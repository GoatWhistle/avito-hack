package app_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/weeklylottery/app"
	"github.com/avito-hack/backend/internal/module/weeklylottery/domain"
	"github.com/avito-hack/backend/internal/shared/domainerr"
)

var fixedNow = time.Date(2026, time.August, 12, 12, 0, 0, 0, time.UTC)

type stubClock struct{ now time.Time }

func (c stubClock) Now() time.Time { return c.now }

type stubTx struct{}

func (stubTx) WithTx(ctx context.Context, fn func(context.Context) error) error { return fn(ctx) }

type stubRepository struct {
	runs      map[string]domain.Run
	lastLocks []bool
}

func newRepository(runs ...domain.Run) *stubRepository {
	repo := &stubRepository{runs: make(map[string]domain.Run)}
	for _, run := range runs {
		repo.runs[run.DisplayID] = run
	}

	return repo
}

func (r *stubRepository) ByUserWeek(
	_ context.Context,
	userID uuid.UUID,
	weekStart time.Time,
	lock bool,
) (domain.Run, error) {
	r.lastLocks = append(r.lastLocks, lock)
	for _, run := range r.runs {
		if run.UserID == userID && run.WeekStart.Equal(weekStart) {
			return run, nil
		}
	}

	return domain.Run{}, domain.ErrRunNotFound
}

func (r *stubRepository) ByDisplayID(_ context.Context, id string, lock bool) (domain.Run, error) {
	r.lastLocks = append(r.lastLocks, lock)
	run, ok := r.runs[id]
	if !ok {
		return domain.Run{}, domain.ErrRunNotFound
	}

	return run, nil
}

func (r *stubRepository) Create(_ context.Context, run domain.Run) (bool, error) {
	for _, existing := range r.runs {
		if existing.UserID == run.UserID && existing.WeekStart.Equal(run.WeekStart) {
			return false, nil
		}
	}
	r.runs[run.DisplayID] = run

	return true, nil
}

func (r *stubRepository) Save(_ context.Context, run domain.Run) error {
	if _, ok := r.runs[run.DisplayID]; !ok {
		return domain.ErrRunNotFound
	}
	r.runs[run.DisplayID] = run

	return nil
}

type fixedRandom struct{ first bool }

func (r *fixedRandom) Intn(limit int) (int, error) {
	if !r.first {
		r.first = true
		return 50 % limit, nil
	}

	return 0, nil
}

type stubSigner struct{}

func (stubSigner) Issue(_ uuid.UUID, rewardID string) (string, error) {
	return "signed-" + rewardID, nil
}

func (stubSigner) Verify(_ uuid.UUID, rewardID, code string) error {
	if code != "signed-"+rewardID {
		return fmt.Errorf("invalid signature")
	}

	return nil
}

func TestStartCreatesOnlyOneRunForCurrentWeek(t *testing.T) {
	t.Parallel()

	repo := newRepository()
	handler := app.NewStartHandler(
		repo, stubTx{}, stubClock{fixedNow}, domain.NewGenerator(&fixedRandom{}),
	)
	userID := uuid.New()

	first, created, err := handler.Handle(context.Background(), userID)
	require.NoError(t, err)
	assert.True(t, created)
	second, created, err := handler.Handle(context.Background(), userID)
	require.NoError(t, err)
	assert.False(t, created)
	assert.Equal(t, first.ID, second.ID)
	assert.Len(t, repo.runs, 1)
	assert.Equal(t, []bool{true, true}, repo.lastLocks)
}

func TestStateDoesNotExposeClosedSymbols(t *testing.T) {
	t.Parallel()

	run := winningRun(uuid.New(), fixedNow)
	_, err := run.Reveal(0, fixedNow)
	require.NoError(t, err)
	handler := app.NewGetStateHandler(newRepository(run), stubClock{fixedNow})

	state, err := handler.Handle(context.Background(), run.UserID)
	require.NoError(t, err)
	require.NotNil(t, state.Run)
	assert.True(t, state.Available)
	assert.Equal(t, domain.SymbolBicycle, state.Run.Slots[0].Symbol)
	for _, slot := range state.Run.Slots[1:] {
		assert.Empty(t, slot.Symbol)
	}
}

func TestStateBecomesUnavailableAfterRunFinishes(t *testing.T) {
	t.Parallel()

	run := winningRun(uuid.New(), fixedNow)
	for _, slot := range []int{0, 1, 2} {
		_, err := run.Reveal(slot, fixedNow)
		require.NoError(t, err)
	}
	run.AttachReward("signed", fixedNow.Add(30*24*time.Hour))
	handler := app.NewGetStateHandler(newRepository(run), stubClock{fixedNow})

	state, err := handler.Handle(context.Background(), run.UserID)
	require.NoError(t, err)
	assert.False(t, state.Available)
	require.NotNil(t, state.Run)
	assert.Equal(t, domain.StateWon, state.Run.State)
}

func TestRevealThirdMatchIssuesAndPersistsReward(t *testing.T) {
	t.Parallel()

	run := winningRun(uuid.New(), fixedNow)
	for _, slot := range []int{0, 1} {
		_, err := run.Reveal(slot, fixedNow)
		require.NoError(t, err)
	}
	repo := newRepository(run)
	handler := app.NewRevealHandler(repo, stubTx{}, stubClock{fixedNow}, stubSigner{})

	result, err := handler.Handle(context.Background(), app.RevealCommand{
		UserID: run.UserID, RunID: run.DisplayID, Slot: 2,
	})
	require.NoError(t, err)
	assert.Equal(t, domain.StateWon, result.Run.State)
	require.NotNil(t, result.Run.Prize)
	assert.Equal(t, "weekly_bicycle_5", result.Run.Prize.ID)
	assert.Equal(t, "signed-weekly_bicycle_5:"+run.DisplayID, result.Run.Prize.Code)
	assert.Equal(t, 5, result.Run.Prize.BenefitValue)
	assert.Equal(t, domain.ScopeTypeCategory, result.Run.Prize.ScopeType)
	assert.Equal(t, domain.CategorySport, result.Run.Prize.ScopeValue)
	assert.Equal(t, true, repo.lastLocks[0])
}

func TestRevealRejectsAnotherUsersAndExpiredRuns(t *testing.T) {
	t.Parallel()

	run := winningRun(uuid.New(), fixedNow)
	repo := newRepository(run)
	handler := app.NewRevealHandler(repo, stubTx{}, stubClock{fixedNow}, stubSigner{})

	_, err := handler.Handle(context.Background(), app.RevealCommand{
		UserID: uuid.New(), RunID: run.DisplayID, Slot: 0,
	})
	require.ErrorIs(t, err, domainerr.ErrNotFound)

	handler = app.NewRevealHandler(repo, stubTx{}, stubClock{fixedNow.AddDate(0, 0, 7)}, stubSigner{})
	_, err = handler.Handle(context.Background(), app.RevealCommand{
		UserID: run.UserID, RunID: run.DisplayID, Slot: 0,
	})
	require.ErrorIs(t, err, domainerr.ErrConflict)
}

func winningRun(userID uuid.UUID, now time.Time) domain.Run {
	board := [domain.BoardSize]domain.Symbol{
		domain.SymbolBicycle, domain.SymbolBicycle, domain.SymbolBicycle,
		domain.SymbolSofa, domain.SymbolSofa, domain.SymbolSneakers,
		domain.SymbolSneakers, domain.SymbolDelivery, domain.SymbolPromotion,
	}
	prize, _ := domain.PrizeByID("weekly_bicycle_5")

	return domain.NewRun(userID, domain.WeekAt(now), board, &prize, now)
}

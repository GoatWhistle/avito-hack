package app_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/games/app"
	"github.com/avito-hack/backend/internal/module/games/domain"
	"github.com/avito-hack/backend/internal/shared/domainerr"
)

func TestStartRoundReturnsFirstPrompt(t *testing.T) {
	t.Parallel()

	rounds := newStubRounds()
	registry := domain.NewRegistry(&scriptedGame{slug: "g", target: 7})

	handler := app.NewStartRoundHandler(registry, rounds, &stubTx{}, stubClock{now: testNow})

	result, err := handler.Handle(context.Background(), app.StartRoundCommand{UserID: uuid.New(), GameSlug: "g"})

	require.NoError(t, err)
	assert.Len(t, result.RoundID, 12)
	assert.Equal(t, 0, result.Streak)
	assert.Equal(t, 7, result.TargetStreak)
	assert.JSONEq(t, `{"question":"first"}`, string(result.Prompt))
	assert.Len(t, rounds.saved, 1)
}

func TestStartRoundAbandonsPreviousActiveRound(t *testing.T) {
	t.Parallel()

	rounds := newStubRounds()
	userID := uuid.New()
	previous := domain.NewRound(userID, "g", testNow)
	rounds.active = previous

	registry := domain.NewRegistry(&scriptedGame{slug: "g", target: 7})
	handler := app.NewStartRoundHandler(registry, rounds, &stubTx{}, stubClock{now: testNow})

	result, err := handler.Handle(context.Background(), app.StartRoundCommand{UserID: userID, GameSlug: "g"})

	require.NoError(t, err)
	assert.Equal(t, domain.StateLost, previous.State())
	assert.NotEqual(t, previous.DisplayID(), result.RoundID)
}

func TestStartRoundUnknownGameIsNotFound(t *testing.T) {
	t.Parallel()

	registry := domain.NewRegistry(&scriptedGame{slug: "g", target: 7})
	handler := app.NewStartRoundHandler(registry, newStubRounds(), &stubTx{}, stubClock{now: testNow})

	_, err := handler.Handle(context.Background(), app.StartRoundCommand{UserID: uuid.New(), GameSlug: "nope"})

	require.ErrorIs(t, err, domain.ErrGameNotFound)
}

func TestListGamesReportsPerGameDailyAndOneGlobalStreak(t *testing.T) {
	t.Parallel()

	progress := newStubProgress()
	userID := uuid.New()

	require.NoError(t, progress.SaveDaily(context.Background(), userID, "a", testDay,
		domain.DailyProgress{Attempts: 2, BestStreak: 7}))
	require.NoError(t, progress.SaveDailyAny(context.Background(), userID, testDay,
		domain.DailyProgress{Attempts: 2, BestStreak: 7}))
	require.NoError(t, progress.SaveStreak(context.Background(), userID,
		domain.Streak{CurrentDays: 7, BestDays: 9}))

	registry := domain.NewRegistry(&scriptedGame{slug: "a", target: 7}, &scriptedGame{slug: "b", target: 7})
	handler := app.NewListGamesHandler(registry, progress)

	view, err := handler.Handle(context.Background(), app.ListGamesQuery{UserID: userID, ClientDay: testDay})

	require.NoError(t, err)
	require.Len(t, view.Games, 2)
	assert.Equal(t, "a", view.Games[0].Slug)
	assert.True(t, view.Games[0].DailyDone)
	assert.False(t, view.Games[1].DailyDone, "the checkmark stays per game")
	assert.True(t, view.DailyDone)
	assert.Equal(t, 7, view.Streak.CurrentDays)
	assert.Equal(t, 9, view.Streak.BestDays)
	assert.True(t, view.Streak.RewardReady)
}

func TestGetStateWithoutActiveRound(t *testing.T) {
	t.Parallel()

	registry := domain.NewRegistry(&scriptedGame{slug: "g", target: 7})
	handler := app.NewGetStateHandler(registry, newStubRounds(), newStubProgress())

	view, err := handler.Handle(context.Background(), app.GetStateQuery{
		UserID: uuid.New(), GameSlug: "g", ClientDay: testDay,
	})

	require.NoError(t, err)
	assert.Nil(t, view.ActiveRound)
	assert.Equal(t, 7, view.TargetStreak)
}

func TestGetStateResumesActiveRound(t *testing.T) {
	t.Parallel()

	rounds := newStubRounds()
	userID := uuid.New()
	round := domain.NewRound(userID, "g", testNow)
	round.Advance(7, testNow)
	rounds.active = round

	registry := domain.NewRegistry(&scriptedGame{slug: "g", target: 7})
	handler := app.NewGetStateHandler(registry, rounds, newStubProgress())

	view, err := handler.Handle(context.Background(), app.GetStateQuery{
		UserID: userID, GameSlug: "g", ClientDay: testDay,
	})

	require.NoError(t, err)
	require.NotNil(t, view.ActiveRound)
	assert.Equal(t, round.DisplayID(), view.ActiveRound.RoundID)
	assert.Equal(t, 1, view.ActiveRound.Streak)
	assert.JSONEq(t, `{"question":"resumed"}`, string(view.ActiveRound.Prompt))
}

func TestClaimRewardIssuesCodeAndResetsCycle(t *testing.T) {
	t.Parallel()

	progress := newStubProgress()
	userID := uuid.New()
	require.NoError(t, progress.SaveStreak(context.Background(), userID,
		domain.Streak{CurrentDays: 7, BestDays: 7}))

	handler := app.NewClaimRewardHandler(progress, &stubTx{}, stubClock{now: testNow})

	result, err := handler.Handle(context.Background(), app.ClaimRewardCommand{UserID: userID})

	require.NoError(t, err)
	assert.Regexp(t, `^PROMO-[0-9A-Z]{8}$`, result.Code)

	streak, err := progress.Streak(context.Background(), userID)
	require.NoError(t, err)
	assert.Equal(t, 0, streak.CurrentDays)
	require.NotNil(t, streak.RewardClaimedAt)
	assert.Equal(t, testNow, streak.RewardClaimedAt.UTC())
}

func TestClaimRewardBeforeSevenDaysConflicts(t *testing.T) {
	t.Parallel()

	progress := newStubProgress()
	userID := uuid.New()
	require.NoError(t, progress.SaveStreak(context.Background(), userID, domain.Streak{CurrentDays: 6}))

	handler := app.NewClaimRewardHandler(progress, &stubTx{}, stubClock{now: testNow})

	_, err := handler.Handle(context.Background(), app.ClaimRewardCommand{UserID: userID})

	require.ErrorIs(t, err, domainerr.ErrConflict)
}

func TestClaimRewardTwiceConflicts(t *testing.T) {
	t.Parallel()

	progress := newStubProgress()
	userID := uuid.New()
	require.NoError(t, progress.SaveStreak(context.Background(), userID, domain.Streak{CurrentDays: 7}))

	handler := app.NewClaimRewardHandler(progress, &stubTx{}, stubClock{now: time.Now().UTC()})

	_, err := handler.Handle(context.Background(), app.ClaimRewardCommand{UserID: userID})
	require.NoError(t, err)

	_, err = handler.Handle(context.Background(), app.ClaimRewardCommand{UserID: userID})
	require.ErrorIs(t, err, domainerr.ErrConflict)
}

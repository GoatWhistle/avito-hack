package app_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/games/app"
	"github.com/avito-hack/backend/internal/module/games/domain"
	"github.com/avito-hack/backend/internal/shared/domainerr"
)

func winOnce(t *testing.T, progress *stubProgress, userID uuid.UUID, slug string, day domain.Day) {
	t.Helper()

	rounds := newStubRounds()
	round := seedRound(t, rounds, userID, slug)
	handler := newGuessHandler(&scriptedGame{slug: slug, target: 1, correct: true}, rounds, progress)

	_, err := handler.Handle(context.Background(), app.GuessCommand{
		UserID: userID, GameSlug: slug, RoundID: round.DisplayID(), Move: move, ClientDay: day,
	})
	require.NoError(t, err)
}

func TestTwoGamesOnTheSameDayAdvanceTheStreakOnce(t *testing.T) {
	t.Parallel()

	progress := newStubProgress()
	userID := uuid.New()

	winOnce(t, progress, userID, "a", testDay)
	winOnce(t, progress, userID, "b", testDay)

	streak, err := progress.Streak(context.Background(), userID)
	require.NoError(t, err)
	assert.Equal(t, 1, streak.CurrentDays,
		"playing two different games on one day must count as a single streak day")
}

func TestConsecutiveDaysInDifferentGamesAdvanceTheSameStreak(t *testing.T) {
	t.Parallel()

	progress := newStubProgress()
	userID := uuid.New()

	winOnce(t, progress, userID, "a", testDay)
	winOnce(t, progress, userID, "b", testDay.AddDays(1))

	streak, err := progress.Streak(context.Background(), userID)
	require.NoError(t, err)
	assert.Equal(t, 2, streak.CurrentDays, "the streak spans every mini-game, not one of them")
}

func TestStateReportsTheSameStreakWhicheverGameIsQueried(t *testing.T) {
	t.Parallel()

	progress := newStubProgress()
	userID := uuid.New()

	winOnce(t, progress, userID, "a", testDay)

	registry := domain.NewRegistry(&scriptedGame{slug: "a", target: 7}, &scriptedGame{slug: "b", target: 7})
	handler := app.NewGetStateHandler(registry, newStubRounds(), progress, newStubScores())

	first, err := handler.Handle(context.Background(), app.GetStateQuery{
		UserID: userID, GameSlug: "a", ClientDay: testDay,
	})
	require.NoError(t, err)

	second, err := handler.Handle(context.Background(), app.GetStateQuery{
		UserID: userID, GameSlug: "b", ClientDay: testDay,
	})
	require.NoError(t, err)

	assert.Equal(t, first.Streak, second.Streak, "the streak must not depend on which game asks for it")
	assert.Equal(t, 1, second.Streak.CurrentDays)
}

func TestClaimMintsOneCodePerCycleAcrossGames(t *testing.T) {
	t.Parallel()

	progress := newStubProgress()
	userID := uuid.New()

	day := testDay
	require.NoError(t, progress.SaveStreak(context.Background(), userID, domain.Streak{
		CurrentDays: domain.RewardTargetDays, BestDays: domain.RewardTargetDays, LastDay: &day,
	}))

	handler := app.NewClaimRewardHandler(progress, &stubTx{}, stubClock{now: testNow})

	first, err := handler.Handle(context.Background(), app.ClaimRewardCommand{UserID: userID})
	require.NoError(t, err)
	assert.NotEmpty(t, first.Code)

	_, err = handler.Handle(context.Background(), app.ClaimRewardCommand{UserID: userID})
	require.ErrorIs(t, err, domainerr.ErrConflict,
		"the cycle is global, so a second claim must not mint another code from another game")

	assert.Equal(t, []string{userID.String(), userID.String()}, progress.lockedStreaks,
		"every claim must take the same single streak row lock")
}

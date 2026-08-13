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

func TestGuessLocksTheRoundRowBeforeReadingIt(t *testing.T) {
	t.Parallel()

	rounds := newStubRounds()
	userID := uuid.New()
	round := seedRound(t, rounds, userID, "g")

	handler := newGuessHandler(&scriptedGame{slug: "g", target: 7, correct: true}, rounds, newStubProgress())

	_, err := handler.Handle(context.Background(), app.GuessCommand{
		UserID: userID, GameSlug: "g", RoundID: round.DisplayID(), Move: move, ClientDay: testDay,
	})

	require.NoError(t, err)
	assert.Equal(t, []string{round.DisplayID()}, rounds.lockedRounds,
		"a guess must take a row lock, otherwise concurrent replays each read the same streak and all advance it")
}

func TestStartRoundSerialisesConcurrentStarts(t *testing.T) {
	t.Parallel()

	rounds := newStubRounds()
	userID := uuid.New()

	handler := app.NewStartRoundHandler(
		domain.NewRegistry(&scriptedGame{slug: "g", target: 7}), rounds, &stubTx{}, stubClock{now: testNow},
	)

	_, err := handler.Handle(context.Background(), app.StartRoundCommand{UserID: userID, GameSlug: "g"})

	require.NoError(t, err)
	assert.Equal(t, []string{userID.String() + "/g"}, rounds.lockedGames,
		"starting a round must serialise per user and game, otherwise parallel starts leave several rounds active")
}

func TestWinningGuessLocksTheStreakRow(t *testing.T) {
	t.Parallel()

	rounds := newStubRounds()
	progress := newStubProgress()
	userID := uuid.New()
	round := seedRound(t, rounds, userID, "g")

	handler := newGuessHandler(&scriptedGame{slug: "g", target: 1, correct: true}, rounds, progress)

	_, err := handler.Handle(context.Background(), app.GuessCommand{
		UserID: userID, GameSlug: "g", RoundID: round.DisplayID(), Move: move, ClientDay: testDay,
	})

	require.NoError(t, err)
	assert.Equal(t, []string{userID.String()}, progress.lockedStreaks,
		"completing a daily attempt must lock the weekly streak row it then overwrites")
}

func TestClaimRewardLocksTheStreakRow(t *testing.T) {
	t.Parallel()

	progress := newStubProgress()
	userID := uuid.New()
	day := testDay
	progress.streaks[userID] = domain.Streak{
		CurrentDays: domain.RewardTargetDays, BestDays: domain.RewardTargetDays, LastDay: &day,
	}

	handler := app.NewClaimRewardHandler(progress, &stubTx{}, stubClock{now: testNow})

	result, err := handler.Handle(context.Background(), app.ClaimRewardCommand{UserID: userID})

	require.NoError(t, err)
	assert.NotEmpty(t, result.Code)
	assert.Equal(t, []string{userID.String()}, progress.lockedStreaks,
		"claiming must lock the streak row, otherwise concurrent claims mint several codes from one cycle")
}

func TestClaimRewardIsSingleUsePerCycle(t *testing.T) {
	t.Parallel()

	progress := newStubProgress()
	userID := uuid.New()
	day := testDay
	progress.streaks[userID] = domain.Streak{
		CurrentDays: domain.RewardTargetDays, BestDays: domain.RewardTargetDays, LastDay: &day,
	}

	handler := app.NewClaimRewardHandler(progress, &stubTx{}, stubClock{now: testNow})

	cmd := app.ClaimRewardCommand{UserID: userID}

	_, err := handler.Handle(context.Background(), cmd)
	require.NoError(t, err)

	_, err = handler.Handle(context.Background(), cmd)
	require.ErrorIs(t, err, domainerr.ErrConflict, "a replayed claim must not mint a second promo code")
}

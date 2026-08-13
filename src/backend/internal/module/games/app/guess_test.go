package app_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/games/app"
	"github.com/avito-hack/backend/internal/module/games/domain"
	"github.com/avito-hack/backend/internal/shared/domainerr"
)

var (
	testNow = time.Date(2026, time.March, 14, 12, 0, 0, 0, time.UTC)
	testDay = domain.DayOf(testNow)
	move    = json.RawMessage(`{"choice":"higher"}`)
)

func newGuessHandler(
	game domain.Game,
	rounds *stubRounds,
	progress *stubProgress,
) *app.GuessHandler {
	return app.NewGuessHandler(
		domain.NewRegistry(game), rounds, progress, &stubTx{}, stubClock{now: testNow},
	)
}

func seedRound(t *testing.T, rounds *stubRounds, userID uuid.UUID, slug string) *domain.Round {
	t.Helper()

	round := domain.NewRound(userID, slug, testNow)
	require.NoError(t, rounds.Save(context.Background(), round))

	return round
}

func TestGuessUnknownRoundIsNotFound(t *testing.T) {
	t.Parallel()

	handler := newGuessHandler(&scriptedGame{slug: "g", target: 7, correct: true}, newStubRounds(), newStubProgress())

	_, err := handler.Handle(context.Background(), app.GuessCommand{
		UserID: uuid.New(), GameSlug: "g", RoundID: "aaaaaaaaaaaa", Move: move, ClientDay: testDay,
	})

	require.ErrorIs(t, err, domainerr.ErrNotFound)
}

func TestGuessRoundOfAnotherUserIsNotFound(t *testing.T) {
	t.Parallel()

	rounds := newStubRounds()
	round := seedRound(t, rounds, uuid.New(), "g")

	handler := newGuessHandler(&scriptedGame{slug: "g", target: 7, correct: true}, rounds, newStubProgress())

	_, err := handler.Handle(context.Background(), app.GuessCommand{
		UserID: uuid.New(), GameSlug: "g", RoundID: round.DisplayID(), Move: move, ClientDay: testDay,
	})

	require.ErrorIs(t, err, domainerr.ErrNotFound)
}

func TestGuessRoundOfAnotherGameIsNotFound(t *testing.T) {
	t.Parallel()

	rounds := newStubRounds()
	userID := uuid.New()
	round := seedRound(t, rounds, userID, "other")

	handler := newGuessHandler(&scriptedGame{slug: "g", target: 7, correct: true}, rounds, newStubProgress())

	_, err := handler.Handle(context.Background(), app.GuessCommand{
		UserID: userID, GameSlug: "g", RoundID: round.DisplayID(), Move: move, ClientDay: testDay,
	})

	require.ErrorIs(t, err, domainerr.ErrNotFound)
}

func TestGuessOnFinishedRoundConflicts(t *testing.T) {
	t.Parallel()

	for _, finish := range []struct {
		name   string
		finish func(*domain.Round)
	}{
		{"lost", func(r *domain.Round) { r.Lose(testNow) }},
		{"won", func(r *domain.Round) { r.Advance(1, testNow) }},
	} {
		t.Run(finish.name, func(t *testing.T) {
			t.Parallel()

			rounds := newStubRounds()
			userID := uuid.New()
			round := seedRound(t, rounds, userID, "g")
			finish.finish(round)

			handler := newGuessHandler(&scriptedGame{slug: "g", target: 7, correct: true}, rounds, newStubProgress())

			_, err := handler.Handle(context.Background(), app.GuessCommand{
				UserID: userID, GameSlug: "g", RoundID: round.DisplayID(), Move: move, ClientDay: testDay,
			})

			require.ErrorIs(t, err, domainerr.ErrConflict)
		})
	}
}

func TestGuessReplayAfterLossConflicts(t *testing.T) {
	t.Parallel()

	rounds := newStubRounds()
	progress := newStubProgress()
	userID := uuid.New()
	round := seedRound(t, rounds, userID, "g")

	handler := newGuessHandler(&scriptedGame{slug: "g", target: 7}, rounds, progress)

	cmd := app.GuessCommand{
		UserID: userID, GameSlug: "g", RoundID: round.DisplayID(), Move: move, ClientDay: testDay,
	}

	result, err := handler.Handle(context.Background(), cmd)
	require.NoError(t, err)
	assert.False(t, result.Correct)
	assert.Equal(t, domain.StateLost, result.State)

	_, err = handler.Handle(context.Background(), cmd)
	require.ErrorIs(t, err, domainerr.ErrConflict)
}

func TestGuessCorrectAdvancesStreakAndReturnsNextPrompt(t *testing.T) {
	t.Parallel()

	rounds := newStubRounds()
	userID := uuid.New()
	round := seedRound(t, rounds, userID, "g")

	handler := newGuessHandler(&scriptedGame{slug: "g", target: 7, correct: true}, rounds, newStubProgress())

	result, err := handler.Handle(context.Background(), app.GuessCommand{
		UserID: userID, GameSlug: "g", RoundID: round.DisplayID(), Move: move, ClientDay: testDay,
	})

	require.NoError(t, err)
	assert.True(t, result.Correct)
	assert.Equal(t, 1, result.Streak)
	assert.Equal(t, domain.StateActive, result.State)
	assert.JSONEq(t, `{"question":"next"}`, string(result.Prompt))
	assert.False(t, result.AttemptCompleted)
	assert.Nil(t, result.StreakAfter)
}

func TestGuessWinningCompletesAttemptAndStartsWeeklyStreak(t *testing.T) {
	t.Parallel()

	rounds := newStubRounds()
	progress := newStubProgress()
	userID := uuid.New()
	round := seedRound(t, rounds, userID, "g")

	handler := newGuessHandler(&scriptedGame{slug: "g", target: 1, correct: true}, rounds, progress)

	result, err := handler.Handle(context.Background(), app.GuessCommand{
		UserID: userID, GameSlug: "g", RoundID: round.DisplayID(), Move: move, ClientDay: testDay,
	})

	require.NoError(t, err)
	assert.Equal(t, domain.StateWon, result.State)
	assert.True(t, result.AttemptCompleted)
	assert.Empty(t, result.Prompt, "a won round must not hand out another question")

	require.NotNil(t, result.StreakAfter)
	assert.Equal(t, 1, result.StreakAfter.CurrentDays)
	assert.False(t, result.StreakAfter.RewardReady)

	daily, err := progress.Daily(context.Background(), userID, "g", testDay)
	require.NoError(t, err)
	assert.Equal(t, 1, daily.Attempts)
}

func TestSecondAttemptSameDayDoesNotAdvanceWeeklyStreak(t *testing.T) {
	t.Parallel()

	rounds := newStubRounds()
	progress := newStubProgress()
	userID := uuid.New()

	handler := newGuessHandler(&scriptedGame{slug: "g", target: 1, correct: true}, rounds, progress)

	for i := 0; i < 2; i++ {
		round := seedRound(t, rounds, userID, "g")

		_, err := handler.Handle(context.Background(), app.GuessCommand{
			UserID: userID, GameSlug: "g", RoundID: round.DisplayID(), Move: move, ClientDay: testDay,
		})
		require.NoError(t, err)
	}

	streak, err := progress.Streak(context.Background(), userID)
	require.NoError(t, err)
	assert.Equal(t, 1, streak.CurrentDays, "a second attempt on the same day must not double count")

	daily, err := progress.Daily(context.Background(), userID, "g", testDay)
	require.NoError(t, err)
	assert.Equal(t, 2, daily.Attempts)
}

func TestGuessContinueKeepsRoundActiveAndStreakUnchanged(t *testing.T) {
	t.Parallel()

	rounds := newStubRounds()
	progress := newStubProgress()
	userID := uuid.New()
	round := seedRound(t, rounds, userID, "g")

	game := &attemptGame{
		scriptedGame: scriptedGame{slug: "g", progress: domain.ProgressContinue, next: true},
		maxAttempts:  6,
		attemptsUsed: 1,
	}

	result, err := newGuessHandler(game, rounds, progress).Handle(context.Background(), app.GuessCommand{
		UserID: userID, GameSlug: "g", RoundID: round.DisplayID(), Move: move, ClientDay: testDay,
	})

	require.NoError(t, err)
	assert.False(t, result.Correct)
	assert.Equal(t, domain.ProgressContinue, result.Progress)
	assert.Equal(t, domain.StateActive, result.State)
	assert.Equal(t, 0, result.Streak, "a continue must never move the streak")
	assert.Equal(t, 0, round.BestStreak())
	assert.False(t, result.AttemptCompleted)
	assert.Nil(t, result.StreakAfter)
	assert.JSONEq(t, `{"question":"next"}`, string(result.Prompt))
	assert.Equal(t, 1, result.AttemptsUsed)
	assert.Equal(t, 6, result.MaxAttempts)
}

func TestGuessWinCompletesAttemptWithoutTargetStreak(t *testing.T) {
	t.Parallel()

	rounds := newStubRounds()
	progress := newStubProgress()
	userID := uuid.New()
	round := seedRound(t, rounds, userID, "g")

	game := &attemptGame{
		scriptedGame: scriptedGame{slug: "g", correct: true, progress: domain.ProgressWin},
		maxAttempts:  6,
		attemptsUsed: 3,
	}

	result, err := newGuessHandler(game, rounds, progress).Handle(context.Background(), app.GuessCommand{
		UserID: userID, GameSlug: "g", RoundID: round.DisplayID(), Move: move, ClientDay: testDay,
	})

	require.NoError(t, err)
	assert.Equal(t, domain.ProgressWin, result.Progress)
	assert.Equal(t, domain.StateWon, result.State)
	assert.True(t, result.AttemptCompleted)
	assert.Equal(t, 0, result.Streak, "winning without a streak target must leave the streak alone")
	assert.Equal(t, 3, result.AttemptsUsed)
	assert.Empty(t, result.Prompt, "a won round must not hand out another question")

	require.NotNil(t, result.StreakAfter)
	assert.Equal(t, 1, result.StreakAfter.CurrentDays)

	daily, err := progress.Daily(context.Background(), userID, "g", testDay)
	require.NoError(t, err)
	assert.Equal(t, 1, daily.Attempts)
}

func TestGuessLoseDoesNotCompleteAttempt(t *testing.T) {
	t.Parallel()

	rounds := newStubRounds()
	progress := newStubProgress()
	userID := uuid.New()
	round := seedRound(t, rounds, userID, "g")

	game := &attemptGame{
		scriptedGame: scriptedGame{slug: "g", progress: domain.ProgressLose},
		maxAttempts:  6,
		attemptsUsed: 6,
	}

	result, err := newGuessHandler(game, rounds, progress).Handle(context.Background(), app.GuessCommand{
		UserID: userID, GameSlug: "g", RoundID: round.DisplayID(), Move: move, ClientDay: testDay,
	})

	require.NoError(t, err)
	assert.Equal(t, domain.ProgressLose, result.Progress)
	assert.Equal(t, domain.StateLost, result.State)
	assert.False(t, result.AttemptCompleted)
	assert.Equal(t, 6, result.AttemptsUsed)

	daily, err := progress.Daily(context.Background(), userID, "g", testDay)
	require.NoError(t, err)
	assert.Equal(t, 0, daily.Attempts)
}

func TestGuessAttemptFieldsAreZeroForStreakGames(t *testing.T) {
	t.Parallel()

	rounds := newStubRounds()
	userID := uuid.New()
	round := seedRound(t, rounds, userID, "g")

	handler := newGuessHandler(
		&scriptedGame{slug: "g", target: 7, correct: true, progress: domain.ProgressAdvance},
		rounds, newStubProgress(),
	)

	result, err := handler.Handle(context.Background(), app.GuessCommand{
		UserID: userID, GameSlug: "g", RoundID: round.DisplayID(), Move: move, ClientDay: testDay,
	})

	require.NoError(t, err)
	assert.Equal(t, 1, result.Streak)
	assert.Equal(t, 0, result.AttemptsUsed)
	assert.Equal(t, 0, result.MaxAttempts)
}

func TestGuessUnknownGameIsNotFound(t *testing.T) {
	t.Parallel()

	handler := newGuessHandler(&scriptedGame{slug: "g", target: 7}, newStubRounds(), newStubProgress())

	_, err := handler.Handle(context.Background(), app.GuessCommand{
		UserID: uuid.New(), GameSlug: "missing", RoundID: "aaaaaaaaaaaa", Move: move, ClientDay: testDay,
	})

	require.ErrorIs(t, err, domain.ErrGameNotFound)
}

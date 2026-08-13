package app_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/games/app"
	"github.com/avito-hack/backend/internal/module/games/domain"
	"github.com/avito-hack/backend/internal/shared/domainerr"
)

type scoredGame struct {
	scriptedGame

	score       int
	minForStrek int
}

func (g *scoredGame) RoundScore(_ *domain.Round) int { return g.score }

func (g *scoredGame) CountsTowardStreak(score int) bool { return score >= g.minForStrek }

func newScoredGame(score int) *scoredGame {
	return &scoredGame{
		scriptedGame: scriptedGame{slug: "raccoonjump", target: 0, correct: true, progress: domain.ProgressWin},
		score:        score,
		minForStrek:  25,
	}
}

func submitScore(
	t *testing.T,
	game domain.Game,
	rounds *stubRounds,
	progress *stubProgress,
	scores *stubScores,
	userID uuid.UUID,
	round *domain.Round,
) app.GuessResult {
	t.Helper()

	handler := app.NewGuessHandler(
		domain.NewRegistry(game), rounds, progress, scores, &stubTx{}, stubClock{now: testNow},
	)

	result, err := handler.Handle(context.Background(), app.GuessCommand{
		UserID:    userID,
		GameSlug:  "raccoonjump",
		RoundID:   round.DisplayID(),
		Move:      json.RawMessage(`{"score":400}`),
		ClientDay: testDay,
	})
	require.NoError(t, err)

	return result
}

func TestRaccoonJumpPersistsBestScore(t *testing.T) {
	t.Parallel()

	rounds := newStubRounds()
	scores := newStubScores()
	progress := newStubProgress()
	userID := uuid.New()
	round := seedRound(t, rounds, userID, "raccoonjump")

	result := submitScore(t, newScoredGame(400), rounds, progress, scores, userID, round)

	assert.Equal(t, 400, result.BestScore)

	best, err := scores.BestScore(context.Background(), userID, "raccoonjump")
	require.NoError(t, err)
	assert.Equal(t, 400, best)
}

func TestRaccoonJumpKeepsHighestBestScore(t *testing.T) {
	t.Parallel()

	rounds := newStubRounds()
	scores := newStubScores()
	progress := newStubProgress()
	userID := uuid.New()

	first := seedRound(t, rounds, userID, "raccoonjump")
	submitScore(t, newScoredGame(900), rounds, progress, scores, userID, first)

	second := seedRound(t, rounds, userID, "raccoonjump")
	result := submitScore(t, newScoredGame(120), rounds, progress, scores, userID, second)

	assert.Equal(t, 900, result.BestScore, "a weaker run must not lower the stored best score")
}

type revealingScoredGame struct {
	scoredGame
}

func (g *revealingScoredGame) Guess(
	_ context.Context,
	_ *domain.Round,
	_ json.RawMessage,
) (domain.GuessOutcome, error) {
	reveal := fmt.Sprintf(`{"score":%d,"best_score":%d,"new_best":true}`, g.score, g.score)

	return domain.GuessOutcome{
		Correct:  true,
		Progress: domain.ProgressWin,
		Reveal:   json.RawMessage(reveal),
	}, nil
}

func newRevealingScoredGame(score int) *revealingScoredGame {
	return &revealingScoredGame{scoredGame: *newScoredGame(score)}
}

func TestRaccoonJumpRevealCarriesPersistedBestScore(t *testing.T) {
	t.Parallel()

	rounds := newStubRounds()
	scores := newStubScores()
	progress := newStubProgress()
	userID := uuid.New()

	first := seedRound(t, rounds, userID, "raccoonjump")
	submitScore(t, newRevealingScoredGame(900), rounds, progress, scores, userID, first)

	second := seedRound(t, rounds, userID, "raccoonjump")
	result := submitScore(t, newRevealingScoredGame(120), rounds, progress, scores, userID, second)

	var reveal struct {
		Score     int  `json:"score"`
		BestScore int  `json:"best_score"`
		NewBest   bool `json:"new_best"`
	}
	require.NoError(t, json.Unmarshal(result.Reveal, &reveal))

	assert.Equal(t, 120, reveal.Score)
	assert.Equal(t, 900, reveal.BestScore, "reveal must report the persisted best, not the round score")
	assert.False(t, reveal.NewBest, "a weaker run must not be announced as a new record")
}

func TestRaccoonJumpLowScoreDoesNotAdvanceStreak(t *testing.T) {
	t.Parallel()

	rounds := newStubRounds()
	scores := newStubScores()
	progress := newStubProgress()
	userID := uuid.New()
	round := seedRound(t, rounds, userID, "raccoonjump")

	result := submitScore(t, newScoredGame(3), rounds, progress, scores, userID, round)

	assert.False(t, result.AttemptCompleted, "an instant death must not farm the global streak")
	assert.Nil(t, result.StreakAfter)

	streak, err := progress.Streak(context.Background(), userID)
	require.NoError(t, err)
	assert.Equal(t, 0, streak.CurrentDays)
}

func TestRaccoonJumpQualifyingScoreAdvancesGlobalStreak(t *testing.T) {
	t.Parallel()

	rounds := newStubRounds()
	scores := newStubScores()
	progress := newStubProgress()
	userID := uuid.New()
	round := seedRound(t, rounds, userID, "raccoonjump")

	result := submitScore(t, newScoredGame(400), rounds, progress, scores, userID, round)

	require.NotNil(t, result.StreakAfter)
	assert.True(t, result.AttemptCompleted)
	assert.Equal(t, 1, result.StreakAfter.CurrentDays)

	streak, err := progress.Streak(context.Background(), userID)
	require.NoError(t, err)
	assert.Equal(t, 1, streak.CurrentDays)
}

func TestRaccoonJumpLowScoreStillRecordsBestScore(t *testing.T) {
	t.Parallel()

	rounds := newStubRounds()
	scores := newStubScores()
	progress := newStubProgress()
	userID := uuid.New()
	round := seedRound(t, rounds, userID, "raccoonjump")

	result := submitScore(t, newScoredGame(10), rounds, progress, scores, userID, round)

	assert.Equal(t, 10, result.BestScore,
		"a sub-threshold run must still count for the personal best")
}

func TestRaccoonJumpDuplicateSubmissionIsRejected(t *testing.T) {
	t.Parallel()

	rounds := newStubRounds()
	scores := newStubScores()
	progress := newStubProgress()
	userID := uuid.New()
	round := seedRound(t, rounds, userID, "raccoonjump")

	game := newScoredGame(400)

	submitScore(t, game, rounds, progress, scores, userID, round)

	handler := app.NewGuessHandler(
		domain.NewRegistry(game), rounds, progress, scores, &stubTx{}, stubClock{now: testNow},
	)

	_, err := handler.Handle(context.Background(), app.GuessCommand{
		UserID:    userID,
		GameSlug:  "raccoonjump",
		RoundID:   round.DisplayID(),
		Move:      json.RawMessage(`{"score":400}`),
		ClientDay: testDay,
	})

	require.ErrorIs(t, err, domainerr.ErrConflict,
		"a replayed submission for a finished round must be refused")
	assert.Len(t, scores.saves, 1, "the score must not be counted twice")
}

func TestRaccoonJumpScoreOfAnotherUserIsNotFound(t *testing.T) {
	t.Parallel()

	rounds := newStubRounds()
	scores := newStubScores()
	round := seedRound(t, rounds, uuid.New(), "raccoonjump")

	handler := app.NewGuessHandler(
		domain.NewRegistry(newScoredGame(400)), rounds, newStubProgress(), scores,
		&stubTx{}, stubClock{now: testNow},
	)

	_, err := handler.Handle(context.Background(), app.GuessCommand{
		UserID:    uuid.New(),
		GameSlug:  "raccoonjump",
		RoundID:   round.DisplayID(),
		Move:      json.RawMessage(`{"score":400}`),
		ClientDay: testDay,
	})

	require.ErrorIs(t, err, domainerr.ErrNotFound)
	assert.Empty(t, scores.saves, "a foreign round must never write a score")
}

func TestRaccoonJumpStartRoundClosesPreviousActiveRound(t *testing.T) {
	t.Parallel()

	rounds := newStubRounds()
	userID := uuid.New()

	previous := domain.NewRound(userID, "raccoonjump", testNow)
	rounds.active = previous

	handler := app.NewStartRoundHandler(
		domain.NewRegistry(newScoredGame(0)), rounds, &stubTx{}, stubClock{now: testNow},
	)

	_, err := handler.Handle(context.Background(), app.StartRoundCommand{
		UserID: userID, GameSlug: "raccoonjump",
	})
	require.NoError(t, err)

	assert.Equal(t, domain.StateLost, previous.State(),
		"only one raccoonjump round may be active per user")
	assert.Contains(t, rounds.lockedGames, userID.String()+"/raccoonjump",
		"the start must take the per-user game lock")
}

func TestRaccoonJumpStateExposesBestScore(t *testing.T) {
	t.Parallel()

	scores := newStubScores()
	userID := uuid.New()

	_, err := scores.SaveBestScore(context.Background(), userID, "raccoonjump", 777)
	require.NoError(t, err)

	handler := app.NewGetStateHandler(
		domain.NewRegistry(newScoredGame(0)), newStubRounds(), newStubProgress(), scores,
	)

	view, err := handler.Handle(context.Background(), app.GetStateQuery{
		UserID: userID, GameSlug: "raccoonjump", ClientDay: testDay,
	})
	require.NoError(t, err)

	assert.Equal(t, 777, view.BestScore)
}

func TestNonScoredGameDoesNotTouchScoreRepository(t *testing.T) {
	t.Parallel()

	rounds := newStubRounds()
	scores := newStubScores()
	userID := uuid.New()
	round := seedRound(t, rounds, userID, "g")

	handler := app.NewGuessHandler(
		domain.NewRegistry(&scriptedGame{slug: "g", target: 1, correct: true}),
		rounds, newStubProgress(), scores, &stubTx{}, stubClock{now: testNow},
	)

	_, err := handler.Handle(context.Background(), app.GuessCommand{
		UserID: userID, GameSlug: "g", RoundID: round.DisplayID(), Move: move, ClientDay: testDay,
	})
	require.NoError(t, err)

	assert.Empty(t, scores.saves)
}

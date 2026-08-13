package raccoonjump_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/games/domain"
	"github.com/avito-hack/backend/internal/module/games/domain/raccoonjump"
	"github.com/avito-hack/backend/internal/shared/domainerr"
)

var testNow = time.Date(2026, time.March, 14, 12, 0, 0, 0, time.UTC)

type stubClock struct {
	now time.Time
}

func (c *stubClock) Now() time.Time { return c.now }

func (c *stubClock) advance(d time.Duration) { c.now = c.now.Add(d) }

type stubListings struct {
	listings []raccoonjump.Listing
	err      error
}

func (s *stubListings) RandomListings(_ context.Context, limit int) ([]raccoonjump.Listing, error) {
	if s.err != nil {
		return nil, s.err
	}

	if len(s.listings) > limit {
		return s.listings[:limit], nil
	}

	return s.listings, nil
}

func newListings(n int) *stubListings {
	out := make([]raccoonjump.Listing, 0, n)

	for i := range n {
		out = append(out, raccoonjump.Listing{
			DisplayID:   fmt.Sprintf("item%02d", i),
			Title:       fmt.Sprintf("Товар %d", i),
			PriceKopeks: int64(1000 * (i + 1)),
			PhotoURL:    fmt.Sprintf("https://cdn.test/%d.jpg", i),
		})
	}

	return &stubListings{listings: out}
}

func newGame(clock *stubClock, listings raccoonjump.ListingPool) *raccoonjump.Game {
	return raccoonjump.New(listings, clock)
}

func startRound(t *testing.T, game *raccoonjump.Game) (*domain.Round, json.RawMessage) {
	t.Helper()

	round := domain.NewRound(uuid.New(), raccoonjump.Slug, testNow)

	view, err := game.Start(context.Background(), round)
	require.NoError(t, err)

	return round, view.Prompt
}

func submit(game *raccoonjump.Game, round *domain.Round, score int, collected ...int) (domain.GuessOutcome, error) {
	if collected == nil {
		collected = []int{}
	}

	move, _ := json.Marshal(map[string]any{"score": score, "collected": collected})

	return game.Guess(context.Background(), round, move)
}

func decodeReveal(t *testing.T, raw json.RawMessage) map[string]any {
	t.Helper()

	var out map[string]any
	require.NoError(t, json.Unmarshal(raw, &out))

	return out
}

func TestSlugAndTargetStreak(t *testing.T) {
	t.Parallel()

	game := newGame(&stubClock{now: testNow}, newListings(4))

	assert.Equal(t, "raccoonjump", game.Slug())
	assert.Equal(t, 0, game.TargetStreak())
}

func TestStartIssuesSeedAndCollectibles(t *testing.T) {
	t.Parallel()

	game := newGame(&stubClock{now: testNow}, newListings(4))

	_, prompt := startRound(t, game)

	var payload struct {
		Seed         uint32 `json:"seed"`
		Collectibles []struct {
			Index    int    `json:"index"`
			Title    string `json:"title"`
			PhotoURL string `json:"photo_url"`
		} `json:"collectibles"`
	}

	require.NoError(t, json.Unmarshal(prompt, &payload))

	assert.NotZero(t, payload.Seed, "server must issue a non-zero seed")
	assert.Len(t, payload.Collectibles, 4)
}

func TestPromptNeverLeaksInternalIdentifiers(t *testing.T) {
	t.Parallel()

	game := newGame(&stubClock{now: testNow}, newListings(4))

	_, prompt := startRound(t, game)

	assert.NotContains(t, string(prompt), "display_id",
		"the prompt must not reveal listing identity before collection")
	assert.NotContains(t, string(prompt), "price_kopeks")
}

func TestSeedsDifferBetweenRounds(t *testing.T) {
	t.Parallel()

	game := newGame(&stubClock{now: testNow}, newListings(1))

	seeds := make(map[uint32]bool)

	for range 16 {
		_, prompt := startRound(t, game)

		var payload struct {
			Seed uint32 `json:"seed"`
		}
		require.NoError(t, json.Unmarshal(prompt, &payload))

		seeds[payload.Seed] = true
	}

	assert.Greater(t, len(seeds), 1, "seeds must not be constant across rounds")
}

func TestScoreWithinPlausibleRateIsAccepted(t *testing.T) {
	t.Parallel()

	clock := &stubClock{now: testNow}
	game := newGame(clock, newListings(4))

	round, _ := startRound(t, game)

	clock.advance(30 * time.Second)

	outcome, err := submit(game, round, 900)
	require.NoError(t, err)

	assert.Equal(t, domain.ProgressWin, outcome.Progress)
	assert.EqualValues(t, 900, decodeReveal(t, outcome.Reveal)["score"])
}

func TestImpossibleScoreForElapsedTimeIsRejected(t *testing.T) {
	t.Parallel()

	clock := &stubClock{now: testNow}
	game := newGame(clock, newListings(4))

	round, _ := startRound(t, game)

	clock.advance(2 * time.Second)

	_, err := submit(game, round, 50000)

	require.ErrorIs(t, err, raccoonjump.ErrImplausibleScore)
}

func TestInstantHugeScoreIsRejected(t *testing.T) {
	t.Parallel()

	clock := &stubClock{now: testNow}
	game := newGame(clock, newListings(4))

	round, _ := startRound(t, game)

	_, err := submit(game, round, 9999)

	require.ErrorIs(t, err, raccoonjump.ErrImplausibleScore)
}

func TestHonestFastPlayerIsNeverRejected(t *testing.T) {
	t.Parallel()

	const theoreticalMaxScorePerSecond = 63

	for _, seconds := range []int{5, 15, 30, 60, 120, 300} {
		clock := &stubClock{now: testNow}
		game := newGame(clock, newListings(4))

		round, _ := startRound(t, game)

		clock.advance(time.Duration(seconds) * time.Second)

		best := seconds * theoreticalMaxScorePerSecond

		_, err := submit(game, round, best)
		require.NoErrorf(t, err, "a physically-perfect run of %ds (score %d) must be accepted", seconds, best)
	}
}

func TestNegativeScoreIsRejected(t *testing.T) {
	t.Parallel()

	clock := &stubClock{now: testNow}
	game := newGame(clock, newListings(4))

	round, _ := startRound(t, game)

	_, err := submit(game, round, -1)

	require.ErrorAs(t, err, new(*domainerr.InvalidError))
}

func TestNonNumericAndMissingScoreAreRejected(t *testing.T) {
	t.Parallel()

	cases := map[string]string{
		"missing score": `{"collected":[]}`,
		"nan score":     `{"score":"NaN"}`,
		"null score":    `{"score":null}`,
		"float string":  `{"score":"120"}`,
		"empty move":    ``,
	}

	for name, move := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			clock := &stubClock{now: testNow}
			game := newGame(clock, newListings(4))

			round, _ := startRound(t, game)

			_, err := game.Guess(context.Background(), round, json.RawMessage(move))

			require.ErrorAs(t, err, new(*domainerr.InvalidError))
		})
	}
}

func TestAbsurdScoreIsCapped(t *testing.T) {
	t.Parallel()

	clock := &stubClock{now: testNow}
	game := newGame(clock, newListings(4))

	round, _ := startRound(t, game)

	clock.advance(365 * 24 * time.Hour)

	_, err := submit(game, round, 1<<30)

	require.ErrorIs(t, err, raccoonjump.ErrImplausibleScore,
		"a score above the absolute cap must be rejected even after a long round")
}

func TestDuplicateSubmissionIsIdempotent(t *testing.T) {
	t.Parallel()

	clock := &stubClock{now: testNow}
	game := newGame(clock, newListings(4))

	round, _ := startRound(t, game)

	clock.advance(30 * time.Second)

	first, err := submit(game, round, 400, 30)
	require.NoError(t, err)

	replay, err := submit(game, round, 5000, 30, 70, 110)
	require.NoError(t, err)

	assert.EqualValues(t, 400, decodeReveal(t, first.Reveal)["score"])
	assert.EqualValues(t, 400, decodeReveal(t, replay.Reveal)["score"],
		"a replayed submission must not overwrite the recorded score")
	assert.Equal(t, 400, game.RoundScore(round))
}

func TestReplayDoesNotInflateBestScore(t *testing.T) {
	t.Parallel()

	clock := &stubClock{now: testNow}
	game := newGame(clock, newListings(4))

	round, _ := startRound(t, game)

	clock.advance(30 * time.Second)

	_, err := submit(game, round, 400)
	require.NoError(t, err)

	replay, err := submit(game, round, 9000)
	require.NoError(t, err)

	assert.EqualValues(t, 400, decodeReveal(t, replay.Reveal)["best_score"])
}

func TestRoundScoreIsZeroBeforeSubmission(t *testing.T) {
	t.Parallel()

	game := newGame(&stubClock{now: testNow}, newListings(4))

	round, _ := startRound(t, game)

	assert.Equal(t, 0, game.RoundScore(round),
		"an unfinished round must not contribute a score")
}

func TestCollectedListingsAreRevealedOnlyWhenReached(t *testing.T) {
	t.Parallel()

	clock := &stubClock{now: testNow}
	game := newGame(clock, newListings(4))

	round, _ := startRound(t, game)

	clock.advance(30 * time.Second)

	outcome, err := submit(game, round, 100, 30, 70, 110)
	require.NoError(t, err)

	var reveal struct {
		Listings []struct {
			DisplayID string `json:"display_id"`
			Title     string `json:"title"`
		} `json:"listings"`
	}
	require.NoError(t, json.Unmarshal(outcome.Reveal, &reveal))

	assert.Len(t, reveal.Listings, 2,
		"a collectible above the reached score must not be claimable")
	assert.Equal(t, "item00", reveal.Listings[0].DisplayID)
}

func TestUnknownCollectibleIsIgnored(t *testing.T) {
	t.Parallel()

	clock := &stubClock{now: testNow}
	game := newGame(clock, newListings(4))

	round, _ := startRound(t, game)

	clock.advance(30 * time.Second)

	outcome, err := submit(game, round, 500, 999, 1234)
	require.NoError(t, err)

	var reveal struct {
		Listings []json.RawMessage `json:"listings"`
	}
	require.NoError(t, json.Unmarshal(outcome.Reveal, &reveal))

	assert.Empty(t, reveal.Listings)
}

func TestTooManyCollectiblesAreRejected(t *testing.T) {
	t.Parallel()

	clock := &stubClock{now: testNow}
	game := newGame(clock, newListings(4))

	round, _ := startRound(t, game)

	_, err := submit(game, round, 10, 1, 2, 3, 4, 5, 6, 7, 8)

	require.ErrorAs(t, err, new(*domainerr.InvalidError))
}

func TestStartSurvivesListingPoolFailure(t *testing.T) {
	t.Parallel()

	game := newGame(&stubClock{now: testNow}, &stubListings{err: assert.AnError})

	round := domain.NewRound(uuid.New(), raccoonjump.Slug, testNow)

	_, err := game.Start(context.Background(), round)

	require.NoError(t, err, "a listings outage must not block the game")
}

func TestStreakRule(t *testing.T) {
	t.Parallel()

	game := newGame(&stubClock{now: testNow}, newListings(4))

	assert.False(t, game.CountsTowardStreak(0), "an instant death must not farm the streak")
	assert.False(t, game.CountsTowardStreak(raccoonjump.MinStreakScore-1))
	assert.True(t, game.CountsTowardStreak(raccoonjump.MinStreakScore))
	assert.True(t, game.CountsTowardStreak(raccoonjump.MinStreakScore+100))
}

func TestMaxPlausibleScoreGrowsWithTimeAndIsCapped(t *testing.T) {
	t.Parallel()

	short := raccoonjump.MaxPlausibleScore(10 * time.Second)
	long := raccoonjump.MaxPlausibleScore(60 * time.Second)

	assert.Greater(t, long, short)
	assert.LessOrEqual(t, raccoonjump.MaxPlausibleScore(10*time.Hour), raccoonjump.MaxScore)
	assert.GreaterOrEqual(t, raccoonjump.MaxPlausibleScore(-5*time.Second), 0,
		"a skewed clock must not produce a negative bound")
}

func TestResumeKeepsSeedStable(t *testing.T) {
	t.Parallel()

	game := newGame(&stubClock{now: testNow}, newListings(4))

	round, prompt := startRound(t, game)

	resumed, err := game.Resume(context.Background(), round)
	require.NoError(t, err)

	assert.JSONEq(t, string(prompt), string(resumed.Prompt),
		"resuming must not reroll the layout seed")
}

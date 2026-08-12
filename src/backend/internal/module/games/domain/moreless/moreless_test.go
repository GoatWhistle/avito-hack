package moreless_test

import (
	"context"
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/games/domain"
	"github.com/avito-hack/backend/internal/module/games/domain/moreless"
	"github.com/avito-hack/backend/internal/shared/domainerr"
)

func TestSlugAndTargetStreak(t *testing.T) {
	t.Parallel()

	game := moreless.New(poolOf(100, 200))

	assert.Equal(t, "moreless", game.Slug())
	assert.Equal(t, 7, game.TargetStreak())
}

func TestPromptNeverContainsHiddenRightPrice(t *testing.T) {
	t.Parallel()

	game := moreless.New(poolOf(1320000, hiddenPrice, 4242424242))
	round := newRound()

	view, err := game.Start(context.Background(), round)
	require.NoError(t, err)

	marshalled := string(view.Prompt)

	assert.NotContains(t, marshalled, strconv.FormatInt(hiddenPrice, 10),
		"the right item price must never be serialized into the prompt")
	assert.Contains(t, marshalled, "1320000")

	prompt := decodePrompt(t, view.Prompt)
	right, ok := prompt["right"].(map[string]any)
	require.True(t, ok)

	_, hasPrice := right["price"]
	assert.False(t, hasPrice, "right side must have no price field at all")
	assert.ElementsMatch(t, []string{"item_id", "title", "photo_url"}, keysOf(right))

	left, ok := prompt["left"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, float64(1320000), left["price"])
}

func TestPayloadKeepsHiddenPriceServerSide(t *testing.T) {
	t.Parallel()

	game := moreless.New(poolOf(100, hiddenPrice, 300))
	round := newRound()

	view, err := game.Start(context.Background(), round)
	require.NoError(t, err)

	assert.NotContains(t, string(view.Prompt), strconv.FormatInt(hiddenPrice, 10))
	assert.Contains(t, string(round.Payload()), strconv.FormatInt(hiddenPrice, 10))
}

func TestStartNeverRepeatsAnItem(t *testing.T) {
	t.Parallel()

	pool := poolOf(100, 200)
	round := newRound()

	_, err := moreless.New(pool).Start(context.Background(), round)
	require.NoError(t, err)

	require.Len(t, pool.excludes, 1)
	assert.Empty(t, pool.excludes[0], "the very first draw has no reference price yet")

	require.Len(t, pool.nearCalls, 1)
	assert.Len(t, pool.nearCalls[0].seen, 1)
	assert.Equal(t, int64(100), pool.nearCalls[0].reference,
		"the opening right item must be matched against the left price")
}

func TestGuessHigherCorrectAdvancesToFreshPair(t *testing.T) {
	t.Parallel()

	pool := poolOf(100, 500, 900)
	game := moreless.New(pool)
	round := newRound()

	_, err := game.Start(context.Background(), round)
	require.NoError(t, err)

	outcome, err := game.Guess(context.Background(), round, json.RawMessage(`{"choice":"higher"}`))
	require.NoError(t, err)

	assert.True(t, outcome.Correct)
	assert.JSONEq(t, `{"right_price":500}`, string(outcome.Reveal))
	require.NotNil(t, outcome.Next)

	prompt := decodePrompt(t, outcome.Next.Prompt)
	left, ok := prompt["left"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, float64(500), left["price"], "the revealed right item becomes the new known left")

	assert.NotContains(t, string(outcome.Next.Prompt), "900")
}

func TestGuessWrongEndsWithoutNextView(t *testing.T) {
	t.Parallel()

	game := moreless.New(poolOf(500, 100, 900))
	round := newRound()

	_, err := game.Start(context.Background(), round)
	require.NoError(t, err)

	outcome, err := game.Guess(context.Background(), round, json.RawMessage(`{"choice":"higher"}`))
	require.NoError(t, err)

	assert.False(t, outcome.Correct)
	assert.Nil(t, outcome.Next)
	assert.JSONEq(t, `{"right_price":100}`, string(outcome.Reveal))
}

func TestEqualPricesCountAsCorrectForEitherChoice(t *testing.T) {
	t.Parallel()

	for _, choice := range []string{"higher", "lower"} {
		t.Run(choice, func(t *testing.T) {
			t.Parallel()

			game := moreless.New(poolOf(700, 700, 900))
			round := newRound()

			_, err := game.Start(context.Background(), round)
			require.NoError(t, err)

			outcome, err := game.Guess(context.Background(), round, json.RawMessage(`{"choice":"`+choice+`"}`))
			require.NoError(t, err)

			assert.True(t, outcome.Correct, "equal prices must never lose")
		})
	}
}

func TestGuessAtTargetStreakReturnsNoNextView(t *testing.T) {
	t.Parallel()

	game := moreless.New(poolOf(100, 200, 300))
	round := newRound()

	_, err := game.Start(context.Background(), round)
	require.NoError(t, err)

	for i := 0; i < 6; i++ {
		round.Advance(7, time.Now())
	}

	outcome, err := game.Guess(context.Background(), round, json.RawMessage(`{"choice":"higher"}`))
	require.NoError(t, err)

	assert.True(t, outcome.Correct)
	assert.Nil(t, outcome.Next, "the winning guess must not draw another pair")
}

func TestGuessRejectsMalformedMoves(t *testing.T) {
	t.Parallel()

	cases := map[string]string{
		"empty":         ``,
		"not an object": `"higher"`,
		"unknown choice": `{"choice":"equal"}`,
		"missing choice": `{}`,
	}

	for name, move := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			game := moreless.New(poolOf(100, 200, 300))
			round := newRound()

			_, err := game.Start(context.Background(), round)
			require.NoError(t, err)

			var invalid *domainerr.InvalidError

			_, err = game.Guess(context.Background(), round, json.RawMessage(move))
			require.ErrorAs(t, err, &invalid)
		})
	}
}

func TestStartPropagatesEmptyPool(t *testing.T) {
	t.Parallel()

	game := moreless.New(&fakePool{err: domain.ErrNoItemsInPool})

	_, err := game.Start(context.Background(), newRound())

	require.ErrorIs(t, err, domain.ErrNoItemsInPool)
}

func TestSeenGrowsWithEveryDrawnItem(t *testing.T) {
	t.Parallel()

	pool := poolOf(100, 200, 300, 400)
	game := moreless.New(pool)
	round := newRound()

	_, err := game.Start(context.Background(), round)
	require.NoError(t, err)

	_, err = game.Guess(context.Background(), round, json.RawMessage(`{"choice":"higher"}`))
	require.NoError(t, err)

	var payload struct {
		Seen []string `json:"seen"`
	}
	require.NoError(t, json.Unmarshal(round.Payload(), &payload))

	assert.Len(t, payload.Seen, 3)

	require.Len(t, pool.nearCalls, 2)
	assert.Len(t, pool.nearCalls[1].seen, 2, "the third draw must exclude both already seen items")
}

func TestAdvanceDrawsNearTheRevealedPrice(t *testing.T) {
	t.Parallel()

	pool := poolOf(100, 200, 300, 400)
	game := moreless.New(pool)
	round := newRound()

	_, err := game.Start(context.Background(), round)
	require.NoError(t, err)

	_, err = game.Guess(context.Background(), round, json.RawMessage(`{"choice":"higher"}`))
	require.NoError(t, err)

	require.Len(t, pool.nearCalls, 2)
	assert.Equal(t, int64(200), pool.nearCalls[1].reference,
		"the new reference is the price that just became visible")
	assert.Equal(t, 0.5, pool.nearCalls[1].minRatio)
	assert.Equal(t, 2.0, pool.nearCalls[1].maxRatio)
}

func TestEmptyBandFallsBackToUnrestrictedDraw(t *testing.T) {
	t.Parallel()

	pool := poolOf(100, 200, 300, 400, 500, 600, 700, 800)
	pool.nearEmpty = true

	game := moreless.New(pool)
	round := newRound()

	view, err := game.Start(context.Background(), round)
	require.NoError(t, err, "an empty price band must never fail the round")
	require.NotEmpty(t, view.Prompt)

	assert.Len(t, pool.nearCalls, 3*3, "every band is tried on each opening attempt")
	assert.NotEmpty(t, pool.excludes, "the unrestricted random draw is the final fallback")

	var payload struct {
		Seen []string `json:"seen"`
	}
	require.NoError(t, json.Unmarshal(round.Payload(), &payload))
	assert.Len(t, payload.Seen, 2)
}

func TestEmptyBandFallbackKeepsTheGameAdvancing(t *testing.T) {
	t.Parallel()

	pool := poolOf(100, 500, 900, 1300, 1700, 2100, 2500, 2900)
	pool.nearEmpty = true

	game := moreless.New(pool)
	round := newRound()

	_, err := game.Start(context.Background(), round)
	require.NoError(t, err)

	outcome, err := game.Guess(context.Background(), round, json.RawMessage(`{"choice":"higher"}`))
	require.NoError(t, err, "advancing must survive an empty band too")

	assert.True(t, outcome.Correct)
	require.NotNil(t, outcome.Next)
	require.NotEmpty(t, outcome.Next.Prompt)
}

func TestBandsWidenUntilACandidateIsFound(t *testing.T) {
	t.Parallel()

	pool := bandedPoolOf(1000, 2500)
	game := moreless.New(pool)
	round := newRound()

	_, err := game.Start(context.Background(), round)
	require.NoError(t, err)

	require.Len(t, pool.nearCalls, 2, "2500 misses the 0.5-2.0 band and is caught by the 0.33-3.0 one")
	assert.Equal(t, 0.5, pool.nearCalls[0].minRatio)
	assert.Equal(t, 0.33, pool.nearCalls[1].minRatio)
}

func TestNearTieCandidateIsRedrawnWhenAnAlternativeExists(t *testing.T) {
	t.Parallel()

	pool := bandedPoolOf(100000, 100500, 150000)
	game := moreless.New(pool)
	round := newRound()

	_, err := game.Start(context.Background(), round)
	require.NoError(t, err)

	var payload struct {
		Right struct {
			PriceKopeks int64 `json:"PriceKopeks"`
		} `json:"right"`
	}
	require.NoError(t, json.Unmarshal(round.Payload(), &payload))

	assert.Equal(t, int64(150000), payload.Right.PriceKopeks,
		"a candidate within 2% of the reference is a coin flip and must be redrawn")
}

func TestNearTieIsAcceptedWhenNothingElseIsAvailable(t *testing.T) {
	t.Parallel()

	pool := bandedPoolOf(100000, 100500)
	game := moreless.New(pool)
	round := newRound()

	_, err := game.Start(context.Background(), round)
	require.NoError(t, err, "a near tie is still better than failing the round")

	var payload struct {
		Right struct {
			PriceKopeks int64 `json:"PriceKopeks"`
		} `json:"right"`
	}
	require.NoError(t, json.Unmarshal(round.Payload(), &payload))

	assert.Equal(t, int64(100500), payload.Right.PriceKopeks)
}

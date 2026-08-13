package bukovki

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/games/domain"
)

type mockWordPool struct {
	word  string
	valid bool
}

func (m *mockWordPool) RandomWord(ctx context.Context, n int) (string, error) {
	return m.word, nil
}

func (m *mockWordPool) IsValidWord(ctx context.Context, word string) bool {
	return m.valid
}

func TestEvaluateGuess(t *testing.T) {
	tests := []struct {
		name     string
		guess    string
		secret   string
		expected []LetterStatus
	}{
		{
			name:     "all correct",
			guess:    "apple",
			secret:   "apple",
			expected: []LetterStatus{StatusCorrect, StatusCorrect, StatusCorrect, StatusCorrect, StatusCorrect},
		},
		{
			name:     "all absent",
			guess:    "ghost",
			secret:   "apple",
			expected: []LetterStatus{StatusAbsent, StatusAbsent, StatusAbsent, StatusAbsent, StatusAbsent},
		},
		{
			name:     "mixed with duplicates",
			guess:    "papel",
			secret:   "apple",
			expected: []LetterStatus{StatusPresent, StatusPresent, StatusCorrect, StatusPresent, StatusPresent},
		},
		{
			name:     "too many of one letter",
			guess:    "ppppp",
			secret:   "apple",
			expected: []LetterStatus{StatusAbsent, StatusCorrect, StatusCorrect, StatusAbsent, StatusAbsent},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evaluateGuess(tt.guess, tt.secret)
			require.Len(t, result, 5)
			for i, r := range result {
				assert.Equal(t, string(tt.guess[i]), r.Char)
				assert.Equal(t, tt.expected[i], r.Status)
			}
		})
	}
}

func TestGuess_WinningMoveAdvancesStreak(t *testing.T) {
	pool := &mockWordPool{word: "apple", valid: true}
	game := New(pool)

	round := domain.NewRound(uuid.New(), Slug, time.Now())
	_, err := game.Start(context.Background(), round)
	require.NoError(t, err)

	assert.Equal(t, 0, round.Streak())

	move := json.RawMessage(`{"guess":"apple"}`)
	outcome, err := game.Guess(context.Background(), round, move)
	require.NoError(t, err)

	assert.True(t, outcome.Correct, "winning guess should be correct")
	assert.Nil(t, outcome.Next, "winning guess should not have a next view")

	assert.Equal(t, game.TargetStreak()-1, round.Streak())
}

func TestGuess_WrongMoveWithTriesRemaining(t *testing.T) {
	pool := &mockWordPool{word: "apple", valid: true}
	game := New(pool)

	round := domain.NewRound(uuid.New(), Slug, time.Now())
	_, err := game.Start(context.Background(), round)
	require.NoError(t, err)

	move := json.RawMessage(`{"guess":"ghost"}`)
	outcome, err := game.Guess(context.Background(), round, move)
	require.NoError(t, err)

	assert.True(t, outcome.Correct, "wrong guess with tries remaining should return true so engine doesn't kill the game")
	assert.NotNil(t, outcome.Next, "should have next view")

	assert.Equal(t, 0, round.Streak())
}

func TestGuess_GameOverOnMaxTries(t *testing.T) {
	pool := &mockWordPool{word: "apple", valid: true}
	game := New(pool)

	round := domain.NewRound(uuid.New(), Slug, time.Now())
	_, err := game.Start(context.Background(), round)
	require.NoError(t, err)

	for i := 0; i < 5; i++ {
		move := json.RawMessage(`{"guess":"ghost"}`)
		_, err := game.Guess(context.Background(), round, move)
		require.NoError(t, err)
	}

	move := json.RawMessage(`{"guess":"ghost"}`)
	outcome, err := game.Guess(context.Background(), round, move)
	require.NoError(t, err)

	assert.False(t, outcome.Correct, "final wrong guess should return false so engine kills the game")
	assert.Nil(t, outcome.Next, "game over, no next view")
}

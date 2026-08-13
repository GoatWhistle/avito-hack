package bukovki

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/avito-hack/backend/internal/module/games/domain"
)

const Slug = "bukovki"

const maxTries = 6
const wordLength = 5

type WordPool interface {
	RandomWord(ctx context.Context, n int) (string, error)
	IsValidWord(ctx context.Context, word string) bool
}

type Game struct {
	pool WordPool
}

func New(pool WordPool) *Game {
	return &Game{pool: pool}
}

func (g *Game) Slug() string { return Slug }

func (g *Game) TargetStreak() int { return maxTries }

func (g *Game) Start(ctx context.Context, r *domain.Round) (domain.View, error) {
	secretWord, err := g.pool.RandomWord(ctx, wordLength)
	if err != nil {
		return domain.View{}, err
	}

	state := state{
		Secret:   secretWord,
		Guesses:  make([]string, 0),
		MaxTries: maxTries,
	}

	return g.commit(r, state)
}

func (g *Game) Resume(ctx context.Context, r *domain.Round) (domain.View, error) {
	s, err := decodeState(r.Payload())
	if err != nil {
		return domain.View{}, err
	}

	prompt, err := s.prompt()
	if err != nil {
		return domain.View{}, err
	}
	return domain.View{Prompt: prompt}, nil
}

func (g *Game) Guess(ctx context.Context, r *domain.Round, move json.RawMessage) (domain.GuessOutcome, error) {

	wordGuess, err := parseMove(move)
	if err != nil {
		return domain.GuessOutcome{}, err
	}

	if !g.pool.IsValidWord(ctx, wordGuess) {
		return domain.GuessOutcome{}, domain.ErrInvalidMove("word is not valid")
	}

	s, err := decodeState(r.Payload())
	if err != nil {
		return domain.GuessOutcome{}, domain.ErrInvalidMove("failed to decode state")
	}

	s.Guesses = append(s.Guesses, wordGuess)

	lettersFeedback := evaluateGuess(wordGuess, s.Secret)
	win := isCorrect(lettersFeedback)
	gameOver := win || len(s.Guesses) == s.MaxTries

	reveal := revealPayload{
		Feedback: lettersFeedback,
		GameOver: gameOver,
		Win:      win,
	}
	if gameOver {
		reveal.Secret = &s.Secret
	}

	revealJSON, err := json.Marshal(reveal)
	if err != nil {
		return domain.GuessOutcome{}, err
	}

	if win {
		for r.Streak() < g.TargetStreak()-1 {
			r.Advance(g.TargetStreak(), time.Now())
		}
		return domain.GuessOutcome{Correct: true, Reveal: revealJSON, Next: nil}, nil
	}

	if gameOver {
		return domain.GuessOutcome{Correct: false, Reveal: revealJSON, Next: nil}, nil
	}

	view, err := g.commit(r, s)
	if err != nil {
		return domain.GuessOutcome{}, err
	}
	return domain.GuessOutcome{Correct: true, Reveal: revealJSON, Next: &view}, nil
}

func isCorrect(feedback []LetterFeedback) bool {
	for _, f := range feedback {
		if f.Status != StatusCorrect {
			return false
		}
	}
	return true
}

func (g *Game) commit(r *domain.Round, s state) (domain.View, error) {
	payload, err := json.Marshal(s)
	if err != nil {
		return domain.View{}, fmt.Errorf("marshal bukovki state: %w", err)
	}

	prompt, err := s.prompt()
	if err != nil {
		return domain.View{}, err
	}

	r.SetPayload(payload)

	return domain.View{Prompt: prompt}, nil
}

func parseMove(raw json.RawMessage) (string, error) {
	if len(raw) == 0 {
		return "", domain.ErrInvalidMove("move is required")
	}

	var move movePayload
	if err := json.Unmarshal(raw, &move); err != nil {
		return "", domain.ErrInvalidMove("move must be an object with a guess field")
	}

	if move.Guess == "" {
		return "", domain.ErrInvalidMove("guess is required")
	}

	return move.Guess, nil
}

func evaluateGuess(guess, secret string) []LetterFeedback {
	secretCounts := make(map[rune]int)
	feedback := make([]LetterFeedback, len(guess))

	for _, r := range secret {
		secretCounts[r]++
	}
	for i, r := range guess {
		if r == rune(secret[i]) {
			feedback[i] = LetterFeedback{Char: string(r), Status: StatusCorrect}
			secretCounts[r]--
		}
	}

	for i, r := range guess {
		if feedback[i].Status == StatusCorrect {
			continue
		}

		if count, ok := secretCounts[r]; ok && count > 0 {
			feedback[i] = LetterFeedback{Char: string(r), Status: StatusPresent}
			secretCounts[r]--
		} else {
			feedback[i] = LetterFeedback{Char: string(r), Status: StatusAbsent}
		}
	}

	return feedback
}

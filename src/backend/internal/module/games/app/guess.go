package app

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/module/games/domain"
)

type GuessCommand struct {
	UserID    uuid.UUID
	GameSlug  string
	RoundID   string
	Move      json.RawMessage
	ClientDay domain.Day
}

type GuessHandler struct {
	registry *domain.Registry
	rounds   RoundRepository
	progress ProgressRepository
	tx       TxManager
	clock    Clock
}

func NewGuessHandler(
	registry *domain.Registry,
	rounds RoundRepository,
	progress ProgressRepository,
	tx TxManager,
	clock Clock,
) *GuessHandler {
	return &GuessHandler{registry: registry, rounds: rounds, progress: progress, tx: tx, clock: clock}
}

func (h *GuessHandler) Handle(ctx context.Context, cmd GuessCommand) (GuessResult, error) {
	game, err := h.registry.Get(cmd.GameSlug)
	if err != nil {
		return GuessResult{}, err
	}

	var result GuessResult

	err = h.tx.WithTx(ctx, func(ctx context.Context) error {
		round, loadErr := h.load(ctx, cmd, game.Slug())
		if loadErr != nil {
			return loadErr
		}

		outcome, guessErr := game.Guess(ctx, round, cmd.Move)
		if guessErr != nil {
			return guessErr
		}

		now := h.clock.Now()

		if !outcome.Correct {
			round.Lose(now)
		} else {
			round.Advance(game.TargetStreak(), now)
		}

		if saveErr := h.rounds.Save(ctx, round); saveErr != nil {
			return saveErr
		}

		result = GuessResult{
			Correct: outcome.Correct,
			Reveal:  outcome.Reveal,
			Streak:  round.Streak(),
			State:   round.State(),
		}

		if outcome.Next != nil && round.IsActive() {
			result.Prompt = outcome.Next.Prompt
		}

		if round.State() != domain.StateWon {
			return nil
		}

		result.AttemptCompleted = true

		streak, progressErr := h.completeAttempt(ctx, cmd, round)
		if progressErr != nil {
			return progressErr
		}

		view := toStreakView(streak)
		result.StreakAfter = &view

		return nil
	})
	if err != nil {
		return GuessResult{}, err
	}

	return result, nil
}

func (h *GuessHandler) load(ctx context.Context, cmd GuessCommand, slug string) (*domain.Round, error) {
	round, err := h.rounds.ByDisplayID(ctx, cmd.RoundID)
	if err != nil {
		return nil, err
	}

	if round.UserID() != cmd.UserID || round.GameSlug() != slug {
		return nil, domain.ErrRoundNotFound
	}

	if !round.IsActive() {
		return nil, domain.ErrRoundFinished(round.State())
	}

	return round, nil
}

func (h *GuessHandler) completeAttempt(
	ctx context.Context,
	cmd GuessCommand,
	round *domain.Round,
) (domain.Streak, error) {
	today := cmd.ClientDay

	daily, err := h.progress.Daily(ctx, cmd.UserID, cmd.GameSlug, today)
	if err != nil {
		return domain.Streak{}, err
	}

	firstOfDay := domain.IsFirstAttemptOfDay(daily)

	if err := h.progress.SaveDaily(
		ctx, cmd.UserID, cmd.GameSlug, today, domain.AdvanceDaily(daily, round.BestStreak()),
	); err != nil {
		return domain.Streak{}, err
	}

	streak, err := h.progress.Streak(ctx, cmd.UserID, cmd.GameSlug)
	if err != nil {
		return domain.Streak{}, err
	}

	if !firstOfDay {
		return streak, nil
	}

	streak = domain.AdvanceStreak(streak, today)

	if err := h.progress.SaveStreak(ctx, cmd.UserID, cmd.GameSlug, streak); err != nil {
		return domain.Streak{}, err
	}

	return streak, nil
}

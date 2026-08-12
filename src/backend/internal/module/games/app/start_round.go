package app

import (
	"context"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/module/games/domain"
)

type StartRoundCommand struct {
	UserID   uuid.UUID
	GameSlug string
}

type StartRoundHandler struct {
	registry *domain.Registry
	rounds   RoundRepository
	tx       TxManager
	clock    Clock
}

func NewStartRoundHandler(
	registry *domain.Registry,
	rounds RoundRepository,
	tx TxManager,
	clock Clock,
) *StartRoundHandler {
	return &StartRoundHandler{registry: registry, rounds: rounds, tx: tx, clock: clock}
}

func (h *StartRoundHandler) Handle(ctx context.Context, cmd StartRoundCommand) (StartRoundResult, error) {
	game, err := h.registry.Get(cmd.GameSlug)
	if err != nil {
		return StartRoundResult{}, err
	}

	var result StartRoundResult

	err = h.tx.WithTx(ctx, func(ctx context.Context) error {
		if active, activeErr := h.rounds.ActiveByUser(ctx, cmd.UserID, cmd.GameSlug); activeErr == nil {
			active.Lose(h.clock.Now())

			if saveErr := h.rounds.Save(ctx, active); saveErr != nil {
				return saveErr
			}
		}

		round := domain.NewRound(cmd.UserID, game.Slug(), h.clock.Now())

		view, startErr := game.Start(ctx, round)
		if startErr != nil {
			return startErr
		}

		if saveErr := h.rounds.Save(ctx, round); saveErr != nil {
			return saveErr
		}

		result = StartRoundResult{
			RoundID:      round.DisplayID(),
			Streak:       round.Streak(),
			TargetStreak: game.TargetStreak(),
			Prompt:       view.Prompt,
		}

		return nil
	})
	if err != nil {
		return StartRoundResult{}, err
	}

	return result, nil
}

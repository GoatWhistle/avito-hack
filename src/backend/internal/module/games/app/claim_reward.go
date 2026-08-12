package app

import (
	"context"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/module/games/domain"
)

type ClaimRewardCommand struct {
	UserID   uuid.UUID
	GameSlug string
}

type ClaimRewardHandler struct {
	registry *domain.Registry
	progress ProgressRepository
	tx       TxManager
	clock    Clock
}

func NewClaimRewardHandler(
	registry *domain.Registry,
	progress ProgressRepository,
	tx TxManager,
	clock Clock,
) *ClaimRewardHandler {
	return &ClaimRewardHandler{registry: registry, progress: progress, tx: tx, clock: clock}
}

func (h *ClaimRewardHandler) Handle(ctx context.Context, cmd ClaimRewardCommand) (ClaimRewardResult, error) {
	if _, err := h.registry.Get(cmd.GameSlug); err != nil {
		return ClaimRewardResult{}, err
	}

	var result ClaimRewardResult

	err := h.tx.WithTx(ctx, func(ctx context.Context) error {
		streak, err := h.progress.Streak(ctx, cmd.UserID, cmd.GameSlug)
		if err != nil {
			return err
		}

		claimed, ok := domain.ClaimStreak(streak, h.clock.Now())
		if !ok {
			return domain.ErrRewardNotReady()
		}

		if err := h.progress.SaveStreak(ctx, cmd.UserID, cmd.GameSlug, claimed); err != nil {
			return err
		}

		result = ClaimRewardResult{Code: domain.NewPromoCode()}

		return nil
	})
	if err != nil {
		return ClaimRewardResult{}, err
	}

	return result, nil
}

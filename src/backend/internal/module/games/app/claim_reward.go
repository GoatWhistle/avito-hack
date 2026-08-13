package app

import (
	"context"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/module/games/domain"
)

type ClaimRewardCommand struct {
	UserID uuid.UUID
}

type ClaimRewardHandler struct {
	progress ProgressRepository
	tx       TxManager
	clock    Clock
}

func NewClaimRewardHandler(
	progress ProgressRepository,
	tx TxManager,
	clock Clock,
) *ClaimRewardHandler {
	return &ClaimRewardHandler{progress: progress, tx: tx, clock: clock}
}

func (h *ClaimRewardHandler) Handle(ctx context.Context, cmd ClaimRewardCommand) (ClaimRewardResult, error) {
	var result ClaimRewardResult

	err := h.tx.WithTx(ctx, func(ctx context.Context) error {
		streak, err := h.progress.StreakForUpdate(ctx, cmd.UserID)
		if err != nil {
			return err
		}

		claimed, ok := domain.ClaimStreak(streak, h.clock.Now())
		if !ok {
			return domain.ErrRewardNotReady()
		}

		if err := h.progress.SaveStreak(ctx, cmd.UserID, claimed); err != nil {
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

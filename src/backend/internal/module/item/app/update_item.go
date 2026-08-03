package app

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/module/item/domain"
	"github.com/avito-hack/backend/internal/shared/domainerr"
	"github.com/avito-hack/backend/internal/shared/vo"
)

type UpdateItemCommand struct {
	ItemID      uuid.UUID
	ActorID     uuid.UUID
	Title       *string
	Description *string
	PriceKopeks *int64
	Attributes  *map[string]string
}

type UpdateItemHandler struct {
	items domain.Repository
	tx    TxManager
	clock Clock
}

func NewUpdateItemHandler(items domain.Repository, tx TxManager, clock Clock) *UpdateItemHandler {
	return &UpdateItemHandler{items: items, tx: tx, clock: clock}
}

func (h *UpdateItemHandler) Handle(ctx context.Context, cmd UpdateItemCommand) (*domain.Item, error) {
	params, err := toUpdateParams(cmd, h.clock.Now())
	if err != nil {
		return nil, err
	}

	var updated *domain.Item

	err = h.tx.WithTx(ctx, func(ctx context.Context) error {
		item, loadErr := h.items.ByIDForUpdate(ctx, cmd.ItemID)
		if loadErr != nil {
			return fmt.Errorf("load item: %w", loadErr)
		}

		if !item.IsOwnedBy(cmd.ActorID) {
			return domainerr.ErrForbidden
		}

		if updateErr := item.Update(params); updateErr != nil {
			return updateErr
		}

		if saveErr := h.items.Save(ctx, item); saveErr != nil {
			return saveErr
		}

		updated = item

		return nil
	})
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func toUpdateParams(cmd UpdateItemCommand, now time.Time) (domain.UpdateItemParams, error) {
	params := domain.UpdateItemParams{
		Title:       cmd.Title,
		Description: cmd.Description,
		Now:         now,
	}

	if cmd.PriceKopeks != nil {
		price, err := vo.NewMoney(*cmd.PriceKopeks)
		if err != nil {
			return domain.UpdateItemParams{}, domainerr.NewInvalid("price", "value must be greater than or equal to 0")
		}
		params.Price = &price
	}

	if cmd.Attributes != nil {
		attributes := domain.NewAttributes(*cmd.Attributes)
		params.Attributes = &attributes
	}

	return params, nil
}

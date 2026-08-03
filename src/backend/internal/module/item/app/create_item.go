package app

import (
	"context"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/module/item/domain"
	"github.com/avito-hack/backend/internal/shared/domainerr"
	"github.com/avito-hack/backend/internal/shared/vo"
)

type CreateItemCommand struct {
	OwnerID     uuid.UUID
	Title       string
	Description string
	PriceKopeks int64
	Attributes  map[string]string
}

type CreateItemHandler struct {
	items domain.Repository
	tx    TxManager
	clock Clock
}

func NewCreateItemHandler(items domain.Repository, tx TxManager, clock Clock) *CreateItemHandler {
	return &CreateItemHandler{items: items, tx: tx, clock: clock}
}

func (h *CreateItemHandler) Handle(ctx context.Context, cmd CreateItemCommand) (*domain.Item, error) {
	price, err := vo.NewMoney(cmd.PriceKopeks)
	if err != nil {
		return nil, domainerr.NewInvalid("price", "value must be greater than or equal to 0")
	}

	item, err := domain.NewItem(domain.NewItemParams{
		OwnerID:     cmd.OwnerID,
		Title:       cmd.Title,
		Description: cmd.Description,
		Price:       price,
		Attributes:  domain.NewAttributes(cmd.Attributes),
		Now:         h.clock.Now(),
	})
	if err != nil {
		return nil, err
	}

	if err := h.tx.WithTx(ctx, func(ctx context.Context) error {
		return h.items.Save(ctx, item)
	}); err != nil {
		return nil, err
	}

	return item, nil
}

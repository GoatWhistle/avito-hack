package app

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/module/item/domain"
	"github.com/avito-hack/backend/internal/shared/events"
)

type ViewItemCommand struct {
	ItemID   uuid.UUID
	ViewerID uuid.UUID
}

type ViewItemHandler struct {
	items domain.Repository
	clock Clock
	bus   events.Publisher
}

func NewViewItemHandler(items domain.Repository, clock Clock, bus events.Publisher) *ViewItemHandler {
	return &ViewItemHandler{items: items, clock: clock, bus: bus}
}

func (h *ViewItemHandler) Handle(ctx context.Context, cmd ViewItemCommand) error {
	item, err := h.items.ByID(ctx, cmd.ItemID)
	if err != nil {
		return fmt.Errorf("load item: %w", err)
	}

	if item.Status() != domain.StatusPublished && item.Status() != domain.StatusSold {
		return domain.ErrItemNotFound
	}

	if item.IsOwnedBy(cmd.ViewerID) {
		return nil
	}

	h.bus.Publish(ctx, events.New(events.TypeItemViewed, cmd.ViewerID, cmd.ItemID, h.clock.Now()))

	return nil
}

package app

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/module/item/domain"
	"github.com/avito-hack/backend/internal/shared/auth"
	"github.com/avito-hack/backend/internal/shared/domainerr"
)

type GetItemQuery struct {
	ItemID uuid.UUID
	Actor  auth.Actor
}

type GetItemHandler struct {
	items domain.Repository
}

func NewGetItemHandler(items domain.Repository) *GetItemHandler {
	return &GetItemHandler{items: items}
}

func (h *GetItemHandler) Handle(ctx context.Context, q GetItemQuery) (*domain.Item, error) {
	item, err := h.items.ByID(ctx, q.ItemID)
	if err != nil {
		return nil, fmt.Errorf("load item: %w", err)
	}

	if item.Status() == domain.StatusPublished {
		return item, nil
	}

	if q.Actor.IsZero() {
		return nil, domain.ErrItemNotFound
	}

	if item.IsOwnedBy(q.Actor.ID) || q.Actor.HasRole(auth.RoleModerator, auth.RoleAdmin) {
		return item, nil
	}

	return nil, domainerr.ErrForbidden
}

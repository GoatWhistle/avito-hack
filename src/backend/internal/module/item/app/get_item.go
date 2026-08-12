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

type ItemView struct {
	Item             *domain.Item
	OwnerDisplayID   string
	ModerationReason string
}

type GetItemHandler struct {
	items         domain.Repository
	moderationLog domain.ModerationLogRepository
	owners        OwnerProvider
}

func NewGetItemHandler(
	items domain.Repository,
	moderationLog domain.ModerationLogRepository,
	owners OwnerProvider,
) *GetItemHandler {
	return &GetItemHandler{items: items, moderationLog: moderationLog, owners: owners}
}

func (h *GetItemHandler) Handle(ctx context.Context, q GetItemQuery) (ItemView, error) {
	item, err := h.items.ByID(ctx, q.ItemID)
	if err != nil {
		return ItemView{}, fmt.Errorf("load item: %w", err)
	}

	if item.Status() == domain.StatusPublished || item.Status() == domain.StatusSold {
		return ItemView{Item: item, OwnerDisplayID: h.ownerDisplayID(ctx, item.OwnerID())}, nil
	}

	if q.Actor.IsZero() {
		return ItemView{}, domain.ErrItemNotFound
	}

	if !item.IsOwnedBy(q.Actor.ID) && !q.Actor.HasRole(auth.RoleModerator, auth.RoleAdmin) {
		return ItemView{}, domainerr.ErrForbidden
	}

	view := ItemView{Item: item, OwnerDisplayID: h.ownerDisplayID(ctx, item.OwnerID())}

	if item.Status() == domain.StatusModeration && h.moderationLog != nil {
		entry, logErr := h.moderationLog.LatestByItemID(ctx, item.ID())
		if logErr == nil && entry != nil && entry.Verdict == domain.ModerationRejected {
			view.ModerationReason = entry.Reason
		}
	}

	return view, nil
}

func (h *GetItemHandler) ownerDisplayID(ctx context.Context, ownerID uuid.UUID) string {
	if h.owners == nil {
		return ""
	}

	owners, err := h.owners.ByIDs(ctx, []uuid.UUID{ownerID})
	if err != nil {
		return ""
	}

	return owners[ownerID].DisplayID
}

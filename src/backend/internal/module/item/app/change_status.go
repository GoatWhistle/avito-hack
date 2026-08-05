package app

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/module/item/domain"
	"github.com/avito-hack/backend/internal/shared/auth"
	"github.com/avito-hack/backend/internal/shared/domainerr"
)

type StatusAction string

const (
	ActionSubmit  StatusAction = "submit"
	ActionPublish StatusAction = "publish"
	ActionArchive StatusAction = "archive"
	ActionRestore StatusAction = "restore"
)

type ChangeStatusCommand struct {
	ItemID uuid.UUID
	Actor  auth.Actor
	Action StatusAction
}

type ChangeStatusHandler struct {
	items domain.Repository
	tx    TxManager
	clock Clock
}

func NewChangeStatusHandler(items domain.Repository, tx TxManager, clock Clock) *ChangeStatusHandler {
	return &ChangeStatusHandler{items: items, tx: tx, clock: clock}
}

func (h *ChangeStatusHandler) Handle(ctx context.Context, cmd ChangeStatusCommand) (*domain.Item, error) {
	var updated *domain.Item

	err := h.tx.WithTx(ctx, func(ctx context.Context) error {
		item, err := h.items.ByIDForUpdate(ctx, cmd.ItemID)
		if err != nil {
			return fmt.Errorf("load item: %w", err)
		}

		if err := authorize(item, cmd); err != nil {
			return err
		}

		if err := applyAction(item, cmd.Action, h.clock.Now()); err != nil {
			return err
		}

		if err := h.items.Save(ctx, item); err != nil {
			return err
		}

		updated = item

		return nil
	})
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func authorize(item *domain.Item, cmd ChangeStatusCommand) error {
	isModerator := cmd.Actor.HasRole(auth.RoleModerator, auth.RoleAdmin)

	if cmd.Action == ActionPublish && !isModerator {
		return domainerr.ErrForbidden
	}

	if !isModerator && !item.IsOwnedBy(cmd.Actor.ID) {
		return domainerr.ErrForbidden
	}

	return nil
}

func applyAction(item *domain.Item, action StatusAction, now time.Time) error {
	switch action {
	case ActionSubmit:
		return item.SubmitForModeration(now)
	case ActionPublish:
		return item.Publish(now)
	case ActionArchive:
		return item.Archive(now)
	case ActionRestore:
		return item.Restore(now)
	default:
		return domainerr.NewInvalid("action", "unknown status action")
	}
}

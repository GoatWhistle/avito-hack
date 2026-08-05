package app

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/module/item/domain"
	"github.com/avito-hack/backend/internal/shared/auth"
	"github.com/avito-hack/backend/internal/shared/domainerr"
	"github.com/avito-hack/backend/internal/shared/events"
)

type StatusAction string

const (
	ActionSubmit  StatusAction = "submit"
	ActionPublish StatusAction = "publish"
	ActionSell    StatusAction = "sell"
	ActionArchive StatusAction = "archive"
	ActionRestore StatusAction = "restore"
)

type ChangeStatusCommand struct {
	ItemID uuid.UUID
	Actor  auth.Actor
	Action StatusAction
}

type ChangeStatusHandler struct {
	items  domain.Repository
	photos domain.PhotoRepository
	tx     TxManager
	clock  Clock
	bus    events.Publisher
}

func NewChangeStatusHandler(
	items domain.Repository,
	photos domain.PhotoRepository,
	tx TxManager,
	clock Clock,
	bus events.Publisher,
) *ChangeStatusHandler {
	return &ChangeStatusHandler{items: items, photos: photos, tx: tx, clock: clock, bus: bus}
}

func (h *ChangeStatusHandler) Handle(ctx context.Context, cmd ChangeStatusCommand) (*domain.Item, error) {
	var updated *domain.Item

	err := events.PublishAfterCommit(ctx, h.bus, h.tx, func(ctx context.Context, out *events.Outbox) error {
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

		if event, ok := h.statusEvent(ctx, item, cmd.Action); ok {
			out.Add(event)
		}

		updated = item

		return nil
	})
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (h *ChangeStatusHandler) statusEvent(
	ctx context.Context,
	item *domain.Item,
	action StatusAction,
) (events.Event, bool) {
	var eventType events.Type

	switch action {
	case ActionPublish:
		eventType = events.TypeItemPublished
	case ActionSell:
		eventType = events.TypeItemSold
	case ActionSubmit, ActionArchive, ActionRestore:
		return events.Event{}, false
	default:
		return events.Event{}, false
	}

	return events.New(eventType, item.OwnerID(), item.ID(), h.clock.Now()).
		WithPayload(itemPayload(item, photoCount(ctx, h.photos, item.ID()))), true
}

func photoCount(ctx context.Context, photos domain.PhotoRepository, itemID uuid.UUID) int {
	if photos == nil {
		return 0
	}

	count, err := photos.CountByItemID(ctx, itemID)
	if err != nil {
		return 0
	}

	return count
}

func itemPayload(item *domain.Item, photoCount int) events.Payload {
	return events.Payload{
		Title:       item.Title(),
		Description: item.Description(),
		PriceKopeks: item.Price().Kopeks(),
		PhotoCount:  photoCount,
		Attributes:  map[string]string(item.Attributes()),
	}
}

func authorize(item *domain.Item, cmd ChangeStatusCommand) error {
	if item.IsOwnedBy(cmd.Actor.ID) {
		return nil
	}

	if cmd.Action == ActionSell {
		return domainerr.ErrForbidden
	}

	if !cmd.Actor.HasRole(auth.RoleModerator, auth.RoleAdmin) {
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
	case ActionSell:
		return item.MarkSold(now)
	case ActionArchive:
		return item.Archive(now)
	case ActionRestore:
		return item.Restore(now)
	default:
		return domainerr.NewInvalid("action", "unknown status action")
	}
}

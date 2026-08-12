package app

import (
	"context"
	"fmt"
	"slices"
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
	items    domain.Repository
	photos   domain.PhotoRepository
	tx       TxManager
	clock    Clock
	bus      events.Publisher
	moderate *ModerateItemHandler
}

func NewChangeStatusHandler(
	items domain.Repository,
	photos domain.PhotoRepository,
	tx TxManager,
	clock Clock,
	bus events.Publisher,
	moderate *ModerateItemHandler,
) *ChangeStatusHandler {
	return &ChangeStatusHandler{items: items, photos: photos, tx: tx, clock: clock, bus: bus, moderate: moderate}
}

func (h *ChangeStatusHandler) Handle(ctx context.Context, cmd ChangeStatusCommand) (*domain.Item, error) {
	var updated *domain.Item

	err := events.PublishAfterCommit(ctx, h.bus, h.tx, func(ctx context.Context, out *events.Outbox) error {
		item, loadErr := h.items.ByIDForUpdate(ctx, cmd.ItemID)
		if loadErr != nil {
			return fmt.Errorf("load item: %w", loadErr)
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

		event, ok, err := h.statusEvent(ctx, item, cmd.Action)
		if err != nil {
			return err
		}

		if ok {
			out.Add(event)
		}

		updated = item

		return nil
	})
	if err != nil {
		return nil, err
	}

	if cmd.Action == ActionSubmit && h.moderate != nil {
		if moderated := h.moderate.HandleOrLog(ctx, updated.ID(), cmd.Actor.ID); moderated != nil {
			updated = moderated
		}
	}

	return updated, nil
}

func (h *ChangeStatusHandler) statusEvent(
	ctx context.Context,
	item *domain.Item,
	action StatusAction,
) (events.Event, bool, error) {
	var eventType events.Type

	switch action {
	case ActionSell:
		eventType = events.TypeItemSold
	case ActionSubmit, ActionArchive, ActionRestore:
		return events.Event{}, false, nil
	default:
		return events.Event{}, false, nil
	}

	count, err := photoCount(ctx, h.photos, item.ID())
	if err != nil {
		return events.Event{}, false, err
	}

	event := events.New(eventType, item.OwnerID(), item.ID(), h.clock.Now()).
		WithPayload(itemPayload(item, count))

	return event, true, nil
}

func photoCount(ctx context.Context, photos domain.PhotoRepository, itemID uuid.UUID) (int, error) {
	if photos == nil {
		return 0, nil
	}

	count, err := photos.CountByItemID(ctx, itemID)
	if err != nil {
		return 0, fmt.Errorf("count item photos: %w", err)
	}

	return count, nil
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

var moderatorActions = []StatusAction{ActionArchive}

func authorize(item *domain.Item, cmd ChangeStatusCommand) error {
	if item.IsOwnedBy(cmd.Actor.ID) {
		return nil
	}

	if !cmd.Actor.HasRole(auth.RoleModerator, auth.RoleAdmin) {
		return domainerr.ErrForbidden
	}

	if !slices.Contains(moderatorActions, cmd.Action) {
		return domainerr.ErrForbidden
	}

	return nil
}

func applyAction(item *domain.Item, action StatusAction, now time.Time) error {
	switch action {
	case ActionSubmit:
		return item.SubmitForModeration(now)
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

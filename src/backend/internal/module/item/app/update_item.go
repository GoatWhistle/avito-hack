package app

import (
	"context"
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/module/item/domain"
	"github.com/avito-hack/backend/internal/shared/domainerr"
	"github.com/avito-hack/backend/internal/shared/events"
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
	items    domain.Repository
	photos   domain.PhotoRepository
	tx       TxManager
	clock    Clock
	bus      events.Publisher
	moderate *ModerateItemHandler
}

func NewUpdateItemHandler(
	items domain.Repository,
	photos domain.PhotoRepository,
	tx TxManager,
	clock Clock,
	bus events.Publisher,
	moderate *ModerateItemHandler,
) *UpdateItemHandler {
	return &UpdateItemHandler{items: items, photos: photos, tx: tx, clock: clock, bus: bus, moderate: moderate}
}

func (h *UpdateItemHandler) Handle(ctx context.Context, cmd UpdateItemCommand) (*domain.Item, error) {
	params, err := toUpdateParams(cmd, h.clock.Now())
	if err != nil {
		return nil, err
	}

	var updated *domain.Item

	err = events.PublishAfterCommit(ctx, h.bus, h.tx, func(ctx context.Context, out *events.Outbox) error {
		item, loadErr := h.items.ByIDForUpdate(ctx, cmd.ItemID)
		if loadErr != nil {
			return fmt.Errorf("load item: %w", loadErr)
		}

		if !item.IsOwnedBy(cmd.ActorID) {
			return domainerr.ErrForbidden
		}

		photosBefore, photosErr := photoCount(ctx, h.photos, item.ID())
		if photosErr != nil {
			return photosErr
		}

		before := improvementSnapshot{
			descriptionLen: utf8.RuneCountInString(item.Description()),
			priceKopeks:    item.Price().Kopeks(),
			photoCount:     photosBefore,
		}

		if updateErr := item.Update(params); updateErr != nil {
			return updateErr
		}

		if saveErr := h.items.Save(ctx, item); saveErr != nil {
			return saveErr
		}

		count, countErr := photoCount(ctx, h.photos, item.ID())
		if countErr != nil {
			return countErr
		}

		payload := itemPayload(item, count)
		payload.Attributes = withImprovementFlags(payload.Attributes, before, item, count)

		out.Add(events.New(events.TypeItemUpdated, item.OwnerID(), item.ID(), h.clock.Now()).
			WithPayload(payload))

		updated = item

		return nil
	})
	if err != nil {
		return nil, err
	}

	if updated.Status() == domain.StatusModeration && h.moderate != nil {
		if moderated := h.moderate.HandleOrLog(ctx, updated.ID(), cmd.ActorID); moderated != nil {
			updated = moderated
		}
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

type improvementSnapshot struct {
	descriptionLen int
	priceKopeks    int64
	photoCount     int
}

func withImprovementFlags(
	attributes map[string]string,
	before improvementSnapshot,
	item *domain.Item,
	photoCount int,
) map[string]string {
	flags := make(map[string]string, len(attributes)+3)
	for key, value := range attributes {
		flags[key] = value
	}

	if photoCount > before.photoCount {
		flags[events.AttrPhotoAdded] = events.AttrFlagTrue
	}
	if utf8.RuneCountInString(item.Description()) > before.descriptionLen {
		flags[events.AttrDescriptionAdded] = events.AttrFlagTrue
	}
	if before.priceKopeks == 0 && item.Price().Kopeks() > 0 {
		flags[events.AttrPriceSet] = events.AttrFlagTrue
	}

	return flags
}

package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/module/item/domain"
	"github.com/avito-hack/backend/internal/shared/domainerr"
	"github.com/avito-hack/backend/internal/shared/events"
)

type ModerateItemCommand struct {
	ItemID  uuid.UUID
	ActorID uuid.UUID
}

type ModerateItemHandler struct {
	items      domain.Repository
	photos     domain.PhotoRepository
	provider   domain.ModerationProvider
	log        domain.ModerationLogRepository
	photoBytes PhotoBytesLoader
	tx         TxManager
	clock      Clock
	bus        events.Publisher
}

func NewModerateItemHandler(
	items domain.Repository,
	photos domain.PhotoRepository,
	provider domain.ModerationProvider,
	log domain.ModerationLogRepository,
	photoBytes PhotoBytesLoader,
	tx TxManager,
	clock Clock,
	bus events.Publisher,
) *ModerateItemHandler {
	return &ModerateItemHandler{
		items: items, photos: photos, provider: provider, log: log,
		photoBytes: photoBytes, tx: tx, clock: clock, bus: bus,
	}
}

func (h *ModerateItemHandler) HandleOrLog(ctx context.Context, itemID, actorID uuid.UUID) *domain.Item {
	moderated, err := h.Handle(ctx, ModerateItemCommand{ItemID: itemID, ActorID: actorID})
	if err != nil {
		slog.ErrorContext(ctx, "item moderation check failed, item stays in moderation",
			slog.String("item_id", itemID.String()),
			slog.Any("error", err),
		)

		return nil
	}

	return moderated
}

func (h *ModerateItemHandler) Handle(ctx context.Context, cmd ModerateItemCommand) (*domain.Item, error) {
	item, err := h.items.ByID(ctx, cmd.ItemID)
	if err != nil {
		return nil, fmt.Errorf("load item: %w", err)
	}

	if !item.IsOwnedBy(cmd.ActorID) {
		return nil, domainerr.ErrForbidden
	}

	if item.Status() != domain.StatusModeration {
		return item, nil
	}

	subject, err := h.buildSubject(ctx, item)
	if err != nil {
		return nil, fmt.Errorf("build moderation subject: %w", err)
	}

	result := h.provider.Review(ctx, subject)

	return h.apply(ctx, item, result)
}

func (h *ModerateItemHandler) buildSubject(ctx context.Context, item *domain.Item) (domain.ModerationSubject, error) {
	photos, err := h.photos.ByItemID(ctx, item.ID())
	if err != nil {
		return domain.ModerationSubject{}, fmt.Errorf("list photos: %w", err)
	}

	urls := make([]string, 0, len(photos))

	for _, photo := range photos {
		dataURL, loadErr := h.photoBytes.DataURL(ctx, item.ID(), photo.URL())
		if loadErr != nil {
			return domain.ModerationSubject{}, fmt.Errorf("load photo %s: %w", photo.URL(), loadErr)
		}

		urls = append(urls, dataURL)
	}

	return domain.ModerationSubject{
		ItemID:      item.ID(),
		Title:       item.Title(),
		Description: item.Description(),
		PhotoURLs:   urls,
	}, nil
}

func (h *ModerateItemHandler) apply(
	ctx context.Context,
	item *domain.Item,
	result domain.ModerationResult,
) (*domain.Item, error) {
	var published bool

	err := events.PublishAfterCommit(ctx, h.bus, h.tx, func(ctx context.Context, out *events.Outbox) error {
		locked, err := h.items.ByIDForUpdate(ctx, item.ID())
		if err != nil {
			return fmt.Errorf("load item for update: %w", err)
		}

		if locked.Status() != domain.StatusModeration {
			item = locked

			return nil
		}

		if result.Verdict == domain.ModerationApproved {
			if err := locked.Publish(h.clock.Now()); err != nil {
				return err
			}
			locked.MarkAIVerified()
			published = true
		}

		if err := h.items.Save(ctx, locked); err != nil {
			return err
		}

		entry := domain.NewModerationLogEntry(locked.ID(), result, h.clock.Now())
		if err := h.log.Add(ctx, entry); err != nil {
			return fmt.Errorf("log moderation result: %w", err)
		}

		item = locked

		if published {
			count, countErr := photoCount(ctx, h.photos, item.ID())
			if countErr != nil {
				return countErr
			}

			out.Add(events.New(events.TypeItemPublished, item.OwnerID(), item.ID(), h.clock.Now()).
				WithPayload(itemPayload(item, count)))
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return item, nil
}

package app

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"path"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/module/item/domain"
	"github.com/avito-hack/backend/internal/shared/domainerr"
)

type AddPhotoCommand struct {
	ItemID      uuid.UUID
	ActorID     uuid.UUID
	Content     io.Reader
	ContentType string
}

type AddPhotoHandler struct {
	items   domain.Repository
	photos  domain.PhotoRepository
	storage PhotoStorage
	tx      TxManager
	clock   Clock
}

func NewAddPhotoHandler(
	items domain.Repository,
	photos domain.PhotoRepository,
	storage PhotoStorage,
	tx TxManager,
	clock Clock,
) *AddPhotoHandler {
	return &AddPhotoHandler{items: items, photos: photos, storage: storage, tx: tx, clock: clock}
}

func (h *AddPhotoHandler) Handle(ctx context.Context, cmd AddPhotoCommand) (*domain.Photo, error) {
	stored, err := h.storage.Save(ctx, cmd.ItemID, cmd.Content, cmd.ContentType)
	if err != nil {
		return nil, fmt.Errorf("store photo: %w", err)
	}

	photo, err := h.persist(ctx, cmd, stored.Name)
	if err != nil {
		discardPhotoFile(ctx, h.storage, cmd.ItemID, stored.Name)

		return nil, err
	}

	return photo, nil
}

func (h *AddPhotoHandler) persist(ctx context.Context, cmd AddPhotoCommand, name string) (*domain.Photo, error) {
	var photo *domain.Photo

	err := h.tx.WithTx(ctx, func(ctx context.Context) error {
		item, err := h.items.ByIDForUpdate(ctx, cmd.ItemID)
		if err != nil {
			return fmt.Errorf("load item: %w", err)
		}

		if !item.IsOwnedBy(cmd.ActorID) {
			return domainerr.ErrForbidden
		}

		count, err := h.photos.CountByItemID(ctx, cmd.ItemID)
		if err != nil {
			return fmt.Errorf("count photos: %w", err)
		}

		if count >= domain.MaxPhotosPerItem {
			return domain.ErrTooManyPhotos()
		}

		created, err := domain.NewPhoto(domain.NewPhotoParams{
			ItemID:   cmd.ItemID,
			URL:      h.storage.URL(cmd.ItemID, name),
			Position: count,
			Now:      h.clock.Now(),
		})
		if err != nil {
			return err
		}

		if err := h.photos.Add(ctx, created); err != nil {
			return err
		}
		photo = created

		return nil
	})
	if err != nil {
		return nil, err
	}

	return photo, nil
}

func discardPhotoFile(ctx context.Context, storage PhotoStorage, itemID uuid.UUID, name string) {
	if storage == nil || name == "" {
		return
	}

	if err := storage.Delete(ctx, itemID, name); err != nil {
		slog.ErrorContext(ctx, "orphaned photo file left on disk",
			slog.String("item_id", itemID.String()),
			slog.String("name", name),
			slog.Any("error", err),
		)
	}
}

type ListPhotosHandler struct {
	photos domain.PhotoRepository
}

func NewListPhotosHandler(photos domain.PhotoRepository) *ListPhotosHandler {
	return &ListPhotosHandler{photos: photos}
}

func (h *ListPhotosHandler) Handle(ctx context.Context, itemID uuid.UUID) ([]*domain.Photo, error) {
	photos, err := h.photos.ByItemID(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("list photos: %w", err)
	}

	return photos, nil
}

type DeletePhotoCommand struct {
	ItemID  uuid.UUID
	PhotoID uuid.UUID
	ActorID uuid.UUID
}

type DeletePhotoHandler struct {
	items   domain.Repository
	photos  domain.PhotoRepository
	storage PhotoStorage
	tx      TxManager
}

func NewDeletePhotoHandler(
	items domain.Repository,
	photos domain.PhotoRepository,
	storage PhotoStorage,
	tx TxManager,
) *DeletePhotoHandler {
	return &DeletePhotoHandler{items: items, photos: photos, storage: storage, tx: tx}
}

func (h *DeletePhotoHandler) Handle(ctx context.Context, cmd DeletePhotoCommand) error {
	var removed string

	err := h.tx.WithTx(ctx, func(ctx context.Context) error {
		item, err := h.items.ByIDForUpdate(ctx, cmd.ItemID)
		if err != nil {
			return fmt.Errorf("load item: %w", err)
		}

		if !item.IsOwnedBy(cmd.ActorID) {
			return domainerr.ErrForbidden
		}

		url, err := h.photos.DeleteByID(ctx, cmd.ItemID, cmd.PhotoID)
		if err != nil {
			return err
		}
		removed = url

		return nil
	})
	if err != nil {
		return err
	}

	if removed != "" {
		discardPhotoFile(ctx, h.storage, cmd.ItemID, path.Base(removed))
	}

	return nil
}

package app

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/module/favorite/domain"
	"github.com/avito-hack/backend/internal/shared/events"
)

type AddFavoriteCommand struct {
	UserID uuid.UUID
	ItemID uuid.UUID
}

type AddFavoriteHandler struct {
	favorites domain.Repository
	items     ItemChecker
	tx        TxManager
	clock     Clock
	bus       events.Publisher
}

func NewAddFavoriteHandler(
	favorites domain.Repository,
	items ItemChecker,
	tx TxManager,
	clock Clock,
	bus events.Publisher,
) *AddFavoriteHandler {
	return &AddFavoriteHandler{favorites: favorites, items: items, tx: tx, clock: clock, bus: bus}
}

func (h *AddFavoriteHandler) Handle(ctx context.Context, cmd AddFavoriteCommand) error {
	exists, err := h.items.Exists(ctx, cmd.ItemID)
	if err != nil {
		return fmt.Errorf("check item: %w", err)
	}

	if !exists {
		return domain.ErrItemNotFound
	}

	favorite, err := domain.New(cmd.UserID, cmd.ItemID, h.clock.Now())
	if err != nil {
		return err
	}

	return events.PublishAfterCommit(ctx, h.bus, h.tx, func(ctx context.Context, out *events.Outbox) error {
		added, addErr := h.favorites.Add(ctx, favorite)
		if addErr != nil {
			return addErr
		}

		if added {
			out.Add(events.New(events.TypeFavoriteAdded, cmd.UserID, cmd.ItemID, h.clock.Now()))
		}

		return nil
	})
}

type RemoveFavoriteCommand struct {
	UserID uuid.UUID
	ItemID uuid.UUID
}

type RemoveFavoriteHandler struct {
	favorites domain.Repository
	tx        TxManager
}

func NewRemoveFavoriteHandler(favorites domain.Repository, tx TxManager) *RemoveFavoriteHandler {
	return &RemoveFavoriteHandler{favorites: favorites, tx: tx}
}

func (h *RemoveFavoriteHandler) Handle(ctx context.Context, cmd RemoveFavoriteCommand) error {
	return h.tx.WithTx(ctx, func(ctx context.Context) error {
		return h.favorites.Remove(ctx, cmd.UserID, cmd.ItemID)
	})
}

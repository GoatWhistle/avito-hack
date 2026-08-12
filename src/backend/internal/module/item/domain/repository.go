package domain

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Save(ctx context.Context, item *Item) error
	ByID(ctx context.Context, id uuid.UUID) (*Item, error)
	ByIDForUpdate(ctx context.Context, id uuid.UUID) (*Item, error)
	ByDisplayID(ctx context.Context, displayID string) (*Item, error)
}

type PhotoRepository interface {
	Add(ctx context.Context, photo *Photo) error
	ByItemID(ctx context.Context, itemID uuid.UUID) ([]*Photo, error)
	CountByItemID(ctx context.Context, itemID uuid.UUID) (int, error)
	DeleteByDisplayID(ctx context.Context, itemID uuid.UUID, displayID string) (string, error)
}

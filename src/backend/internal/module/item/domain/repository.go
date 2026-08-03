package domain

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Save(ctx context.Context, item *Item) error
	ByID(ctx context.Context, id uuid.UUID) (*Item, error)
	ByIDForUpdate(ctx context.Context, id uuid.UUID) (*Item, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

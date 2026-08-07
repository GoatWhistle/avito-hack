package domain

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Add(ctx context.Context, favorite *Favorite) (bool, error)
	Remove(ctx context.Context, userID, itemID uuid.UUID) error
}

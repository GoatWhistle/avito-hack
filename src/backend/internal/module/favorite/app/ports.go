package app

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/shared/pagination"
)

type TxManager interface {
	WithTx(ctx context.Context, fn func(context.Context) error) error
}

type Clock interface {
	Now() time.Time
}

type ItemChecker interface {
	Exists(ctx context.Context, itemID uuid.UUID) (bool, error)
}

type FavoriteItem struct {
	ItemID         uuid.UUID
	ItemDisplayID  string
	OwnerID        uuid.UUID
	OwnerDisplayID string
	Title          string
	PriceKopeks   int64
	Status        string
	PhotoURL      string
	CreatedAt     time.Time
}

type ListFilter struct {
	UserID uuid.UUID
	Cursor pagination.Cursor
	Limit  int
}

type ReadModel interface {
	List(ctx context.Context, f ListFilter) ([]FavoriteItem, error)
}

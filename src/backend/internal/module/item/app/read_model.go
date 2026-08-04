package app

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/module/item/domain"
	"github.com/avito-hack/backend/internal/shared/pagination"
)

type ListItem struct {
	ID          uuid.UUID
	OwnerID     uuid.UUID
	OwnerName   string
	Title       string
	PriceKopeks int64
	Status      domain.Status
	CreatedAt   time.Time
}

type ListFilter struct {
	Status  domain.Status
	OwnerID uuid.UUID
	Search  string
	Cursor  pagination.Cursor
	Limit   int
}

type ReadModel interface {
	List(ctx context.Context, f ListFilter) ([]ListItem, error)
}

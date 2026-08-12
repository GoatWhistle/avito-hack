package app

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/module/item/domain"
	"github.com/avito-hack/backend/internal/shared/pagination"
)

type ListItem struct {
	ID             uuid.UUID
	DisplayID      string
	OwnerID        uuid.UUID
	OwnerDisplayID string
	OwnerName      string
	Title          string
	PriceKopeks    int64
	Status         domain.Status
	Category       string
	Condition      string
	CreatedAt      time.Time
	IsSeed         bool
	AIVerified     bool
}

type ListSort string

const (
	ListSortNewest    ListSort = "newest"
	ListSortPriceAsc  ListSort = "price_asc"
	ListSortPriceDesc ListSort = "price_desc"
)

func NormalizeListSort(raw string) ListSort {
	switch ListSort(raw) {
	case ListSortPriceAsc:
		return ListSortPriceAsc
	case ListSortPriceDesc:
		return ListSortPriceDesc
	case ListSortNewest:
		return ListSortNewest
	default:
		return ListSortNewest
	}
}

type ListFilter struct {
	Status    domain.Status
	OwnerID   uuid.UUID
	ViewerID  uuid.UUID
	Search    string
	Category  string
	Condition string
	Sort      ListSort
	Cursor    pagination.Cursor
	Limit     int
}

type ReadModel interface {
	List(ctx context.Context, f ListFilter) ([]ListItem, error)
}

package app

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/shared/pagination"
)

type ListFavoritesQuery struct {
	UserID uuid.UUID
	Cursor pagination.Cursor
	Limit  int
}

type ListFavoritesResult struct {
	Items      []FavoriteItem
	NextCursor string
}

type ListFavoritesHandler struct {
	read ReadModel
}

func NewListFavoritesHandler(read ReadModel) *ListFavoritesHandler {
	return &ListFavoritesHandler{read: read}
}

func (h *ListFavoritesHandler) Handle(ctx context.Context, q ListFavoritesQuery) (ListFavoritesResult, error) {
	limit := pagination.NormalizeLimit(q.Limit)

	rows, err := h.read.List(ctx, ListFilter{
		UserID: q.UserID,
		Cursor: q.Cursor,
		Limit:  limit + 1,
	})
	if err != nil {
		return ListFavoritesResult{}, fmt.Errorf("list favorites: %w", err)
	}

	page, next := pagination.Paginate(rows, limit, func(item FavoriteItem) pagination.CursorKey {
		return pagination.CursorKey{CreatedAt: item.CreatedAt, ID: item.ItemID}
	})

	return ListFavoritesResult{Items: page, NextCursor: next}, nil
}

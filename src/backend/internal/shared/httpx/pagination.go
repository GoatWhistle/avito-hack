package httpx

import (
	"net/http"

	"github.com/avito-hack/backend/internal/shared/domainerr"
	"github.com/avito-hack/backend/internal/shared/pagination"
)

type Page struct {
	Limit  int
	Cursor pagination.Cursor
}

func PageFromRequest(r *http.Request) (Page, error) {
	limit, err := IntQuery(r, "limit", pagination.DefaultLimit)
	if err != nil {
		return Page{}, err
	}

	cursor, err := pagination.DecodeCursor(r.URL.Query().Get("cursor"))
	if err != nil {
		return Page{}, domainerr.NewInvalid("cursor", "invalid cursor")
	}

	return Page{Limit: pagination.NormalizeLimit(limit), Cursor: cursor}, nil
}

type ListResponse[T any] struct {
	Items      []T    `json:"items"`
	NextCursor string `json:"next_cursor,omitempty"`
}

func NewListResponse[T any](items []T, next string) ListResponse[T] {
	if items == nil {
		items = []T{}
	}

	return ListResponse[T]{Items: items, NextCursor: next}
}

package app

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/module/item/domain"
	"github.com/avito-hack/backend/internal/shared/auth"
	"github.com/avito-hack/backend/internal/shared/pagination"
)

type ListItemsQuery struct {
	Status  domain.Status
	OwnerID uuid.UUID
	Viewer  auth.Actor
	Search  string
	Cursor  pagination.Cursor
	Limit   int
}

type ListItemsResult struct {
	Items      []ListItem
	NextCursor string
}

type ListItemsHandler struct {
	read   ReadModel
	owners OwnerProvider
}

func NewListItemsHandler(read ReadModel, owners OwnerProvider) *ListItemsHandler {
	return &ListItemsHandler{read: read, owners: owners}
}

func (h *ListItemsHandler) Handle(ctx context.Context, q ListItemsQuery) (ListItemsResult, error) {
	limit := pagination.NormalizeLimit(q.Limit)

	rows, err := h.read.List(ctx, ListFilter{
		Status:   q.Status,
		OwnerID:  q.OwnerID,
		ViewerID: q.Viewer.ID,
		Search:   q.Search,
		Cursor:   q.Cursor,
		Limit:    limit + 1,
	})
	if err != nil {
		return ListItemsResult{}, fmt.Errorf("list items: %w", err)
	}

	var next string
	if len(rows) > limit {
		last := rows[limit-1]
		next = pagination.Cursor{CreatedAt: last.CreatedAt, ID: last.ID}.Encode()
		rows = rows[:limit]
	}

	if err := h.fillOwners(ctx, rows); err != nil {
		return ListItemsResult{}, err
	}

	return ListItemsResult{Items: rows, NextCursor: next}, nil
}

func (h *ListItemsHandler) fillOwners(ctx context.Context, rows []ListItem) error {
	if len(rows) == 0 || h.owners == nil {
		return nil
	}

	ids := make([]uuid.UUID, 0, len(rows))
	seen := make(map[uuid.UUID]struct{}, len(rows))

	for _, row := range rows {
		if _, ok := seen[row.OwnerID]; ok {
			continue
		}
		seen[row.OwnerID] = struct{}{}
		ids = append(ids, row.OwnerID)
	}

	owners, err := h.owners.ByIDs(ctx, ids)
	if err != nil {
		return fmt.Errorf("load owners: %w", err)
	}

	for i := range rows {
		if owner, ok := owners[rows[i].OwnerID]; ok {
			rows[i].OwnerName = owner.DisplayName
		}
	}

	return nil
}

package infra

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/avito-hack/backend/internal/module/item/app"
	"github.com/avito-hack/backend/internal/module/item/domain"
	"github.com/avito-hack/backend/internal/shared/pagination"
)

var cursorTime = time.Date(2024, 2, 2, 0, 0, 0, 0, time.UTC)

func TestBuildListQueryWithoutFilters(t *testing.T) {
	t.Parallel()

	query, args := buildListQuery(app.ListFilter{Limit: 20})

	assert.Contains(t, query, "FROM items")
	assert.Contains(t, query, "WHERE deleted_at IS NULL")
	assert.Contains(t, query, "ORDER BY created_at DESC, id DESC")
	assert.Contains(t, query, "LIMIT $1")
	assert.NotContains(t, query, "status =")
	assert.NotContains(t, query, "owner_id =")
	assert.NotContains(t, query, "ILIKE")
	assert.Equal(t, []any{20}, args)
}

func TestBuildListQueryWithEveryFilter(t *testing.T) {
	t.Parallel()

	ownerID, cursorID := uuid.New(), uuid.New()

	query, args := buildListQuery(app.ListFilter{
		Status:  domain.StatusPublished,
		OwnerID: ownerID,
		Search:  "bike",
		Cursor:  pagination.Cursor{CreatedAt: cursorTime, ID: cursorID},
		Limit:   5,
	})

	assert.Contains(t, query, "status = $1")
	assert.Contains(t, query, "owner_id = $2")
	assert.Contains(t, query, "title ILIKE $3")
	assert.Contains(t, query, "(created_at, id) < ($4, $5)")
	assert.Contains(t, query, "LIMIT $6")
	assert.Equal(t, []any{"published", ownerID, "%bike%", cursorTime, cursorID, 5}, args)
}

func TestBuildListQueryStatusOnly(t *testing.T) {
	t.Parallel()

	query, args := buildListQuery(app.ListFilter{Status: domain.StatusSold, Limit: 3})

	assert.Contains(t, query, "status = $1")
	assert.Contains(t, query, "LIMIT $2")
	assert.Equal(t, []any{"sold", 3}, args)
}

func TestBuildListQueryOwnerOnly(t *testing.T) {
	t.Parallel()

	ownerID := uuid.New()

	query, args := buildListQuery(app.ListFilter{OwnerID: ownerID, Limit: 7})

	assert.Contains(t, query, "owner_id = $1")
	assert.Equal(t, []any{ownerID, 7}, args)
}

func TestBuildListQuerySearchOnly(t *testing.T) {
	t.Parallel()

	query, args := buildListQuery(app.ListFilter{Search: "chair", Limit: 2})

	assert.Contains(t, query, "title ILIKE $1")
	assert.Equal(t, []any{"%chair%", 2}, args)
}

func TestBuildListQueryIgnoresZeroCursor(t *testing.T) {
	t.Parallel()

	query, args := buildListQuery(app.ListFilter{
		Cursor: pagination.Cursor{CreatedAt: cursorTime},
		Limit:  4,
	})

	assert.NotContains(t, query, "(created_at, id) <")
	assert.Equal(t, []any{4}, args)
}

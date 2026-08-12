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
	assert.Contains(t, query, "WHERE i.deleted_at IS NULL")
	assert.Contains(t, query, "status IN ($1, $2)")
	assert.Contains(t, query, "ORDER BY i.created_at DESC, i.id DESC")
	assert.Contains(t, query, "LIMIT $3")
	assert.NotContains(t, query, "status =")
	assert.NotContains(t, query, "owner_id =")
	assert.NotContains(t, query, "ILIKE")
	assert.Equal(t, []any{"published", "sold", 20}, args)
}

func TestBuildListQueryAlwaysRestrictsAnonymousToPublicStatuses(t *testing.T) {
	t.Parallel()

	for _, status := range []domain.Status{domain.StatusDraft, domain.StatusModeration, domain.StatusArchived} {
		query, args := buildListQuery(app.ListFilter{Status: status, Limit: 20})

		assert.Contains(t, query, "status IN ($1, $2)")
		assert.Contains(t, query, "status = $3")
		assert.NotContains(t, query, "owner_id =")
		assert.Equal(t, []any{"published", "sold", status.String(), 20}, args)
	}
}

func TestBuildListQueryViewerSeesOwnItemsInAnyStatus(t *testing.T) {
	t.Parallel()

	viewerID := uuid.New()

	query, args := buildListQuery(app.ListFilter{ViewerID: viewerID, Limit: 20})

	assert.Contains(t, query, "(i.status IN ($1, $2) OR i.owner_id = $3)")
	assert.Equal(t, []any{"published", "sold", viewerID, 20}, args)
}

func TestBuildListQueryOwnerFilterCannotBypassVisibility(t *testing.T) {
	t.Parallel()

	victimID, viewerID := uuid.New(), uuid.New()

	query, args := buildListQuery(app.ListFilter{
		Status: domain.StatusDraft, OwnerID: victimID, ViewerID: viewerID, Limit: 20,
	})

	assert.Contains(t, query, "(i.status IN ($1, $2) OR i.owner_id = $3)")
	assert.Contains(t, query, "status = $4")
	assert.Contains(t, query, "owner_id = $5")
	assert.Equal(t, []any{"published", "sold", viewerID, "draft", victimID, 20}, args)
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

	assert.Contains(t, query, "status = $3")
	assert.Contains(t, query, "owner_id = $4")
	assert.Contains(t, query, "title ILIKE $5")
	assert.Contains(t, query, "(i.created_at, i.id) < ($6, $7)")
	assert.Contains(t, query, "LIMIT $8")
	assert.Equal(t, []any{
		"published", "sold", "published", ownerID, "%bike%", cursorTime, cursorID, 5,
	}, args)
}

func TestBuildListQueryStatusOnly(t *testing.T) {
	t.Parallel()

	query, args := buildListQuery(app.ListFilter{Status: domain.StatusSold, Limit: 3})

	assert.Contains(t, query, "status = $3")
	assert.Contains(t, query, "LIMIT $4")
	assert.Equal(t, []any{"published", "sold", "sold", 3}, args)
}

func TestBuildListQueryOwnerOnly(t *testing.T) {
	t.Parallel()

	ownerID := uuid.New()

	query, args := buildListQuery(app.ListFilter{OwnerID: ownerID, Limit: 7})

	assert.Contains(t, query, "owner_id = $3")
	assert.Equal(t, []any{"published", "sold", ownerID, 7}, args)
}

func TestBuildListQuerySearchOnly(t *testing.T) {
	t.Parallel()

	query, args := buildListQuery(app.ListFilter{Search: "chair", Limit: 2})

	assert.Contains(t, query, "title ILIKE $3")
	assert.Equal(t, []any{"published", "sold", "%chair%", 2}, args)
}

func TestBuildListQueryCategoryAndConditionFilters(t *testing.T) {
	t.Parallel()

	query, args := buildListQuery(app.ListFilter{Category: "electronics", Condition: "used", Limit: 9})

	assert.Contains(t, query, "i.attributes->>'category' = $3")
	assert.Contains(t, query, "i.attributes->>'condition' = $4")
	assert.Equal(t, []any{"published", "sold", "electronics", "used", 9}, args)
}

func TestBuildListQueryPriceAscOrderAndCursor(t *testing.T) {
	t.Parallel()

	cursorID := uuid.New()

	query, args := buildListQuery(app.ListFilter{
		Sort:   app.ListSortPriceAsc,
		Cursor: pagination.Cursor{CreatedAt: cursorTime, ID: cursorID, PriceKopeks: 5000},
		Limit:  6,
	})

	assert.Contains(t, query, "ORDER BY i.price_kopeks ASC, i.id ASC")
	assert.Contains(t, query, "(i.price_kopeks, i.id) > ($3, $4)")
	assert.NotContains(t, query, "(i.created_at, i.id) <")
	assert.NotContains(t, query, "ORDER BY i.created_at")
	assert.Equal(t, []any{"published", "sold", int64(5000), cursorID, 6}, args)
}

func TestBuildListQueryPriceDescOrderAndCursor(t *testing.T) {
	t.Parallel()

	cursorID := uuid.New()

	query, args := buildListQuery(app.ListFilter{
		Sort:   app.ListSortPriceDesc,
		Cursor: pagination.Cursor{CreatedAt: cursorTime, ID: cursorID, PriceKopeks: 7500},
		Limit:  6,
	})

	assert.Contains(t, query, "ORDER BY i.price_kopeks DESC, i.id DESC")
	assert.Contains(t, query, "(i.price_kopeks, i.id) < ($3, $4)")
	assert.Equal(t, []any{"published", "sold", int64(7500), cursorID, 6}, args)
}

func TestBuildListQueryUnknownSortFallsBackToNewest(t *testing.T) {
	t.Parallel()

	cursorID := uuid.New()

	query, args := buildListQuery(app.ListFilter{
		Sort:   app.ListSort("price_asc; DROP TABLE items"),
		Cursor: pagination.Cursor{CreatedAt: cursorTime, ID: cursorID, PriceKopeks: 100},
		Limit:  6,
	})

	assert.Contains(t, query, "ORDER BY i.created_at DESC, i.id DESC")
	assert.Contains(t, query, "(i.created_at, i.id) < ($3, $4)")
	assert.NotContains(t, query, "DROP TABLE")
	assert.NotContains(t, query, "price_kopeks, i.id")
	assert.Equal(t, []any{"published", "sold", cursorTime, cursorID, 6}, args)
}

func TestBuildListQueryExplicitNewestSortMatchesDefault(t *testing.T) {
	t.Parallel()

	defaultQuery, defaultArgs := buildListQuery(app.ListFilter{Limit: 20})
	newestQuery, newestArgs := buildListQuery(app.ListFilter{Sort: app.ListSortNewest, Limit: 20})

	assert.Equal(t, defaultQuery, newestQuery)
	assert.Equal(t, defaultArgs, newestArgs)
}

func TestBuildListQueryIgnoresZeroCursor(t *testing.T) {
	t.Parallel()

	query, args := buildListQuery(app.ListFilter{
		Cursor: pagination.Cursor{CreatedAt: cursorTime},
		Limit:  4,
	})

	assert.NotContains(t, query, "(created_at, id) <")
	assert.Equal(t, []any{"published", "sold", 4}, args)
}

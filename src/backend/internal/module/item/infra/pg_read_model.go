package infra

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/avito-hack/backend/internal/module/item/app"
	"github.com/avito-hack/backend/internal/module/item/domain"
	"github.com/avito-hack/backend/internal/shared/pagination"
	"github.com/avito-hack/backend/internal/shared/postgres"
)

type PgReadModel struct {
	pool *pgxpool.Pool
}

func NewPgReadModel(pool *pgxpool.Pool) *PgReadModel {
	return &PgReadModel{pool: pool}
}

func (m *PgReadModel) List(ctx context.Context, f app.ListFilter) ([]app.ListItem, error) {
	query, args := buildListQuery(f)

	return postgres.QueryAll(ctx, postgres.QuerierFrom(ctx, m.pool), f.Limit, scanListItem, query, args...)
}

func scanListItem(row pgx.Row) (app.ListItem, error) {
	var (
		item   app.ListItem
		status string
	)

	if err := row.Scan(&item.ID, &item.DisplayID, &item.OwnerID, &item.Title,
		&item.PriceKopeks, &status, &item.CreatedAt, &item.IsSeed, &item.AIVerified,
		&item.Category, &item.Condition); err != nil {
		return app.ListItem{}, fmt.Errorf("scan item: %w", err)
	}

	item.Status = statusFrom(status)

	return item, nil
}

func buildListQuery(f app.ListFilter) (query string, args []any) {
	conditions := []string{"i.deleted_at IS NULL"}

	next := func(value any) string {
		args = append(args, value)

		return "$" + strconv.Itoa(len(args))
	}

	conditions = append(conditions, visibilityCondition(f.ViewerID, next))

	if f.Status != "" {
		conditions = append(conditions, "i.status = "+next(f.Status.String()))
	}

	if f.OwnerID != uuid.Nil {
		conditions = append(conditions, "i.owner_id = "+next(f.OwnerID))
	}

	if f.Search != "" {
		conditions = append(conditions, "i.title ILIKE "+next("%"+f.Search+"%"))
	}

	if f.Category != "" {
		conditions = append(conditions, "i.attributes->>'category' = "+next(f.Category))
	}

	if f.Condition != "" {
		conditions = append(conditions, "i.attributes->>'condition' = "+next(f.Condition))
	}

	sort := app.NormalizeListSort(string(f.Sort))

	if !f.Cursor.IsZero() {
		conditions = append(conditions, cursorCondition(sort, f.Cursor, next))
	}

	query = `SELECT i.id, i.display_id, i.owner_id, i.title,
			i.price_kopeks, i.status, i.created_at, i.is_seed, i.ai_verified,
			coalesce(i.attributes->>'category', ''), coalesce(i.attributes->>'condition', '')
		FROM items i
		WHERE ` + strings.Join(conditions, " AND ") + `
		ORDER BY ` + orderByClause(sort) + `
		LIMIT ` + next(f.Limit)

	return query, args
}

func orderByClause(sort app.ListSort) string {
	switch sort {
	case app.ListSortPriceAsc:
		return "i.price_kopeks ASC, i.id ASC"
	case app.ListSortPriceDesc:
		return "i.price_kopeks DESC, i.id DESC"
	case app.ListSortNewest:
		return "i.created_at DESC, i.id DESC"
	default:
		return "i.created_at DESC, i.id DESC"
	}
}

func cursorCondition(sort app.ListSort, cursor pagination.Cursor, next func(any) string) string {
	switch sort {
	case app.ListSortPriceAsc:
		return "(i.price_kopeks, i.id) > (" + next(cursor.PriceKopeks) + ", " + next(cursor.ID) + ")"
	case app.ListSortPriceDesc:
		return "(i.price_kopeks, i.id) < (" + next(cursor.PriceKopeks) + ", " + next(cursor.ID) + ")"
	case app.ListSortNewest:
		return "(i.created_at, i.id) < (" + next(cursor.CreatedAt) + ", " + next(cursor.ID) + ")"
	default:
		return "(i.created_at, i.id) < (" + next(cursor.CreatedAt) + ", " + next(cursor.ID) + ")"
	}
}

func visibilityCondition(viewerID uuid.UUID, next func(any) string) string {
	placeholders := make([]string, 0, len(domain.PublicStatuses))
	for _, status := range domain.PublicStatuses {
		placeholders = append(placeholders, next(status.String()))
	}

	public := "i.status IN (" + strings.Join(placeholders, ", ") + ")"

	if viewerID == uuid.Nil {
		return public
	}

	return "(" + public + " OR i.owner_id = " + next(viewerID) + ")"
}

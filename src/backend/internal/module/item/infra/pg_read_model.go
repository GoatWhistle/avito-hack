package infra

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/avito-hack/backend/internal/module/item/app"
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

	rows, err := postgres.QuerierFrom(ctx, m.pool).Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query items: %w", err)
	}
	defer rows.Close()

	items := make([]app.ListItem, 0, f.Limit)

	for rows.Next() {
		var item app.ListItem
		var status string

		if err := rows.Scan(&item.ID, &item.OwnerID, &item.Title,
			&item.PriceKopeks, &status, &item.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan item: %w", err)
		}

		item.Status = statusFrom(status)
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate items: %w", err)
	}

	return items, nil
}

func buildListQuery(f app.ListFilter) (string, []any) {
	var (
		conditions = []string{"deleted_at IS NULL"}
		args       []any
	)

	next := func(value any) string {
		args = append(args, value)

		return "$" + strconv.Itoa(len(args))
	}

	if f.Status != "" {
		conditions = append(conditions, "status = "+next(f.Status.String()))
	}

	if f.OwnerID != uuid.Nil {
		conditions = append(conditions, "owner_id = "+next(f.OwnerID))
	}

	if f.Search != "" {
		conditions = append(conditions, "title ILIKE "+next("%"+f.Search+"%"))
	}

	if !f.Cursor.IsZero() {
		conditions = append(conditions,
			"(created_at, id) < ("+next(f.Cursor.CreatedAt)+", "+next(f.Cursor.ID)+")")
	}

	query := `SELECT id, owner_id, title, price_kopeks, status, created_at
		FROM items
		WHERE ` + strings.Join(conditions, " AND ") + `
		ORDER BY created_at DESC, id DESC
		LIMIT ` + next(f.Limit)

	return query, args
}

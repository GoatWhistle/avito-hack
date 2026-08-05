package infra

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/avito-hack/backend/internal/module/favorite/app"
	"github.com/avito-hack/backend/internal/shared/postgres"
)

const listQuery = `
	SELECT i.id, i.owner_id, i.title, i.price_kopeks, i.status, f.created_at,
	       coalesce((SELECT p.url FROM item_photos p WHERE p.item_id = i.id ORDER BY p.position LIMIT 1), '')
	FROM favorites f
	JOIN items i ON i.id = f.item_id AND i.deleted_at IS NULL
	WHERE f.user_id = $1 AND ($2::timestamptz IS NULL OR (f.created_at, f.item_id) < ($2::timestamptz, $3::uuid))
	ORDER BY f.created_at DESC, f.item_id DESC
	LIMIT $4`

type PgReadModel struct {
	pool *pgxpool.Pool
}

func NewPgReadModel(pool *pgxpool.Pool) *PgReadModel {
	return &PgReadModel{pool: pool}
}

func (m *PgReadModel) List(ctx context.Context, f app.ListFilter) ([]app.FavoriteItem, error) {
	var cursorAt, cursorID any
	if !f.Cursor.IsZero() {
		cursorAt, cursorID = f.Cursor.CreatedAt, f.Cursor.ID
	}

	rows, err := postgres.QuerierFrom(ctx, m.pool).Query(ctx, listQuery, f.UserID, cursorAt, cursorID, f.Limit)
	if err != nil {
		return nil, fmt.Errorf("query favorites: %w", err)
	}
	defer rows.Close()

	items := make([]app.FavoriteItem, 0, f.Limit)

	for rows.Next() {
		var item app.FavoriteItem

		if err := rows.Scan(&item.ItemID, &item.OwnerID, &item.Title,
			&item.PriceKopeks, &item.Status, &item.CreatedAt, &item.PhotoURL); err != nil {
			return nil, fmt.Errorf("scan favorite: %w", err)
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate favorites: %w", err)
	}

	return items, nil
}

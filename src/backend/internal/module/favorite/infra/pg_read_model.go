package infra

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
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

	return postgres.QueryAll(ctx, postgres.QuerierFrom(ctx, m.pool), f.Limit, scanFavoriteItem,
		listQuery, f.UserID, cursorAt, cursorID, f.Limit)
}

func scanFavoriteItem(row pgx.Row) (app.FavoriteItem, error) {
	var item app.FavoriteItem

	if err := row.Scan(&item.ItemID, &item.OwnerID, &item.Title,
		&item.PriceKopeks, &item.Status, &item.CreatedAt, &item.PhotoURL); err != nil {
		return app.FavoriteItem{}, fmt.Errorf("scan favorite: %w", err)
	}

	return item, nil
}

package infra

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/avito-hack/backend/internal/module/games/domain"
	"github.com/avito-hack/backend/internal/module/games/domain/moreless"
	"github.com/avito-hack/backend/internal/shared/postgres"
)

type PgItemPool struct {
	pool *pgxpool.Pool
}

func NewPgItemPool(pool *pgxpool.Pool) *PgItemPool {
	return &PgItemPool{pool: pool}
}

func (p *PgItemPool) Random(ctx context.Context, exclude []string) (moreless.Item, error) {
	const query = `
		WITH pool AS (
			SELECT i.display_id, i.title, i.price_kopeks,
			       (SELECT ph.url FROM item_photos ph
			         WHERE ph.item_id = i.id
			         ORDER BY ph.position
			         LIMIT 1) AS photo_url
			FROM items i
			WHERE i.status = 'published'
			  AND i.deleted_at IS NULL
			  AND NOT (i.display_id = ANY($1::text[]))
			  AND EXISTS (SELECT 1 FROM item_photos ph WHERE ph.item_id = i.id)
			OFFSET floor(random() * GREATEST((
				SELECT count(*) FROM items i2
				WHERE i2.status = 'published'
				  AND i2.deleted_at IS NULL
				  AND NOT (i2.display_id = ANY($1::text[]))
				  AND EXISTS (SELECT 1 FROM item_photos ph WHERE ph.item_id = i2.id)
			), 1))
			LIMIT 1
		)
		SELECT display_id, title, price_kopeks, coalesce(photo_url, '') FROM pool`

	if exclude == nil {
		exclude = []string{}
	}

	var item moreless.Item

	err := postgres.QuerierFrom(ctx, p.pool).
		QueryRow(ctx, query, exclude).
		Scan(&item.DisplayID, &item.Title, &item.PriceKopeks, &item.PhotoURL)
	if errors.Is(err, pgx.ErrNoRows) {
		return moreless.Item{}, domain.ErrNoItemsInPool
	}
	if err != nil {
		return moreless.Item{}, fmt.Errorf("draw random game item: %w", err)
	}

	return item, nil
}

func (p *PgItemPool) RandomNear(
	ctx context.Context,
	seen []string,
	referenceKopeks int64,
	minRatio, maxRatio float64,
) (moreless.Item, error) {
	const query = `
		SELECT i.display_id, i.title, i.price_kopeks,
		       coalesce((SELECT ph.url FROM item_photos ph
		         WHERE ph.item_id = i.id
		         ORDER BY ph.position
		         LIMIT 1), '') AS photo_url
		FROM items i
		WHERE i.status = 'published'
		  AND i.deleted_at IS NULL
		  AND NOT (i.display_id = ANY($1::text[]))
		  AND i.price_kopeks BETWEEN $2 AND $3
		  AND EXISTS (SELECT 1 FROM item_photos ph WHERE ph.item_id = i.id)
		ORDER BY random()
		LIMIT 1`

	if seen == nil {
		seen = []string{}
	}

	lo := int64(math.Floor(float64(referenceKopeks) * minRatio))
	hi := int64(math.Ceil(float64(referenceKopeks) * maxRatio))

	if lo < 0 {
		lo = 0
	}

	if hi < lo {
		hi = lo
	}

	var item moreless.Item

	err := postgres.QuerierFrom(ctx, p.pool).
		QueryRow(ctx, query, seen, lo, hi).
		Scan(&item.DisplayID, &item.Title, &item.PriceKopeks, &item.PhotoURL)
	if errors.Is(err, pgx.ErrNoRows) {
		return moreless.Item{}, domain.ErrNoItemsInPool
	}
	if err != nil {
		return moreless.Item{}, fmt.Errorf("draw nearby game item: %w", err)
	}

	return item, nil
}

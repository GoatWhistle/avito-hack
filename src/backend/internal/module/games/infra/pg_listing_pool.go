package infra

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/avito-hack/backend/internal/module/games/domain/raccoonjump"
	"github.com/avito-hack/backend/internal/shared/postgres"
)

type PgListingPool struct {
	pool *pgxpool.Pool
}

func NewPgListingPool(pool *pgxpool.Pool) *PgListingPool {
	return &PgListingPool{pool: pool}
}

const randomListingsQuery = `
	SELECT i.display_id, i.title, i.price_kopeks,
	       coalesce((SELECT ph.url FROM item_photos ph
	         WHERE ph.item_id = i.id
	         ORDER BY ph.position
	         LIMIT 1), '') AS photo_url
	FROM items i
	WHERE i.status = 'published'
	  AND i.deleted_at IS NULL
	  AND EXISTS (SELECT 1 FROM item_photos ph WHERE ph.item_id = i.id)
	ORDER BY random()
	LIMIT $1`

func (p *PgListingPool) RandomListings(ctx context.Context, limit int) ([]raccoonjump.Listing, error) {
	rows, err := postgres.QuerierFrom(ctx, p.pool).Query(ctx, randomListingsQuery, limit)
	if err != nil {
		return nil, fmt.Errorf("query raccoonjump listings: %w", err)
	}
	defer rows.Close()

	listings := make([]raccoonjump.Listing, 0, limit)

	for rows.Next() {
		var listing raccoonjump.Listing

		if err := rows.Scan(&listing.DisplayID, &listing.Title, &listing.PriceKopeks, &listing.PhotoURL); err != nil {
			return nil, fmt.Errorf("scan raccoonjump listing: %w", err)
		}

		listings = append(listings, listing)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate raccoonjump listings: %w", err)
	}

	return listings, nil
}

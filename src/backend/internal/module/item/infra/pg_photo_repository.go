package infra

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/avito-hack/backend/internal/module/item/domain"
	"github.com/avito-hack/backend/internal/shared/postgres"
)

type PgPhotoRepository struct {
	pool *pgxpool.Pool
}

func NewPgPhotoRepository(pool *pgxpool.Pool) *PgPhotoRepository {
	return &PgPhotoRepository{pool: pool}
}

func (r *PgPhotoRepository) Add(ctx context.Context, photo *domain.Photo) error {
	const query = `
		INSERT INTO item_photos (id, item_id, url, position, created_at)
		VALUES ($1, $2, $3, $4, $5)`

	_, err := postgres.QuerierFrom(ctx, r.pool).Exec(ctx, query,
		photo.ID(), photo.ItemID(), photo.URL(), photo.Position(), photo.CreatedAt())
	if err != nil {
		return fmt.Errorf("insert photo: %w", err)
	}

	return nil
}

func (r *PgPhotoRepository) ByItemID(ctx context.Context, itemID uuid.UUID) ([]*domain.Photo, error) {
	const query = `
		SELECT id, item_id, url, position, created_at
		FROM item_photos
		WHERE item_id = $1
		ORDER BY position`

	rows, err := postgres.QuerierFrom(ctx, r.pool).Query(ctx, query, itemID)
	if err != nil {
		return nil, fmt.Errorf("query photos: %w", err)
	}
	defer rows.Close()

	photos := make([]*domain.Photo, 0, domain.MaxPhotosPerItem)

	for rows.Next() {
		photo, scanErr := scanPhoto(rows)
		if scanErr != nil {
			return nil, scanErr
		}

		photos = append(photos, photo)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate photos: %w", err)
	}

	return photos, nil
}

func (r *PgPhotoRepository) CountByItemID(ctx context.Context, itemID uuid.UUID) (int, error) {
	const query = `SELECT count(*) FROM item_photos WHERE item_id = $1`

	var count int
	if err := postgres.QuerierFrom(ctx, r.pool).QueryRow(ctx, query, itemID).Scan(&count); err != nil {
		return 0, fmt.Errorf("count photos: %w", err)
	}

	return count, nil
}

func (r *PgPhotoRepository) DeleteByID(ctx context.Context, itemID, photoID uuid.UUID) (string, error) {
	const query = `DELETE FROM item_photos WHERE item_id = $1 AND id = $2 RETURNING url`

	var url string
	err := postgres.QuerierFrom(ctx, r.pool).QueryRow(ctx, query, itemID, photoID).Scan(&url)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", domain.ErrPhotoNotFound
	}
	if err != nil {
		return "", fmt.Errorf("delete photo: %w", err)
	}

	return url, nil
}

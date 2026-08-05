package infra

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/avito-hack/backend/internal/module/favorite/domain"
	"github.com/avito-hack/backend/internal/shared/postgres"
)

type PgRepository struct {
	pool *pgxpool.Pool
}

func NewPgRepository(pool *pgxpool.Pool) *PgRepository {
	return &PgRepository{pool: pool}
}

func (r *PgRepository) Add(ctx context.Context, favorite *domain.Favorite) (bool, error) {
	const query = `
		INSERT INTO favorites (user_id, item_id, created_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, item_id) DO NOTHING`

	tag, err := postgres.QuerierFrom(ctx, r.pool).Exec(ctx, query,
		favorite.UserID(), favorite.ItemID(), favorite.CreatedAt())
	if err != nil {
		return false, fmt.Errorf("insert favorite: %w", err)
	}

	return tag.RowsAffected() > 0, nil
}

func (r *PgRepository) Remove(ctx context.Context, userID, itemID uuid.UUID) error {
	const query = `DELETE FROM favorites WHERE user_id = $1 AND item_id = $2`

	tag, err := postgres.QuerierFrom(ctx, r.pool).Exec(ctx, query, userID, itemID)
	if err != nil {
		return fmt.Errorf("delete favorite: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return domain.ErrFavoriteNotFound
	}

	return nil
}

func (r *PgRepository) Exists(ctx context.Context, userID, itemID uuid.UUID) (bool, error) {
	const query = `SELECT EXISTS (SELECT 1 FROM favorites WHERE user_id = $1 AND item_id = $2)`

	var exists bool
	if err := postgres.QuerierFrom(ctx, r.pool).QueryRow(ctx, query, userID, itemID).Scan(&exists); err != nil {
		return false, fmt.Errorf("check favorite: %w", err)
	}

	return exists, nil
}

type PgItemChecker struct {
	pool *pgxpool.Pool
}

func NewPgItemChecker(pool *pgxpool.Pool) *PgItemChecker {
	return &PgItemChecker{pool: pool}
}

func (c *PgItemChecker) Exists(ctx context.Context, itemID uuid.UUID) (bool, error) {
	const query = `SELECT EXISTS (SELECT 1 FROM items WHERE id = $1 AND deleted_at IS NULL)`

	var exists bool
	if err := postgres.QuerierFrom(ctx, c.pool).QueryRow(ctx, query, itemID).Scan(&exists); err != nil {
		return false, fmt.Errorf("check item: %w", err)
	}

	return exists, nil
}

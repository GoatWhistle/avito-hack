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

const itemColumns = `id, owner_id, title, description, price_kopeks, status, attributes, created_at, updated_at`

type PgRepository struct {
	pool *pgxpool.Pool
}

func NewPgRepository(pool *pgxpool.Pool) *PgRepository {
	return &PgRepository{pool: pool}
}

func (r *PgRepository) Save(ctx context.Context, item *domain.Item) error {
	attributes, err := encodeAttributes(item.Attributes())
	if err != nil {
		return err
	}

	const query = `
		INSERT INTO items (id, owner_id, title, description, price_kopeks, status, attributes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (id) DO UPDATE SET
			title = EXCLUDED.title,
			description = EXCLUDED.description,
			price_kopeks = EXCLUDED.price_kopeks,
			status = EXCLUDED.status,
			attributes = EXCLUDED.attributes,
			updated_at = EXCLUDED.updated_at`

	_, err = postgres.QuerierFrom(ctx, r.pool).Exec(ctx, query,
		item.ID(),
		item.OwnerID(),
		item.Title(),
		item.Description(),
		item.Price().Kopeks(),
		item.Status().String(),
		attributes,
		item.CreatedAt(),
		item.UpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf("save item: %w", err)
	}

	return nil
}

func (r *PgRepository) ByID(ctx context.Context, id uuid.UUID) (*domain.Item, error) {
	query := `SELECT ` + itemColumns + ` FROM items WHERE id = $1 AND deleted_at IS NULL`

	return r.queryOne(ctx, query, id)
}

func (r *PgRepository) ByIDForUpdate(ctx context.Context, id uuid.UUID) (*domain.Item, error) {
	query := `SELECT ` + itemColumns + ` FROM items WHERE id = $1 AND deleted_at IS NULL FOR UPDATE`

	return r.queryOne(ctx, query, id)
}

func (r *PgRepository) Delete(ctx context.Context, id uuid.UUID) error {
	const query = `UPDATE items SET deleted_at = now() WHERE id = $1 AND deleted_at IS NULL`

	tag, err := postgres.QuerierFrom(ctx, r.pool).Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete item: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return domain.ErrItemNotFound
	}

	return nil
}

func (r *PgRepository) queryOne(ctx context.Context, query string, args ...any) (*domain.Item, error) {
	row := postgres.QuerierFrom(ctx, r.pool).QueryRow(ctx, query, args...)

	item, err := scanItem(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrItemNotFound
		}

		return nil, fmt.Errorf("query item: %w", err)
	}

	return item, nil
}

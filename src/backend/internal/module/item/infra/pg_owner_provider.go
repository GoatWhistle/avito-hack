package infra

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/avito-hack/backend/internal/module/item/app"
	"github.com/avito-hack/backend/internal/shared/postgres"
)

type PgOwnerProvider struct {
	pool *pgxpool.Pool
}

func NewPgOwnerProvider(pool *pgxpool.Pool) *PgOwnerProvider {
	return &PgOwnerProvider{pool: pool}
}

func (p *PgOwnerProvider) ByIDs(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]app.OwnerView, error) {
	if len(ids) == 0 {
		return map[uuid.UUID]app.OwnerView{}, nil
	}

	const query = `SELECT id, display_name FROM users WHERE id = ANY($1)`

	rows, err := postgres.QuerierFrom(ctx, p.pool).Query(ctx, query, ids)
	if err != nil {
		return nil, fmt.Errorf("query owners: %w", err)
	}
	defer rows.Close()

	owners := make(map[uuid.UUID]app.OwnerView, len(ids))

	for rows.Next() {
		var owner app.OwnerView
		if err := rows.Scan(&owner.ID, &owner.DisplayName); err != nil {
			return nil, fmt.Errorf("scan owner: %w", err)
		}

		owners[owner.ID] = owner
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate owners: %w", err)
	}

	return owners, nil
}

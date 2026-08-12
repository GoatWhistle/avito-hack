package infra

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

	const query = `SELECT id, display_id, full_name FROM users WHERE id = ANY($1) AND deleted_at IS NULL`

	rows, err := postgres.QueryAll(ctx, postgres.QuerierFrom(ctx, p.pool), len(ids), scanOwner, query, ids)
	if err != nil {
		return nil, err
	}

	owners := make(map[uuid.UUID]app.OwnerView, len(rows))
	for _, owner := range rows {
		owners[owner.ID] = owner
	}

	return owners, nil
}

func scanOwner(row pgx.Row) (app.OwnerView, error) {
	var owner app.OwnerView

	if err := row.Scan(&owner.ID, &owner.DisplayID, &owner.DisplayName); err != nil {
		return app.OwnerView{}, fmt.Errorf("scan owner: %w", err)
	}

	return owner, nil
}

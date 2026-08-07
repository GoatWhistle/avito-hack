package infra

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/avito-hack/backend/internal/shared/events"
	"github.com/avito-hack/backend/internal/shared/postgres"
)

type PgEventDeduplicator struct {
	pool *pgxpool.Pool
}

func NewPgEventDeduplicator(pool *pgxpool.Pool) *PgEventDeduplicator {
	return &PgEventDeduplicator{pool: pool}
}

func (d *PgEventDeduplicator) Claim(ctx context.Context, id uuid.UUID, e events.Event) (bool, error) {
	const query = `INSERT INTO pet_processed_events (
		event_id, event_type, user_id, subject_id, occurred_at
	) VALUES ($1, $2, $3, NULLIF($4, '00000000-0000-0000-0000-000000000000'::uuid), $5)
	ON CONFLICT DO NOTHING`

	tag, err := postgres.QuerierFrom(ctx, d.pool).Exec(ctx, query,
		id, string(e.Type), e.UserID, e.SubjectID, e.OccurredAt,
	)
	if err != nil {
		return false, fmt.Errorf("insert processed pet event: %w", err)
	}

	return tag.RowsAffected() == 1, nil
}

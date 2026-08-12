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

type PgModerationLogRepository struct {
	pool *pgxpool.Pool
}

func NewPgModerationLogRepository(pool *pgxpool.Pool) *PgModerationLogRepository {
	return &PgModerationLogRepository{pool: pool}
}

func (r *PgModerationLogRepository) Add(ctx context.Context, entry domain.ModerationLogEntry) error {
	const query = `
		INSERT INTO item_moderation_log (id, item_id, verdict, reason, provider, checked_at)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := postgres.QuerierFrom(ctx, r.pool).Exec(ctx, query,
		entry.ID, entry.ItemID, string(entry.Verdict), entry.Reason, entry.Provider, entry.CheckedAt)
	if err != nil {
		return fmt.Errorf("insert moderation log entry: %w", err)
	}

	return nil
}

func (r *PgModerationLogRepository) LatestByItemID(
	ctx context.Context,
	itemID uuid.UUID,
) (*domain.ModerationLogEntry, error) {
	const query = `
		SELECT id, item_id, verdict, reason, provider, checked_at
		FROM item_moderation_log
		WHERE item_id = $1
		ORDER BY checked_at DESC
		LIMIT 1`

	var (
		entry   domain.ModerationLogEntry
		verdict string
	)

	err := postgres.QuerierFrom(ctx, r.pool).QueryRow(ctx, query, itemID).Scan(
		&entry.ID, &entry.ItemID, &verdict, &entry.Reason, &entry.Provider, &entry.CheckedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil //nolint:nilnil
	}
	if err != nil {
		return nil, fmt.Errorf("query latest moderation entry: %w", err)
	}

	entry.Verdict = domain.ModerationVerdict(verdict)

	return &entry, nil
}

package infra

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/avito-hack/backend/internal/module/pet/domain"
	"github.com/avito-hack/backend/internal/shared/postgres"
)

const uniqueViolationCode = "23505"

type PgXPEventRepository struct {
	pool *pgxpool.Pool
}

func NewPgXPEventRepository(pool *pgxpool.Pool) *PgXPEventRepository {
	return &PgXPEventRepository{pool: pool}
}

func (r *PgXPEventRepository) Append(ctx context.Context, event domain.XPEvent) error {
	const query = `INSERT INTO xp_events (id, user_id, action, subject_id, amount, created_at)
		VALUES ($1,$2,$3,$4,$5,$6)`

	_, err := postgres.QuerierFrom(ctx, r.pool).Exec(ctx, query,
		event.ID, event.UserID, string(event.Action), event.SubjectID, event.Amount, event.CreatedAt,
	)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == uniqueViolationCode {
		return domain.ErrDuplicateAction
	}
	if err != nil {
		return fmt.Errorf("append xp event: %w", err)
	}

	return nil
}

func (r *PgXPEventRepository) CountSince(
	ctx context.Context,
	userID uuid.UUID,
	action domain.Action,
	since time.Time,
) (int, error) {
	const query = `SELECT count(*) FROM xp_events
		WHERE user_id = $1 AND action = $2 AND created_at >= $3`

	var count int
	err := postgres.QuerierFrom(ctx, r.pool).
		QueryRow(ctx, query, userID, string(action), since).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count xp events: %w", err)
	}

	return count, nil
}

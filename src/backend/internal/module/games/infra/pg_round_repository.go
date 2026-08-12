package infra

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/avito-hack/backend/internal/module/games/domain"
	"github.com/avito-hack/backend/internal/shared/postgres"
)

const roundColumns = `id, display_id, user_id, game_slug, state, streak, best_streak,
		payload, created_at, updated_at`

type PgRoundRepository struct {
	pool *pgxpool.Pool
}

func NewPgRoundRepository(pool *pgxpool.Pool) *PgRoundRepository {
	return &PgRoundRepository{pool: pool}
}

func (r *PgRoundRepository) Save(ctx context.Context, round *domain.Round) error {
	const query = `
		INSERT INTO game_rounds (` + roundColumns + `)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (id) DO UPDATE SET
			state = EXCLUDED.state,
			streak = EXCLUDED.streak,
			best_streak = EXCLUDED.best_streak,
			payload = EXCLUDED.payload,
			updated_at = EXCLUDED.updated_at`

	_, err := postgres.QuerierFrom(ctx, r.pool).Exec(ctx, query,
		round.ID(), round.DisplayID(), round.UserID(), round.GameSlug(), round.State().String(),
		round.Streak(), round.BestStreak(), []byte(round.Payload()), round.CreatedAt(), round.UpdatedAt())
	if err != nil {
		return fmt.Errorf("save game round: %w", err)
	}

	return nil
}

func (r *PgRoundRepository) ByDisplayID(ctx context.Context, displayID string) (*domain.Round, error) {
	const query = `SELECT ` + roundColumns + ` FROM game_rounds WHERE display_id = $1`

	return scanRound(postgres.QuerierFrom(ctx, r.pool).QueryRow(ctx, query, displayID))
}

func (r *PgRoundRepository) ActiveByUser(
	ctx context.Context,
	userID uuid.UUID,
	gameSlug string,
) (*domain.Round, error) {
	const query = `
		SELECT ` + roundColumns + ` FROM game_rounds
		WHERE user_id = $1 AND game_slug = $2 AND state = 'active'
		ORDER BY created_at DESC
		LIMIT 1`

	return scanRound(postgres.QuerierFrom(ctx, r.pool).QueryRow(ctx, query, userID, gameSlug))
}

func scanRound(row pgx.Row) (*domain.Round, error) {
	var (
		id         uuid.UUID
		displayID  string
		userID     uuid.UUID
		gameSlug   string
		state      string
		streak     int
		bestStreak int
		payload    []byte
		createdAt  time.Time
		updatedAt  time.Time
	)

	if err := row.Scan(&id, &displayID, &userID, &gameSlug, &state, &streak, &bestStreak,
		&payload, &createdAt, &updatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrRoundNotFound
		}

		return nil, fmt.Errorf("scan game round: %w", err)
	}

	return domain.RestoreRound(id, displayID, userID, gameSlug, domain.State(state),
		streak, bestStreak, json.RawMessage(payload), createdAt, updatedAt), nil
}

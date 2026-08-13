package infra

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/avito-hack/backend/internal/module/games/domain"
	"github.com/avito-hack/backend/internal/shared/postgres"
)

type PgProgressRepository struct {
	pool *pgxpool.Pool
}

func NewPgProgressRepository(pool *pgxpool.Pool) *PgProgressRepository {
	return &PgProgressRepository{pool: pool}
}

func (r *PgProgressRepository) Daily(
	ctx context.Context,
	userID uuid.UUID,
	gameSlug string,
	day domain.Day,
) (domain.DailyProgress, error) {
	const query = `
		SELECT attempts, best_streak FROM game_daily_progress
		WHERE user_id = $1 AND game_slug = $2 AND day = $3`

	var progress domain.DailyProgress

	err := postgres.QuerierFrom(ctx, r.pool).
		QueryRow(ctx, query, userID, gameSlug, day.Time()).
		Scan(&progress.Attempts, &progress.BestStreak)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.DailyProgress{}, nil
	}
	if err != nil {
		return domain.DailyProgress{}, fmt.Errorf("load daily progress: %w", err)
	}

	return progress, nil
}

func (r *PgProgressRepository) SaveDaily(
	ctx context.Context,
	userID uuid.UUID,
	gameSlug string,
	day domain.Day,
	progress domain.DailyProgress,
) error {
	const query = `
		INSERT INTO game_daily_progress (user_id, game_slug, day, attempts, best_streak)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id, game_slug, day) DO UPDATE SET
			attempts = EXCLUDED.attempts,
			best_streak = GREATEST(game_daily_progress.best_streak, EXCLUDED.best_streak)`

	_, err := postgres.QuerierFrom(ctx, r.pool).Exec(ctx, query,
		userID, gameSlug, day.Time(), progress.Attempts, progress.BestStreak)
	if err != nil {
		return fmt.Errorf("save daily progress: %w", err)
	}

	return nil
}

func (r *PgProgressRepository) DailyAny(
	ctx context.Context,
	userID uuid.UUID,
	day domain.Day,
) (domain.DailyProgress, error) {
	const query = `
		SELECT attempts FROM game_daily_days
		WHERE user_id = $1 AND day = $2`

	var progress domain.DailyProgress

	err := postgres.QuerierFrom(ctx, r.pool).
		QueryRow(ctx, query, userID, day.Time()).
		Scan(&progress.Attempts)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.DailyProgress{}, nil
	}
	if err != nil {
		return domain.DailyProgress{}, fmt.Errorf("load global daily progress: %w", err)
	}

	return progress, nil
}

func (r *PgProgressRepository) SaveDailyAny(
	ctx context.Context,
	userID uuid.UUID,
	day domain.Day,
	progress domain.DailyProgress,
) error {
	const query = `
		INSERT INTO game_daily_days (user_id, day, attempts)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, day) DO UPDATE SET attempts = EXCLUDED.attempts`

	_, err := postgres.QuerierFrom(ctx, r.pool).Exec(ctx, query, userID, day.Time(), progress.Attempts)
	if err != nil {
		return fmt.Errorf("save global daily progress: %w", err)
	}

	return nil
}

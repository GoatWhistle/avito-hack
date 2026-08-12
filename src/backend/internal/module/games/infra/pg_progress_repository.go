package infra

import (
	"context"
	"errors"
	"fmt"
	"time"

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

func (r *PgProgressRepository) Streak(
	ctx context.Context,
	userID uuid.UUID,
	gameSlug string,
) (domain.Streak, error) {
	const query = `
		SELECT current_days, best_days, last_day, reward_claimed_at FROM game_streaks
		WHERE user_id = $1 AND game_slug = $2`

	var (
		streak    domain.Streak
		lastDay   *time.Time
		claimedAt *time.Time
	)

	err := postgres.QuerierFrom(ctx, r.pool).
		QueryRow(ctx, query, userID, gameSlug).
		Scan(&streak.CurrentDays, &streak.BestDays, &lastDay, &claimedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Streak{}, nil
	}
	if err != nil {
		return domain.Streak{}, fmt.Errorf("load game streak: %w", err)
	}

	if lastDay != nil {
		day := domain.DayOf(*lastDay)
		streak.LastDay = &day
	}

	streak.RewardClaimedAt = claimedAt

	return streak, nil
}

func (r *PgProgressRepository) SaveStreak(
	ctx context.Context,
	userID uuid.UUID,
	gameSlug string,
	streak domain.Streak,
) error {
	const query = `
		INSERT INTO game_streaks (user_id, game_slug, current_days, best_days, last_day, reward_claimed_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (user_id, game_slug) DO UPDATE SET
			current_days = EXCLUDED.current_days,
			best_days = EXCLUDED.best_days,
			last_day = EXCLUDED.last_day,
			reward_claimed_at = EXCLUDED.reward_claimed_at`

	var lastDay *time.Time
	if streak.LastDay != nil {
		day := streak.LastDay.Time()
		lastDay = &day
	}

	_, err := postgres.QuerierFrom(ctx, r.pool).Exec(ctx, query,
		userID, gameSlug, streak.CurrentDays, streak.BestDays, lastDay, streak.RewardClaimedAt)
	if err != nil {
		return fmt.Errorf("save game streak: %w", err)
	}

	return nil
}

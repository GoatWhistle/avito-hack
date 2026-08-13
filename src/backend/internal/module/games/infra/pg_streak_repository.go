package infra

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/avito-hack/backend/internal/module/games/domain"
	"github.com/avito-hack/backend/internal/shared/postgres"
)

const selectStreak = `
	SELECT current_days, best_days, last_day, reward_claimed_at FROM game_streaks
	WHERE user_id = $1`

func (r *PgProgressRepository) Streak(ctx context.Context, userID uuid.UUID) (domain.Streak, error) {
	streak, err := scanStreak(postgres.QuerierFrom(ctx, r.pool).QueryRow(ctx, selectStreak, userID))
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Streak{}, nil
	}
	if err != nil {
		return domain.Streak{}, fmt.Errorf("load game streak: %w", err)
	}

	return streak, nil
}

func (r *PgProgressRepository) StreakForUpdate(ctx context.Context, userID uuid.UUID) (domain.Streak, error) {
	const ensure = `
		INSERT INTO game_streaks (user_id, current_days, best_days)
		VALUES ($1, 0, 0)
		ON CONFLICT (user_id) DO NOTHING`

	querier := postgres.QuerierFrom(ctx, r.pool)

	if _, err := querier.Exec(ctx, ensure, userID); err != nil {
		return domain.Streak{}, fmt.Errorf("ensure game streak: %w", err)
	}

	streak, err := scanStreak(querier.QueryRow(ctx, selectStreak+" FOR UPDATE", userID))
	if err != nil {
		return domain.Streak{}, fmt.Errorf("lock game streak: %w", err)
	}

	return streak, nil
}

func (r *PgProgressRepository) SaveStreak(ctx context.Context, userID uuid.UUID, streak domain.Streak) error {
	const query = `
		INSERT INTO game_streaks (user_id, current_days, best_days, last_day, reward_claimed_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id) DO UPDATE SET
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
		userID, streak.CurrentDays, streak.BestDays, lastDay, streak.RewardClaimedAt)
	if err != nil {
		return fmt.Errorf("save game streak: %w", err)
	}

	return nil
}

func scanStreak(row pgx.Row) (domain.Streak, error) {
	var (
		streak    domain.Streak
		lastDay   *time.Time
		claimedAt *time.Time
	)

	if err := row.Scan(&streak.CurrentDays, &streak.BestDays, &lastDay, &claimedAt); err != nil {
		return domain.Streak{}, err
	}

	if lastDay != nil {
		day := domain.DayOf(*lastDay)
		streak.LastDay = &day
	}

	streak.RewardClaimedAt = claimedAt

	return streak, nil
}

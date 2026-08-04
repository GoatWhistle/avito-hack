package infra

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/avito-hack/backend/internal/module/pet/domain"
	"github.com/avito-hack/backend/internal/shared/domainerr"
	"github.com/avito-hack/backend/internal/shared/postgres"
)

const petColumns = `id, user_id, name, stage, level, xp, next_level_xp, satiety, happiness,
	streak_days, last_checkin_date, last_decay_time, updated_at`

type PgRepository struct {
	pool *pgxpool.Pool
}

func NewPgRepository(pool *pgxpool.Pool) *PgRepository {
	return &PgRepository{pool: pool}
}

func (r *PgRepository) ByUserID(ctx context.Context, userID uuid.UUID) (*domain.Pet, error) {
	query := `SELECT ` + petColumns + ` FROM pets WHERE user_id = $1`
	return r.queryOne(ctx, query, userID)
}

func (r *PgRepository) ByUserIDForUpdate(ctx context.Context, userID uuid.UUID) (*domain.Pet, error) {
	querier := postgres.QuerierFrom(ctx, r.pool)
	if _, err := querier.Exec(ctx, `SELECT pg_advisory_xact_lock(hashtextextended($1::text, 0))`, userID); err != nil {
		return nil, fmt.Errorf("lock pet owner: %w", err)
	}

	query := `SELECT ` + petColumns + ` FROM pets WHERE user_id = $1 FOR UPDATE`
	return r.queryOne(ctx, query, userID)
}

func (r *PgRepository) Save(ctx context.Context, pet *domain.Pet) error {
	const query = `INSERT INTO pets (
		id, user_id, name, stage, level, xp, next_level_xp, satiety, happiness,
		streak_days, last_checkin_date, last_decay_time, updated_at
	) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
	ON CONFLICT (user_id) DO UPDATE SET
		name = EXCLUDED.name, stage = EXCLUDED.stage, level = EXCLUDED.level,
		xp = EXCLUDED.xp, next_level_xp = EXCLUDED.next_level_xp,
		satiety = EXCLUDED.satiety, happiness = EXCLUDED.happiness,
		streak_days = EXCLUDED.streak_days, last_checkin_date = EXCLUDED.last_checkin_date,
		last_decay_time = EXCLUDED.last_decay_time,
		updated_at = EXCLUDED.updated_at`

	_, err := postgres.QuerierFrom(ctx, r.pool).Exec(ctx, query,
		pet.ID(), pet.UserID(), pet.Name(), pet.Stage(), pet.Level(), pet.XP(), pet.NextLevelXP(),
		pet.Satiety(), pet.Happiness(), pet.StreakDays(), pet.LastCheckInDate(),
		pet.LastDecayTime(), pet.UpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf("save pet: %w", err)
	}
	return nil
}

func (r *PgRepository) queryOne(ctx context.Context, query string, userID uuid.UUID) (*domain.Pet, error) {
	var p domain.RestoreParams
	var stage string
	err := postgres.QuerierFrom(ctx, r.pool).QueryRow(ctx, query, userID).Scan(
		&p.ID, &p.UserID, &p.Name, &stage, &p.Level, &p.XP, &p.NextLevelXP,
		&p.Satiety, &p.Happiness, &p.StreakDays, &p.LastCheckInDate, &p.LastDecayTime, &p.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domainerr.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("query pet: %w", err)
	}
	p.Stage = domain.Stage(stage)

	return domain.Restore(p), nil
}

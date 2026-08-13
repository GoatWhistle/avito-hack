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

	"github.com/avito-hack/backend/internal/module/weeklylottery/domain"
	"github.com/avito-hack/backend/internal/shared/postgres"
)

const columns = `id, display_id, user_id, week_start, state, board, opened,
	prize_id, reward_code, reward_expires_at, created_at, updated_at`

type PgRepository struct {
	pool *pgxpool.Pool
}

func NewPgRepository(pool *pgxpool.Pool) *PgRepository {
	return &PgRepository{pool: pool}
}

func (r *PgRepository) ByUserWeek(
	ctx context.Context,
	userID uuid.UUID,
	weekStart time.Time,
	lock bool,
) (domain.Run, error) {
	query := `SELECT ` + columns + ` FROM weekly_lottery_runs WHERE user_id = $1 AND week_start = $2`
	if lock {
		query += ` FOR UPDATE`
	}

	return scanRun(postgres.QuerierFrom(ctx, r.pool).QueryRow(ctx, query, userID, weekStart))
}

func (r *PgRepository) ByDisplayID(ctx context.Context, displayID string, lock bool) (domain.Run, error) {
	query := `SELECT ` + columns + ` FROM weekly_lottery_runs WHERE display_id = $1`
	if lock {
		query += ` FOR UPDATE`
	}

	return scanRun(postgres.QuerierFrom(ctx, r.pool).QueryRow(ctx, query, displayID))
}

func (r *PgRepository) Create(ctx context.Context, run domain.Run) (bool, error) {
	const query = `
		INSERT INTO weekly_lottery_runs (` + columns + `)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NULLIF($8, ''), NULLIF($9, ''), $10, $11, $12)
		ON CONFLICT (user_id, week_start) DO NOTHING`

	board, opened, err := encodeState(run)
	if err != nil {
		return false, err
	}

	tag, err := postgres.QuerierFrom(ctx, r.pool).Exec(ctx, query,
		run.ID, run.DisplayID, run.UserID, run.WeekStart, string(run.State), board, opened,
		run.PrizeID, run.RewardCode, run.RewardExpiresAt, run.CreatedAt, run.UpdatedAt)
	if err != nil {
		return false, fmt.Errorf("create weekly lottery run: %w", err)
	}

	return tag.RowsAffected() > 0, nil
}

func (r *PgRepository) Save(ctx context.Context, run domain.Run) error {
	const query = `
		UPDATE weekly_lottery_runs SET
			state = $1, opened = $2, prize_id = NULLIF($3, ''), reward_code = NULLIF($4, ''),
			reward_expires_at = $5, updated_at = $6
		WHERE id = $7`

	_, opened, err := encodeState(run)
	if err != nil {
		return err
	}

	tag, err := postgres.QuerierFrom(ctx, r.pool).Exec(ctx, query,
		string(run.State), opened, run.PrizeID, run.RewardCode, run.RewardExpiresAt, run.UpdatedAt, run.ID)
	if err != nil {
		return fmt.Errorf("save weekly lottery run: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrRunNotFound
	}

	return nil
}

func encodeState(run domain.Run) ([]byte, []byte, error) {
	board, err := json.Marshal(run.Board)
	if err != nil {
		return nil, nil, fmt.Errorf("encode weekly lottery board: %w", err)
	}
	opened, err := json.Marshal(run.Opened)
	if err != nil {
		return nil, nil, fmt.Errorf("encode weekly lottery opened slots: %w", err)
	}

	return board, opened, nil
}

func scanRun(row pgx.Row) (domain.Run, error) {
	var (
		run                 domain.Run
		state               string
		boardRaw, openedRaw []byte
		prizeID, rewardCode *string
	)

	err := row.Scan(
		&run.ID, &run.DisplayID, &run.UserID, &run.WeekStart, &state, &boardRaw, &openedRaw,
		&prizeID, &rewardCode, &run.RewardExpiresAt, &run.CreatedAt, &run.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Run{}, domain.ErrRunNotFound
	}
	if err != nil {
		return domain.Run{}, fmt.Errorf("scan weekly lottery run: %w", err)
	}

	run.State = domain.State(state)
	if prizeID != nil {
		run.PrizeID = *prizeID
	}
	if rewardCode != nil {
		run.RewardCode = *rewardCode
	}
	if err := json.Unmarshal(boardRaw, &run.Board); err != nil {
		return domain.Run{}, fmt.Errorf("decode weekly lottery board: %w", err)
	}
	if err := json.Unmarshal(openedRaw, &run.Opened); err != nil {
		return domain.Run{}, fmt.Errorf("decode weekly lottery opened slots: %w", err)
	}
	if err := domain.ValidateRun(run); err != nil {
		return domain.Run{}, fmt.Errorf("validate stored weekly lottery run: %w", err)
	}

	return run, nil
}

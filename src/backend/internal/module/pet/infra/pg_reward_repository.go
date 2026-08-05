package infra

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/avito-hack/backend/internal/module/pet/domain"
	"github.com/avito-hack/backend/internal/shared/domainerr"
	"github.com/avito-hack/backend/internal/shared/postgres"
)

const rewardColumns = `id, title, description, kind, condition_type, condition_value`

const userRewardColumns = `id, user_id, reward_id, status, code, granted_at, activated_at, expires_at`

type PgRewardRepository struct {
	pool *pgxpool.Pool
}

func NewPgRewardRepository(pool *pgxpool.Pool) *PgRewardRepository {
	return &PgRewardRepository{pool: pool}
}

func (r *PgRewardRepository) Catalog(ctx context.Context) ([]domain.Reward, error) {
	const query = `SELECT ` + rewardColumns + ` FROM rewards ORDER BY condition_type, condition_value, id`

	rows, err := postgres.QuerierFrom(ctx, r.pool).Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query rewards: %w", err)
	}
	defer rows.Close()

	rewards := make([]domain.Reward, 0)
	for rows.Next() {
		reward, scanErr := scanReward(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		rewards = append(rewards, reward)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rewards: %w", err)
	}

	return rewards, nil
}

func (r *PgRewardRepository) ByID(ctx context.Context, rewardID string) (domain.Reward, error) {
	const query = `SELECT ` + rewardColumns + ` FROM rewards WHERE id = $1`

	row := postgres.QuerierFrom(ctx, r.pool).QueryRow(ctx, query, rewardID)

	reward, err := scanReward(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Reward{}, domainerr.ErrNotFound
	}
	if err != nil {
		return domain.Reward{}, err
	}

	return reward, nil
}

func (r *PgRewardRepository) Grant(ctx context.Context, granted domain.UserReward) (bool, error) {
	const query = `
		INSERT INTO user_rewards (id, user_id, reward_id, status, granted_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id, reward_id) DO NOTHING`

	tag, err := postgres.QuerierFrom(ctx, r.pool).Exec(ctx, query,
		granted.ID(), granted.UserID(), granted.RewardID(), string(granted.Status()), granted.GrantedAt())
	if err != nil {
		return false, fmt.Errorf("grant reward: %w", err)
	}

	return tag.RowsAffected() > 0, nil
}

func (r *PgRewardRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]domain.UserReward, error) {
	const query = `SELECT ` + userRewardColumns + `
		FROM user_rewards WHERE user_id = $1 ORDER BY granted_at DESC`

	rows, err := postgres.QuerierFrom(ctx, r.pool).Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("query user rewards: %w", err)
	}
	defer rows.Close()

	granted := make([]domain.UserReward, 0)
	for rows.Next() {
		item, scanErr := scanUserReward(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		granted = append(granted, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate user rewards: %w", err)
	}

	return granted, nil
}

func (r *PgRewardRepository) ByUserAndReward(
	ctx context.Context,
	userID uuid.UUID,
	rewardID string,
) (domain.UserReward, error) {
	const query = `SELECT ` + userRewardColumns + ` FROM user_rewards WHERE user_id = $1 AND reward_id = $2`

	row := postgres.QuerierFrom(ctx, r.pool).QueryRow(ctx, query, userID, rewardID)

	granted, err := scanUserReward(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.UserReward{}, domainerr.ErrNotFound
	}
	if err != nil {
		return domain.UserReward{}, err
	}

	return granted, nil
}

func (r *PgRewardRepository) Activate(ctx context.Context, granted domain.UserReward) (bool, error) {
	const query = `
		UPDATE user_rewards
		SET status = $1, code = $2, activated_at = $3
		WHERE user_id = $4 AND reward_id = $5 AND status = 'granted'`

	tag, err := postgres.QuerierFrom(ctx, r.pool).Exec(ctx, query,
		string(granted.Status()), granted.Code(), granted.ActivatedAt(),
		granted.UserID(), granted.RewardID())
	if err != nil {
		return false, fmt.Errorf("activate reward: %w", err)
	}

	return tag.RowsAffected() > 0, nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanReward(row scanner) (domain.Reward, error) {
	var (
		params        domain.RewardParams
		kind          string
		conditionType string
	)

	if err := row.Scan(&params.ID, &params.Title, &params.Description,
		&kind, &conditionType, &params.ConditionValue); err != nil {
		return domain.Reward{}, err
	}

	params.Kind = domain.RewardKind(kind)
	params.ConditionType = domain.ConditionType(conditionType)

	reward, err := domain.NewReward(params)
	if err != nil {
		return domain.Reward{}, fmt.Errorf("restore reward %q: %w", params.ID, err)
	}

	return reward, nil
}

func scanUserReward(row scanner) (domain.UserReward, error) {
	var (
		params      domain.UserRewardParams
		status      string
		code        *string
		activatedAt *time.Time
		expiresAt   *time.Time
	)

	if err := row.Scan(&params.ID, &params.UserID, &params.RewardID, &status,
		&code, &params.GrantedAt, &activatedAt, &expiresAt); err != nil {
		return domain.UserReward{}, err
	}

	params.Status = domain.RewardStatus(status)
	if code != nil {
		params.Code = *code
	}
	params.ActivatedAt = activatedAt
	params.ExpiresAt = expiresAt

	return domain.RestoreUserReward(params), nil
}

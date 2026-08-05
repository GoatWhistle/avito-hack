package infra

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/avito-hack/backend/internal/module/pet/domain"
	"github.com/avito-hack/backend/internal/shared/postgres"
)

type PgBadgeRepository struct {
	pool *pgxpool.Pool
}

func NewPgBadgeRepository(pool *pgxpool.Pool) *PgBadgeRepository {
	return &PgBadgeRepository{pool: pool}
}

func (r *PgBadgeRepository) Catalog(ctx context.Context) ([]domain.Badge, error) {
	const query = `SELECT id, name, description, icon_url FROM badges ORDER BY id`

	rows, err := postgres.QuerierFrom(ctx, r.pool).Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query badges: %w", err)
	}
	defer rows.Close()

	badges := make([]domain.Badge, 0)
	for rows.Next() {
		var id, name, description, iconURL string
		if err := rows.Scan(&id, &name, &description, &iconURL); err != nil {
			return nil, fmt.Errorf("scan badge: %w", err)
		}
		badges = append(badges, domain.NewBadge(id, name, description, iconURL))
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate badges: %w", err)
	}

	return badges, nil
}

func (r *PgBadgeRepository) EarnedBy(ctx context.Context, userID uuid.UUID) ([]domain.EarnedBadge, error) {
	const query = `
		SELECT b.id, b.name, b.description, b.icon_url, ub.earned_at
		FROM user_badges ub
		JOIN badges b ON b.id = ub.badge_id
		WHERE ub.user_id = $1
		ORDER BY ub.earned_at DESC`

	rows, err := postgres.QuerierFrom(ctx, r.pool).Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("query earned badges: %w", err)
	}
	defer rows.Close()

	earned := make([]domain.EarnedBadge, 0)
	for rows.Next() {
		var (
			id, name, description, iconURL string
			earnedAt                       time.Time
		)
		if err := rows.Scan(&id, &name, &description, &iconURL, &earnedAt); err != nil {
			return nil, fmt.Errorf("scan earned badge: %w", err)
		}
		earned = append(earned,
			domain.NewEarnedBadge(domain.NewBadge(id, name, description, iconURL), userID, earnedAt))
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate earned badges: %w", err)
	}

	return earned, nil
}

func (r *PgBadgeRepository) Award(
	ctx context.Context,
	userID uuid.UUID,
	badgeID string,
	earnedAt time.Time,
) (bool, error) {
	const query = `
		INSERT INTO user_badges (id, user_id, badge_id, earned_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id, badge_id) DO NOTHING`

	tag, err := postgres.QuerierFrom(ctx, r.pool).Exec(ctx, query, uuid.New(), userID, badgeID, earnedAt)
	if err != nil {
		return false, fmt.Errorf("award badge: %w", err)
	}

	return tag.RowsAffected() > 0, nil
}

package infra

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/avito-hack/backend/internal/module/pet/app"
	"github.com/avito-hack/backend/internal/module/pet/domain"
	"github.com/avito-hack/backend/internal/shared/postgres"
)

const staleListingDays = 14

const xpByActionQuery = `
	SELECT action, count(*), COALESCE(sum(amount), 0)
	FROM xp_events
	WHERE user_id = $1 AND created_at >= $2 AND created_at < $3
	GROUP BY action
	ORDER BY 3 DESC, 1`

const xpTotalBeforeQuery = `
	SELECT COALESCE(sum(amount), 0) FROM xp_events
	WHERE user_id = $1 AND created_at < $2`

const rewardsGrantedQuery = `
	SELECT r.title FROM user_rewards ur
	JOIN rewards r ON r.id = ur.reward_id
	WHERE ur.user_id = $1 AND ur.granted_at >= $2 AND ur.granted_at < $3
	ORDER BY ur.granted_at`

const badgesEarnedQuery = `
	SELECT b.name FROM user_badges ub
	JOIN badges b ON b.id = ub.badge_id
	WHERE ub.user_id = $1 AND ub.earned_at >= $2 AND ub.earned_at < $3
	ORDER BY ub.earned_at`

const listingIssuesQuery = `
	SELECT i.display_id, i.title,
		CASE
			WHEN NOT EXISTS (SELECT 1 FROM item_photos p WHERE p.item_id = i.id) THEN 'no_photo'
			WHEN i.price_kopeks = 0 THEN 'no_price'
			WHEN char_length(i.description) < 200 THEN 'short_description'
			ELSE 'stale'
		END,
		GREATEST(0, EXTRACT(DAY FROM (now() - i.updated_at))::int)
	FROM items i
	WHERE i.owner_id = $1 AND i.deleted_at IS NULL AND i.status = 'published'
	  AND (
		NOT EXISTS (SELECT 1 FROM item_photos p WHERE p.item_id = i.id)
		OR i.price_kopeks = 0
		OR char_length(i.description) < 200
		OR i.updated_at < now() - make_interval(days => $3)
	  )
	ORDER BY i.updated_at ASC
	LIMIT $2`

type PgDayActivity struct {
	pool *pgxpool.Pool
}

func NewPgDayActivity(pool *pgxpool.Pool) *PgDayActivity {
	return &PgDayActivity{pool: pool}
}

func (m *PgDayActivity) XPByAction(
	ctx context.Context,
	userID uuid.UUID,
	from, to time.Time,
) ([]app.XPAggregate, error) {
	rows, err := postgres.QuerierFrom(ctx, m.pool).Query(ctx, xpByActionQuery, userID, from, to)
	if err != nil {
		return nil, fmt.Errorf("query xp aggregates: %w", err)
	}
	defer rows.Close()

	aggregates := make([]app.XPAggregate, 0)
	for rows.Next() {
		var (
			action    string
			aggregate app.XPAggregate
		)
		if err := rows.Scan(&action, &aggregate.Count, &aggregate.Amount); err != nil {
			return nil, fmt.Errorf("scan xp aggregate: %w", err)
		}
		aggregate.Action = domain.Action(action)
		aggregates = append(aggregates, aggregate)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate xp aggregates: %w", err)
	}

	return aggregates, nil
}

func (m *PgDayActivity) XPTotalBefore(
	ctx context.Context,
	userID uuid.UUID,
	before time.Time,
) (int, error) {
	var total int

	err := postgres.QuerierFrom(ctx, m.pool).
		QueryRow(ctx, xpTotalBeforeQuery, userID, before).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("sum xp before day: %w", err)
	}

	return total, nil
}

func (m *PgDayActivity) RewardsGranted(
	ctx context.Context,
	userID uuid.UUID,
	from, to time.Time,
) ([]string, error) {
	return m.titles(ctx, rewardsGrantedQuery, "granted rewards", userID, from, to)
}

func (m *PgDayActivity) BadgesEarned(
	ctx context.Context,
	userID uuid.UUID,
	from, to time.Time,
) ([]string, error) {
	return m.titles(ctx, badgesEarnedQuery, "earned badges", userID, from, to)
}

func (m *PgDayActivity) titles(
	ctx context.Context,
	query, label string,
	userID uuid.UUID,
	from, to time.Time,
) ([]string, error) {
	rows, err := postgres.QuerierFrom(ctx, m.pool).Query(ctx, query, userID, from, to)
	if err != nil {
		return nil, fmt.Errorf("query %s: %w", label, err)
	}
	defer rows.Close()

	titles := make([]string, 0)
	for rows.Next() {
		var title string
		if err := rows.Scan(&title); err != nil {
			return nil, fmt.Errorf("scan %s: %w", label, err)
		}
		titles = append(titles, title)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate %s: %w", label, err)
	}

	return titles, nil
}

func (m *PgDayActivity) ListingIssues(
	ctx context.Context,
	userID uuid.UUID,
	limit int,
) ([]domain.ListingIssue, error) {
	rows, err := postgres.QuerierFrom(ctx, m.pool).
		Query(ctx, listingIssuesQuery, userID, limit, staleListingDays)
	if err != nil {
		return nil, fmt.Errorf("query listing issues: %w", err)
	}
	defer rows.Close()

	issues := make([]domain.ListingIssue, 0, limit)
	for rows.Next() {
		var (
			issue domain.ListingIssue
			kind  string
		)
		if err := rows.Scan(&issue.ItemID, &issue.Title, &kind, &issue.StaleDays); err != nil {
			return nil, fmt.Errorf("scan listing issue: %w", err)
		}
		issue.Kind = domain.ListingIssueKind(kind)
		issues = append(issues, issue)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate listing issues: %w", err)
	}

	return issues, nil
}

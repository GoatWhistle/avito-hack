package infra

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/avito-hack/backend/internal/module/pet/app"
	"github.com/avito-hack/backend/internal/shared/postgres"
)

const leaderboardSelect = `
	SELECT p.user_id, COALESCE(u.full_name, ''), p.level, p.xp, p.streak_days
	FROM pets p
	JOIN users u ON u.id = p.user_id AND u.deleted_at IS NULL`

const leaderboardOrder = `
	ORDER BY p.level DESC, p.xp DESC, p.user_id ASC
	LIMIT $%d OFFSET $%d`

const leaderboardKeyset = `
	WHERE (p.level < $1)
	   OR (p.level = $1 AND p.xp < $2)
	   OR (p.level = $1 AND p.xp = $2 AND p.user_id > $3)`

const rankBeforeCursor = `
	SELECT count(*)
	FROM pets p
	JOIN users u ON u.id = p.user_id AND u.deleted_at IS NULL
	WHERE (p.level > $1)
	   OR (p.level = $1 AND p.xp > $2)
	   OR (p.level = $1 AND p.xp = $2 AND p.user_id <= $3)`

const rankOfUser = `
	SELECT count(*) + 1
	FROM pets p
	JOIN users u ON u.id = p.user_id AND u.deleted_at IS NULL
	CROSS JOIN me
	WHERE (p.level > me.level)
	   OR (p.level = me.level AND p.xp > me.xp)
	   OR (p.level = me.level AND p.xp = me.xp AND p.user_id < me.user_id)`

const rankOfUserQuery = `WITH me AS (
		SELECT p.user_id, p.level, p.xp
		FROM pets p
		JOIN users u ON u.id = p.user_id AND u.deleted_at IS NULL
		WHERE p.user_id = $1
	)` + rankOfUser

type PgLeaderboard struct {
	pool *pgxpool.Pool
}

func NewPgLeaderboard(pool *pgxpool.Pool) *PgLeaderboard {
	return &PgLeaderboard{pool: pool}
}

func (m *PgLeaderboard) Page(ctx context.Context, f app.LeaderboardFilter) ([]app.LeaderboardEntry, error) {
	query, args, baseRank, err := m.buildPageQuery(ctx, f)
	if err != nil {
		return nil, err
	}

	rows, err := postgres.QuerierFrom(ctx, m.pool).Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query leaderboard: %w", err)
	}
	defer rows.Close()

	entries := make([]app.LeaderboardEntry, 0, f.Limit)

	for rows.Next() {
		var e app.LeaderboardEntry
		if err := rows.Scan(&e.UserID, &e.Name, &e.Level, &e.XP, &e.StreakDays); err != nil {
			return nil, fmt.Errorf("scan leaderboard entry: %w", err)
		}

		e.Rank = baseRank + len(entries)
		entries = append(entries, e)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate leaderboard: %w", err)
	}

	return entries, nil
}

func (m *PgLeaderboard) buildPageQuery(
	ctx context.Context,
	f app.LeaderboardFilter,
) (query string, args []any, baseRank int, err error) {
	if f.Cursor.IsZero() {
		return leaderboardSelect + fmt.Sprintf(leaderboardOrder, 1, 2),
			[]any{f.Limit, f.Offset}, f.Offset + 1, nil
	}

	baseRank, err = m.rankAfterCursor(ctx, f.Cursor.Level, f.Cursor.XP, f.Cursor.UserID)
	if err != nil {
		return "", nil, 0, err
	}

	args = []any{f.Cursor.Level, f.Cursor.XP, f.Cursor.UserID, f.Limit, f.Offset}

	return leaderboardSelect + leaderboardKeyset + fmt.Sprintf(leaderboardOrder, 4, 5), args, baseRank, nil
}

func (m *PgLeaderboard) rankAfterCursor(ctx context.Context, level, xp int, userID uuid.UUID) (int, error) {
	var passed int

	row := postgres.QuerierFrom(ctx, m.pool).QueryRow(ctx, rankBeforeCursor, level, xp, userID)
	if err := row.Scan(&passed); err != nil {
		return 0, fmt.Errorf("count leaderboard offset: %w", err)
	}

	return passed + 1, nil
}

func (m *PgLeaderboard) RankOf(ctx context.Context, userID uuid.UUID) (rank int, found bool, err error) {
	if userID == uuid.Nil {
		return 0, false, nil
	}

	row := postgres.QuerierFrom(ctx, m.pool).QueryRow(ctx, rankOfUserQuery, userID)

	if scanErr := row.Scan(&rank); scanErr != nil {
		if errors.Is(scanErr, pgx.ErrNoRows) {
			return 0, false, nil
		}

		return 0, false, fmt.Errorf("query user rank: %w", scanErr)
	}

	return rank, true, nil
}

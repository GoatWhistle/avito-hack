package infra

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/avito-hack/backend/internal/shared/postgres"
)

type PgScoreRepository struct {
	pool *pgxpool.Pool
}

func NewPgScoreRepository(pool *pgxpool.Pool) *PgScoreRepository {
	return &PgScoreRepository{pool: pool}
}

func (r *PgScoreRepository) BestScore(ctx context.Context, userID uuid.UUID, gameSlug string) (int, error) {
	const query = `SELECT best_score FROM game_best_scores WHERE user_id = $1 AND game_slug = $2`

	var best int

	err := postgres.QuerierFrom(ctx, r.pool).
		QueryRow(ctx, query, userID, gameSlug).
		Scan(&best)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("load best score: %w", err)
	}

	return best, nil
}

func (r *PgScoreRepository) SaveBestScore(
	ctx context.Context,
	userID uuid.UUID,
	gameSlug string,
	score int,
) (int, error) {
	const query = `
		INSERT INTO game_best_scores (user_id, game_slug, best_score)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, game_slug) DO UPDATE SET
			best_score = GREATEST(game_best_scores.best_score, EXCLUDED.best_score),
			updated_at = now()
		RETURNING best_score`

	var best int

	err := postgres.QuerierFrom(ctx, r.pool).
		QueryRow(ctx, query, userID, gameSlug, score).
		Scan(&best)
	if err != nil {
		return 0, fmt.Errorf("save best score: %w", err)
	}

	return best, nil
}

-- +goose Up
CREATE TABLE game_best_scores (
    user_id    uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    game_slug  text NOT NULL,
    best_score int NOT NULL DEFAULT 0 CHECK (best_score >= 0),
    updated_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, game_slug)
);

CREATE INDEX idx_game_best_scores_game_score ON game_best_scores (game_slug, best_score DESC);

-- +goose Down
DROP INDEX IF EXISTS idx_game_best_scores_game_score;
DROP TABLE IF EXISTS game_best_scores;

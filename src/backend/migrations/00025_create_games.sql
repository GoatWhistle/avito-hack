-- +goose Up
CREATE TABLE game_rounds (
    id          uuid PRIMARY KEY,
    display_id  text NOT NULL,
    user_id     uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    game_slug   text NOT NULL,
    state       text NOT NULL DEFAULT 'active' CHECK (state IN ('active', 'won', 'lost')),
    streak      int NOT NULL DEFAULT 0 CHECK (streak >= 0),
    best_streak int NOT NULL DEFAULT 0 CHECK (best_streak >= 0),
    payload     jsonb NOT NULL DEFAULT '{}'::jsonb,
    created_at  timestamptz NOT NULL DEFAULT now(),
    updated_at  timestamptz NOT NULL DEFAULT now()
);

ALTER TABLE game_rounds ADD CONSTRAINT game_rounds_display_id_unique UNIQUE (display_id);

CREATE INDEX idx_game_rounds_user_game_state ON game_rounds (user_id, game_slug, state);

CREATE TABLE game_daily_progress (
    user_id     uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    game_slug   text NOT NULL,
    day         date NOT NULL,
    attempts    int NOT NULL DEFAULT 0 CHECK (attempts >= 0),
    best_streak int NOT NULL DEFAULT 0 CHECK (best_streak >= 0),
    PRIMARY KEY (user_id, game_slug, day)
);

CREATE TABLE game_streaks (
    user_id           uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    game_slug         text NOT NULL,
    current_days      int NOT NULL DEFAULT 0 CHECK (current_days >= 0),
    best_days         int NOT NULL DEFAULT 0 CHECK (best_days >= 0),
    last_day          date,
    reward_claimed_at timestamptz,
    PRIMARY KEY (user_id, game_slug)
);

-- +goose Down
DROP TABLE IF EXISTS game_streaks;
DROP TABLE IF EXISTS game_daily_progress;
DROP INDEX IF EXISTS idx_game_rounds_user_game_state;
DROP TABLE IF EXISTS game_rounds;

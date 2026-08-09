-- +goose Up
CREATE TABLE pets (
    id                uuid PRIMARY KEY,
    user_id           uuid NOT NULL UNIQUE REFERENCES users (id) ON DELETE CASCADE,
    name              text NOT NULL DEFAULT '' CHECK (char_length(name) <= 100),
    stage             text NOT NULL DEFAULT 'egg'
                      CHECK (stage IN ('egg', 'baby', 'teen', 'adult', 'legend')),
    level             integer NOT NULL DEFAULT 1 CHECK (level >= 1),
    xp                integer NOT NULL DEFAULT 0 CHECK (xp >= 0),
    next_level_xp     integer NOT NULL DEFAULT 5 CHECK (next_level_xp >= 0),
    satiety           integer NOT NULL DEFAULT 70 CHECK (satiety BETWEEN 0 AND 100),
    happiness         integer NOT NULL DEFAULT 70 CHECK (happiness BETWEEN 0 AND 100),
    streak_days       integer NOT NULL DEFAULT 0 CHECK (streak_days >= 0),
    last_checkin_date date,
    last_decay_time   timestamptz NOT NULL DEFAULT now(),
    updated_at        timestamptz NOT NULL DEFAULT now(),
    energy            integer NOT NULL DEFAULT 100 CHECK (energy BETWEEN 0 AND 100),
    freezes           integer NOT NULL DEFAULT 0 CHECK (freezes BETWEEN 0 AND 3),
    hatched_at        timestamptz
);

CREATE INDEX idx_pets_leaderboard ON pets (level DESC, xp DESC, user_id);

-- +goose Down
DROP INDEX IF EXISTS idx_pets_leaderboard;
DROP TABLE IF EXISTS pets;

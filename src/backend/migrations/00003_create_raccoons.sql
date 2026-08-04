-- +goose Up
CREATE TABLE raccoon_profiles (
    id             uuid PRIMARY KEY,
    user_id        uuid NOT NULL UNIQUE REFERENCES users (id) ON DELETE CASCADE,
    name           text NOT NULL DEFAULT 'Rocky',
    level          integer NOT NULL DEFAULT 1 CHECK (level >= 1),
    xp             integer NOT NULL DEFAULT 0 CHECK (xp >= 0),
    current_streak integer NOT NULL DEFAULT 0 CHECK (current_streak >= 0),
    created_at     timestamptz NOT NULL DEFAULT now(),
    updated_at     timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_raccoon_profiles_user_id ON raccoon_profiles (user_id);

CREATE TABLE user_badges (
    id        uuid PRIMARY KEY,
    user_id   uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    badge_id  text NOT NULL,
    earned_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (user_id, badge_id)
);

CREATE INDEX idx_user_badges_user ON user_badges (user_id);

CREATE TABLE claimed_rewards (
    id         uuid PRIMARY KEY,
    user_id    uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    reward_id  text NOT NULL,
    promocode  text NOT NULL,
    claimed_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_claimed_rewards_user ON claimed_rewards (user_id);

-- +goose Down
DROP INDEX IF EXISTS idx_claimed_rewards_user;
DROP TABLE IF EXISTS claimed_rewards;
DROP INDEX IF EXISTS idx_user_badges_user;
DROP TABLE IF EXISTS user_badges;
DROP INDEX IF EXISTS idx_raccoon_profiles_user_id;
DROP TABLE IF EXISTS raccoon_profiles;

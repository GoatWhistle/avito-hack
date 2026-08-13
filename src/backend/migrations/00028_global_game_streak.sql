-- +goose Up
CREATE TABLE game_streaks_global (
    user_id           uuid PRIMARY KEY REFERENCES users (id) ON DELETE CASCADE,
    current_days      int NOT NULL DEFAULT 0 CHECK (current_days >= 0),
    best_days         int NOT NULL DEFAULT 0 CHECK (best_days >= 0),
    last_day          date,
    reward_claimed_at timestamptz
);

INSERT INTO game_streaks_global (user_id, current_days, best_days, last_day, reward_claimed_at)
SELECT user_id,
       MAX(current_days),
       MAX(best_days),
       MAX(last_day),
       MAX(reward_claimed_at)
FROM game_streaks
GROUP BY user_id;

DROP TABLE game_streaks;

ALTER TABLE game_streaks_global RENAME TO game_streaks;

CREATE TABLE game_daily_days (
    user_id  uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    day      date NOT NULL,
    attempts int NOT NULL DEFAULT 0 CHECK (attempts >= 0),
    PRIMARY KEY (user_id, day)
);

INSERT INTO game_daily_days (user_id, day, attempts)
SELECT user_id, day, SUM(attempts)
FROM game_daily_progress
GROUP BY user_id, day;

-- +goose Down
DROP TABLE IF EXISTS game_daily_days;

CREATE TABLE game_streaks_per_game (
    user_id           uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    game_slug         text NOT NULL,
    current_days      int NOT NULL DEFAULT 0 CHECK (current_days >= 0),
    best_days         int NOT NULL DEFAULT 0 CHECK (best_days >= 0),
    last_day          date,
    reward_claimed_at timestamptz,
    PRIMARY KEY (user_id, game_slug)
);

INSERT INTO game_streaks_per_game (user_id, game_slug, current_days, best_days, last_day, reward_claimed_at)
SELECT s.user_id, p.game_slug, s.current_days, s.best_days, s.last_day, s.reward_claimed_at
FROM game_streaks s
JOIN (SELECT DISTINCT user_id, game_slug FROM game_daily_progress) p ON p.user_id = s.user_id;

DROP TABLE game_streaks;

ALTER TABLE game_streaks_per_game RENAME TO game_streaks;

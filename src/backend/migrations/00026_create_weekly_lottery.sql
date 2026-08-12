-- +goose Up
CREATE TABLE weekly_lottery_runs (
    id             uuid PRIMARY KEY,
    display_id     text NOT NULL UNIQUE,
    user_id        uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    week_start     date NOT NULL,
    state          text NOT NULL CHECK (state IN ('active', 'won', 'lost')),
    board          jsonb NOT NULL,
    opened         jsonb NOT NULL DEFAULT '[]'::jsonb,
    prize_id       text,
    reward_code    text,
    reward_expires_at timestamptz,
    created_at     timestamptz NOT NULL DEFAULT now(),
    updated_at     timestamptz NOT NULL DEFAULT now(),
    UNIQUE (user_id, week_start),
    CHECK (jsonb_typeof(board) = 'array' AND jsonb_array_length(board) = 9),
    CHECK (jsonb_typeof(opened) = 'array'),
    CHECK ((state = 'won') = (
        prize_id IS NOT NULL AND reward_code IS NOT NULL AND reward_expires_at IS NOT NULL
    )),
    CHECK (state <> 'lost' OR (
        prize_id IS NULL AND reward_code IS NULL AND reward_expires_at IS NULL
    ))
);

CREATE UNIQUE INDEX idx_weekly_lottery_reward_code
    ON weekly_lottery_runs (reward_code) WHERE reward_code IS NOT NULL;

-- +goose Down
DROP INDEX IF EXISTS idx_weekly_lottery_reward_code;
DROP TABLE IF EXISTS weekly_lottery_runs;

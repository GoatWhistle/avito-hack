-- +goose Up
CREATE TABLE daily_summaries (
    id           uuid PRIMARY KEY,
    user_id      uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    date         date NOT NULL,
    facts        jsonb NOT NULL,
    message      text NOT NULL CHECK (char_length(message) BETWEEN 1 AND 4000),
    advice       jsonb,
    generated_by text NOT NULL CHECK (generated_by IN ('template', 'llm')),
    created_at   timestamptz NOT NULL DEFAULT now(),
    UNIQUE (user_id, date)
);

CREATE INDEX idx_daily_summaries_user_date ON daily_summaries (user_id, date DESC, id DESC);

-- +goose Down
DROP INDEX IF EXISTS idx_daily_summaries_user_date;
DROP TABLE IF EXISTS daily_summaries;

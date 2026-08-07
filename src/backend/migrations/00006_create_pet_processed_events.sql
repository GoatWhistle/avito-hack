-- +goose Up
CREATE TABLE pet_processed_events (
    event_id     uuid PRIMARY KEY,
    event_type   text NOT NULL,
    user_id      uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    subject_id   uuid,
    occurred_at  timestamptz NOT NULL,
    processed_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_pet_events_user_time
    ON pet_processed_events (user_id, occurred_at DESC);

CREATE INDEX idx_pet_events_processed_at
    ON pet_processed_events (processed_at);

-- +goose Down
DROP INDEX IF EXISTS idx_pet_events_processed_at;
DROP INDEX IF EXISTS idx_pet_events_user_time;
DROP TABLE IF EXISTS pet_processed_events;

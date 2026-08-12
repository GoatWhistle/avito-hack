-- +goose Up
CREATE TABLE item_moderation_log (
    id           uuid PRIMARY KEY,
    item_id      uuid NOT NULL REFERENCES items (id) ON DELETE CASCADE,
    verdict      text NOT NULL CHECK (verdict IN ('approved', 'rejected', 'unavailable')),
    reason       text NOT NULL DEFAULT '',
    provider     text NOT NULL,
    checked_at   timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_item_moderation_log_item ON item_moderation_log (item_id, checked_at DESC);

-- +goose Down
DROP INDEX IF EXISTS idx_item_moderation_log_item;
DROP TABLE IF EXISTS item_moderation_log;

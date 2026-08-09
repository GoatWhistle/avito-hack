-- +goose Up
CREATE INDEX IF NOT EXISTS idx_items_feed_created ON items (created_at DESC, id DESC)
    WHERE deleted_at IS NULL;

DROP INDEX IF EXISTS idx_items_published_created;

-- +goose Down
CREATE INDEX IF NOT EXISTS idx_items_published_created ON items (created_at DESC, id DESC)
    WHERE status = 'published' AND deleted_at IS NULL;

DROP INDEX IF EXISTS idx_items_feed_created;

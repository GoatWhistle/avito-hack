-- +goose Up
CREATE TABLE items (
    id           uuid PRIMARY KEY,
    owner_id     uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    title        text NOT NULL CHECK (char_length(title) BETWEEN 3 AND 200),
    description  text NOT NULL DEFAULT '' CHECK (char_length(description) <= 5000),
    price_kopeks bigint NOT NULL DEFAULT 0 CHECK (price_kopeks >= 0),
    status       text NOT NULL DEFAULT 'draft'
                 CHECK (status IN ('draft', 'moderation', 'published', 'archived')),
    attributes   jsonb NOT NULL DEFAULT '{}'::jsonb,
    created_at   timestamptz NOT NULL DEFAULT now(),
    updated_at   timestamptz NOT NULL DEFAULT now(),
    deleted_at   timestamptz
);

CREATE INDEX idx_items_published_created ON items (created_at DESC, id DESC)
    WHERE status = 'published' AND deleted_at IS NULL;

CREATE INDEX idx_items_owner_created ON items (owner_id, created_at DESC, id DESC)
    WHERE deleted_at IS NULL;

CREATE INDEX idx_items_status_created ON items (status, created_at DESC, id DESC)
    WHERE deleted_at IS NULL;

CREATE INDEX idx_items_title_trgm ON items USING gin (title gin_trgm_ops);
CREATE INDEX idx_items_attributes ON items USING gin (attributes);

-- +goose Down
DROP INDEX IF EXISTS idx_items_attributes;
DROP INDEX IF EXISTS idx_items_title_trgm;
DROP INDEX IF EXISTS idx_items_status_created;
DROP INDEX IF EXISTS idx_items_owner_created;
DROP INDEX IF EXISTS idx_items_published_created;
DROP TABLE IF EXISTS items;

-- +goose Up
ALTER TABLE items ADD COLUMN display_id text;

UPDATE items SET display_id = substr(encode(gen_random_bytes(8), 'hex'), 1, 12);

ALTER TABLE items ALTER COLUMN display_id SET NOT NULL;
ALTER TABLE items ADD CONSTRAINT items_display_id_unique UNIQUE (display_id);

-- +goose Down
ALTER TABLE items DROP CONSTRAINT IF EXISTS items_display_id_unique;
ALTER TABLE items DROP COLUMN IF EXISTS display_id;

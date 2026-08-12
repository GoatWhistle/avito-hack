-- +goose Up
ALTER TABLE items ADD COLUMN is_seed boolean NOT NULL DEFAULT false;
ALTER TABLE items ADD COLUMN ai_verified boolean NOT NULL DEFAULT false;

UPDATE items SET is_seed = true;

-- +goose Down
ALTER TABLE items DROP COLUMN ai_verified;
ALTER TABLE items DROP COLUMN is_seed;

-- +goose Up
ALTER TABLE users ADD COLUMN display_id text;

UPDATE users SET display_id = substr(encode(gen_random_bytes(8), 'hex'), 1, 12);

ALTER TABLE users ALTER COLUMN display_id SET NOT NULL;
ALTER TABLE users ADD CONSTRAINT users_display_id_unique UNIQUE (display_id);

-- +goose Down
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_display_id_unique;
ALTER TABLE users DROP COLUMN IF EXISTS display_id;

-- +goose Up
ALTER TABLE pets
    ADD COLUMN interaction_version bigint NOT NULL DEFAULT 0
    CHECK (interaction_version >= 0);

-- +goose Down
ALTER TABLE pets DROP COLUMN IF EXISTS interaction_version;

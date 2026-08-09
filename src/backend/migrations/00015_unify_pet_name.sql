-- +goose Up
UPDATE pets
SET name = 'Ноти'
WHERE name IS NULL OR btrim(name) = '' OR name <> 'Ноти';

ALTER TABLE pets ALTER COLUMN name SET DEFAULT 'Ноти';

-- +goose Down
ALTER TABLE pets ALTER COLUMN name DROP DEFAULT;

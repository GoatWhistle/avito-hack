-- +goose Up
UPDATE pets
SET stage = 'baby',
    hatched_at = COALESCE(hatched_at, updated_at, now())
WHERE stage = 'egg';

UPDATE pets
SET hatched_at = COALESCE(updated_at, now())
WHERE hatched_at IS NULL;

ALTER TABLE pets DROP CONSTRAINT IF EXISTS pets_stage_check;

ALTER TABLE pets
    ADD CONSTRAINT pets_stage_check
    CHECK (stage IN ('baby', 'teen', 'adult', 'legend'));

ALTER TABLE pets ALTER COLUMN stage SET DEFAULT 'baby';

-- +goose Down
ALTER TABLE pets ALTER COLUMN stage SET DEFAULT 'egg';

ALTER TABLE pets DROP CONSTRAINT IF EXISTS pets_stage_check;

ALTER TABLE pets
    ADD CONSTRAINT pets_stage_check
    CHECK (stage IN ('egg', 'baby', 'teen', 'adult', 'legend'));

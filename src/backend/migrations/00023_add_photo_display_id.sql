-- +goose Up
ALTER TABLE item_photos ADD COLUMN display_id text;

UPDATE item_photos SET display_id = substr(encode(gen_random_bytes(8), 'hex'), 1, 12);

ALTER TABLE item_photos ALTER COLUMN display_id SET NOT NULL;
ALTER TABLE item_photos ADD CONSTRAINT item_photos_display_id_unique UNIQUE (display_id);

UPDATE item_photos p
SET url = '/uploads/' || i.display_id || '/' || substr(p.url, length('/uploads/' || i.id::text || '/') + 1)
FROM items i
WHERE p.item_id = i.id AND p.url LIKE '/uploads/' || i.id::text || '/%';

-- +goose Down
UPDATE item_photos p
SET url = '/uploads/' || i.id::text || '/' || substr(p.url, length('/uploads/' || i.display_id || '/') + 1)
FROM items i
WHERE p.item_id = i.id AND p.url LIKE '/uploads/' || i.display_id || '/%';

ALTER TABLE item_photos DROP CONSTRAINT IF EXISTS item_photos_display_id_unique;
ALTER TABLE item_photos DROP COLUMN IF EXISTS display_id;

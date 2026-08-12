-- +goose Up
INSERT INTO item_photos (id, item_id, url, position, created_at) VALUES
    (gen_random_uuid(), 'a0000000-0000-4000-a000-00000000000d',
     '/uploads/a0000000-0000-4000-a000-00000000000d/demo-1.jpg', 0, now()),
    (gen_random_uuid(), 'a0000000-0000-4000-a000-00000000000e',
     '/uploads/a0000000-0000-4000-a000-00000000000e/demo-1.jpg', 0, now())
ON CONFLICT (item_id, position) DO NOTHING;

-- +goose Down
DELETE FROM item_photos
WHERE item_id IN (
    'a0000000-0000-4000-a000-00000000000d',
    'a0000000-0000-4000-a000-00000000000e'
);

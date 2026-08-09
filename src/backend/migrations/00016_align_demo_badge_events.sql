-- +goose Up

INSERT INTO xp_events (id, user_id, action, subject_id, amount, created_at)
SELECT
    md5('demo-badge-evt:' || b.user_id::text || ':' || b.action || ':' || i.item_id::text)::uuid,
    b.user_id,
    b.action,
    i.item_id,
    b.amount,
    b.earned_at - ((row_number() OVER (PARTITION BY b.user_id, b.action ORDER BY i.item_id)) || ' hours')::interval
FROM (
    SELECT ub.user_id, ub.earned_at,
           CASE ub.badge_id
               WHEN 'nothing_hidden' THEN 'item_improved'
               WHEN 'explorer'       THEN 'favorite'
               WHEN 'green_planet'   THEN 'item_sold'
               WHEN 'curious_eye'    THEN 'item_viewed'
           END AS action,
           CASE ub.badge_id
               WHEN 'nothing_hidden' THEN 3
               WHEN 'explorer'       THEN 5
               WHEN 'green_planet'   THEN 5
               WHEN 'curious_eye'    THEN 10
           END AS target,
           CASE ub.badge_id
               WHEN 'nothing_hidden' THEN 15
               WHEN 'explorer'       THEN 1
               WHEN 'green_planet'   THEN 100
               WHEN 'curious_eye'    THEN 1
           END AS amount
    FROM user_badges ub
    WHERE ub.badge_id IN ('nothing_hidden', 'explorer', 'green_planet', 'curious_eye')
      AND ub.user_id::text LIKE 'd0000000-0000-4000-a000-%'
) AS b
CROSS JOIN LATERAL (
    SELECT it.id AS item_id
    FROM items it
    WHERE it.deleted_at IS NULL
    ORDER BY it.id
    LIMIT b.target
) AS i
WHERE b.action IS NOT NULL
ON CONFLICT (user_id, action, subject_id) WHERE subject_id IS NOT NULL DO NOTHING;

DELETE FROM user_badges
WHERE user_id::text LIKE 'd0000000-0000-4000-a000-%'
  AND badge_id NOT IN (
      'explorer', 'curious_eye', 'first_listing', 'nothing_hidden',
      'green_planet', 'raccoon_friend', 'charged_streak', 'quest_runner'
  );

INSERT INTO user_badges (id, user_id, badge_id, earned_at)
SELECT
    md5('demo-first-listing:' || u.id::text)::uuid,
    u.id,
    'first_listing',
    min(it.created_at)
FROM users u
JOIN items it ON it.owner_id = u.id AND it.deleted_at IS NULL
WHERE u.id::text LIKE 'd0000000-0000-4000-a000-%'
GROUP BY u.id
ON CONFLICT (user_id, badge_id) DO NOTHING;

INSERT INTO xp_events (id, user_id, action, subject_id, amount, created_at)
SELECT
    md5('demo-published:' || it.owner_id::text || ':' || it.id::text)::uuid,
    it.owner_id,
    'item_published',
    it.id,
    50,
    it.created_at
FROM items it
WHERE it.deleted_at IS NULL
  AND it.owner_id::text LIKE 'd0000000-0000-4000-a000-%'
ON CONFLICT (user_id, action, subject_id) WHERE subject_id IS NOT NULL DO NOTHING;

-- +goose Down

DELETE FROM xp_events
WHERE user_id::text LIKE 'd0000000-0000-4000-a000-%'
  AND action IN ('item_improved', 'item_viewed', 'item_sold', 'item_published');

DELETE FROM user_badges
WHERE user_id::text LIKE 'd0000000-0000-4000-a000-%'
  AND badge_id = 'first_listing';

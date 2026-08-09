-- +goose Up

INSERT INTO item_photos (id, item_id, url, position, created_at)
SELECT
    md5('demo-photo:' || i.item_id::text || ':' || p.position)::uuid,
    i.item_id,
    '/uploads/' || i.item_id || '/demo-' || (p.position + 1) || '.jpg',
    p.position,
    it.created_at + (p.position || ' minutes')::interval
FROM (VALUES
    ('a0000000-0000-4000-a000-000000000001'::uuid, 3),
    ('a0000000-0000-4000-a000-000000000002'::uuid, 3),
    ('a0000000-0000-4000-a000-000000000003'::uuid, 2),
    ('a0000000-0000-4000-a000-000000000004'::uuid, 2),
    ('a0000000-0000-4000-a000-000000000005'::uuid, 3),
    ('a0000000-0000-4000-a000-000000000006'::uuid, 2),
    ('a0000000-0000-4000-a000-000000000007'::uuid, 2),
    ('a0000000-0000-4000-a000-000000000008'::uuid, 3),
    ('a0000000-0000-4000-a000-000000000009'::uuid, 2),
    ('a0000000-0000-4000-a000-00000000000a'::uuid, 2),
    ('a0000000-0000-4000-a000-00000000000b'::uuid, 2),
    ('a0000000-0000-4000-a000-00000000000c'::uuid, 3)
) AS i(item_id, photo_count)
CROSS JOIN LATERAL generate_series(0, i.photo_count - 1) AS p(position)
JOIN items it ON it.id = i.item_id
ON CONFLICT (item_id, position) DO NOTHING;

UPDATE user_rewards
SET status       = 'granted',
    code         = NULL,
    activated_at = NULL
WHERE user_id::text LIKE 'd0000000-0000-4000-a000-%'
  AND code LIKE '%DEMOSIGNATURE%'
  AND status <> 'expired';

UPDATE user_rewards
SET code = NULL
WHERE user_id::text LIKE 'd0000000-0000-4000-a000-%'
  AND code LIKE '%DEMOSIGNATURE%';

UPDATE pets
SET satiety           = 88,
    happiness         = 90,
    energy            = 82,
    last_decay_time   = now() - interval '30 minutes',
    last_checkin_date = current_date,
    updated_at        = now()
WHERE user_id = 'd0000000-0000-4000-a000-000000000007';

INSERT INTO xp_events (id, user_id, action, subject_id, amount, created_at)
SELECT
    md5('demo-xp-backfill:' || b.user_id::text || ':' || d.n)::uuid,
    b.user_id,
    'daily_login',
    NULL,
    b.per_day + CASE WHEN d.n <= b.remainder THEN 1 ELSE 0 END,
    date_trunc('day', now()) - ((b.span_days - d.n + 2) || ' days')::interval + interval '11 hours'
FROM (
    SELECT
        p.user_id,
        GREATEST(p.xp - COALESCE(e.total, 0), 0)                    AS backfill,
        GREATEST(LEAST(EXTRACT(DAY FROM now() - p.hatched_at)::int, 45), 1) AS span_days,
        GREATEST(p.xp - COALESCE(e.total, 0), 0)
            / GREATEST(LEAST(EXTRACT(DAY FROM now() - p.hatched_at)::int, 45), 1) AS per_day,
        GREATEST(p.xp - COALESCE(e.total, 0), 0)
            % GREATEST(LEAST(EXTRACT(DAY FROM now() - p.hatched_at)::int, 45), 1) AS remainder
    FROM pets p
    LEFT JOIN (
        SELECT user_id, SUM(amount)::int AS total FROM xp_events GROUP BY user_id
    ) e ON e.user_id = p.user_id
    WHERE p.user_id::text LIKE 'd0000000-0000-4000-a000-%'
      AND p.hatched_at IS NOT NULL
) AS b
CROSS JOIN LATERAL generate_series(1, b.span_days) AS d(n)
WHERE b.backfill > 0
  AND b.per_day + CASE WHEN d.n <= b.remainder THEN 1 ELSE 0 END > 0
ON CONFLICT (id) DO NOTHING;

INSERT INTO user_rewards (id, user_id, reward_id, status, code, granted_at, activated_at, expires_at)
SELECT
    md5('demo-reward:anna:' || r.reward_id)::uuid,
    'd0000000-0000-4000-a000-000000000001'::uuid,
    r.reward_id,
    'granted',
    NULL,
    now() - ((16 - r.level) || ' days')::interval,
    NULL,
    CASE WHEN r.promo THEN now() + interval '45 days' END
FROM (VALUES
    ('starter_status',                     1,  false),
    ('raccoon_accessory',                  2,  false),
    ('services_discount_10',               3,  true),
    ('attentive_badge',                    4,  false),
    ('listing_badge_or_delivery_discount', 5,  true),
    ('raccoon_avatar_frame',               6,  false),
    ('favorites_personal_offer',           7,  true),
    ('dialogue_master_badge',              8,  false),
    ('autoteka_discount_30',               9,  true),
    ('free_delivery_500',                  10, true),
    ('xl_listing_discount_50',             11, true),
    ('treasure_hunter_badge',              12, false)
) AS r(reward_id, level, promo)
ON CONFLICT (user_id, reward_id) DO NOTHING;

DELETE FROM daily_summaries WHERE user_id::text LIKE 'd0000000-0000-4000-a000-%';

-- +goose Down
DELETE FROM daily_summaries WHERE user_id::text LIKE 'd0000000-0000-4000-a000-%';

DELETE FROM user_rewards
WHERE user_id = 'd0000000-0000-4000-a000-000000000001'
  AND id IN (
      SELECT md5('demo-reward:anna:' || reward_id)::uuid
      FROM (VALUES
          ('starter_status'), ('raccoon_accessory'), ('services_discount_10'),
          ('attentive_badge'), ('listing_badge_or_delivery_discount'),
          ('raccoon_avatar_frame'), ('favorites_personal_offer'),
          ('dialogue_master_badge'), ('autoteka_discount_30'),
          ('free_delivery_500'), ('xl_listing_discount_50'), ('treasure_hunter_badge')
      ) AS r(reward_id)
  );

DELETE FROM xp_events
WHERE user_id::text LIKE 'd0000000-0000-4000-a000-%'
  AND id IN (
      SELECT md5('demo-xp-backfill:' || user_id::text || ':' || n)::uuid
      FROM xp_events, generate_series(1, 45) AS g(n)
      WHERE user_id::text LIKE 'd0000000-0000-4000-a000-%'
  );

DELETE FROM item_photos
WHERE id IN (
    SELECT md5('demo-photo:' || item_id::text || ':' || position)::uuid
    FROM item_photos
    WHERE item_id::text LIKE 'a0000000-0000-4000-a000-%'
);

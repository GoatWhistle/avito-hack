-- +goose Up

UPDATE user_rewards
SET user_id = 'd0000000-0000-4000-a000-000000000003'
WHERE user_id = 'd0000000-0000-4000-a000-000000000007'
  AND status = 'expired'
  AND NOT EXISTS (
    SELECT 1
    FROM user_rewards existing
    WHERE existing.user_id = 'd0000000-0000-4000-a000-000000000003'
      AND existing.reward_id = user_rewards.reward_id
  );

INSERT INTO user_rewards (id, user_id, reward_id, status, granted_at)
SELECT
  md5('zhanna-granted-' || r.id)::uuid,
  'd0000000-0000-4000-a000-000000000007',
  r.id,
  'granted',
  now() - interval '2 days'
FROM rewards r
WHERE r.condition_type = 'level'
  AND r.condition_value <= 7
  AND NOT EXISTS (
    SELECT 1
    FROM user_rewards ur
    WHERE ur.user_id = 'd0000000-0000-4000-a000-000000000007'
      AND ur.reward_id = r.id
  )
ON CONFLICT (id) DO NOTHING;

-- +goose Down
DELETE FROM user_rewards
WHERE user_id = 'd0000000-0000-4000-a000-000000000007'
  AND id IN (
    SELECT md5('zhanna-granted-' || r.id)::uuid
    FROM rewards r
    WHERE r.condition_type = 'level'
      AND r.condition_value <= 7
  );

UPDATE user_rewards
SET user_id = 'd0000000-0000-4000-a000-000000000007'
WHERE user_id = 'd0000000-0000-4000-a000-000000000003'
  AND status = 'expired';

package infra

import "github.com/redis/go-redis/v9"

var initializeHotState = redis.NewScript(`
local version = tonumber(redis.call('HGET', KEYS[1], 'version') or '-1')
local initial_version = tonumber(ARGV[4])
if version < initial_version then
  redis.call('HSET', KEYS[1],
    'happiness', ARGV[2], 'satiety', ARGV[3],
    'version', ARGV[4], 'updated_at', ARGV[5])
  redis.call('SREM', KEYS[2], ARGV[1])
end
redis.call('EXPIRE', KEYS[1], ARGV[6])
return redis.call('HMGET', KEYS[1], 'happiness', 'satiety', 'version', 'updated_at')
`)

var strokeHotState = redis.NewScript(`
local version = tonumber(redis.call('HGET', KEYS[1], 'version') or '-1')
local initial_version = tonumber(ARGV[4])
if version < initial_version then
  redis.call('HSET', KEYS[1],
    'happiness', ARGV[2], 'satiety', ARGV[3],
    'version', ARGV[4], 'updated_at', ARGV[5])
  redis.call('SREM', KEYS[3], ARGV[1])
end
local applied = redis.call('SET', KEYS[2], '1', 'NX', 'EX', ARGV[8])
if applied then
  local current_happiness = tonumber(redis.call('HGET', KEYS[1], 'happiness'))
  local happiness = math.min(tonumber(ARGV[10]), current_happiness + tonumber(ARGV[9]))
  redis.call('HINCRBY', KEYS[1], 'version', 1)
  redis.call('HSET', KEYS[1], 'happiness', happiness, 'updated_at', ARGV[6])
  redis.call('SADD', KEYS[3], ARGV[1])
end
redis.call('EXPIRE', KEYS[1], ARGV[7])
local values = redis.call('HMGET', KEYS[1], 'happiness', 'satiety', 'version', 'updated_at')
table.insert(values, applied and 1 or 0)
return values
`)

var acknowledgeHotState = redis.NewScript(`
local version = tonumber(redis.call('HGET', KEYS[1], 'version') or '-1')
if version == tonumber(ARGV[2]) then
  return redis.call('SREM', KEYS[2], ARGV[1])
end
return 0
`)

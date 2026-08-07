package infra

import (
	"context"
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func hashReply(happiness, satiety, version, updatedAt string) string {
	return respArray(
		respBulk("happiness"), respBulk(happiness),
		respBulk("satiety"), respBulk(satiety),
		respBulk("version"), respBulk(version),
		respBulk("updated_at"), respBulk(updatedAt),
	)
}

func TestDirtyBatchEmptySet(t *testing.T) {
	t.Parallel()

	client := startFakeRedis(t, map[string]string{"SRANDMEMBER": respArray()})

	states, err := NewRedisHotStateStore(client).DirtyBatch(context.Background(), 5)
	require.NoError(t, err)
	assert.Nil(t, states)
}

func TestDirtyBatchReturnsStates(t *testing.T) {
	t.Parallel()

	updatedAt := time.Date(2026, time.March, 6, 1, 2, 3, 0, time.UTC)
	client := startFakeRedis(t, map[string]string{
		"SRANDMEMBER": respArray(respBulk(hotTestUser.String())),
		"HGETALL":     hashReply("81", "72", "6", strconv.FormatInt(updatedAt.UnixNano(), 10)),
	})

	states, err := NewRedisHotStateStore(client).DirtyBatch(context.Background(), 5)
	require.NoError(t, err)
	require.Len(t, states, 1)

	assert.Equal(t, hotTestUser, states[0].UserID)
	assert.Equal(t, 81, states[0].Happiness)
	assert.Equal(t, 72, states[0].Satiety)
	assert.Equal(t, int64(6), states[0].Version)
	assert.Equal(t, updatedAt, states[0].UpdatedAt)
}

func TestDirtyBatchInvalidUUID(t *testing.T) {
	t.Parallel()

	client := startFakeRedis(t, map[string]string{
		"SRANDMEMBER": respArray(respBulk("not-a-uuid")),
		"HGETALL":     hashReply("1", "2", "3", "4"),
	})

	states, err := NewRedisHotStateStore(client).DirtyBatch(context.Background(), 5)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "parse dirty pet id")
	assert.Nil(t, states)
}

func TestDirtyBatchSkipsMissingHash(t *testing.T) {
	t.Parallel()

	client := startFakeRedis(t, map[string]string{
		"SRANDMEMBER": respArray(respBulk(hotTestUser.String())),
		"HGETALL":     respArray(),
		"SREM":        respInt(1),
	})

	states, err := NewRedisHotStateStore(client).DirtyBatch(context.Background(), 5)
	require.NoError(t, err)
	assert.Empty(t, states)
}

func TestDirtyBatchMalformedHash(t *testing.T) {
	t.Parallel()

	client := startFakeRedis(t, map[string]string{
		"SRANDMEMBER": respArray(respBulk(hotTestUser.String())),
		"HGETALL":     hashReply("nope", "2", "3", "4"),
	})

	states, err := NewRedisHotStateStore(client).DirtyBatch(context.Background(), 5)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "parse redis integer")
	assert.Nil(t, states)
}

func TestRedisCacheGetMiss(t *testing.T) {
	t.Parallel()

	client := startFakeRedis(t, map[string]string{"GET": respNil()})

	pet, err := NewRedisCache(client).Get(context.Background(), hotTestUser)
	require.NoError(t, err)
	assert.Nil(t, pet)
}

func TestRedisCacheGetMalformedPayload(t *testing.T) {
	t.Parallel()

	client := startFakeRedis(t, map[string]string{"GET": respBulk("{not json")})

	pet, err := NewRedisCache(client).Get(context.Background(), hotTestUser)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "decode cached pet")
	assert.Nil(t, pet)
}

func TestRedisCacheSetAndGetRoundTrip(t *testing.T) {
	t.Parallel()

	stored := fixturePet(t, true)
	raw, err := json.Marshal(toCachedPet(stored))
	require.NoError(t, err)

	client := startFakeRedis(t, map[string]string{
		"SET": "+OK\r\n",
		"GET": respBulk(string(raw)),
		"DEL": respInt(1),
	})
	cache := NewRedisCache(client)
	ctx := context.Background()

	require.NoError(t, cache.Set(ctx, stored))

	loaded, err := cache.Get(ctx, stored.UserID())
	require.NoError(t, err)
	require.NotNil(t, loaded)
	assertSamePet(t, stored, loaded)

	require.NoError(t, cache.Delete(ctx, uuid.New()))
}

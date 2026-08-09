package infra

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/pet/app"
)

func unreachableClient(t *testing.T) *redis.Client {
	t.Helper()

	var config net.ListenConfig
	listener, err := config.Listen(t.Context(), "tcp", "127.0.0.1:0")
	require.NoError(t, err)

	addr := listener.Addr().String()
	require.NoError(t, listener.Close())

	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		DialTimeout:  200 * time.Millisecond,
		ReadTimeout:  200 * time.Millisecond,
		WriteTimeout: 200 * time.Millisecond,
		MaxRetries:   -1,
	})
	t.Cleanup(func() { _ = client.Close() })

	return client
}

func testHotState() app.HotState {
	return app.HotState{
		UserID:    hotTestUser,
		Happiness: 50,
		Satiety:   50,
		Version:   1,
		UpdatedAt: time.Date(2026, time.March, 1, 0, 0, 0, 0, time.UTC),
	}
}

func TestRedisHotStateStoreUnreachable(t *testing.T) {
	t.Parallel()

	store := NewRedisHotStateStore(unreachableClient(t))
	require.NotNil(t, store)

	ctx := context.Background()
	initial := testHotState()

	t.Run("get or initialize", func(t *testing.T) {
		t.Parallel()

		state, err := store.GetOrInitialize(ctx, initial)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "initialize redis hot state")
		assert.Zero(t, state)
	})

	t.Run("stroke", func(t *testing.T) {
		t.Parallel()

		now := time.Date(2026, time.March, 1, 1, 0, 0, 0, time.UTC)
		state, applied, err := store.Stroke(ctx, initial, now)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "stroke redis hot state")
		assert.False(t, applied)
		assert.Zero(t, state)
	})

	t.Run("feed", func(t *testing.T) {
		t.Parallel()

		now := time.Date(2026, time.March, 1, 1, 0, 0, 0, time.UTC)
		state, applied, err := store.Feed(ctx, initial, now)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "feed redis hot state")
		assert.False(t, applied)
		assert.Zero(t, state)
	})

	t.Run("dirty batch", func(t *testing.T) {
		t.Parallel()

		states, err := store.DirtyBatch(ctx, 10)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "read dirty pet ids")
		assert.Nil(t, states)
	})

	t.Run("acknowledge", func(t *testing.T) {
		t.Parallel()

		err := store.Acknowledge(ctx, initial)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "acknowledge redis hot state")
	})

	t.Run("read states", func(t *testing.T) {
		t.Parallel()

		values, err := store.readStates(ctx, []string{hotTestUser.String()})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "read dirty pet states")
		assert.Nil(t, values)
	})
}

func TestRedisHotStateStoreReadStatesEmpty(t *testing.T) {
	t.Parallel()

	values, err := NewRedisHotStateStore(unreachableClient(t)).
		readStates(context.Background(), nil)
	require.NoError(t, err)
	assert.Empty(t, values)
}

func TestRedisCacheUnreachable(t *testing.T) {
	t.Parallel()

	cache := NewRedisCache(unreachableClient(t))
	require.NotNil(t, cache)

	ctx := context.Background()
	userID := uuid.MustParse("88888888-8888-8888-8888-888888888888")

	t.Run("get", func(t *testing.T) {
		t.Parallel()

		pet, err := cache.Get(ctx, userID)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "get cached pet")
		assert.Nil(t, pet)
	})

	t.Run("set", func(t *testing.T) {
		t.Parallel()

		err := cache.Set(ctx, fixturePet(t, true))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "set cached pet")
	})

	t.Run("delete", func(t *testing.T) {
		t.Parallel()

		err := cache.Delete(ctx, userID)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "delete cached pet")
	})
}

func TestRedisCacheCancelledContext(t *testing.T) {
	t.Parallel()

	cache := NewRedisCache(unreachableClient(t))
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := cache.Get(ctx, hotTestUser)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "get cached pet")
}

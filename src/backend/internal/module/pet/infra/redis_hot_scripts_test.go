package infra

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/pet/app"
)

func TestHotScriptsAreLoadable(t *testing.T) {
	t.Parallel()

	scripts := map[string]struct {
		hash   string
		source string
	}{
		"initialize":  {hash: initializeHotState.Hash(), source: "HMGET"},
		"stroke":      {hash: strokeHotState.Hash(), source: "HINCRBY"},
		"acknowledge": {hash: acknowledgeHotState.Hash(), source: "SREM"},
	}

	seen := make(map[string]string, len(scripts))
	for name, script := range scripts {
		require.Len(t, script.hash, 40, "script %s hash must be a sha1 hex digest", name)
		assert.NotContains(t, seen, script.hash, "script %s hash collides", name)
		seen[script.hash] = name
	}
}

func TestHotStateConstants(t *testing.T) {
	t.Parallel()

	assert.Equal(t, 24*time.Hour, hotStateTTL)
	assert.Equal(t, time.Second, strokeCooldown)
	assert.Equal(t, 86400, int(hotStateTTL.Seconds()))
	assert.Equal(t, 1, int(strokeCooldown.Seconds()))
	assert.Equal(t, 5*time.Minute, petCacheTTL)
}

func TestHotArgs(t *testing.T) {
	t.Parallel()

	updatedAt := time.Date(2026, time.January, 2, 3, 4, 5, 0, time.UTC)
	args := hotArgs(app.HotState{
		UserID: hotTestUser, Happiness: 61, Satiety: 42, Version: 8, UpdatedAt: updatedAt,
	})

	require.Len(t, args, 6)
	assert.Equal(t, hotTestUser.String(), args[0])
	assert.Equal(t, 61, args[1])
	assert.Equal(t, 42, args[2])
	assert.Equal(t, int64(8), args[3])
	assert.Equal(t, updatedAt.UnixNano(), args[4])
	assert.Equal(t, 86400, args[5])
}

func TestHotArgsRoundTripsThroughHotStateFromValues(t *testing.T) {
	t.Parallel()

	updatedAt := time.Date(2026, time.July, 7, 7, 7, 7, 0, time.UTC)
	initial := app.HotState{
		UserID: hotTestUser, Happiness: 10, Satiety: 20, Version: 3, UpdatedAt: updatedAt,
	}
	args := hotArgs(initial)

	state, err := hotStateFromValues(hotTestUser, []any{
		int64(args[1].(int)), int64(args[2].(int)), args[3], args[4],
	})
	require.NoError(t, err)
	assert.Equal(t, initial, state)
}

func TestCacheKey(t *testing.T) {
	t.Parallel()

	key := cacheKey(hotTestUser)
	assert.Equal(t, "pet:user:"+hotTestUser.String(), key)
	assert.True(t, strings.HasPrefix(key, "pet:user:"))
	assert.NotEqual(t, hotKey(hotTestUser), key)
}

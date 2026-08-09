package infra

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func feedCooldownNow() time.Time {
	return time.Date(2026, time.March, 5, 12, 0, 0, 0, time.UTC)
}

func TestFeedAvailableAtReturnsDeadlineWhileCooldownHolds(t *testing.T) {
	t.Parallel()

	client := startFakeRedis(t, map[string]string{"TTL": respInt(9000)})
	now := feedCooldownNow()

	available, err := NewRedisHotStateStore(client).
		FeedAvailableAt(context.Background(), hotTestUser, now)

	require.NoError(t, err)
	require.NotNil(t, available)
	assert.Equal(t, now.Add(9000*time.Second).UTC(), *available)
}

func TestFeedAvailableAtIsNilWhenCooldownExpired(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		ttl  int64
	}{
		{name: "missing key", ttl: -2},
		{name: "no expiry", ttl: -1},
		{name: "zero", ttl: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := startFakeRedis(t, map[string]string{"TTL": respInt(tt.ttl)})

			available, err := NewRedisHotStateStore(client).
				FeedAvailableAt(context.Background(), hotTestUser, feedCooldownNow())

			require.NoError(t, err)
			assert.Nil(t, available)
		})
	}
}

func TestFeedAvailableAtPropagatesRedisError(t *testing.T) {
	t.Parallel()

	client := startFakeRedis(t, map[string]string{"TTL": "-ERR redis exploded\r\n"})

	available, err := NewRedisHotStateStore(client).
		FeedAvailableAt(context.Background(), hotTestUser, feedCooldownNow())

	require.Error(t, err)
	assert.Nil(t, available)
}

func TestFeedAvailableAtCoversFullCooldown(t *testing.T) {
	t.Parallel()

	client := startFakeRedis(t, map[string]string{
		"TTL": respInt(int64(feedCooldown.Seconds())),
	})
	now := feedCooldownNow()

	available, err := NewRedisHotStateStore(client).
		FeedAvailableAt(context.Background(), hotTestUser, now)

	require.NoError(t, err)
	require.NotNil(t, available)
	assert.Equal(t, now.Add(feedCooldown).UTC(), *available)
}

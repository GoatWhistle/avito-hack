package infra

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetOrInitializeDecodesScriptReply(t *testing.T) {
	t.Parallel()

	updatedAt := time.Date(2026, time.March, 5, 6, 7, 8, 0, time.UTC)
	client := startFakeRedis(t, map[string]string{
		"EVALSHA": respArray(
			respBulk("70"), respBulk("60"), respBulk("4"),
			respBulk(strconv.FormatInt(updatedAt.UnixNano(), 10)),
		),
	})

	state, err := NewRedisHotStateStore(client).
		GetOrInitialize(context.Background(), testHotState())
	require.NoError(t, err)

	assert.Equal(t, hotTestUser, state.UserID)
	assert.Equal(t, 70, state.Happiness)
	assert.Equal(t, 60, state.Satiety)
	assert.Equal(t, int64(4), state.Version)
	assert.Equal(t, updatedAt, state.UpdatedAt)
}

func TestGetOrInitializeRejectsShortReply(t *testing.T) {
	t.Parallel()

	client := startFakeRedis(t, map[string]string{
		"EVALSHA": respArray(respBulk("70"), respBulk("60")),
	})

	state, err := NewRedisHotStateStore(client).
		GetOrInitialize(context.Background(), testHotState())
	require.ErrorIs(t, err, errInvalidHotState)
	assert.Zero(t, state)
}

func TestStrokeApplied(t *testing.T) {
	t.Parallel()

	client := startFakeRedis(t, map[string]string{
		"EVALSHA": respArray(
			respBulk("55"), respBulk("50"), respBulk("2"), respBulk("1000"), respInt(1),
		),
	})

	now := time.Date(2026, time.March, 5, 0, 0, 0, 0, time.UTC)
	state, applied, err := NewRedisHotStateStore(client).
		Stroke(context.Background(), testHotState(), now)
	require.NoError(t, err)

	assert.True(t, applied)
	assert.Equal(t, 55, state.Happiness)
	assert.Equal(t, 50, state.Satiety)
	assert.Equal(t, int64(2), state.Version)
	assert.Equal(t, time.Unix(0, 1000).UTC(), state.UpdatedAt)
}

func TestStrokeNotApplied(t *testing.T) {
	t.Parallel()

	client := startFakeRedis(t, map[string]string{
		"EVALSHA": respArray(
			respBulk("50"), respBulk("50"), respBulk("1"), respBulk("0"), respInt(0),
		),
	})

	now := time.Date(2026, time.March, 5, 0, 0, 0, 0, time.UTC)
	state, applied, err := NewRedisHotStateStore(client).
		Stroke(context.Background(), testHotState(), now)
	require.NoError(t, err)

	assert.False(t, applied)
	assert.Equal(t, 50, state.Happiness)
}

func TestStrokeReplyErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		reply string
	}{
		{name: "too short", reply: respArray(respBulk("1"), respBulk("2"), respBulk("3"), respBulk("4"))},
		{name: "bad state", reply: respArray(
			respBulk("x"), respBulk("2"), respBulk("3"), respBulk("4"), respInt(1),
		)},
		{name: "bad applied flag", reply: respArray(
			respBulk("1"), respBulk("2"), respBulk("3"), respBulk("4"), respBulk("nope"),
		)},
		{name: "nil applied flag", reply: respArray(
			respBulk("1"), respBulk("2"), respBulk("3"), respBulk("4"), respNil(),
		)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := startFakeRedis(t, map[string]string{"EVALSHA": tt.reply})
			now := time.Date(2026, time.March, 5, 0, 0, 0, 0, time.UTC)

			state, applied, err := NewRedisHotStateStore(client).
				Stroke(context.Background(), testHotState(), now)
			require.Error(t, err)
			assert.False(t, applied)
			assert.Zero(t, state)
		})
	}
}

func TestFeedApplied(t *testing.T) {
	t.Parallel()

	client := startFakeRedis(t, map[string]string{
		"EVALSHA": respArray(
			respBulk("50"), respBulk("65"), respBulk("2"), respBulk("1000"), respInt(1),
		),
	})

	now := time.Date(2026, time.March, 5, 0, 0, 0, 0, time.UTC)
	state, applied, err := NewRedisHotStateStore(client).
		Feed(context.Background(), testHotState(), now)
	require.NoError(t, err)

	assert.True(t, applied)
	assert.Equal(t, 65, state.Satiety)
	assert.Equal(t, 50, state.Happiness)
	assert.Equal(t, int64(2), state.Version)
	assert.Equal(t, time.Unix(0, 1000).UTC(), state.UpdatedAt)
}

func TestFeedNotAppliedWithinCooldown(t *testing.T) {
	t.Parallel()

	client := startFakeRedis(t, map[string]string{
		"EVALSHA": respArray(
			respBulk("50"), respBulk("50"), respBulk("1"), respBulk("0"), respInt(0),
		),
	})

	now := time.Date(2026, time.March, 5, 0, 0, 0, 0, time.UTC)
	state, applied, err := NewRedisHotStateStore(client).
		Feed(context.Background(), testHotState(), now)
	require.NoError(t, err)

	assert.False(t, applied)
	assert.Equal(t, 50, state.Satiety)
}

func TestFeedReplyErrors(t *testing.T) {
	t.Parallel()

	client := startFakeRedis(t, map[string]string{
		"EVALSHA": respArray(respBulk("1"), respBulk("2"), respBulk("3"), respBulk("4")),
	})
	now := time.Date(2026, time.March, 5, 0, 0, 0, 0, time.UTC)

	state, applied, err := NewRedisHotStateStore(client).
		Feed(context.Background(), testHotState(), now)

	require.ErrorIs(t, err, errInvalidHotState)
	assert.False(t, applied)
	assert.Zero(t, state)
}

func TestAcknowledgeSucceeds(t *testing.T) {
	t.Parallel()

	client := startFakeRedis(t, map[string]string{"EVALSHA": respInt(1)})

	require.NoError(t, NewRedisHotStateStore(client).
		Acknowledge(context.Background(), testHotState()))
}

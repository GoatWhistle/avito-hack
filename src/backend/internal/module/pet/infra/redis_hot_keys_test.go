package infra

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var hotTestUser = uuid.MustParse("55555555-5555-5555-5555-555555555555")

func TestHotKeys(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "pet:hot:"+hotTestUser.String(), hotKey(hotTestUser))
	assert.Equal(t, "pet:hot:abc", hotKeyString("abc"))
	assert.Equal(t, "pet:hot:", hotKeyString(""))
	assert.Equal(t, "pet:cooldown:stroke:"+hotTestUser.String(), cooldownKey(hotTestUser))
	assert.Equal(t, "pet:hot:dirty", dirtyKey())
	assert.Equal(t, hotKeyString(hotTestUser.String()), hotKey(hotTestUser))
}

func TestValueInt64(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		value   any
		want    int64
		wantErr string
	}{
		{name: "int64 passthrough", value: int64(42), want: 42},
		{name: "negative int64", value: int64(-7), want: -7},
		{name: "numeric string", value: "123", want: 123},
		{name: "negative string", value: "-9", want: -9},
		{name: "zero string", value: "0", want: 0},
		{name: "malformed string", value: "abc", wantErr: "parse redis integer"},
		{name: "empty string", value: "", wantErr: "parse redis integer"},
		{name: "nil value", value: nil, wantErr: "invalid redis hot state response"},
		{name: "unexpected type", value: 3.5, wantErr: "unexpected redis value type float64"},
		{name: "unexpected int type", value: 3, wantErr: "unexpected redis value type int"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := valueInt64(tt.value)
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				assert.Zero(t, got)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestValueInt64NilWrapsSentinel(t *testing.T) {
	t.Parallel()

	_, err := valueInt64(nil)
	require.ErrorIs(t, err, errInvalidHotState)
}

func TestHotStateFromValues(t *testing.T) {
	t.Parallel()

	updatedAt := time.Date(2026, time.February, 3, 4, 5, 6, 0, time.UTC)

	state, err := hotStateFromValues(hotTestUser, []any{
		int64(70), int64(55), int64(9), updatedAt.UnixNano(),
	})
	require.NoError(t, err)

	assert.Equal(t, hotTestUser, state.UserID)
	assert.Equal(t, 70, state.Happiness)
	assert.Equal(t, 55, state.Satiety)
	assert.Equal(t, int64(9), state.Version)
	assert.Equal(t, updatedAt, state.UpdatedAt)
	assert.Equal(t, time.UTC, state.UpdatedAt.Location())
}

func TestHotStateFromValuesStrings(t *testing.T) {
	t.Parallel()

	state, err := hotStateFromValues(hotTestUser, []any{"70", "55", "9", "0"})
	require.NoError(t, err)

	assert.Equal(t, 70, state.Happiness)
	assert.Equal(t, 55, state.Satiety)
	assert.Equal(t, int64(9), state.Version)
	assert.Equal(t, time.Unix(0, 0).UTC(), state.UpdatedAt)
}

func TestHotStateFromValuesIgnoresExtraValues(t *testing.T) {
	t.Parallel()

	state, err := hotStateFromValues(hotTestUser, []any{"1", "2", "3", "4", "5", "6"})
	require.NoError(t, err)
	assert.Equal(t, 1, state.Happiness)
	assert.Equal(t, int64(3), state.Version)
}

func TestHotStateFromValuesErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		values []any
	}{
		{name: "nil slice", values: nil},
		{name: "empty slice", values: []any{}},
		{name: "too short", values: []any{"1", "2", "3"}},
		{name: "bad happiness", values: []any{"x", "2", "3", "4"}},
		{name: "bad satiety", values: []any{"1", "x", "3", "4"}},
		{name: "bad version", values: []any{"1", "2", "x", "4"}},
		{name: "bad updated_at", values: []any{"1", "2", "3", "x"}},
		{name: "nil happiness", values: []any{nil, "2", "3", "4"}},
		{name: "nil updated_at", values: []any{"1", "2", "3", nil}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			state, err := hotStateFromValues(hotTestUser, tt.values)
			require.Error(t, err)
			assert.Zero(t, state)
		})
	}
}

func TestHotStateFromMap(t *testing.T) {
	t.Parallel()

	state, err := hotStateFromMap(hotTestUser, map[string]string{
		"happiness": "80", "satiety": "60", "version": "12", "updated_at": "1000",
	})
	require.NoError(t, err)

	assert.Equal(t, hotTestUser, state.UserID)
	assert.Equal(t, 80, state.Happiness)
	assert.Equal(t, 60, state.Satiety)
	assert.Equal(t, int64(12), state.Version)
	assert.Equal(t, time.Unix(0, 1000).UTC(), state.UpdatedAt)
}

func TestHotStateFromMapMissingFields(t *testing.T) {
	t.Parallel()

	for _, values := range []map[string]string{
		nil,
		{},
		{"happiness": "1"},
		{"happiness": "1", "satiety": "2", "version": "3"},
	} {
		state, err := hotStateFromMap(hotTestUser, values)
		require.Error(t, err)
		assert.Zero(t, state)
	}
}

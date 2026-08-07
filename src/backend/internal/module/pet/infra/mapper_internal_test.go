package infra

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/pet/domain"
	"github.com/avito-hack/backend/internal/shared/postgres/pgtest"
)

var errScan = errors.New("scan failed")

var scanTime = time.Date(2024, 3, 15, 10, 0, 0, 0, time.UTC)

func TestScanRewardMapsColumns(t *testing.T) {
	t.Parallel()

	row := pgtest.Row{Values: []any{"promo10", "Promo", "Desc", "cosmetic", "achievement", 3}}

	reward, err := scanReward(row)

	require.NoError(t, err)
	assert.Equal(t, "promo10", reward.ID())
	assert.Equal(t, "Promo", reward.Title())
	assert.Equal(t, "Desc", reward.Description())
	assert.Equal(t, domain.RewardKindCosmetic, reward.Kind())
	assert.Equal(t, domain.ConditionAchievement, reward.ConditionType())
	assert.Equal(t, 3, reward.ConditionValue())
}

func TestScanRewardPropagatesScanError(t *testing.T) {
	t.Parallel()

	_, err := scanReward(pgtest.Row{Err: errScan})

	require.ErrorIs(t, err, errScan)
}

func TestScanRewardRejectsInvalidValues(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		values []any
	}{
		{name: "empty id", values: []any{"", "T", "D", "promo", "level", 1}},
		{name: "bad kind", values: []any{"r", "T", "D", "unknown", "level", 1}},
		{name: "bad condition", values: []any{"r", "T", "D", "promo", "unknown", 1}},
		{name: "negative value", values: []any{"r", "T", "D", "promo", "level", -1}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := scanReward(pgtest.Row{Values: tc.values})

			require.ErrorIs(t, err, domain.ErrInvalidReward)
		})
	}
}

func TestScanUserRewardMapsNullables(t *testing.T) {
	t.Parallel()

	id := uuid.New()
	userID := uuid.New()
	code := "CODE"
	activated := scanTime.Add(time.Hour)
	expires := scanTime.Add(48 * time.Hour)

	row := pgtest.Row{Values: []any{
		id, userID, "promo10", "activated", &code, scanTime, &activated, &expires,
	}}

	granted, err := scanUserReward(row)

	require.NoError(t, err)
	assert.Equal(t, id, granted.ID())
	assert.Equal(t, userID, granted.UserID())
	assert.Equal(t, domain.RewardActivated, granted.Status())
	assert.Equal(t, "CODE", granted.Code())
	assert.Equal(t, scanTime, granted.GrantedAt())
	require.NotNil(t, granted.ActivatedAt())
	assert.Equal(t, activated, *granted.ActivatedAt())
	require.NotNil(t, granted.ExpiresAt())
	assert.Equal(t, expires, *granted.ExpiresAt())
}

func TestScanUserRewardHandlesNulls(t *testing.T) {
	t.Parallel()

	row := pgtest.Row{Values: []any{
		uuid.New(), uuid.New(), "promo10", "granted", nil, scanTime, nil, nil,
	}}

	granted, err := scanUserReward(row)

	require.NoError(t, err)
	assert.Empty(t, granted.Code())
	assert.Nil(t, granted.ActivatedAt())
	assert.Nil(t, granted.ExpiresAt())
	assert.False(t, granted.IsActivated())
}

func TestScanUserRewardPropagatesScanError(t *testing.T) {
	t.Parallel()

	_, err := scanUserReward(pgtest.Row{Err: errScan})

	require.ErrorIs(t, err, errScan)
}

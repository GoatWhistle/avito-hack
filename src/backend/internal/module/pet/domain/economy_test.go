package domain

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLevelProgressionAndRewards(t *testing.T) {
	t.Parallel()

	pet := New(uuid.New(), testTime())
	progress := pet.awardXP(420)

	assert.Equal(t, MaxLevel, pet.Level())
	assert.Equal(t, StageLegend, pet.Stage())
	assert.Equal(t, 420, pet.NextLevelXP())
	assert.Len(t, progress.UnlockedRewards, 14)
	assert.Equal(t, "free_delivery_and_avito_guru_badge", progress.UnlockedRewards[13])
}

func TestDailyCheckIn(t *testing.T) {
	t.Parallel()

	pet := New(uuid.New(), testTime())
	for dayNumber := range 7 {
		progress, err := pet.DailyCheckIn(testTime().AddDate(0, 0, dayNumber))
		require.NoError(t, err)
		if dayNumber == 6 {
			assert.Equal(t, 2, progress.XPGranted)
		}
	}

	assert.Equal(t, 7, pet.StreakDays())
	_, err := pet.DailyCheckIn(testTime().AddDate(0, 0, 6).Add(time.Hour))
	assert.ErrorIs(t, err, ErrDuplicateAction)
}

func TestDailyCheckInResetsBrokenStreak(t *testing.T) {
	t.Parallel()

	pet := New(uuid.New(), testTime())
	_, err := pet.DailyCheckIn(testTime())
	require.NoError(t, err)
	_, err = pet.DailyCheckIn(testTime().AddDate(0, 0, 2))
	require.NoError(t, err)

	assert.Equal(t, 1, pet.StreakDays())
}

func TestLimitedActions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		action LimitedAction
		call   func(*Pet, LimitedAction) (Progress, error)
		err    error
		xp     int
	}{
		{name: "favorite", action: validLimitedAction(), call: (*Pet).RewardFavorite, xp: 1},
		{name: "favorite duplicate", action: actionWithUnique(false), call: (*Pet).RewardFavorite,
			err: ErrDuplicateAction},
		{name: "favorite limit", action: actionWithCount(5), call: (*Pet).RewardFavorite,
			err: ErrLimitReached},
		{name: "search", action: validLimitedAction(), call: (*Pet).RewardSearchSubscription, xp: 3},
		{name: "search limit", action: actionWithCount(3), call: (*Pet).RewardSearchSubscription,
			err: ErrLimitReached},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			pet := New(uuid.New(), testTime())
			progress, err := tt.call(pet, tt.action)
			assert.ErrorIs(t, err, tt.err)
			assert.Equal(t, tt.xp, progress.XPGranted)
		})
	}
}

func TestContentAndTrustRewards(t *testing.T) {
	t.Parallel()

	pet := New(uuid.New(), testTime())
	progress, err := pet.RewardQualityListing(QualityListing{
		HasVideo: true, Description: strings.Repeat("я", 201), Now: testTime(),
	})
	require.NoError(t, err)
	assert.Equal(t, 2, progress.XPGranted)

	progress, err = pet.RewardPhotoReview(Review{
		ConfirmedDeal: true, HasAttachment: true, Now: testTime(),
	})
	require.NoError(t, err)
	assert.Equal(t, 2, progress.XPGranted)

	_, err = pet.RewardVideoReview(Review{
		ConfirmedDeal: true, VideoDuration: 30 * time.Second, Now: testTime(),
	})
	assert.ErrorIs(t, err, ErrConditionNotMet)
}

func TestQuickReply(t *testing.T) {
	t.Parallel()

	pet := New(uuid.New(), testTime())
	progress, err := pet.RewardQuickReply(59*time.Minute, testTime())
	require.NoError(t, err)
	assert.Equal(t, 3, progress.XPGranted)

	_, err = pet.RewardQuickReply(61*time.Minute, testTime())
	assert.True(t, errors.Is(err, ErrConditionNotMet))
}

func validLimitedAction() LimitedAction {
	return LimitedAction{SubjectID: uuid.New(), Unique: true, Now: testTime()}
}

func actionWithUnique(unique bool) LimitedAction {
	action := validLimitedAction()
	action.Unique = unique
	return action
}

func actionWithCount(count int) LimitedAction {
	action := validLimitedAction()
	action.RewardedCount = count
	return action
}

func testTime() time.Time {
	return time.Date(2026, time.August, 4, 12, 0, 0, 0, time.UTC)
}

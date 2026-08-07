package app_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/raccoon/app"
	"github.com/avito-hack/backend/internal/shared/domainerr"
)

var fixedTime = time.Date(2026, time.March, 3, 9, 0, 0, 0, time.UTC)

type stubPetReader struct {
	view     app.RaccoonProfileView
	err      error
	lastUser uuid.UUID
}

func (s *stubPetReader) Profile(_ context.Context, userID uuid.UUID) (app.RaccoonProfileView, error) {
	s.lastUser = userID
	if s.err != nil {
		return app.RaccoonProfileView{}, s.err
	}

	return s.view, nil
}

type stubBadgeReader struct {
	badges []app.BadgeView
	err    error
}

func (s *stubBadgeReader) Earned(context.Context, uuid.UUID) ([]app.BadgeView, error) {
	if s.err != nil {
		return nil, s.err
	}

	return s.badges, nil
}

type stubActivator struct {
	code       string
	err        error
	lastReward string
	calls      int
}

func (s *stubActivator) Activate(_ context.Context, _ uuid.UUID, rewardID string) (string, error) {
	s.calls++
	s.lastReward = rewardID
	if s.err != nil {
		return "", s.err
	}

	return s.code, nil
}

func sampleProfile(userID uuid.UUID) app.RaccoonProfileView {
	return app.RaccoonProfileView{
		ID: uuid.New(), UserID: userID, Name: "Enot", Level: 4, XP: 120,
		XPToNextLevel: 200, CurrentStreak: 3, Stage: "baby", State: "idle",
	}
}

func TestGetProfileReturnsViewWithBadges(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	earned := fixedTime
	pets := &stubPetReader{view: sampleProfile(userID)}
	badges := &stubBadgeReader{badges: []app.BadgeView{
		{ID: "b1", Name: "First", Description: "desc", IconURL: "http://i/1.png", EarnedAt: &earned},
	}}

	view, err := app.NewGetRaccoonProfileUseCase(pets, badges).Execute(t.Context(), userID)

	require.NoError(t, err)
	assert.Equal(t, userID, pets.lastUser)
	assert.Equal(t, "Enot", view.Name)
	assert.Equal(t, 4, view.Level)
	assert.Equal(t, 3, view.CurrentStreak)
	require.Len(t, view.Badges, 1)
	assert.Equal(t, "b1", view.Badges[0].ID)
}

func TestGetProfileOverwritesBadgesFromReader(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	profile := sampleProfile(userID)
	profile.Badges = []app.BadgeView{{ID: "stale"}}
	pets := &stubPetReader{view: profile}

	view, err := app.NewGetRaccoonProfileUseCase(pets, &stubBadgeReader{}).Execute(t.Context(), userID)

	require.NoError(t, err)
	assert.Empty(t, view.Badges)
}

func TestGetProfileRejectsNilUser(t *testing.T) {
	t.Parallel()

	pets := &stubPetReader{}

	_, err := app.NewGetRaccoonProfileUseCase(pets, &stubBadgeReader{}).Execute(t.Context(), uuid.Nil)

	var invalid *domainerr.InvalidError
	require.ErrorAs(t, err, &invalid)
	assert.Equal(t, "user_id", invalid.Field)
	assert.Equal(t, uuid.Nil, pets.lastUser)
}

func TestGetProfilePropagatesErrors(t *testing.T) {
	t.Parallel()

	sentinel := errors.New("infra down")

	tests := map[string]struct {
		pets   *stubPetReader
		badges *stubBadgeReader
	}{
		"pet reader fails":   {pets: &stubPetReader{err: sentinel}, badges: &stubBadgeReader{}},
		"badge reader fails": {pets: &stubPetReader{}, badges: &stubBadgeReader{err: sentinel}},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := app.NewGetRaccoonProfileUseCase(tt.pets, tt.badges).Execute(t.Context(), uuid.New())

			require.ErrorIs(t, err, sentinel)
		})
	}
}

func TestListBadgesWithoutReaderReturnsEmpty(t *testing.T) {
	t.Parallel()

	badges, err := app.NewGetRaccoonProfileUseCase(&stubPetReader{}, nil).ListBadges(t.Context(), uuid.New())

	require.NoError(t, err)
	assert.NotNil(t, badges)
	assert.Empty(t, badges)
}

func TestListBadgesRejectsNilUser(t *testing.T) {
	t.Parallel()

	_, err := app.NewGetRaccoonProfileUseCase(&stubPetReader{}, &stubBadgeReader{}).
		ListBadges(t.Context(), uuid.Nil)

	var invalid *domainerr.InvalidError
	require.ErrorAs(t, err, &invalid)
	assert.Equal(t, "user_id", invalid.Field)
}

func TestClaimRewardReturnsPromocode(t *testing.T) {
	t.Parallel()

	activator := &stubActivator{code: "PROMO-42"}

	result, err := app.NewClaimRewardUseCase(activator).Execute(t.Context(), uuid.New(), "r1")

	require.NoError(t, err)
	assert.Equal(t, "r1", result.RewardID)
	assert.Equal(t, "PROMO-42", result.Promocode)
	assert.Equal(t, "r1", activator.lastReward)
}

func TestClaimRewardValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		userID    uuid.UUID
		rewardID  string
		wantField string
	}{
		{name: "nil user", userID: uuid.Nil, rewardID: "r1", wantField: "user_id"},
		{name: "empty reward", userID: uuid.New(), rewardID: "", wantField: "reward_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			activator := &stubActivator{}

			_, err := app.NewClaimRewardUseCase(activator).Execute(t.Context(), tt.userID, tt.rewardID)

			var invalid *domainerr.InvalidError
			require.ErrorAs(t, err, &invalid)
			assert.Equal(t, tt.wantField, invalid.Field)
			assert.Zero(t, activator.calls)
		})
	}
}

func TestClaimRewardPropagatesActivationError(t *testing.T) {
	t.Parallel()

	sentinel := errors.New("already activated")

	result, err := app.NewClaimRewardUseCase(&stubActivator{err: sentinel}).
		Execute(t.Context(), uuid.New(), "r1")

	require.ErrorIs(t, err, sentinel)
	assert.Empty(t, result.Promocode)
}

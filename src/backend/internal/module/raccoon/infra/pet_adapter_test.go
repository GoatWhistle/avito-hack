package infra_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	petdomain "github.com/avito-hack/backend/internal/module/pet/domain"
	"github.com/avito-hack/backend/internal/module/raccoon/infra"
)

var errAdapter = errors.New("adapter failure")

type stubPetState struct {
	pet    *petdomain.Pet
	err    error
	called uuid.UUID
}

func (s *stubPetState) State(_ context.Context, userID uuid.UUID) (*petdomain.Pet, error) {
	s.called = userID
	if s.err != nil {
		return nil, s.err
	}

	return s.pet, nil
}

type stubBadges struct {
	badges []petdomain.EarnedBadge
	err    error
	called uuid.UUID
}

func (s *stubBadges) EarnedBy(_ context.Context, userID uuid.UUID) ([]petdomain.EarnedBadge, error) {
	s.called = userID
	if s.err != nil {
		return nil, s.err
	}

	return s.badges, nil
}

func newTestPet(t *testing.T, userID uuid.UUID) *petdomain.Pet {
	t.Helper()

	return petdomain.Restore(petdomain.RestoreParams{
		ID:          uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		UserID:      userID,
		Name:        "Enot",
		Stage:       petdomain.StageBaby,
		Level:       4,
		XP:          120,
		NextLevelXP: 200,
		Satiety:     90,
		Happiness:   95,
		Energy:      80,
		StreakDays:  7,
		Freezes:     1,
		UpdatedAt:   time.Date(2026, time.March, 1, 12, 0, 0, 0, time.UTC),
	})
}

func TestPetAdapterProfile(t *testing.T) {
	t.Parallel()

	userID := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	pet := newTestPet(t, userID)
	pets := &stubPetState{pet: pet}

	view, err := infra.NewPetAdapter(pets).Profile(context.Background(), userID)
	require.NoError(t, err)

	assert.Equal(t, userID, pets.called)
	assert.Equal(t, pet.ID(), view.ID)
	assert.Equal(t, userID, view.UserID)
	assert.Equal(t, "Enot", view.Name)
	assert.Equal(t, 4, view.Level)
	assert.Equal(t, 120, view.XP)
	assert.Equal(t, pet.NextLevelXP(), view.XPToNextLevel)
	assert.Equal(t, 7, view.CurrentStreak)
	assert.Equal(t, string(pet.Stage()), view.Stage)
	assert.Equal(t, string(pet.State()), view.State)
	assert.Nil(t, view.Badges)
}

func TestPetAdapterProfileError(t *testing.T) {
	t.Parallel()

	adapter := infra.NewPetAdapter(&stubPetState{err: errAdapter})

	view, err := adapter.Profile(context.Background(), uuid.New())
	require.ErrorIs(t, err, errAdapter)
	assert.Zero(t, view)
}

func TestBadgeAdapterEarned(t *testing.T) {
	t.Parallel()

	userID := uuid.MustParse("33333333-3333-3333-3333-333333333333")
	earnedAt := time.Date(2026, time.April, 2, 9, 30, 0, 0, time.UTC)
	badges := &stubBadges{badges: []petdomain.EarnedBadge{
		petdomain.NewEarnedBadge(
			petdomain.NewBadge("first_deal", "First deal", "Closed the first deal", "/icons/first.png"),
			userID, earnedAt,
		),
		petdomain.NewEarnedBadge(
			petdomain.NewBadge("raccoon_friend", "Friend", "Reached level 10", "/icons/friend.png"),
			userID, earnedAt.Add(time.Hour),
		),
	}}

	views, err := infra.NewBadgeAdapter(badges).Earned(context.Background(), userID)
	require.NoError(t, err)
	require.Len(t, views, 2)

	assert.Equal(t, userID, badges.called)
	assert.Equal(t, "first_deal", views[0].ID)
	assert.Equal(t, "First deal", views[0].Name)
	assert.Equal(t, "Closed the first deal", views[0].Description)
	assert.Equal(t, "/icons/first.png", views[0].IconURL)
	require.NotNil(t, views[0].EarnedAt)
	assert.Equal(t, earnedAt, *views[0].EarnedAt)

	require.NotNil(t, views[1].EarnedAt)
	assert.Equal(t, earnedAt.Add(time.Hour), *views[1].EarnedAt)
	assert.NotSame(t, views[0].EarnedAt, views[1].EarnedAt)
}

func TestBadgeAdapterEarnedEmpty(t *testing.T) {
	t.Parallel()

	views, err := infra.NewBadgeAdapter(&stubBadges{}).Earned(context.Background(), uuid.New())
	require.NoError(t, err)
	assert.Empty(t, views)
	assert.NotNil(t, views)
}

func TestBadgeAdapterEarnedError(t *testing.T) {
	t.Parallel()

	views, err := infra.NewBadgeAdapter(&stubBadges{err: errAdapter}).Earned(context.Background(), uuid.New())
	require.ErrorIs(t, err, errAdapter)
	assert.Nil(t, views)
}

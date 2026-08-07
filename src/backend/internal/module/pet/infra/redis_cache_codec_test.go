package infra

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/pet/domain"
)

func fixturePet(t *testing.T, withOptionals bool) *domain.Pet {
	t.Helper()

	params := domain.RestoreParams{
		ID:                 uuid.MustParse("66666666-6666-6666-6666-666666666666"),
		UserID:             uuid.MustParse("77777777-7777-7777-7777-777777777777"),
		Name:               "Enot",
		Stage:              domain.StageTeen,
		Level:              12,
		XP:                 4200,
		NextLevelXP:        5000,
		Satiety:            64,
		Happiness:          88,
		Energy:             51,
		StreakDays:         21,
		Freezes:            2,
		LastDecayTime:      time.Date(2026, time.March, 10, 8, 0, 0, 0, time.UTC),
		UpdatedAt:          time.Date(2026, time.March, 10, 9, 0, 0, 0, time.UTC),
		InteractionVersion: 17,
	}
	if withOptionals {
		checkIn := time.Date(2026, time.March, 9, 0, 0, 0, 0, time.UTC)
		hatched := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
		params.LastCheckInDate = &checkIn
		params.HatchedAt = &hatched
	}

	return domain.Restore(params)
}

func assertSamePet(t *testing.T, want, got *domain.Pet) {
	t.Helper()

	assert.Equal(t, want.ID(), got.ID())
	assert.Equal(t, want.UserID(), got.UserID())
	assert.Equal(t, want.Name(), got.Name())
	assert.Equal(t, want.Stage(), got.Stage())
	assert.Equal(t, want.Level(), got.Level())
	assert.Equal(t, want.XP(), got.XP())
	assert.Equal(t, want.NextLevelXP(), got.NextLevelXP())
	assert.Equal(t, want.Satiety(), got.Satiety())
	assert.Equal(t, want.Happiness(), got.Happiness())
	assert.Equal(t, want.Energy(), got.Energy())
	assert.Equal(t, want.StreakDays(), got.StreakDays())
	assert.Equal(t, want.Freezes(), got.Freezes())
	assert.Equal(t, want.LastCheckInDate(), got.LastCheckInDate())
	assert.Equal(t, want.HatchedAt(), got.HatchedAt())
	assert.True(t, want.LastDecayTime().Equal(got.LastDecayTime()))
	assert.True(t, want.UpdatedAt().Equal(got.UpdatedAt()))
	assert.Equal(t, want.InteractionVersion(), got.InteractionVersion())
}

func TestCachedPetRoundTripWithOptionals(t *testing.T) {
	t.Parallel()

	pet := fixturePet(t, true)

	raw, err := json.Marshal(toCachedPet(pet))
	require.NoError(t, err)

	var decoded cachedPet
	require.NoError(t, json.Unmarshal(raw, &decoded))
	assertSamePet(t, pet, decoded.toDomain())
}

func TestCachedPetRoundTripWithoutOptionals(t *testing.T) {
	t.Parallel()

	pet := fixturePet(t, false)

	raw, err := json.Marshal(toCachedPet(pet))
	require.NoError(t, err)
	assert.NotContains(t, string(raw), "last_checkin_date")
	assert.NotContains(t, string(raw), "hatched_at")

	var decoded cachedPet
	require.NoError(t, json.Unmarshal(raw, &decoded))
	assert.Nil(t, decoded.LastCheckInDate)
	assert.Nil(t, decoded.HatchedAt)
	assertSamePet(t, pet, decoded.toDomain())
}

func TestToCachedPetFields(t *testing.T) {
	t.Parallel()

	pet := fixturePet(t, true)
	value := toCachedPet(pet)

	assert.Equal(t, pet.ID(), value.ID)
	assert.Equal(t, pet.UserID(), value.UserID)
	assert.Equal(t, domain.StageTeen, value.Stage)
	assert.Equal(t, 12, value.Level)
	assert.Equal(t, 4200, value.XP)
	assert.Equal(t, pet.NextLevelXP(), value.NextLevelXP)
	assert.Equal(t, 64, value.Satiety)
	assert.Equal(t, 88, value.Happiness)
	assert.Equal(t, 51, value.Energy)
	assert.Equal(t, 21, value.StreakDays)
	assert.Equal(t, 2, value.Freezes)
	assert.Equal(t, int64(17), value.InteractionVersion)
}

func TestCachedPetJSONKeys(t *testing.T) {
	t.Parallel()

	raw, err := json.Marshal(toCachedPet(fixturePet(t, true)))
	require.NoError(t, err)

	var decoded map[string]any
	require.NoError(t, json.Unmarshal(raw, &decoded))

	for _, key := range []string{
		"id", "user_id", "name", "stage", "level", "xp", "next_level_xp",
		"satiety", "happiness", "energy", "streak_days", "freezes",
		"last_checkin_date", "hatched_at", "last_decay_time", "updated_at",
		"interaction_version",
	} {
		assert.Contains(t, decoded, key)
	}
}

func TestCachedPetDecodeMalformed(t *testing.T) {
	t.Parallel()

	var decoded cachedPet
	require.Error(t, json.Unmarshal([]byte(`{"level":"twelve"}`), &decoded))
	require.Error(t, json.Unmarshal([]byte(`not json`), &decoded))
}

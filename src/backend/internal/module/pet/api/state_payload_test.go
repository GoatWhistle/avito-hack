package api_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/pet/api"
	"github.com/avito-hack/backend/internal/module/pet/app"
	"github.com/avito-hack/backend/internal/module/pet/domain"
	"github.com/avito-hack/backend/internal/shared/auth"
)

func decodeBody(t *testing.T, raw []byte) map[string]any {
	t.Helper()

	var body map[string]any
	require.NoError(t, json.Unmarshal(raw, &body))

	return body
}

func TestPetGetExposesFeedAvailableAt(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	actor := auth.Actor{ID: userID, Role: auth.RoleUser}
	available := petTestTime.Add(5 * time.Hour)
	service := &stubPetService{pet: testPet(userID), feedAvailableAt: &available}

	router := petRouter(t, api.PetDeps{Service: service}, &actor)
	rec := doRequest(t, router, http.MethodGet, "/pet")

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, available.Format(time.RFC3339Nano), decodeBody(t, rec.Body.Bytes())["feed_available_at"])
}

func TestPetGetReportsNullFeedAvailableAt(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	actor := auth.Actor{ID: userID, Role: auth.RoleUser}
	router := petRouter(t, api.PetDeps{Service: &stubPetService{pet: testPet(userID)}}, &actor)

	rec := doRequest(t, router, http.MethodGet, "/pet")

	require.Equal(t, http.StatusOK, rec.Code)

	body := decodeBody(t, rec.Body.Bytes())
	require.Contains(t, body, "feed_available_at")
	assert.Nil(t, body["feed_available_at"])
	assert.Equal(t, false, body["checkin_applied"])
	assert.NotContains(t, body, "checkin")
}

func TestPetFeedExposesFeedAvailableAt(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	actor := auth.Actor{ID: userID, Role: auth.RoleUser}
	available := petTestTime.Add(5 * time.Hour)
	service := &stubPetService{pet: testPet(userID), feedAvailableAt: &available}

	router := petRouter(t, api.PetDeps{Service: service}, &actor)
	rec := doRequest(t, router, http.MethodPost, "/pet/actions/feed")

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, available.Format(time.RFC3339Nano), decodeBody(t, rec.Body.Bytes())["feed_available_at"])
}

func TestPetGetReportsAutomaticCheckIn(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	actor := auth.Actor{ID: userID, Role: auth.RoleUser}
	granter := &spyGranter{}
	service := &stubPetService{
		pet:            testPet(userID),
		checkInApplied: true,
		result: app.ActionResult{
			Progress: domain.Progress{
				XPGranted: 2, PreviousLevel: 3, Level: 3, NextLevelXP: 22,
				UnlockedRewards: []string{"attentive_badge"},
			},
			Streak: domain.StreakOutcome{Days: 5, Continued: true, FreezesLeft: 1},
		},
	}

	router := petRouter(t, api.PetDeps{Service: service, Rewards: granter}, &actor)
	rec := doRequest(t, router, http.MethodGet, "/pet")

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, 1, granter.calls)

	var body struct {
		CheckInApplied bool `json:"checkin_applied"`
		CheckIn        *struct {
			XPGranted       int      `json:"xp_granted"`
			Level           int      `json:"level"`
			UnlockedRewards []string `json:"unlocked_rewards"`
			Streak          struct {
				Days      int  `json:"days"`
				Continued bool `json:"continued"`
			} `json:"streak"`
		} `json:"checkin"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))

	assert.True(t, body.CheckInApplied)
	require.NotNil(t, body.CheckIn)
	assert.Equal(t, 2, body.CheckIn.XPGranted)
	assert.Equal(t, 3, body.CheckIn.Level)
	assert.Equal(t, []string{"attentive_badge"}, body.CheckIn.UnlockedRewards)
	assert.Equal(t, 5, body.CheckIn.Streak.Days)
	assert.True(t, body.CheckIn.Streak.Continued)
}

func TestPetGetSkipsRewardGrantWithoutCheckIn(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	actor := auth.Actor{ID: userID, Role: auth.RoleUser}
	granter := &spyGranter{}
	service := &stubPetService{pet: testPet(userID)}

	router := petRouter(t, api.PetDeps{Service: service, Rewards: granter}, &actor)
	rec := doRequest(t, router, http.MethodGet, "/pet")

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Zero(t, granter.calls)
}

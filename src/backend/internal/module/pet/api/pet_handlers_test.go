package api_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/pet/api"
	"github.com/avito-hack/backend/internal/module/pet/app"
	"github.com/avito-hack/backend/internal/module/pet/domain"
	"github.com/avito-hack/backend/internal/shared/auth"
)

func TestPetGetReturnsState(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	actor := auth.Actor{ID: userID, Role: auth.RoleUser}
	router := petRouter(t, api.PetDeps{Service: &stubPetService{pet: testPet(userID)}}, &actor)

	rec := doRequest(t, router, http.MethodGet, "/pet")

	require.Equal(t, http.StatusOK, rec.Code)

	var body map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.Equal(t, "Enot", body["name"])
	assert.InDelta(t, 3, body["level"], 0)
	assert.Equal(t, true, body["is_hatched"])
	assert.Equal(t, userID.String(), body["user_id"])
}

func TestPetEndpointsRequireActor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/pet"},
		{http.MethodPost, "/pet/actions/stroke"},
		{http.MethodPost, "/pet/actions/feed"},
		{http.MethodPost, "/checkin"},
		{http.MethodGet, "/progress"},
	}

	for _, tt := range tests {
		t.Run(tt.method+" "+tt.path, func(t *testing.T) {
			t.Parallel()

			router := petRouter(t, api.PetDeps{Service: &stubPetService{pet: testPet(uuid.New())}}, nil)

			rec := doRequest(t, router, tt.method, tt.path)

			assert.Equal(t, http.StatusUnauthorized, rec.Code)
			assert.Contains(t, rec.Body.String(), "unauthorized")
		})
	}
}

func TestPetStrokeReturnsUpdatedPet(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	actor := auth.Actor{ID: userID, Role: auth.RoleUser}
	router := petRouter(t, api.PetDeps{Service: &stubPetService{pet: testPet(userID)}}, &actor)

	rec := doRequest(t, router, http.MethodPost, "/pet/actions/stroke")

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"happiness":80`)
}

func TestPetCheckInGrantsRewardsAndReturnsStreak(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	actor := auth.Actor{ID: userID, Role: auth.RoleUser}
	granter := &spyGranter{}
	service := &stubPetService{result: app.ActionResult{
		Pet: testPet(userID),
		Progress: domain.Progress{
			XPGranted: 25, PreviousLevel: 2, Level: 3, NextLevelXP: 100,
			UnlockedRewards: []string{"r3"},
		},
		Streak: domain.StreakOutcome{Days: 4, Continued: true, MilestoneBonus: 10, FreezesLeft: 1},
	}}

	router := petRouter(t, api.PetDeps{Service: service, Rewards: granter}, &actor)
	rec := doRequest(t, router, http.MethodPost, "/checkin")

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, 1, granter.calls)

	var body struct {
		XPGranted       int      `json:"xp_granted"`
		UnlockedRewards []string `json:"unlocked_rewards"`
		Streak          struct {
			Days           int  `json:"days"`
			Continued      bool `json:"continued"`
			MilestoneBonus int  `json:"milestone_bonus"`
		} `json:"streak"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.Equal(t, 25, body.XPGranted)
	assert.Equal(t, []string{"r3"}, body.UnlockedRewards)
	assert.Equal(t, 4, body.Streak.Days)
	assert.True(t, body.Streak.Continued)
	assert.Equal(t, 10, body.Streak.MilestoneBonus)
}

func TestPetCheckInWithoutGranterStillSucceeds(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	actor := auth.Actor{ID: userID, Role: auth.RoleUser}
	service := &stubPetService{result: app.ActionResult{Pet: testPet(userID)}}

	router := petRouter(t, api.PetDeps{Service: service}, &actor)
	rec := doRequest(t, router, http.MethodPost, "/checkin")

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"unlocked_rewards":[]`)
}

func TestPetCheckInMapsDomainErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantCode   string
	}{
		{"duplicate", domain.ErrDuplicateAction, http.StatusConflict, "conflict"},
		{"limit reached", domain.ErrLimitReached, http.StatusConflict, "conflict"},
		{"condition not met", domain.ErrConditionNotMet, http.StatusBadRequest, "validation_error"},
		{"invalid action", domain.ErrInvalidAction, http.StatusBadRequest, "validation_error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actor := auth.Actor{ID: uuid.New(), Role: auth.RoleUser}
			router := petRouter(t, api.PetDeps{Service: &stubPetService{checkInErr: tt.err}}, &actor)

			rec := doRequest(t, router, http.MethodPost, "/checkin")

			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.Contains(t, rec.Body.String(), tt.wantCode)
		})
	}
}

func TestPetProgressReturnsView(t *testing.T) {
	t.Parallel()

	actor := auth.Actor{ID: uuid.New(), Role: auth.RoleUser}
	service := &stubPetService{progress: app.ProgressView{
		Level: 5, XP: 210, NextLevelXP: 300, XPToNext: 90,
		Stage: domain.StageTeen, StreakDays: 7, Freezes: 2,
	}}

	router := petRouter(t, api.PetDeps{Service: service}, &actor)
	rec := doRequest(t, router, http.MethodGet, "/progress")

	require.Equal(t, http.StatusOK, rec.Code)

	var body struct {
		Level      int    `json:"level"`
		XPToNext   int    `json:"xp_to_next_level"`
		Stage      string `json:"stage"`
		IsMaxLevel bool   `json:"is_max_level"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.Equal(t, 5, body.Level)
	assert.Equal(t, 90, body.XPToNext)
	assert.Equal(t, string(domain.StageTeen), body.Stage)
	assert.False(t, body.IsMaxLevel)
}

package api_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/raccoon/api"
	"github.com/avito-hack/backend/internal/module/raccoon/app"
	"github.com/avito-hack/backend/internal/shared/auth"
	"github.com/avito-hack/backend/internal/shared/domainerr"
)

func TestGetProfileEndpoint(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	actor := auth.Actor{ID: userID, Role: auth.RoleUser}
	earned := raccoonTime
	badges := &stubBadgeReader{badges: []app.BadgeView{
		{ID: "b1", Name: "First", Description: "d", IconURL: "http://i/1.png", EarnedAt: &earned},
	}}

	rec := newFixture(t, &stubPetReader{view: sampleView()}, badges, &actor).
		do(http.MethodGet, "/raccoon/profile", "")

	require.Equal(t, http.StatusOK, rec.Code)

	var body api.RaccoonProfileResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.Equal(t, "Enot", body.Name)
	assert.Equal(t, userID.String(), body.UserID)
	assert.Equal(t, 4, body.Level)
	assert.Equal(t, 3, body.CurrentStreak)
	require.Len(t, body.Badges, 1)
	assert.Equal(t, "http://i/1.png", body.Badges[0].Icon)
	assert.Equal(t, raccoonTime.Format(time.RFC3339), body.Badges[0].EarnedAt)
}

func TestGetBadgesEndpointReturnsArray(t *testing.T) {
	t.Parallel()

	actor := auth.Actor{ID: uuid.New(), Role: auth.RoleUser}

	rec := newFixture(t, &stubPetReader{view: sampleView()}, &stubBadgeReader{}, &actor).
		do(http.MethodGet, "/badges", "")

	require.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `[]`, rec.Body.String())
}

func TestGetBadgesOmitsUnearnedTimestamp(t *testing.T) {
	t.Parallel()

	actor := auth.Actor{ID: uuid.New(), Role: auth.RoleUser}
	badges := &stubBadgeReader{badges: []app.BadgeView{{ID: "b1", Name: "Locked"}}}

	rec := newFixture(t, &stubPetReader{view: sampleView()}, badges, &actor).
		do(http.MethodGet, "/badges", "")

	require.Equal(t, http.StatusOK, rec.Code)
	assert.NotContains(t, rec.Body.String(), "earned_at")
}

func TestClaimRewardEndpoint(t *testing.T) {
	t.Parallel()

	actor := auth.Actor{ID: uuid.New(), Role: auth.RoleUser}
	f := newFixture(t, &stubPetReader{view: sampleView()}, &stubBadgeReader{}, &actor)

	rec := f.do(http.MethodPost, "/rewards/claim", `{"reward_id":"r1"}`)

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "r1", f.activator.lastReward)

	var body api.ClaimRewardResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.Equal(t, "r1", body.RewardID)
	assert.Equal(t, "PROMO-1", body.Promocode)
}

func TestClaimRewardRejectsInvalidPayload(t *testing.T) {
	t.Parallel()

	tests := map[string]string{
		"missing reward id": `{}`,
		"empty reward id":   `{"reward_id":""}`,
		"malformed json":    `{"reward_id":`,
		"unknown field":     `{"reward_id":"r1","evil":true}`,
	}

	for name, body := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			actor := auth.Actor{ID: uuid.New(), Role: auth.RoleUser}
			f := newFixture(t, &stubPetReader{view: sampleView()}, &stubBadgeReader{}, &actor)

			rec := f.do(http.MethodPost, "/rewards/claim", body)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Empty(t, f.activator.lastReward)
		})
	}
}

func TestRaccoonEndpointsRequireActor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		method string
		path   string
		body   string
	}{
		{http.MethodGet, "/raccoon/profile", ""},
		{http.MethodGet, "/badges", ""},
		{http.MethodPost, "/rewards/claim", `{"reward_id":"r1"}`},
	}

	for _, tt := range tests {
		t.Run(tt.method+" "+tt.path, func(t *testing.T) {
			t.Parallel()

			f := newFixture(t, &stubPetReader{view: sampleView()}, &stubBadgeReader{}, nil)

			rec := f.do(tt.method, tt.path, tt.body)

			assert.Equal(t, http.StatusUnauthorized, rec.Code)
			assert.Contains(t, rec.Body.String(), "unauthorized")
		})
	}
}

func TestRaccoonProfileMapsDomainErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{"not found", domainerr.ErrNotFound, http.StatusNotFound},
		{"forbidden", domainerr.ErrForbidden, http.StatusForbidden},
		{"conflict", domainerr.ErrConflict, http.StatusConflict},
		{"unexpected", errors.New("db down"), http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actor := auth.Actor{ID: uuid.New(), Role: auth.RoleUser}
			f := newFixture(t, &stubPetReader{err: tt.err}, &stubBadgeReader{}, &actor)

			rec := f.do(http.MethodGet, "/raccoon/profile", "")

			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestClaimRewardHidesInternalDetails(t *testing.T) {
	t.Parallel()

	actor := auth.Actor{ID: uuid.New(), Role: auth.RoleUser}
	f := newFixture(t, &stubPetReader{view: sampleView()}, &stubBadgeReader{}, &actor)
	f.activator.err = errors.New("pgx: dial tcp 10.0.0.7:5432 refused")

	rec := f.do(http.MethodPost, "/rewards/claim", `{"reward_id":"r1"}`)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.NotContains(t, rec.Body.String(), "10.0.0.7")
}

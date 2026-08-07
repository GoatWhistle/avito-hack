package api_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/raccoon/app"
	"github.com/avito-hack/backend/internal/shared/auth"
	"github.com/avito-hack/backend/internal/shared/domainerr"
)

func TestGetBadgesMapsErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{name: "not found", err: domainerr.ErrNotFound, wantStatus: http.StatusNotFound},
		{name: "forbidden", err: domainerr.ErrForbidden, wantStatus: http.StatusForbidden},
		{name: "unexpected", err: errors.New("badge table locked"), wantStatus: http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actor := auth.Actor{ID: uuid.New(), Role: auth.RoleUser}
			f := newFixture(t, &stubPetReader{view: sampleView()}, &stubBadgeReader{err: tt.err}, &actor)

			rec := f.do(http.MethodGet, "/badges", "")

			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestGetProfilePropagatesBadgeFailure(t *testing.T) {
	t.Parallel()

	actor := auth.Actor{ID: uuid.New(), Role: auth.RoleUser}
	badges := &stubBadgeReader{err: errors.New("badge read model exploded")}
	f := newFixture(t, &stubPetReader{view: sampleView()}, badges, &actor)

	rec := f.do(http.MethodGet, "/raccoon/profile", "")

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestGetBadgesReturnsEarnedBadges(t *testing.T) {
	t.Parallel()

	actor := auth.Actor{ID: uuid.New(), Role: auth.RoleUser}
	earned := raccoonTime
	badges := &stubBadgeReader{badges: []app.BadgeView{
		{ID: "b1", Name: "First", Description: "d1", IconURL: "http://i/1.png", EarnedAt: &earned},
		{ID: "b2", Name: "Second", Description: "d2"},
	}}

	rec := newFixture(t, &stubPetReader{view: sampleView()}, badges, &actor).
		do(http.MethodGet, "/badges", "")

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"id":"b1"`)
	assert.Contains(t, rec.Body.String(), `"id":"b2"`)
}

func TestClaimRewardMapsDomainErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{name: "not found", err: domainerr.ErrNotFound, wantStatus: http.StatusNotFound},
		{name: "conflict", err: domainerr.ErrConflict, wantStatus: http.StatusConflict},
		{name: "forbidden", err: domainerr.ErrForbidden, wantStatus: http.StatusForbidden},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actor := auth.Actor{ID: uuid.New(), Role: auth.RoleUser}
			f := newFixture(t, &stubPetReader{view: sampleView()}, &stubBadgeReader{}, &actor)
			f.activator.err = tt.err

			rec := f.do(http.MethodPost, "/rewards/claim", `{"reward_id":"r1"}`)

			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

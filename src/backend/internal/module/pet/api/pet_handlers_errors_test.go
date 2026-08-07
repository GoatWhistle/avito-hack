package api_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/pet/api"
	"github.com/avito-hack/backend/internal/module/pet/app"
	"github.com/avito-hack/backend/internal/shared/auth"
	"github.com/avito-hack/backend/internal/shared/domainerr"
)

func testActor() auth.Actor {
	return testActorFor(uuid.New())
}

func testActorFor(userID uuid.UUID) auth.Actor {
	return auth.Actor{ID: userID, Role: auth.RoleUser}
}

func checkInResult(userID uuid.UUID) app.ActionResult {
	return app.ActionResult{Pet: testPet(userID)}
}

func TestPetGetPropagatesServiceErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{name: "not found", err: domainerr.ErrNotFound, wantStatus: http.StatusNotFound},
		{name: "unauthorized", err: domainerr.ErrUnauthorized, wantStatus: http.StatusUnauthorized},
		{name: "unexpected", err: errors.New("pg connection lost"), wantStatus: http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actor := testActor()
			router := petRouter(t, api.PetDeps{Service: &stubPetService{stateErr: tt.err}}, &actor)

			rec := doRequest(t, router, http.MethodGet, "/pet")

			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestPetStrokePropagatesServiceError(t *testing.T) {
	t.Parallel()

	actor := testActor()
	router := petRouter(t, api.PetDeps{
		Service: &stubPetService{strokeErr: domainerr.ErrNotFound},
	}, &actor)

	rec := doRequest(t, router, http.MethodPost, "/pet/actions/stroke")

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestPetProgressPropagatesServiceError(t *testing.T) {
	t.Parallel()

	actor := testActor()
	router := petRouter(t, api.PetDeps{
		Service: &stubPetService{progErr: errors.New("redis unreachable")},
	}, &actor)

	rec := doRequest(t, router, http.MethodGet, "/progress")

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestPetCheckInPassesThroughUnmappedErrors(t *testing.T) {
	t.Parallel()

	actor := testActor()
	router := petRouter(t, api.PetDeps{
		Service: &stubPetService{checkInErr: errors.New("transaction aborted")},
	}, &actor)

	rec := doRequest(t, router, http.MethodPost, "/checkin")

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestPetCheckInIgnoresRewardGrantFailure(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	actor := testActorFor(userID)
	granter := &spyGranter{err: errors.New("rewards catalog unavailable")}
	router := petRouter(t, api.PetDeps{
		Service: &stubPetService{result: checkInResult(userID)},
		Rewards: granter,
	}, &actor)

	rec := doRequest(t, router, http.MethodPost, "/checkin")

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, 1, granter.calls)
}

package api_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/pet/api"
	"github.com/avito-hack/backend/internal/module/pet/domain"
	"github.com/avito-hack/backend/internal/shared/auth"
	"github.com/avito-hack/backend/internal/shared/domainerr"
)

type spyPetNotifier struct {
	calls   int
	userID  uuid.UUID
	lastPet *domain.Pet
}

func (s *spyPetNotifier) PetUpdated(userID uuid.UUID, pet *domain.Pet) {
	s.calls++
	s.userID = userID
	s.lastPet = pet
}

func TestPetFeedReturnsUpdatedPet(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	actor := auth.Actor{ID: userID, Role: auth.RoleUser}
	service := &stubPetService{pet: testPet(userID)}
	router := petRouter(t, api.PetDeps{Service: service}, &actor)

	rec := doRequest(t, router, http.MethodPost, "/pet/actions/feed")

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, 1, service.feedCalls)

	var body map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.InDelta(t, 70, body["satiety"], 0)
	assert.Equal(t, userID.String(), body["user_id"])
}

func TestPetFeedRequiresActor(t *testing.T) {
	t.Parallel()

	service := &stubPetService{pet: testPet(uuid.New())}
	router := petRouter(t, api.PetDeps{Service: service}, nil)

	rec := doRequest(t, router, http.MethodPost, "/pet/actions/feed")

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), "unauthorized")
	assert.Zero(t, service.feedCalls)
}

func TestPetFeedBroadcastsPetUpdated(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	actor := auth.Actor{ID: userID, Role: auth.RoleUser}
	notifier := &spyPetNotifier{}
	router := petRouter(t, api.PetDeps{
		Service: &stubPetService{pet: testPet(userID)}, Notifier: notifier,
	}, &actor)

	rec := doRequest(t, router, http.MethodPost, "/pet/actions/feed")

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, 1, notifier.calls)
	assert.Equal(t, userID, notifier.userID)
	require.NotNil(t, notifier.lastPet)
	assert.Equal(t, userID, notifier.lastPet.UserID())
}

func TestPetFeedWithoutNotifierStillSucceeds(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	actor := auth.Actor{ID: userID, Role: auth.RoleUser}
	router := petRouter(t, api.PetDeps{Service: &stubPetService{pet: testPet(userID)}}, &actor)

	rec := doRequest(t, router, http.MethodPost, "/pet/actions/feed")

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestPetFeedMapsServiceErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
		want int
	}{
		{name: "not found", err: domainerr.ErrNotFound, want: http.StatusNotFound},
		{name: "limit reached", err: domain.ErrLimitReached, want: http.StatusConflict},
		{name: "invalid action", err: domain.ErrInvalidAction, want: http.StatusBadRequest},
		{name: "unexpected", err: errors.New("redis exploded"), want: http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userID := uuid.New()
			actor := auth.Actor{ID: userID, Role: auth.RoleUser}
			notifier := &spyPetNotifier{}
			router := petRouter(t, api.PetDeps{
				Service: &stubPetService{feedErr: tt.err}, Notifier: notifier,
			}, &actor)

			rec := doRequest(t, router, http.MethodPost, "/pet/actions/feed")

			assert.Equal(t, tt.want, rec.Code)
			assert.Zero(t, notifier.calls)
		})
	}
}

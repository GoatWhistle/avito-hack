package api_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/module/pet/api"
	"github.com/avito-hack/backend/internal/module/pet/app"
	"github.com/avito-hack/backend/internal/module/pet/domain"
	"github.com/avito-hack/backend/internal/shared/auth"
)

var petTestTime = time.Date(2026, time.March, 3, 9, 0, 0, 0, time.UTC)

type stubPetService struct {
	pet             *domain.Pet
	result          app.ActionResult
	progress        app.ProgressView
	feedAvailableAt *time.Time
	checkInApplied  bool
	stateErr        error
	strokeErr       error
	feedErr         error
	checkInErr      error
	progErr         error
	feedCalls       int
	stateCalls      int
}

func (s *stubPetService) StateView(context.Context, uuid.UUID) (app.StateView, error) {
	s.stateCalls++
	if s.stateErr != nil {
		return app.StateView{}, s.stateErr
	}

	return app.StateView{
		Pet: s.pet, FeedAvailableAt: s.feedAvailableAt,
		CheckInApplied: s.checkInApplied,
		Progress:       s.result.Progress, Streak: s.result.Streak,
	}, nil
}

func (s *stubPetService) Stroke(context.Context, uuid.UUID) (*domain.Pet, error) {
	if s.strokeErr != nil {
		return nil, s.strokeErr
	}

	return s.pet, nil
}

func (s *stubPetService) FeedView(context.Context, uuid.UUID) (app.StateView, error) {
	s.feedCalls++
	if s.feedErr != nil {
		return app.StateView{}, s.feedErr
	}

	return app.StateView{Pet: s.pet, FeedAvailableAt: s.feedAvailableAt}, nil
}

func (s *stubPetService) CheckIn(context.Context, uuid.UUID) (app.ActionResult, error) {
	if s.checkInErr != nil {
		return app.ActionResult{}, s.checkInErr
	}

	return s.result, nil
}

func (s *stubPetService) Progress(context.Context, uuid.UUID) (app.ProgressView, error) {
	if s.progErr != nil {
		return app.ProgressView{}, s.progErr
	}

	return s.progress, nil
}

type spyGranter struct {
	calls int
	err   error
}

func (g *spyGranter) GrantEligible(context.Context, uuid.UUID) ([]domain.Reward, error) {
	g.calls++

	return nil, g.err
}

func testPet(userID uuid.UUID) *domain.Pet {
	hatched := petTestTime

	return domain.Restore(domain.RestoreParams{
		ID: uuid.New(), UserID: userID, Name: "Enot", Stage: domain.StageBaby,
		Level: 3, XP: 40, Satiety: 70, Happiness: 80, Energy: 90,
		StreakDays: 4, Freezes: 1, HatchedAt: &hatched,
		LastDecayTime: petTestTime, UpdatedAt: petTestTime,
	})
}

func petRouter(t *testing.T, deps api.PetDeps, actor *auth.Actor) http.Handler {
	t.Helper()

	deps.Authenticate = injectActor(actor)
	router := chi.NewRouter()
	api.NewPetHandlers(deps).RegisterRoutes(router)

	return router
}

func injectActor(actor *auth.Actor) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if actor != nil {
				r = r.WithContext(auth.WithActor(r.Context(), *actor))
			}
			next.ServeHTTP(w, r)
		})
	}
}

func doRequest(t *testing.T, h http.Handler, method, path string) *httptest.ResponseRecorder {
	t.Helper()

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(method, path, http.NoBody))

	return rec
}

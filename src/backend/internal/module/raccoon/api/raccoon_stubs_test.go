package api_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/module/raccoon/api"
	"github.com/avito-hack/backend/internal/module/raccoon/app"
	"github.com/avito-hack/backend/internal/shared/auth"
	"github.com/avito-hack/backend/internal/shared/validate"
)

var raccoonTime = time.Date(2026, time.March, 3, 9, 0, 0, 0, time.UTC)

type stubPetReader struct {
	view app.RaccoonProfileView
	err  error
}

func (s *stubPetReader) Profile(_ context.Context, userID uuid.UUID) (app.RaccoonProfileView, error) {
	if s.err != nil {
		return app.RaccoonProfileView{}, s.err
	}

	view := s.view
	view.UserID = userID

	return view, nil
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
}

func (s *stubActivator) Activate(_ context.Context, _ uuid.UUID, rewardID string) (string, error) {
	s.lastReward = rewardID
	if s.err != nil {
		return "", s.err
	}

	return s.code, nil
}

type fixture struct {
	router    http.Handler
	activator *stubActivator
}

func newFixture(t *testing.T, pets *stubPetReader, badges *stubBadgeReader, actor *auth.Actor) fixture {
	t.Helper()

	activator := &stubActivator{code: "PROMO-1"}
	handlers := api.NewHandlers(api.Deps{
		GetProfile:   app.NewGetRaccoonProfileUseCase(pets, badges),
		ClaimReward:  app.NewClaimRewardUseCase(activator),
		Validator:    validate.New(),
		MaxBodyBytes: 1 << 20,
		Authenticate: func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if actor != nil {
					r = r.WithContext(auth.WithActor(r.Context(), *actor))
				}
				next.ServeHTTP(w, r)
			})
		},
	})

	router := chi.NewRouter()
	handlers.RegisterRoutes(router)

	return fixture{router: router, activator: activator}
}

func (f fixture) do(method, path, body string) *httptest.ResponseRecorder {
	var req *http.Request
	if body == "" {
		req = httptest.NewRequest(method, path, http.NoBody)
	} else {
		req = httptest.NewRequest(method, path, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
	}

	rec := httptest.NewRecorder()
	f.router.ServeHTTP(rec, req)

	return rec
}

func sampleView() app.RaccoonProfileView {
	return app.RaccoonProfileView{
		ID: uuid.New(), Name: "Enot", Level: 4, XP: 120, XPToNextLevel: 200,
		CurrentStreak: 3, Stage: "baby", State: "idle",
	}
}

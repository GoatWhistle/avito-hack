package api_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/user/api"
	"github.com/avito-hack/backend/internal/module/user/app"
	"github.com/avito-hack/backend/internal/module/user/domain"
	"github.com/avito-hack/backend/internal/shared/auth"
	"github.com/avito-hack/backend/internal/shared/validate"
	"github.com/avito-hack/backend/internal/shared/vo"
)

const testPassword = "correct horse battery"

var fixedNow = time.Date(2026, time.April, 1, 10, 0, 0, 0, time.UTC)

type fixedClock struct{}

func (fixedClock) Now() time.Time { return fixedNow }

type passthroughTx struct{}

func (passthroughTx) WithTx(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

type stubTokens struct{}

func (stubTokens) Issue(auth.Actor) (string, time.Time, error) {
	return "signed.jwt.token", fixedNow.Add(time.Hour), nil
}

type memoryUsers struct {
	byID    map[uuid.UUID]*domain.User
	byEmail map[string]*domain.User
}

func newMemoryUsers() *memoryUsers {
	return &memoryUsers{byID: map[uuid.UUID]*domain.User{}, byEmail: map[string]*domain.User{}}
}

func (m *memoryUsers) Save(_ context.Context, user *domain.User) error {
	m.byID[user.ID()] = user
	m.byEmail[user.Email().String()] = user

	return nil
}

func (m *memoryUsers) ByID(_ context.Context, id uuid.UUID) (*domain.User, error) {
	user, ok := m.byID[id]
	if !ok {
		return nil, domain.ErrUserNotFound
	}

	return user, nil
}

func (m *memoryUsers) ByEmail(_ context.Context, email vo.Email) (*domain.User, error) {
	user, ok := m.byEmail[email.String()]
	if !ok {
		return nil, domain.ErrUserNotFound
	}

	return user, nil
}

func (m *memoryUsers) ExistsByEmail(_ context.Context, email vo.Email) (bool, error) {
	_, ok := m.byEmail[email.String()]

	return ok, nil
}

func newRouter(t *testing.T, repo *memoryUsers, actor *auth.Actor) http.Handler {
	t.Helper()

	handlers := api.NewHandlers(api.Deps{
		Register:      app.NewRegisterUserHandler(repo, passthroughTx{}, fixedClock{}, nil, stubTokens{}),
		Login:         app.NewLoginUserHandler(repo, stubTokens{}),
		GetProfile:    app.NewGetProfileHandler(repo),
		UpdateProfile: app.NewUpdateProfileHandler(repo, passthroughTx{}, fixedClock{}),
		Validator:     validate.New(),
		MaxBodyBytes:  1 << 20,
		Authenticate: func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if actor == nil {
					w.WriteHeader(http.StatusUnauthorized)

					return
				}

				next.ServeHTTP(w, r.WithContext(auth.WithActor(r.Context(), *actor)))
			})
		},
	})

	router := chi.NewRouter()
	handlers.RegisterRoutes(router)

	return router
}

func do(t *testing.T, router http.Handler, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	return rec
}

func decodeBody(t *testing.T, rec *httptest.ResponseRecorder) map[string]any {
	t.Helper()

	var body map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))

	return body
}

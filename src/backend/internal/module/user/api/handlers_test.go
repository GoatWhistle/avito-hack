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
	"github.com/stretchr/testify/assert"
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
		Register:      app.NewRegisterUserHandler(repo, passthroughTx{}, fixedClock{}, nil),
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

func TestRegisterEndpoint(t *testing.T) {
	t.Parallel()

	repo := newMemoryUsers()
	router := newRouter(t, repo, nil)

	rec := do(t, router, http.MethodPost, "/auth/register",
		`{"email":"new@example.com","password":"`+testPassword+`","full_name":"Ivan"}`)

	require.Equal(t, http.StatusCreated, rec.Code)

	body := decodeBody(t, rec)
	assert.Equal(t, "new@example.com", body["email"])
	assert.Equal(t, "Ivan", body["full_name"])
	assert.Equal(t, "user", body["role"])
	assert.NotEmpty(t, body["id"])
	assert.Len(t, repo.byID, 1)
}

func TestRegisterEndpointRejectsBadPayloads(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{name: "malformed json", body: `{`, wantStatus: http.StatusBadRequest},
		{name: "unknown field", body: `{"email":"a@b.co","password":"x","full_name":"N","x":1}`, wantStatus: 400},
		{name: "missing email", body: `{"password":"` + testPassword + `","full_name":"Ivan"}`, wantStatus: 400},
		{name: "invalid email", body: `{"email":"nope","password":"` + testPassword + `","full_name":"I"}`, wantStatus: 400},
		{name: "missing full name", body: `{"email":"a@b.co","password":"` + testPassword + `"}`, wantStatus: 400},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := newMemoryUsers()
			rec := do(t, newRouter(t, repo, nil), http.MethodPost, "/auth/register", tc.body)

			assert.Equal(t, tc.wantStatus, rec.Code)
			assert.Empty(t, repo.byID)
		})
	}
}

func TestRegisterEndpointConflictOnDuplicate(t *testing.T) {
	t.Parallel()

	repo := newMemoryUsers()
	router := newRouter(t, repo, nil)
	body := `{"email":"dup@example.com","password":"` + testPassword + `","full_name":"Ivan"}`

	require.Equal(t, http.StatusCreated, do(t, router, http.MethodPost, "/auth/register", body).Code)

	rec := do(t, router, http.MethodPost, "/auth/register", body)

	assert.Equal(t, http.StatusConflict, rec.Code)
}

func TestLoginEndpoint(t *testing.T) {
	t.Parallel()

	repo := newMemoryUsers()
	router := newRouter(t, repo, nil)
	registerBody := `{"email":"login@example.com","password":"` + testPassword + `","full_name":"Ivan"}`
	require.Equal(t, http.StatusCreated, do(t, router, http.MethodPost, "/auth/register", registerBody).Code)

	rec := do(t, router, http.MethodPost, "/auth/login",
		`{"email":"login@example.com","password":"`+testPassword+`"}`)

	require.Equal(t, http.StatusOK, rec.Code)

	body := decodeBody(t, rec)
	assert.Equal(t, "signed.jwt.token", body["token"])
	assert.NotEmpty(t, body["expires_at"])
	require.Contains(t, body, "user")
}

func TestLoginEndpointRejectsWrongCredentials(t *testing.T) {
	t.Parallel()

	repo := newMemoryUsers()
	router := newRouter(t, repo, nil)

	rec := do(t, router, http.MethodPost, "/auth/login",
		`{"email":"ghost@example.com","password":"`+testPassword+`"}`)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

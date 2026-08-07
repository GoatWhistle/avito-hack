package raccoon_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/raccoon"
	"github.com/avito-hack/backend/internal/shared/auth"
	"github.com/avito-hack/backend/internal/shared/validate"
)

func passthroughAuth(next http.Handler) http.Handler { return next }

func TestNewBuildsHandlers(t *testing.T) {
	t.Parallel()

	module := raccoon.New(raccoon.Options{
		Validator:    validate.New(),
		MaxBodyBytes: 1 << 20,
		Authenticate: passthroughAuth,
	})

	require.NotNil(t, module)
	assert.NotNil(t, module.Handlers)
}

func TestNewToleratesMissingDependencies(t *testing.T) {
	t.Parallel()

	module := raccoon.New(raccoon.Options{Authenticate: passthroughAuth})

	require.NotNil(t, module)
	assert.NotNil(t, module.Handlers)
}

func TestNewRegistersRoutes(t *testing.T) {
	t.Parallel()

	module := raccoon.New(raccoon.Options{
		Validator:    validate.New(),
		MaxBodyBytes: 1 << 20,
		Authenticate: passthroughAuth,
	})

	router := chi.NewRouter()
	require.NotPanics(t, func() { module.Handlers.RegisterRoutes(router) })

	paths := []struct {
		method string
		path   string
	}{
		{method: http.MethodGet, path: "/raccoon/profile"},
		{method: http.MethodGet, path: "/badges"},
		{method: http.MethodPost, path: "/rewards/claim"},
	}

	for _, tc := range paths {
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, httptest.NewRequest(tc.method, tc.path, http.NoBody))

		assert.NotEqual(t, http.StatusNotFound, rec.Code, tc.path)
		assert.NotEqual(t, http.StatusMethodNotAllowed, rec.Code, tc.path)
	}
}

func TestModuleRoutesRequireAuthenticatedActor(t *testing.T) {
	t.Parallel()

	module := raccoon.New(raccoon.Options{
		Validator:    validate.New(),
		MaxBodyBytes: 1 << 20,
		Authenticate: passthroughAuth,
	})

	router := chi.NewRouter()
	module.Handlers.RegisterRoutes(router)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/badges", http.NoBody))

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestModuleRejectsMalformedClaimPayload(t *testing.T) {
	t.Parallel()

	module := raccoon.New(raccoon.Options{
		Validator:    validate.New(),
		MaxBodyBytes: 1 << 20,
		Authenticate: func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r.WithContext(auth.WithActor(r.Context(), testActor())))
			})
		},
	})

	router := chi.NewRouter()
	module.Handlers.RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodPost, "/rewards/claim", strings.NewReader(`{"reward_id":""}`))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func testActor() auth.Actor {
	return auth.Actor{ID: uuid.New(), Role: auth.RoleUser}
}

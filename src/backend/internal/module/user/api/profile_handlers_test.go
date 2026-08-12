package api_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/user/app"
	"github.com/avito-hack/backend/internal/module/user/domain"
	"github.com/avito-hack/backend/internal/shared/auth"
)

func TestMeEndpoint(t *testing.T) {
	t.Parallel()

	repo := newMemoryUsers()
	registered := seed(t, repo)
	actor := registered.Actor()
	router := newRouter(t, repo, &actor)

	rec := do(t, router, http.MethodGet, "/users/me", "")

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, registered.DisplayID(), decodeBody(t, rec)["id"])
}

func TestMeEndpointRequiresAuth(t *testing.T) {
	t.Parallel()

	rec := do(t, newRouter(t, newMemoryUsers(), nil), http.MethodGet, "/users/me", "")

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestMeEndpointNotFoundForUnknownActor(t *testing.T) {
	t.Parallel()

	actor := auth.Actor{ID: uuid.New(), Role: auth.RoleUser}

	rec := do(t, newRouter(t, newMemoryUsers(), &actor), http.MethodGet, "/users/me", "")

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestUpdateMeEndpoint(t *testing.T) {
	t.Parallel()

	repo := newMemoryUsers()
	registered := seed(t, repo)
	actor := registered.Actor()
	router := newRouter(t, repo, &actor)

	rec := do(t, router, http.MethodPatch, "/users/me", `{"full_name":"  Renamed  "}`)

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "Renamed", decodeBody(t, rec)["full_name"])
}

func TestUpdateMeEndpointValidation(t *testing.T) {
	t.Parallel()

	repo := newMemoryUsers()
	registered := seed(t, repo)
	actor := registered.Actor()
	router := newRouter(t, repo, &actor)

	assert.Equal(t, http.StatusBadRequest, do(t, router, http.MethodPatch, "/users/me", `{}`).Code)
	assert.Equal(t, http.StatusBadRequest, do(t, router, http.MethodPatch, "/users/me", `{`).Code)
}

func seed(t *testing.T, repo *memoryUsers) *domain.User {
	t.Helper()

	result, err := app.NewRegisterUserHandler(repo, passthroughTx{}, fixedClock{}, nil, stubTokens{}).
		Handle(context.Background(), app.RegisterUserCommand{
			Email: "seed@example.com", Password: testPassword, FullName: "Seeded",
		})
	require.NoError(t, err)

	return result.User
}

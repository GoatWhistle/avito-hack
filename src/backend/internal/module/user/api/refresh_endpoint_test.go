package api_test

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/shared/auth"
)

func TestRefreshEndpointIssuesNewToken(t *testing.T) {
	t.Parallel()

	repo := newMemoryUsers()
	guest := newRouter(t, repo, nil)

	registerRec := do(t, guest, http.MethodPost, "/auth/register",
		`{"email":"refresh@example.com","password":"`+testPassword+`","full_name":"Ivan"}`)
	require.Equal(t, http.StatusCreated, registerRec.Code)

	var stored *auth.Actor
	for _, user := range repo.byID {
		actor := user.Actor()
		stored = &actor
	}
	require.NotNil(t, stored)

	router := newRouter(t, repo, stored)
	rec := do(t, router, http.MethodPost, "/auth/refresh", "")

	require.Equal(t, http.StatusOK, rec.Code)

	body := decodeBody(t, rec)
	assert.Equal(t, "signed.jwt.token", body["token"])
	assert.NotContains(t, body, "expires_at")

	user, ok := body["user"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "refresh@example.com", user["email"])
	assert.Equal(t, stored.ID.String(), user["id"])
}

func TestRefreshEndpointMatchesLoginEnvelope(t *testing.T) {
	t.Parallel()

	repo := newMemoryUsers()
	guest := newRouter(t, repo, nil)

	require.Equal(t, http.StatusCreated, do(t, guest, http.MethodPost, "/auth/register",
		`{"email":"shape@example.com","password":"`+testPassword+`","full_name":"Ivan"}`).Code)

	loginRec := do(t, guest, http.MethodPost, "/auth/login",
		`{"email":"shape@example.com","password":"`+testPassword+`"}`)
	require.Equal(t, http.StatusOK, loginRec.Code)

	var stored *auth.Actor
	for _, user := range repo.byID {
		actor := user.Actor()
		stored = &actor
	}
	require.NotNil(t, stored)

	refreshRec := do(t, newRouter(t, repo, stored), http.MethodPost, "/auth/refresh", "")
	require.Equal(t, http.StatusOK, refreshRec.Code)

	assert.Equal(t, keysOf(decodeBody(t, loginRec)), keysOf(decodeBody(t, refreshRec)))
}

func TestRefreshEndpointRejectsGuest(t *testing.T) {
	t.Parallel()

	router := newRouter(t, newMemoryUsers(), nil)

	rec := do(t, router, http.MethodPost, "/auth/refresh", "")

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestRefreshEndpointRejectsDeletedUser(t *testing.T) {
	t.Parallel()

	repo := newMemoryUsers()
	actor := auth.Actor{ID: uuid.New(), Role: auth.RoleUser}

	rec := do(t, newRouter(t, repo, &actor), http.MethodPost, "/auth/refresh", "")

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

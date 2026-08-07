package api_test

import (
	"net/http"
	"sort"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/shared/auth"
)

func TestRegisterEndpoint(t *testing.T) {
	t.Parallel()

	repo := newMemoryUsers()
	router := newRouter(t, repo, nil)

	rec := do(t, router, http.MethodPost, "/auth/register",
		`{"email":"new@example.com","password":"`+testPassword+`","full_name":"Ivan"}`)

	require.Equal(t, http.StatusCreated, rec.Code)

	body := decodeBody(t, rec)
	assert.Equal(t, "signed.jwt.token", body["token"])

	user, ok := body["user"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "new@example.com", user["email"])
	assert.Equal(t, "Ivan", user["full_name"])
	assert.Equal(t, "user", user["role"])
	assert.NotEmpty(t, user["id"])
	assert.Len(t, repo.byID, 1)
}

func TestRegisterEndpointReturnsSameEnvelopeAsLogin(t *testing.T) {
	t.Parallel()

	repo := newMemoryUsers()
	router := newRouter(t, repo, nil)

	registerRec := do(t, router, http.MethodPost, "/auth/register",
		`{"email":"envelope@example.com","password":"`+testPassword+`","full_name":"Ivan"}`)
	require.Equal(t, http.StatusCreated, registerRec.Code)

	loginRec := do(t, router, http.MethodPost, "/auth/login",
		`{"email":"envelope@example.com","password":"`+testPassword+`"}`)
	require.Equal(t, http.StatusOK, loginRec.Code)

	registerBody := decodeBody(t, registerRec)
	loginBody := decodeBody(t, loginRec)

	assert.Equal(t, keysOf(loginBody), keysOf(registerBody))
	assert.NotEmpty(t, registerBody["token"])
	assert.NotContains(t, registerBody, "expires_at")
}

func TestRegisterTokenAuthenticatesFollowUpRequest(t *testing.T) {
	t.Parallel()

	repo := newMemoryUsers()
	rec := do(t, newRouter(t, repo, nil), http.MethodPost, "/auth/register",
		`{"email":"session@example.com","password":"`+testPassword+`","full_name":"Ivan"}`)
	require.Equal(t, http.StatusCreated, rec.Code)

	body := decodeBody(t, rec)
	require.NotEmpty(t, body["token"])

	user, ok := body["user"].(map[string]any)
	require.True(t, ok)

	id, err := uuid.Parse(user["id"].(string))
	require.NoError(t, err)

	actor := auth.Actor{ID: id, Role: auth.RoleUser}
	meRec := do(t, newRouter(t, repo, &actor), http.MethodGet, "/users/me", "")

	require.Equal(t, http.StatusOK, meRec.Code)
	assert.Equal(t, "session@example.com", decodeBody(t, meRec)["email"])
}

func keysOf(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	return keys
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
	assert.NotContains(t, body, "expires_at")
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

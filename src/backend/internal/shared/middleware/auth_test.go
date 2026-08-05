package middleware_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/shared/auth"
	"github.com/avito-hack/backend/internal/shared/middleware"
)

type stubParser struct {
	actor auth.Actor
	err   error
}

func (s stubParser) Parse(string) (auth.Actor, error) {
	if s.err != nil {
		return auth.Actor{}, s.err
	}

	return s.actor, nil
}

func probe(seen *auth.Actor, called *bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*called = true
		if actor, err := auth.ActorFrom(r.Context()); err == nil {
			*seen = actor
		}
		w.WriteHeader(http.StatusOK)
	})
}

func TestAuthenticate(t *testing.T) {
	t.Parallel()

	actor := auth.Actor{ID: uuid.New(), Role: auth.RoleUser}

	tests := []struct {
		name       string
		header     string
		parser     stubParser
		wantStatus int
		wantCalled bool
	}{
		{
			name:       "valid bearer token",
			header:     "Bearer good.token",
			parser:     stubParser{actor: actor},
			wantStatus: http.StatusOK,
			wantCalled: true,
		},
		{
			name:       "missing header",
			parser:     stubParser{actor: actor},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "wrong scheme",
			header:     "Basic dXNlcjpwYXNz",
			parser:     stubParser{actor: actor},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "empty token after prefix",
			header:     "Bearer    ",
			parser:     stubParser{actor: actor},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "lowercase scheme is rejected",
			header:     "bearer good.token",
			parser:     stubParser{actor: actor},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "parser rejects token",
			header:     "Bearer forged.token",
			parser:     stubParser{err: errors.New("invalid token")},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var (
				called bool
				seen   auth.Actor
			)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tc.header != "" {
				req.Header.Set("Authorization", tc.header)
			}

			rec := httptest.NewRecorder()
			middleware.Authenticate(tc.parser)(probe(&seen, &called)).ServeHTTP(rec, req)

			assert.Equal(t, tc.wantStatus, rec.Code)
			assert.Equal(t, tc.wantCalled, called)

			if tc.wantCalled {
				assert.Equal(t, actor, seen)
			}
		})
	}
}

func TestOptionalAuthenticate(t *testing.T) {
	t.Parallel()

	actor := auth.Actor{ID: uuid.New(), Role: auth.RoleModerator}

	tests := []struct {
		name      string
		header    string
		parser    stubParser
		wantActor auth.Actor
	}{
		{name: "no header still passes through", parser: stubParser{actor: actor}},
		{name: "invalid token still passes through", header: "Bearer x", parser: stubParser{err: errors.New("bad")}},
		{name: "valid token attaches actor", header: "Bearer ok", parser: stubParser{actor: actor}, wantActor: actor},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var (
				called bool
				seen   auth.Actor
			)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tc.header != "" {
				req.Header.Set("Authorization", tc.header)
			}

			rec := httptest.NewRecorder()
			middleware.OptionalAuthenticate(tc.parser)(probe(&seen, &called)).ServeHTTP(rec, req)

			require.True(t, called)
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, tc.wantActor, seen)
		})
	}
}

func TestRequireRole(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		actor      *auth.Actor
		allowed    []auth.Role
		wantStatus int
	}{
		{
			name:       "no actor in context",
			allowed:    []auth.Role{auth.RoleUser},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "matching role",
			actor:      &auth.Actor{ID: uuid.New(), Role: auth.RoleModerator},
			allowed:    []auth.Role{auth.RoleModerator, auth.RoleAdmin},
			wantStatus: http.StatusOK,
		},
		{
			name:       "insufficient role",
			actor:      &auth.Actor{ID: uuid.New(), Role: auth.RoleUser},
			allowed:    []auth.Role{auth.RoleAdmin},
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "empty allow list forbids everything",
			actor:      &auth.Actor{ID: uuid.New(), Role: auth.RoleAdmin},
			wantStatus: http.StatusForbidden,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var (
				called bool
				seen   auth.Actor
			)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tc.actor != nil {
				req = req.WithContext(auth.WithActor(req.Context(), *tc.actor))
			}

			rec := httptest.NewRecorder()
			middleware.RequireRole(tc.allowed...)(probe(&seen, &called)).ServeHTTP(rec, req)

			assert.Equal(t, tc.wantStatus, rec.Code)
			assert.Equal(t, tc.wantStatus == http.StatusOK, called)
		})
	}
}

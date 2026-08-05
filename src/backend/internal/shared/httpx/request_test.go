package httpx_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/shared/domainerr"
	"github.com/avito-hack/backend/internal/shared/httpx"
)

type payload struct {
	Title string `json:"title"`
}

func TestDecodeJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		body    string
		max     int64
		wantErr bool
		want    string
	}{
		{name: "valid body", body: `{"title":"chair"}`, max: 1024, want: "chair"},
		{name: "unknown field is rejected", body: `{"title":"x","extra":1}`, max: 1024, wantErr: true},
		{name: "malformed json", body: `{"title":`, max: 1024, wantErr: true},
		{name: "empty body", body: ``, max: 1024, wantErr: true},
		{name: "body over limit", body: `{"title":"` + strings.Repeat("a", 200) + `"}`, max: 32, wantErr: true},
		{name: "wrong type", body: `{"title":42}`, max: 1024, wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tc.body))

			var dst payload
			err := httpx.DecodeJSON(rec, req, tc.max, &dst)

			if tc.wantErr {
				require.ErrorIs(t, err, httpx.ErrMalformedBody)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.want, dst.Title)
		})
	}
}

func TestUUIDParam(t *testing.T) {
	t.Parallel()

	id := uuid.New()

	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		parsed, err := httpx.UUIDParam(requestWithParam(t, "id", id.String()), "id")

		require.NoError(t, err)
		assert.Equal(t, id, parsed)
	})

	t.Run("invalid", func(t *testing.T) {
		t.Parallel()

		_, err := httpx.UUIDParam(requestWithParam(t, "id", "abc"), "id")

		var invalid *domainerr.InvalidError
		require.ErrorAs(t, err, &invalid)
		assert.Equal(t, "id", invalid.Field)
	})

	t.Run("missing", func(t *testing.T) {
		t.Parallel()

		_, err := httpx.UUIDParam(httptest.NewRequest(http.MethodGet, "/", nil), "id")

		require.Error(t, err)
	})
}

func TestOptionalUUIDQuery(t *testing.T) {
	t.Parallel()

	id := uuid.New()

	t.Run("present", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodGet, "/?owner_id="+id.String(), nil)

		parsed, ok, err := httpx.OptionalUUIDQuery(req, "owner_id")

		require.NoError(t, err)
		assert.True(t, ok)
		assert.Equal(t, id, parsed)
	})

	t.Run("absent", func(t *testing.T) {
		t.Parallel()

		parsed, ok, err := httpx.OptionalUUIDQuery(httptest.NewRequest(http.MethodGet, "/", nil), "owner_id")

		require.NoError(t, err)
		assert.False(t, ok)
		assert.Equal(t, uuid.Nil, parsed)
	})

	t.Run("malformed", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodGet, "/?owner_id=zzz", nil)

		_, ok, err := httpx.OptionalUUIDQuery(req, "owner_id")

		require.Error(t, err)
		assert.False(t, ok)
	})
}

func TestIntQuery(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		url     string
		want    int
		wantErr bool
	}{
		{name: "absent uses fallback", url: "/", want: 7},
		{name: "empty uses fallback", url: "/?limit=", want: 7},
		{name: "parses positive", url: "/?limit=25", want: 25},
		{name: "parses negative", url: "/?limit=-5", want: -5},
		{name: "parses zero", url: "/?limit=0", want: 0},
		{name: "rejects non numeric", url: "/?limit=abc", wantErr: true},
		{name: "rejects float", url: "/?limit=1.5", wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			value, err := httpx.IntQuery(httptest.NewRequest(http.MethodGet, tc.url, nil), "limit", 7)

			if tc.wantErr {
				var invalid *domainerr.InvalidError
				require.ErrorAs(t, err, &invalid)
				assert.Equal(t, "limit", invalid.Field)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.want, value)
		})
	}
}

func requestWithParam(t *testing.T, name, value string) *http.Request {
	t.Helper()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add(name, value)

	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
}

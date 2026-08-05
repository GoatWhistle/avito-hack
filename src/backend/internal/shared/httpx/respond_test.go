package httpx_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/shared/httpx"
	"github.com/avito-hack/backend/internal/shared/pagination"
)

func TestJSONResponders(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		respond    func(w http.ResponseWriter)
		wantStatus int
		wantBody   string
	}{
		{
			name:       "ok",
			respond:    func(w http.ResponseWriter) { httpx.OK(w, payload{Title: "chair"}) },
			wantStatus: http.StatusOK,
			wantBody:   `{"title":"chair"}`,
		},
		{
			name:       "created",
			respond:    func(w http.ResponseWriter) { httpx.Created(w, payload{Title: "table"}) },
			wantStatus: http.StatusCreated,
			wantBody:   `{"title":"table"}`,
		},
		{
			name:       "no content",
			respond:    httpx.NoContent,
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "nil payload writes no body",
			respond:    func(w http.ResponseWriter) { httpx.JSON(w, http.StatusAccepted, nil) },
			wantStatus: http.StatusAccepted,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			tc.respond(rec)

			require.Equal(t, tc.wantStatus, rec.Code)

			if tc.wantBody == "" {
				assert.Empty(t, rec.Body.String())

				return
			}

			assert.JSONEq(t, tc.wantBody, rec.Body.String())
			assert.Equal(t, "application/json; charset=utf-8", rec.Header().Get("Content-Type"))
		})
	}
}

func TestNewListResponse(t *testing.T) {
	t.Parallel()

	t.Run("nil items become empty array", func(t *testing.T) {
		t.Parallel()

		encoded, err := json.Marshal(httpx.NewListResponse[payload](nil, ""))

		require.NoError(t, err)
		assert.JSONEq(t, `{"items":[]}`, string(encoded))
	})

	t.Run("next cursor is included", func(t *testing.T) {
		t.Parallel()

		encoded, err := json.Marshal(httpx.NewListResponse([]payload{{Title: "a"}}, "cur"))

		require.NoError(t, err)
		assert.JSONEq(t, `{"items":[{"title":"a"}],"next_cursor":"cur"}`, string(encoded))
	})
}

func TestPageFromRequest(t *testing.T) {
	t.Parallel()

	cursor := pagination.Cursor{CreatedAt: time.Now().UTC().Truncate(time.Second), ID: uuid.New()}

	t.Run("defaults", func(t *testing.T) {
		t.Parallel()

		page, err := httpx.PageFromRequest(httptest.NewRequest(http.MethodGet, "/", nil))

		require.NoError(t, err)
		assert.Equal(t, pagination.DefaultLimit, page.Limit)
		assert.True(t, page.Cursor.IsZero())
	})

	t.Run("parses cursor and clamps limit", func(t *testing.T) {
		t.Parallel()

		url := "/?limit=1000&cursor=" + cursor.Encode()

		page, err := httpx.PageFromRequest(httptest.NewRequest(http.MethodGet, url, nil))

		require.NoError(t, err)
		assert.Equal(t, pagination.MaxLimit, page.Limit)
		assert.Equal(t, cursor.ID, page.Cursor.ID)
	})

	t.Run("rejects malformed limit", func(t *testing.T) {
		t.Parallel()

		_, err := httpx.PageFromRequest(httptest.NewRequest(http.MethodGet, "/?limit=many", nil))

		require.Error(t, err)
	})

	t.Run("rejects malformed cursor", func(t *testing.T) {
		t.Parallel()

		_, err := httpx.PageFromRequest(httptest.NewRequest(http.MethodGet, "/?cursor=%21%21%21", nil))

		require.Error(t, err)
	})
}

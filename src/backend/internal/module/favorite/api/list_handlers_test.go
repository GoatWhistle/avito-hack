package api_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/favorite/app"
)

func TestListFavoritesEndpoint(t *testing.T) {
	t.Parallel()

	current := actor()
	rows := []app.FavoriteItem{
		{
			ItemID: uuid.New(), OwnerID: uuid.New(), Title: "Chair",
			PriceKopeks: 1500, Status: "published", PhotoURL: "/media/a.jpg", CreatedAt: fixedTime,
		},
	}

	router := newRouter(t, routerOptions{actor: &current, read: stubRead{rows: rows}})

	rec := do(t, router, http.MethodGet, "/favorites?limit=10")
	require.Equal(t, http.StatusOK, rec.Code)

	var body struct {
		Items []struct {
			ItemID      string `json:"item_id"`
			Title       string `json:"title"`
			PriceKopeks int64  `json:"price"`
			Status      string `json:"status"`
			PhotoURL    string `json:"photo_url"`
		} `json:"items"`
		NextCursor string `json:"next_cursor"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))

	require.Len(t, body.Items, 1)
	assert.Equal(t, rows[0].ItemID.String(), body.Items[0].ItemID)
	assert.Equal(t, "Chair", body.Items[0].Title)
	assert.Equal(t, int64(1500), body.Items[0].PriceKopeks)
	assert.Equal(t, "/media/a.jpg", body.Items[0].PhotoURL)
	assert.Empty(t, body.NextCursor)
}

func TestListFavoritesEndpointErrors(t *testing.T) {
	t.Parallel()

	current := actor()

	t.Run("requires auth", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, http.StatusUnauthorized, do(t, newRouter(t, routerOptions{}), http.MethodGet, "/favorites").Code)
	})

	t.Run("rejects malformed cursor", func(t *testing.T) {
		t.Parallel()

		router := newRouter(t, routerOptions{actor: &current})

		assert.Equal(t, http.StatusBadRequest, do(t, router, http.MethodGet, "/favorites?cursor=%21%21").Code)
	})

	t.Run("rejects malformed limit", func(t *testing.T) {
		t.Parallel()

		router := newRouter(t, routerOptions{actor: &current})

		assert.Equal(t, http.StatusBadRequest, do(t, router, http.MethodGet, "/favorites?limit=lots").Code)
	})

	t.Run("read model failure", func(t *testing.T) {
		t.Parallel()

		router := newRouter(t, routerOptions{actor: &current, read: stubRead{err: errors.New("db down")}})

		assert.Equal(t, http.StatusInternalServerError, do(t, router, http.MethodGet, "/favorites").Code)
	})
}

func TestListFavoritesEndpointEmitsEmptyArray(t *testing.T) {
	t.Parallel()

	current := actor()
	router := newRouter(t, routerOptions{actor: &current})

	rec := do(t, router, http.MethodGet, "/favorites")

	require.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"items":[]}`, rec.Body.String())
}

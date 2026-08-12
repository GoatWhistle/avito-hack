package api_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/item/domain"
)

func body(t *testing.T, rec *http.Response) map[string]any {
	t.Helper()

	var decoded map[string]any
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&decoded))

	return decoded
}

func TestCreateItemEndpoint(t *testing.T) {
	t.Parallel()

	actor := userActor()
	f := newFixture(t, &actor)

	rec := f.do(t, http.MethodPost, "/items/",
		`{"title":"Mountain bike","description":"as new","price":150000,"attributes":{"color":"red"}}`)

	require.Equal(t, http.StatusCreated, rec.Code)

	decoded := body(t, rec.Result())
	assert.Equal(t, "Mountain bike", decoded["title"])
	assert.Equal(t, "draft", decoded["status"])
	assert.InDelta(t, float64(150000), decoded["price"], 0.1)
	assert.NotContains(t, decoded["owner_id"], actor.ID.String())
	assert.Len(t, f.items.items, 1)
}

func TestCreateItemEndpointValidation(t *testing.T) {
	t.Parallel()

	actor := userActor()

	tests := []struct {
		name string
		body string
	}{
		{name: "malformed json", body: `{`},
		{name: "missing title", body: `{"price":100}`},
		{name: "short title", body: `{"title":"ab","price":100}`},
		{name: "negative price", body: `{"title":"Bicycle","price":-5}`},
		{name: "unknown field", body: `{"title":"Bicycle","price":100,"color":"red"}`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			f := newFixture(t, &actor)

			assert.Equal(t, http.StatusBadRequest, f.do(t, http.MethodPost, "/items/", tc.body).Code)
			assert.Empty(t, f.items.items)
		})
	}
}

func TestCreateItemEndpointRequiresAuth(t *testing.T) {
	t.Parallel()

	f := newFixture(t, nil)

	assert.Equal(t, http.StatusUnauthorized,
		f.do(t, http.MethodPost, "/items/", `{"title":"Bicycle","price":100}`).Code)
}

func TestUpdateItemEndpoint(t *testing.T) {
	t.Parallel()

	actor := userActor()
	f := newFixture(t, &actor)
	item := f.seedItem(t, actor.ID, domain.StatusDraft)

	rec := f.do(t, http.MethodPatch, "/items/"+item.DisplayID(), `{"title":"Updated bike","price":99}`)

	require.Equal(t, http.StatusOK, rec.Code)

	decoded := body(t, rec.Result())
	assert.Equal(t, "Updated bike", decoded["title"])
	assert.InDelta(t, float64(99), decoded["price"], 0.1)
}

func TestUpdateItemEndpointErrors(t *testing.T) {
	t.Parallel()

	owner := userActor()
	stranger := userActor()

	t.Run("unknown id", func(t *testing.T) {
		t.Parallel()

		f := newFixture(t, &owner)

		assert.Equal(t, http.StatusNotFound, f.do(t, http.MethodPatch, "/items/abc", `{"title":"New"}`).Code)
	})

	t.Run("unknown item", func(t *testing.T) {
		t.Parallel()

		f := newFixture(t, &owner)

		assert.Equal(t, http.StatusNotFound,
			f.do(t, http.MethodPatch, "/items/"+uuid.NewString(), `{"title":"New"}`).Code)
	})

	t.Run("not the owner", func(t *testing.T) {
		t.Parallel()

		f := newFixture(t, &stranger)
		item := f.seedItem(t, owner.ID, domain.StatusDraft)

		assert.Equal(t, http.StatusForbidden,
			f.do(t, http.MethodPatch, "/items/"+item.DisplayID(), `{"title":"Hijack"}`).Code)
	})

	t.Run("invalid payload", func(t *testing.T) {
		t.Parallel()

		f := newFixture(t, &owner)
		item := f.seedItem(t, owner.ID, domain.StatusDraft)

		assert.Equal(t, http.StatusBadRequest,
			f.do(t, http.MethodPatch, "/items/"+item.DisplayID(), `{"title":"ab"}`).Code)
	})

	t.Run("requires auth", func(t *testing.T) {
		t.Parallel()

		f := newFixture(t, nil)

		assert.Equal(t, http.StatusUnauthorized,
			f.do(t, http.MethodPatch, "/items/"+uuid.NewString(), `{"title":"New"}`).Code)
	})
}

func TestChangeStatusEndpoint(t *testing.T) {
	t.Parallel()

	actor := userActor()

	tests := []struct {
		name     string
		from     domain.Status
		action   string
		wantCode int
		wantTo   string
	}{
		{name: "submit draft", from: domain.StatusDraft, action: "submit", wantCode: 200, wantTo: "moderation"},
		{name: "sell published", from: domain.StatusPublished, action: "sell", wantCode: 200, wantTo: "sold"},
		{name: "archive published", from: domain.StatusPublished, action: "archive", wantCode: 200, wantTo: "archived"},
		{name: "restore archived", from: domain.StatusArchived, action: "restore", wantCode: 200, wantTo: "draft"},
		{name: "sell draft is a conflict", from: domain.StatusDraft, action: "sell", wantCode: 409},
		{name: "submit sold is a conflict", from: domain.StatusSold, action: "submit", wantCode: 409},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			f := newFixture(t, &actor)
			item := f.seedItem(t, actor.ID, tc.from)

			rec := f.do(t, http.MethodPost, "/items/"+item.DisplayID()+"/status", `{"action":"`+tc.action+`"}`)

			require.Equal(t, tc.wantCode, rec.Code)

			if tc.wantTo != "" {
				assert.Equal(t, tc.wantTo, body(t, rec.Result())["status"])
			}
		})
	}
}

func TestChangeStatusEndpointRejectsUnknownAction(t *testing.T) {
	t.Parallel()

	actor := userActor()
	f := newFixture(t, &actor)
	item := f.seedItem(t, actor.ID, domain.StatusDraft)

	rec := f.do(t, http.MethodPost, "/items/"+item.DisplayID()+"/status", `{"action":"delete"}`)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

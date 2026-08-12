package api_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/item/app"
	"github.com/avito-hack/backend/internal/module/item/domain"
)

type listBody struct {
	Items []struct {
		ID        string `json:"id"`
		OwnerName string `json:"owner_name"`
		Title     string `json:"title"`
		Status    string `json:"status"`
	} `json:"items"`
	NextCursor string `json:"next_cursor"`
}

func decodeList(t *testing.T, rec *http.Response) listBody {
	t.Helper()

	var decoded listBody
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&decoded))

	return decoded
}

func TestGetItemEndpointPublicAccess(t *testing.T) {
	t.Parallel()

	owner := userActor()
	f := newFixture(t, nil)
	item := f.seedItem(t, owner.ID, domain.StatusPublished)

	rec := f.do(t, http.MethodGet, "/items/"+item.DisplayID(), "")

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, item.DisplayID(), body(t, rec.Result())["id"])
}

func TestGetItemEndpointHidesDraftsFromAnonymous(t *testing.T) {
	t.Parallel()

	owner := userActor()
	f := newFixture(t, nil)
	item := f.seedItem(t, owner.ID, domain.StatusDraft)

	assert.Equal(t, http.StatusNotFound, f.do(t, http.MethodGet, "/items/"+item.DisplayID(), "").Code)
}

func TestGetItemEndpointOwnerSeesDraft(t *testing.T) {
	t.Parallel()

	owner := userActor()
	f := newFixture(t, &owner)
	item := f.seedItem(t, owner.ID, domain.StatusDraft)

	assert.Equal(t, http.StatusOK, f.do(t, http.MethodGet, "/items/"+item.DisplayID(), "").Code)
}

func TestGetItemEndpointErrors(t *testing.T) {
	t.Parallel()

	f := newFixture(t, nil)

	assert.Equal(t, http.StatusNotFound, f.do(t, http.MethodGet, "/items/not-a-uuid", "").Code)
	assert.Equal(t, http.StatusNotFound, f.do(t, http.MethodGet, "/items/"+uuid.NewString(), "").Code)
}

func TestListItemsEndpoint(t *testing.T) {
	t.Parallel()

	f := newFixture(t, nil)
	ownerID := uuid.New()
	f.read.rows = []app.ListItem{
		{ID: uuid.New(), DisplayID: domain.NewDisplayID(), OwnerID: ownerID, Title: "Bike",
			Status: domain.StatusPublished, CreatedAt: fixedTime},
		{ID: uuid.New(), DisplayID: domain.NewDisplayID(), OwnerID: uuid.New(), Title: "Chair",
			Status: domain.StatusDraft, CreatedAt: fixedTime},
	}

	rec := f.do(t, http.MethodGet, "/items/?limit=10", "")

	require.Equal(t, http.StatusOK, rec.Code)

	decoded := decodeList(t, rec.Result())
	require.Len(t, decoded.Items, 1)
	assert.Equal(t, "Bike", decoded.Items[0].Title)
	assert.Empty(t, decoded.NextCursor)
}

func TestListItemsEndpointFilters(t *testing.T) {
	t.Parallel()

	f := newFixture(t, nil)
	ownerID := uuid.New()
	f.read.rows = []app.ListItem{
		{ID: uuid.New(), DisplayID: domain.NewDisplayID(), OwnerID: ownerID, Title: "Bike",
			Status: domain.StatusPublished, CreatedAt: fixedTime},
		{ID: uuid.New(), DisplayID: domain.NewDisplayID(), OwnerID: uuid.New(), Title: "Chair",
			Status: domain.StatusDraft, CreatedAt: fixedTime},
	}

	byStatus := decodeList(t, f.do(t, http.MethodGet, "/items/?status=published", "").Result())
	require.Len(t, byStatus.Items, 1)
	assert.Equal(t, "Bike", byStatus.Items[0].Title)

	byOwner := decodeList(t, f.do(t, http.MethodGet, "/items/?owner_id="+ownerID.String(), "").Result())
	require.Len(t, byOwner.Items, 1)
	assert.Equal(t, "Bike", byOwner.Items[0].Title)

	assert.Empty(t, decodeList(t, f.do(t, http.MethodGet, "/items/?status=draft", "").Result()).Items)
}

func TestListItemsEndpointRejectsBadQuery(t *testing.T) {
	t.Parallel()

	f := newFixture(t, nil)

	tests := []struct {
		name string
		path string
	}{
		{name: "unknown status", path: "/items/?status=deleted"},
		{name: "malformed owner", path: "/items/?owner_id=nope"},
		{name: "malformed limit", path: "/items/?limit=many"},
		{name: "malformed cursor", path: "/items/?cursor=%21%21%21"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, http.StatusBadRequest, f.do(t, http.MethodGet, tc.path, "").Code)
		})
	}
}

func TestListMineEndpoint(t *testing.T) {
	t.Parallel()

	actor := userActor()
	f := newFixture(t, &actor)
	f.read.rows = []app.ListItem{
		{ID: uuid.New(), DisplayID: domain.NewDisplayID(), OwnerID: actor.ID, Title: "Mine",
			Status: domain.StatusDraft, CreatedAt: fixedTime},
		{ID: uuid.New(), DisplayID: domain.NewDisplayID(), OwnerID: uuid.New(), Title: "Theirs",
			Status: domain.StatusDraft, CreatedAt: fixedTime},
	}

	rec := f.do(t, http.MethodGet, "/items/mine", "")

	require.Equal(t, http.StatusOK, rec.Code)

	decoded := decodeList(t, rec.Result())
	require.Len(t, decoded.Items, 1)
	assert.Equal(t, "Mine", decoded.Items[0].Title)
}

func TestListMineEndpointErrors(t *testing.T) {
	t.Parallel()

	actor := userActor()

	assert.Equal(t, http.StatusUnauthorized, newFixture(t, nil).do(t, http.MethodGet, "/items/mine", "").Code)

	f := newFixture(t, &actor)
	assert.Equal(t, http.StatusBadRequest, f.do(t, http.MethodGet, "/items/mine?status=deleted", "").Code)
	assert.Equal(t, http.StatusBadRequest, f.do(t, http.MethodGet, "/items/mine?limit=many", "").Code)
}

func TestListPhotosEndpoint(t *testing.T) {
	t.Parallel()

	owner := userActor()
	f := newFixture(t, &owner)
	item := f.seedItem(t, owner.ID, domain.StatusPublished)

	photo, err := domain.NewPhoto(domain.NewPhotoParams{
		ItemID: item.ID(), URL: "/uploads/a.jpg", Position: 0, Now: fixedTime,
	})
	require.NoError(t, err)
	require.NoError(t, f.photos.Add(t.Context(), photo))

	rec := f.do(t, http.MethodGet, "/items/"+item.DisplayID()+"/photos", "")

	require.Equal(t, http.StatusOK, rec.Code)

	var decoded []map[string]any
	require.NoError(t, json.NewDecoder(rec.Result().Body).Decode(&decoded))
	require.Len(t, decoded, 1)
	assert.Equal(t, "/uploads/a.jpg", decoded[0]["url"])
}

func TestListPhotosEndpointRejectsMalformedID(t *testing.T) {
	t.Parallel()

	f := newFixture(t, nil)

	assert.Equal(t, http.StatusNotFound, f.do(t, http.MethodGet, "/items/nope/photos", "").Code)
}

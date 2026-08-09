package api_test

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/item/app"
	"github.com/avito-hack/backend/internal/module/item/domain"
)

func TestListItemsEndpointHidesForeignNonPublicStatusesFromAnonymous(t *testing.T) {
	t.Parallel()

	f := newFixture(t, nil)
	ownerID := uuid.New()
	f.read.rows = []app.ListItem{
		{ID: uuid.New(), OwnerID: ownerID, Title: "Draft", Status: domain.StatusDraft, CreatedAt: fixedTime},
		{ID: uuid.New(), OwnerID: ownerID, Title: "Moderation", Status: domain.StatusModeration, CreatedAt: fixedTime},
		{ID: uuid.New(), OwnerID: ownerID, Title: "Archived", Status: domain.StatusArchived, CreatedAt: fixedTime},
		{ID: uuid.New(), OwnerID: ownerID, Title: "Sold", Status: domain.StatusSold, CreatedAt: fixedTime},
	}

	decoded := decodeList(t, f.do(t, http.MethodGet, "/items/?limit=10", "").Result())

	require.Len(t, decoded.Items, 1)
	assert.Equal(t, "Sold", decoded.Items[0].Title)

	for _, status := range []string{"draft", "moderation", "archived"} {
		byStatus := decodeList(t, f.do(t, http.MethodGet, "/items/?status="+status, "").Result())
		assert.Emptyf(t, byStatus.Items, "anonymous must not see %s items", status)
	}
}

func TestListItemsEndpointOwnerSeesOwnDrafts(t *testing.T) {
	t.Parallel()

	owner := userActor()
	f := newFixture(t, &owner)
	f.read.rows = []app.ListItem{
		{ID: uuid.New(), OwnerID: owner.ID, Title: "My draft", Status: domain.StatusDraft, CreatedAt: fixedTime},
		{ID: uuid.New(), OwnerID: uuid.New(), Title: "Their draft", Status: domain.StatusDraft, CreatedAt: fixedTime},
	}

	decoded := decodeList(t, f.do(t, http.MethodGet, "/items/?limit=10", "").Result())

	require.Len(t, decoded.Items, 1)
	assert.Equal(t, "My draft", decoded.Items[0].Title)
}

func TestListItemsEndpointOwnerFilterCannotRevealForeignDrafts(t *testing.T) {
	t.Parallel()

	victimID := uuid.New()
	viewer := userActor()
	f := newFixture(t, &viewer)
	f.read.rows = []app.ListItem{
		{ID: uuid.New(), OwnerID: victimID, Title: "Secret", Status: domain.StatusDraft, CreatedAt: fixedTime},
		{ID: uuid.New(), OwnerID: victimID, Title: "Public", Status: domain.StatusPublished, CreatedAt: fixedTime},
	}

	path := "/items/?owner_id=" + victimID.String() + "&status=draft"
	assert.Empty(t, decodeList(t, f.do(t, http.MethodGet, path, "").Result()).Items)

	all := decodeList(t, f.do(t, http.MethodGet, "/items/?owner_id="+victimID.String(), "").Result())
	require.Len(t, all.Items, 1)
	assert.Equal(t, "Public", all.Items[0].Title)
}

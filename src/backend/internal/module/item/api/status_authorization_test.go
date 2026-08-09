package api_test

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/avito-hack/backend/internal/module/item/domain"
	"github.com/avito-hack/backend/internal/shared/auth"
)

func TestChangeStatusEndpointAuthorization(t *testing.T) {
	t.Parallel()

	owner := userActor()
	moderator := auth.Actor{ID: uuid.New(), Role: auth.RoleModerator}
	stranger := userActor()

	t.Run("moderator may archive", func(t *testing.T) {
		t.Parallel()

		f := newFixture(t, &moderator)
		item := f.seedItem(t, owner.ID, domain.StatusPublished)

		assert.Equal(t, http.StatusOK,
			f.do(t, http.MethodPost, "/items/"+item.ID().String()+"/status", `{"action":"archive"}`).Code)
	})

	t.Run("moderator may not publish", func(t *testing.T) {
		t.Parallel()

		f := newFixture(t, &moderator)
		item := f.seedItem(t, owner.ID, domain.StatusModeration)

		assert.Equal(t, http.StatusForbidden,
			f.do(t, http.MethodPost, "/items/"+item.ID().String()+"/status", `{"action":"publish"}`).Code)
	})

	t.Run("moderator may not restore", func(t *testing.T) {
		t.Parallel()

		f := newFixture(t, &moderator)
		item := f.seedItem(t, owner.ID, domain.StatusArchived)

		assert.Equal(t, http.StatusForbidden,
			f.do(t, http.MethodPost, "/items/"+item.ID().String()+"/status", `{"action":"restore"}`).Code)
	})

	t.Run("moderator may not sell", func(t *testing.T) {
		t.Parallel()

		f := newFixture(t, &moderator)
		item := f.seedItem(t, owner.ID, domain.StatusPublished)

		assert.Equal(t, http.StatusForbidden,
			f.do(t, http.MethodPost, "/items/"+item.ID().String()+"/status", `{"action":"sell"}`).Code)
	})

	t.Run("stranger is forbidden", func(t *testing.T) {
		t.Parallel()

		f := newFixture(t, &stranger)
		item := f.seedItem(t, owner.ID, domain.StatusDraft)

		assert.Equal(t, http.StatusForbidden,
			f.do(t, http.MethodPost, "/items/"+item.ID().String()+"/status", `{"action":"publish"}`).Code)
	})
}

package domain_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/item/domain"
	"github.com/avito-hack/backend/internal/shared/vo"
)

const defaultTitle = "MacBook Pro"

var fixedTime = time.Date(2026, time.January, 1, 12, 0, 0, 0, time.UTC)

func newItem(t *testing.T) *domain.Item {
	t.Helper()

	item, err := domain.NewItem(domain.NewItemParams{
		OwnerID:     uuid.New(),
		Title:       defaultTitle,
		Description: "description",
		Price:       vo.MustMoney(10_000),
		Attributes:  domain.NewAttributes(map[string]string{"color": "black"}),
		Now:         fixedTime,
	})
	require.NoError(t, err)

	return item
}

func TestNewItem_Validation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		title     string
		ownerID   uuid.UUID
		wantField string
	}{
		{name: "valid", title: defaultTitle, ownerID: uuid.New()},
		{name: "title too short", title: "ab", ownerID: uuid.New(), wantField: "title"},
		{name: "title blank", title: "   ", ownerID: uuid.New(), wantField: "title"},
		{name: "missing owner", title: defaultTitle, ownerID: uuid.Nil, wantField: "owner_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			item, err := domain.NewItem(domain.NewItemParams{
				OwnerID: tt.ownerID,
				Title:   tt.title,
				Price:   vo.ZeroMoney(),
				Now:     fixedTime,
			})

			if tt.wantField == "" {
				require.NoError(t, err)
				require.Equal(t, domain.StatusDraft, item.Status())
				return
			}

			require.Error(t, err)
			require.Contains(t, err.Error(), tt.wantField)
		})
	}
}

func TestItem_StatusTransitions(t *testing.T) {
	t.Parallel()

	t.Run("owner publishes draft directly", func(t *testing.T) {
		t.Parallel()

		item := newItem(t)

		require.NoError(t, item.Publish(fixedTime))
		require.Equal(t, domain.StatusPublished, item.Status())
	})

	t.Run("moderation stays an optional path", func(t *testing.T) {
		t.Parallel()

		item := newItem(t)

		require.NoError(t, item.SubmitForModeration(fixedTime))
		require.Equal(t, domain.StatusModeration, item.Status())

		require.NoError(t, item.Publish(fixedTime))
		require.Equal(t, domain.StatusPublished, item.Status())

		require.NoError(t, item.Archive(fixedTime))
		require.Equal(t, domain.StatusArchived, item.Status())
	})

	t.Run("published item can be sold", func(t *testing.T) {
		t.Parallel()

		item := newItem(t)

		require.NoError(t, item.Publish(fixedTime))
		require.NoError(t, item.MarkSold(fixedTime))
		require.Equal(t, domain.StatusSold, item.Status())
	})

	t.Run("draft cannot be sold", func(t *testing.T) {
		t.Parallel()

		item := newItem(t)

		require.Error(t, item.MarkSold(fixedTime))
		require.Equal(t, domain.StatusDraft, item.Status())
	})

	t.Run("sold is terminal", func(t *testing.T) {
		t.Parallel()

		item := newItem(t)
		require.NoError(t, item.Publish(fixedTime))
		require.NoError(t, item.MarkSold(fixedTime))

		require.True(t, domain.StatusSold.IsTerminal())
		require.Error(t, item.Archive(fixedTime))
		require.Error(t, item.Restore(fixedTime))
		require.Error(t, item.Publish(fixedTime))
		require.Equal(t, domain.StatusSold, item.Status())
	})

	t.Run("sold item cannot be updated", func(t *testing.T) {
		t.Parallel()

		item := newItem(t)
		require.NoError(t, item.Publish(fixedTime))
		require.NoError(t, item.MarkSold(fixedTime))

		newTitle := "Another title"
		require.Error(t, item.Update(domain.UpdateItemParams{Title: &newTitle, Now: fixedTime}))
	})

	t.Run("archived item cannot be updated", func(t *testing.T) {
		t.Parallel()

		item := newItem(t)
		require.NoError(t, item.SubmitForModeration(fixedTime))
		require.NoError(t, item.Publish(fixedTime))
		require.NoError(t, item.Archive(fixedTime))

		newTitle := "Another title"
		require.Error(t, item.Update(domain.UpdateItemParams{Title: &newTitle, Now: fixedTime}))
	})
}

func TestItem_AttributesAreIsolated(t *testing.T) {
	t.Parallel()

	item := newItem(t)

	attrs := item.Attributes()
	attrs["color"] = "white"

	stored, ok := item.Attributes().Get("color")
	require.True(t, ok)
	require.Equal(t, "black", stored)
}

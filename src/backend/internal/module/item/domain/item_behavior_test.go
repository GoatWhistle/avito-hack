package domain_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/item/domain"
	"github.com/avito-hack/backend/internal/shared/domainerr"
	"github.com/avito-hack/backend/internal/shared/vo"
)

var behaviorTestTime = time.Date(2026, time.August, 1, 10, 0, 0, 0, time.UTC)

func seedItem(t *testing.T, status domain.Status) *domain.Item {
	t.Helper()

	return domain.RestoreItem(domain.RestoreItemParams{
		ID: uuid.New(), DisplayID: domain.NewDisplayID(), OwnerID: uuid.New(), Title: "Seed listing",
		Description: "provided by platform",
		Price:       vo.MustMoney(10000), Status: status, Attributes: domain.NewAttributes(nil),
		CreatedAt: behaviorTestTime, UpdatedAt: behaviorTestTime, IsSeed: true, AIVerified: true,
	})
}

func TestSeedItemRejectsMutatingActions(t *testing.T) {
	t.Parallel()

	later := behaviorTestTime.Add(time.Hour)
	newTitle := "Renamed"

	tests := []struct {
		name  string
		from  domain.Status
		apply func(*domain.Item) error
	}{
		{name: "update", from: domain.StatusDraft, apply: func(i *domain.Item) error {
			return i.Update(domain.UpdateItemParams{Title: &newTitle, Now: later})
		}},
		{name: "submit for moderation", from: domain.StatusDraft, apply: func(i *domain.Item) error {
			return i.SubmitForModeration(later)
		}},
		{name: "mark sold", from: domain.StatusPublished, apply: func(i *domain.Item) error {
			return i.MarkSold(later)
		}},
		{name: "archive", from: domain.StatusPublished, apply: func(i *domain.Item) error {
			return i.Archive(later)
		}},
		{name: "restore", from: domain.StatusArchived, apply: func(i *domain.Item) error {
			return i.Restore(later)
		}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			item := seedItem(t, tc.from)

			err := tc.apply(item)

			require.ErrorIs(t, err, domainerr.ErrConflict)
			assert.Equal(t, tc.from, item.Status())
			assert.Equal(t, behaviorTestTime, item.UpdatedAt())
			assert.Equal(t, "Seed listing", item.Title())
		})
	}
}

func TestUpdateContentChangeTriggersReverification(t *testing.T) {
	t.Parallel()

	item := itemInStatus(t, domain.StatusPublished)
	item.MarkAIVerified()
	later := statusTestTime.Add(time.Hour)
	newTitle := "New title"

	err := item.Update(domain.UpdateItemParams{Title: &newTitle, Now: later})

	require.NoError(t, err)
	assert.Equal(t, domain.StatusModeration, item.Status())
	assert.False(t, item.AIVerified())
}

func TestUpdatePriceOnlyDoesNotTriggerReverification(t *testing.T) {
	t.Parallel()

	item := itemInStatus(t, domain.StatusPublished)
	item.MarkAIVerified()
	later := statusTestTime.Add(time.Hour)
	price := vo.MustMoney(999900)

	err := item.Update(domain.UpdateItemParams{Price: &price, Now: later})

	require.NoError(t, err)
	assert.Equal(t, domain.StatusPublished, item.Status())
	assert.True(t, item.AIVerified())
}

func TestRequireReverificationRejectsSeedItem(t *testing.T) {
	t.Parallel()

	item := seedItem(t, domain.StatusPublished)

	err := item.RequireReverification(behaviorTestTime.Add(time.Hour))

	require.ErrorIs(t, err, domainerr.ErrConflict)
	assert.Equal(t, domain.StatusPublished, item.Status())
	assert.True(t, item.AIVerified())
}

func TestRequireReverificationMovesPublishedItemToModeration(t *testing.T) {
	t.Parallel()

	item := itemInStatus(t, domain.StatusPublished)
	item.MarkAIVerified()
	later := statusTestTime.Add(time.Hour)

	err := item.RequireReverification(later)

	require.NoError(t, err)
	assert.Equal(t, domain.StatusModeration, item.Status())
	assert.False(t, item.AIVerified())
	assert.Equal(t, later, item.UpdatedAt())
}

func TestMarkAndClearAIVerification(t *testing.T) {
	t.Parallel()

	item := itemInStatus(t, domain.StatusDraft)
	assert.False(t, item.AIVerified())

	item.MarkAIVerified()
	assert.True(t, item.AIVerified())

	item.ClearAIVerification()
	assert.False(t, item.AIVerified())
}

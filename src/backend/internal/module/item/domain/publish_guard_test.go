package domain_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/item/domain"
	"github.com/avito-hack/backend/internal/shared/domainerr"
)

func TestPublishRequiresApprovedVerdict(t *testing.T) {
	t.Parallel()

	verdicts := []domain.ModerationVerdict{
		domain.ModerationRejected,
		domain.ModerationUnavailable,
		domain.ModerationVerdict(""),
		domain.ModerationVerdict("APPROVED"),
		domain.ModerationVerdict("approved "),
	}

	for _, verdict := range verdicts {
		t.Run(string(verdict), func(t *testing.T) {
			t.Parallel()

			item := itemInModeration(t)

			err := item.Publish(fixedTime, verdict)

			require.ErrorIs(t, err, domainerr.ErrConflict)
			assert.Equal(t, domain.StatusModeration, item.Status())
			assert.False(t, item.AIVerified())
		})
	}
}

func TestPublishWithApprovedVerdictMarksVerified(t *testing.T) {
	t.Parallel()

	item := itemInModeration(t)

	require.NoError(t, item.Publish(fixedTime, domain.ModerationApproved))
	assert.Equal(t, domain.StatusPublished, item.Status())
	assert.True(t, item.AIVerified())
}

func TestAttributeChangeClearsVerificationAndDemotes(t *testing.T) {
	t.Parallel()

	item := itemInModeration(t)
	require.NoError(t, item.Publish(fixedTime, domain.ModerationApproved))

	attributes := domain.NewAttributes(map[string]string{"note": "new unmoderated text"})
	require.NoError(t, item.Update(domain.UpdateItemParams{Attributes: &attributes, Now: fixedTime}))

	assert.Equal(t, domain.StatusModeration, item.Status())
	assert.False(t, item.AIVerified())
}

func itemInModeration(t *testing.T) *domain.Item {
	t.Helper()

	item, err := domain.NewItem(domain.NewItemParams{
		OwnerID: uuid.New(), Title: "Чайник", Description: "обычный чайник", Now: fixedTime,
	})
	require.NoError(t, err)
	require.NoError(t, item.SubmitForModeration(fixedTime))

	return item
}

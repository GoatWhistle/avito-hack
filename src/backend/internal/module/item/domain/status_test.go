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

var statusTestTime = time.Date(2026, time.July, 1, 12, 0, 0, 0, time.UTC)

func allStatuses() []domain.Status {
	return []domain.Status{
		domain.StatusDraft, domain.StatusModeration,
		domain.StatusPublished, domain.StatusSold, domain.StatusArchived,
	}
}

func itemInStatus(t *testing.T, status domain.Status) *domain.Item {
	t.Helper()

	return domain.RestoreItem(domain.RestoreItemParams{
		ID: uuid.New(), DisplayID: domain.NewDisplayID(), OwnerID: uuid.New(), Title: "Bicycle",
		Description: "good one",
		Price:       vo.MustMoney(50000), Status: status, Attributes: domain.NewAttributes(nil),
		CreatedAt: statusTestTime, UpdatedAt: statusTestTime,
	})
}

func TestStatusValidity(t *testing.T) {
	t.Parallel()

	for _, status := range allStatuses() {
		assert.True(t, status.Valid(), status.String())
		assert.Equal(t, string(status), status.String())
	}

	assert.False(t, domain.Status("").Valid())
	assert.False(t, domain.Status("deleted").Valid())
}

func TestStatusIsTerminal(t *testing.T) {
	t.Parallel()

	assert.True(t, domain.StatusSold.IsTerminal())
	assert.True(t, domain.Status("unknown").IsTerminal())

	for _, status := range []domain.Status{
		domain.StatusDraft, domain.StatusModeration, domain.StatusPublished, domain.StatusArchived,
	} {
		assert.False(t, status.IsTerminal(), status.String())
	}
}

func TestStatusTransitionMatrix(t *testing.T) {
	t.Parallel()

	allowed := map[domain.Status]map[domain.Status]bool{
		domain.StatusDraft: {
			domain.StatusModeration: true, domain.StatusArchived: true,
		},
		domain.StatusModeration: {
			domain.StatusPublished: true, domain.StatusDraft: true, domain.StatusArchived: true,
		},
		domain.StatusPublished: {domain.StatusSold: true, domain.StatusArchived: true},
		domain.StatusSold:      {},
		domain.StatusArchived:  {domain.StatusDraft: true},
	}

	for from, targets := range allowed {
		for _, to := range allStatuses() {
			t.Run(from.String()+"_to_"+to.String(), func(t *testing.T) {
				t.Parallel()

				assert.Equal(t, targets[to], from.CanTransitionTo(to))
			})
		}
	}
}

func TestStatusUnknownSourceAllowsNothing(t *testing.T) {
	t.Parallel()

	unknown := domain.Status("mystery")

	for _, to := range allStatuses() {
		assert.False(t, unknown.CanTransitionTo(to))
	}
}

func TestItemBehaviorAllowedTransitions(t *testing.T) {
	t.Parallel()

	later := statusTestTime.Add(time.Hour)

	tests := []struct {
		name  string
		from  domain.Status
		apply func(*domain.Item) error
		want  domain.Status
	}{
		{
			name: "draft to moderation",
			from: domain.StatusDraft,
			apply: func(i *domain.Item) error {
				return i.SubmitForModeration(later)
			},
			want: domain.StatusModeration,
		},
		{
			name:  "moderation to published",
			from:  domain.StatusModeration,
			apply: func(i *domain.Item) error { return i.Publish(later) },
			want:  domain.StatusPublished,
		},
		{
			name:  "moderation back to draft",
			from:  domain.StatusModeration,
			apply: func(i *domain.Item) error { return i.Restore(later) },
			want:  domain.StatusDraft,
		},
		{
			name:  "published to sold",
			from:  domain.StatusPublished,
			apply: func(i *domain.Item) error { return i.MarkSold(later) },
			want:  domain.StatusSold,
		},
		{
			name:  "published to archived",
			from:  domain.StatusPublished,
			apply: func(i *domain.Item) error { return i.Archive(later) },
			want:  domain.StatusArchived,
		},
		{
			name:  "archived back to draft",
			from:  domain.StatusArchived,
			apply: func(i *domain.Item) error { return i.Restore(later) },
			want:  domain.StatusDraft,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			item := itemInStatus(t, tc.from)

			require.NoError(t, tc.apply(item))
			assert.Equal(t, tc.want, item.Status())
			assert.Equal(t, later, item.UpdatedAt())
		})
	}
}

func TestItemBehaviorForbiddenTransitions(t *testing.T) {
	t.Parallel()

	later := statusTestTime.Add(time.Hour)

	tests := []struct {
		name  string
		from  domain.Status
		apply func(*domain.Item) error
	}{
		{name: "sold cannot be published", from: domain.StatusSold, apply: func(i *domain.Item) error {
			return i.Publish(later)
		}},
		{name: "sold cannot be archived", from: domain.StatusSold, apply: func(i *domain.Item) error {
			return i.Archive(later)
		}},
		{name: "sold cannot be restored", from: domain.StatusSold, apply: func(i *domain.Item) error {
			return i.Restore(later)
		}},
		{name: "sold cannot be resold", from: domain.StatusSold, apply: func(i *domain.Item) error {
			return i.MarkSold(later)
		}},
		{name: "draft cannot be sold", from: domain.StatusDraft, apply: func(i *domain.Item) error {
			return i.MarkSold(later)
		}},
		{name: "draft cannot be published directly", from: domain.StatusDraft, apply: func(i *domain.Item) error {
			return i.Publish(later)
		}},
		{name: "draft cannot be restored", from: domain.StatusDraft, apply: func(i *domain.Item) error {
			return i.Restore(later)
		}},
		{name: "moderation cannot be sold", from: domain.StatusModeration, apply: func(i *domain.Item) error {
			return i.MarkSold(later)
		}},
		{name: "published cannot go to moderation", from: domain.StatusPublished, apply: func(i *domain.Item) error {
			return i.SubmitForModeration(later)
		}},
		{name: "published cannot be republished", from: domain.StatusPublished, apply: func(i *domain.Item) error {
			return i.Publish(later)
		}},
		{name: "archived cannot be published", from: domain.StatusArchived, apply: func(i *domain.Item) error {
			return i.Publish(later)
		}},
		{name: "archived cannot be sold", from: domain.StatusArchived, apply: func(i *domain.Item) error {
			return i.MarkSold(later)
		}},
		{name: "archived cannot be re-archived", from: domain.StatusArchived, apply: func(i *domain.Item) error {
			return i.Archive(later)
		}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			item := itemInStatus(t, tc.from)

			err := tc.apply(item)

			require.ErrorIs(t, err, domainerr.ErrConflict)
			assert.Equal(t, tc.from, item.Status())
			assert.Equal(t, statusTestTime, item.UpdatedAt())
		})
	}
}

func TestSubmitForModerationRequiresDescription(t *testing.T) {
	t.Parallel()

	item := domain.RestoreItem(domain.RestoreItemParams{
		ID: uuid.New(), DisplayID: domain.NewDisplayID(), OwnerID: uuid.New(), Title: "Bicycle", Description: "",
		Price: vo.MustMoney(1), Status: domain.StatusDraft,
		CreatedAt: statusTestTime, UpdatedAt: statusTestTime,
	})

	err := item.SubmitForModeration(statusTestTime.Add(time.Hour))

	var invalid *domainerr.InvalidError
	require.ErrorAs(t, err, &invalid)
	assert.Equal(t, "description", invalid.Field)
	assert.Equal(t, domain.StatusDraft, item.Status())
}

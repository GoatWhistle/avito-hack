package domain_test

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/item/domain"
	"github.com/avito-hack/backend/internal/shared/domainerr"
	"github.com/avito-hack/backend/internal/shared/vo"
)

func ptr[T any](value T) *T { return &value }

func TestItemUpdateAppliesPartialChanges(t *testing.T) {
	t.Parallel()

	later := statusTestTime.Add(time.Hour)
	item := itemInStatus(t, domain.StatusDraft)

	require.NoError(t, item.Update(domain.UpdateItemParams{
		Title:       ptr("  Mountain bike  "),
		Description: ptr("  fresh description  "),
		Price:       ptr(vo.MustMoney(99900)),
		Attributes:  ptr(domain.NewAttributes(map[string]string{"color": "red"})),
		Now:         later,
	}))

	assert.Equal(t, "Mountain bike", item.Title())
	assert.Equal(t, "fresh description", item.Description())
	assert.Equal(t, int64(99900), item.Price().Kopeks())
	assert.Equal(t, later, item.UpdatedAt())

	color, ok := item.Attributes().Get("color")
	require.True(t, ok)
	assert.Equal(t, "red", color)
}

func TestItemUpdateKeepsUntouchedFields(t *testing.T) {
	t.Parallel()

	item := itemInStatus(t, domain.StatusPublished)
	title, description, price := item.Title(), item.Description(), item.Price()

	require.NoError(t, item.Update(domain.UpdateItemParams{Now: statusTestTime.Add(time.Hour)}))

	assert.Equal(t, title, item.Title())
	assert.Equal(t, description, item.Description())
	assert.Equal(t, price, item.Price())
	assert.Equal(t, domain.StatusPublished, item.Status())
}

func TestItemUpdateRejectsFinalStatuses(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		status domain.Status
	}{
		{name: "archived", status: domain.StatusArchived},
		{name: "sold", status: domain.StatusSold},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			item := itemInStatus(t, tc.status)
			original := item.Title()

			err := item.Update(domain.UpdateItemParams{
				Title: ptr("New title"), Now: statusTestTime.Add(time.Hour),
			})

			require.ErrorIs(t, err, domainerr.ErrConflict)
			assert.Equal(t, original, item.Title())
			assert.Equal(t, statusTestTime, item.UpdatedAt())
		})
	}
}

func TestItemUpdateValidatesFields(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		params    domain.UpdateItemParams
		wantField string
	}{
		{name: "title too short", params: domain.UpdateItemParams{Title: ptr("ab")}, wantField: "title"},
		{name: "title blank", params: domain.UpdateItemParams{Title: ptr("     ")}, wantField: "title"},
		{
			name:      "title too long",
			params:    domain.UpdateItemParams{Title: ptr(strings.Repeat("a", 201))},
			wantField: "title",
		},
		{
			name:      "description too long",
			params:    domain.UpdateItemParams{Description: ptr(strings.Repeat("b", 5001))},
			wantField: "description",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			item := itemInStatus(t, domain.StatusDraft)
			params := tc.params
			params.Now = statusTestTime.Add(time.Hour)

			err := item.Update(params)

			var invalid *domainerr.InvalidError
			require.ErrorAs(t, err, &invalid)
			assert.Equal(t, tc.wantField, invalid.Field)
			assert.Equal(t, statusTestTime, item.UpdatedAt())
		})
	}
}

func TestItemUpdateBoundaryLengths(t *testing.T) {
	t.Parallel()

	item := itemInStatus(t, domain.StatusDraft)

	require.NoError(t, item.Update(domain.UpdateItemParams{
		Title:       ptr(strings.Repeat("a", 200)),
		Description: ptr(strings.Repeat("b", 5000)),
		Now:         statusTestTime.Add(time.Hour),
	}))

	assert.Len(t, item.Title(), 200)
	assert.Len(t, item.Description(), 5000)

	require.NoError(t, item.Update(domain.UpdateItemParams{
		Title: ptr("abc"), Description: ptr(""), Now: statusTestTime.Add(2 * time.Hour),
	}))

	assert.Equal(t, "abc", item.Title())
	assert.Empty(t, item.Description())
}

func TestAttributes(t *testing.T) {
	t.Parallel()

	assert.Empty(t, domain.NewAttributes(nil))
	assert.Zero(t, domain.NewAttributes(nil).Len())
	assert.Empty(t, domain.NewAttributes(map[string]string{}))

	raw := map[string]string{"size": "L", "brand": "Avito"}
	attributes := domain.NewAttributes(raw)

	assert.Equal(t, 2, attributes.Len())

	raw["size"] = "XL"
	size, ok := attributes.Get("size")
	require.True(t, ok)
	assert.Equal(t, "L", size)

	_, missing := attributes.Get("absent")
	assert.False(t, missing)

	var nilAttributes domain.Attributes
	assert.NotNil(t, nilAttributes.Clone())
	assert.Empty(t, nilAttributes.Clone())
}

func TestItemOwnership(t *testing.T) {
	t.Parallel()

	item := itemInStatus(t, domain.StatusDraft)

	assert.True(t, item.IsOwnedBy(item.OwnerID()))
	assert.False(t, item.IsOwnedBy(item.ID()))
	assert.Equal(t, statusTestTime, item.CreatedAt())
}

func TestErrTooManyPhotos(t *testing.T) {
	t.Parallel()

	err := domain.ErrTooManyPhotos()

	require.ErrorIs(t, err, domainerr.ErrConflict)
	assert.Contains(t, err.Error(), "10")
}

package app_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/item/app"
	"github.com/avito-hack/backend/internal/module/item/domain"
	"github.com/avito-hack/backend/internal/shared/auth"
	"github.com/avito-hack/backend/internal/shared/domainerr"
	"github.com/avito-hack/backend/internal/shared/vo"
)

func itemWithStatus(t *testing.T, ownerID uuid.UUID, status domain.Status) *domain.Item {
	t.Helper()

	return domain.RestoreItem(domain.RestoreItemParams{
		ID: uuid.New(), DisplayID: domain.NewDisplayID(), OwnerID: ownerID, Title: "Chair",
		Description: "solid oak",
		Price:       vo.MustMoney(1000), Status: status, Attributes: domain.NewAttributes(nil),
		CreatedAt: fixedTime, UpdatedAt: fixedTime,
	})
}

func TestGetItemVisibility(t *testing.T) {
	t.Parallel()

	ownerID := uuid.New()
	strangerID := uuid.New()

	tests := []struct {
		name    string
		status  domain.Status
		actor   auth.Actor
		wantErr error
	}{
		{name: "published is public", status: domain.StatusPublished},
		{name: "sold is public", status: domain.StatusSold},
		{
			name:   "owner sees own draft",
			status: domain.StatusDraft,
			actor:  auth.Actor{ID: ownerID, Role: auth.RoleUser},
		},
		{
			name:   "moderator sees draft",
			status: domain.StatusDraft,
			actor:  auth.Actor{ID: strangerID, Role: auth.RoleModerator},
		},
		{
			name:   "admin sees archived",
			status: domain.StatusArchived,
			actor:  auth.Actor{ID: strangerID, Role: auth.RoleAdmin},
		},
		{
			name:    "anonymous gets not found for draft",
			status:  domain.StatusDraft,
			wantErr: domain.ErrItemNotFound,
		},
		{
			name:    "anonymous gets not found for moderation",
			status:  domain.StatusModeration,
			wantErr: domain.ErrItemNotFound,
		},
		{
			name:    "stranger is forbidden",
			status:  domain.StatusDraft,
			actor:   auth.Actor{ID: strangerID, Role: auth.RoleUser},
			wantErr: domainerr.ErrForbidden,
		},
		{
			name:    "stranger is forbidden for archived",
			status:  domain.StatusArchived,
			actor:   auth.Actor{ID: strangerID, Role: auth.RoleUser},
			wantErr: domainerr.ErrForbidden,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			item := itemWithStatus(t, ownerID, tc.status)
			handler := app.NewGetItemHandler(&stubRepository{item: item}, nil, nil)

			view, err := handler.Handle(t.Context(), app.GetItemQuery{ItemID: item.ID(), Actor: tc.actor})

			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
				assert.Nil(t, view.Item)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, item.ID(), view.Item.ID())
		})
	}
}

func TestGetItemPropagatesRepositoryError(t *testing.T) {
	t.Parallel()

	sentinel := errors.New("db unreachable")
	handler := app.NewGetItemHandler(&failingRepository{err: sentinel}, nil, nil)

	_, err := handler.Handle(t.Context(), app.GetItemQuery{ItemID: uuid.New()})

	require.ErrorIs(t, err, sentinel)
}

func TestCreateItemHandler(t *testing.T) {
	t.Parallel()

	repo := &stubRepository{}
	handler := app.NewCreateItemHandler(repo, passthroughTx{}, fakeClock{})
	ownerID := uuid.New()

	item, err := handler.Handle(t.Context(), app.CreateItemCommand{
		OwnerID:     ownerID,
		Title:       "  Sofa  ",
		Description: "  comfy  ",
		PriceKopeks: 250000,
		Attributes:  map[string]string{"color": "grey"},
	})

	require.NoError(t, err)
	assert.Equal(t, "Sofa", item.Title())
	assert.Equal(t, "comfy", item.Description())
	assert.Equal(t, domain.StatusDraft, item.Status())
	assert.Equal(t, ownerID, item.OwnerID())
	assert.Equal(t, int64(250000), item.Price().Kopeks())
	assert.Equal(t, fixedTime, item.CreatedAt())
	require.NotNil(t, repo.saved)
}

func TestCreateItemHandlerValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		cmd       app.CreateItemCommand
		wantField string
	}{
		{
			name:      "negative price",
			cmd:       app.CreateItemCommand{OwnerID: uuid.New(), Title: "Sofa", PriceKopeks: -1},
			wantField: "price",
		},
		{
			name:      "short title",
			cmd:       app.CreateItemCommand{OwnerID: uuid.New(), Title: "so", PriceKopeks: 100},
			wantField: "title",
		},
		{
			name:      "missing owner",
			cmd:       app.CreateItemCommand{Title: "Sofa", PriceKopeks: 100},
			wantField: "owner_id",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := &stubRepository{}
			handler := app.NewCreateItemHandler(repo, passthroughTx{}, fakeClock{})

			_, err := handler.Handle(t.Context(), tc.cmd)

			var invalid *domainerr.InvalidError
			require.ErrorAs(t, err, &invalid)
			assert.Equal(t, tc.wantField, invalid.Field)
			assert.Nil(t, repo.saved)
		})
	}
}

func TestCreateItemHandlerPropagatesSaveError(t *testing.T) {
	t.Parallel()

	sentinel := errors.New("insert failed")
	handler := app.NewCreateItemHandler(&failingRepository{err: sentinel}, passthroughTx{}, fakeClock{})

	_, err := handler.Handle(t.Context(), app.CreateItemCommand{
		OwnerID: uuid.New(), Title: "Sofa", PriceKopeks: 100,
	})

	require.ErrorIs(t, err, sentinel)
}

type failingRepository struct{ err error }

func (f *failingRepository) Save(_ context.Context, _ *domain.Item) error { return f.err }

func (f *failingRepository) ByID(_ context.Context, _ uuid.UUID) (*domain.Item, error) {
	return nil, f.err
}

func (f *failingRepository) ByIDForUpdate(_ context.Context, _ uuid.UUID) (*domain.Item, error) {
	return nil, f.err
}

func (f *failingRepository) ByDisplayID(_ context.Context, _ string) (*domain.Item, error) {
	return nil, f.err
}

func (f *failingRepository) Delete(_ context.Context, _ uuid.UUID) error { return f.err }

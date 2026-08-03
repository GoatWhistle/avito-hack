package app_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/item/app"
	"github.com/avito-hack/backend/internal/module/item/domain"
	"github.com/avito-hack/backend/internal/shared/auth"
	"github.com/avito-hack/backend/internal/shared/domainerr"
	"github.com/avito-hack/backend/internal/shared/vo"
)

var fixedTime = time.Date(2026, time.January, 1, 12, 0, 0, 0, time.UTC)

type fakeClock struct{}

func (fakeClock) Now() time.Time { return fixedTime }

type passthroughTx struct{}

func (passthroughTx) WithTx(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

type stubRepository struct {
	item  *domain.Item
	saved *domain.Item
}

func (s *stubRepository) Save(_ context.Context, item *domain.Item) error {
	s.saved = item

	return nil
}

func (s *stubRepository) ByID(_ context.Context, _ uuid.UUID) (*domain.Item, error) {
	return s.item, nil
}

func (s *stubRepository) ByIDForUpdate(_ context.Context, _ uuid.UUID) (*domain.Item, error) {
	return s.item, nil
}

func (s *stubRepository) Delete(_ context.Context, _ uuid.UUID) error { return nil }

func newModeratedItem(t *testing.T, ownerID uuid.UUID) *domain.Item {
	t.Helper()

	item, err := domain.NewItem(domain.NewItemParams{
		OwnerID:     ownerID,
		Title:       "MacBook Pro",
		Description: "description",
		Price:       vo.MustMoney(1000),
		Now:         fixedTime,
	})
	require.NoError(t, err)
	require.NoError(t, item.SubmitForModeration(fixedTime))

	return item
}

func TestChangeStatusHandler_Publish(t *testing.T) {
	t.Parallel()

	ownerID := uuid.New()

	tests := []struct {
		name    string
		actor   auth.Actor
		wantErr error
	}{
		{
			name:    "moderator publishes",
			actor:   auth.Actor{ID: uuid.New(), Role: auth.RoleModerator},
			wantErr: nil,
		},
		{
			name:    "owner cannot publish",
			actor:   auth.Actor{ID: ownerID, Role: auth.RoleUser},
			wantErr: domainerr.ErrForbidden,
		},
		{
			name:    "stranger cannot publish",
			actor:   auth.Actor{ID: uuid.New(), Role: auth.RoleUser},
			wantErr: domainerr.ErrForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &stubRepository{item: newModeratedItem(t, ownerID)}
			handler := app.NewChangeStatusHandler(repo, passthroughTx{}, fakeClock{})

			result, err := handler.Handle(context.Background(), app.ChangeStatusCommand{
				ItemID: uuid.New(),
				Actor:  tt.actor,
				Action: app.ActionPublish,
			})

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.Nil(t, repo.saved)
				return
			}

			require.NoError(t, err)
			require.Equal(t, domain.StatusPublished, result.Status())
			require.NotNil(t, repo.saved)
		})
	}
}

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
	"github.com/avito-hack/backend/internal/shared/events"
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

type stubPhotos struct{ count int }

func (s *stubPhotos) Add(_ context.Context, _ *domain.Photo) error { return nil }

func (s *stubPhotos) ByItemID(_ context.Context, _ uuid.UUID) ([]*domain.Photo, error) {
	return nil, nil
}

func (s *stubPhotos) CountByItemID(_ context.Context, _ uuid.UUID) (int, error) {
	return s.count, nil
}

func (s *stubPhotos) DeleteByID(_ context.Context, _, _ uuid.UUID) (string, error) {
	return "", nil
}

func newDraftItem(t *testing.T, ownerID uuid.UUID) *domain.Item {
	t.Helper()

	item, err := domain.NewItem(domain.NewItemParams{
		OwnerID:     ownerID,
		Title:       "MacBook Pro",
		Description: "description",
		Price:       vo.MustMoney(1000),
		Now:         fixedTime,
	})
	require.NoError(t, err)

	return item
}

func newPublishedItem(t *testing.T, ownerID uuid.UUID) *domain.Item {
	t.Helper()

	item := newDraftItem(t, ownerID)
	require.NoError(t, item.Publish(fixedTime))

	return item
}

func newHandler(repo *stubRepository, bus events.Publisher) *app.ChangeStatusHandler {
	return app.NewChangeStatusHandler(repo, &stubPhotos{count: 2}, passthroughTx{}, fakeClock{}, bus)
}

func TestChangeStatusHandler_OwnerPublishesDraft(t *testing.T) {
	t.Parallel()

	ownerID := uuid.New()

	tests := []struct {
		name    string
		actor   auth.Actor
		wantErr error
	}{
		{name: "owner publishes", actor: auth.Actor{ID: ownerID, Role: auth.RoleUser}},
		{name: "moderator publishes", actor: auth.Actor{ID: uuid.New(), Role: auth.RoleModerator}},
		{
			name:    "stranger cannot publish",
			actor:   auth.Actor{ID: uuid.New(), Role: auth.RoleUser},
			wantErr: domainerr.ErrForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &stubRepository{item: newDraftItem(t, ownerID)}
			handler := newHandler(repo, events.NopPublisher{})

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

func TestChangeStatusHandler_Sell(t *testing.T) {
	t.Parallel()

	ownerID := uuid.New()

	tests := []struct {
		name    string
		actor   auth.Actor
		wantErr error
	}{
		{name: "owner sells", actor: auth.Actor{ID: ownerID, Role: auth.RoleUser}},
		{
			name:    "moderator cannot sell someone else item",
			actor:   auth.Actor{ID: uuid.New(), Role: auth.RoleModerator},
			wantErr: domainerr.ErrForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &stubRepository{item: newPublishedItem(t, ownerID)}
			handler := newHandler(repo, events.NopPublisher{})

			result, err := handler.Handle(context.Background(), app.ChangeStatusCommand{
				ItemID: uuid.New(),
				Actor:  tt.actor,
				Action: app.ActionSell,
			})

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)

				return
			}

			require.NoError(t, err)
			require.Equal(t, domain.StatusSold, result.Status())
		})
	}
}

func TestChangeStatusHandler_EmitsEvents(t *testing.T) {
	t.Parallel()

	ownerID := uuid.New()

	tests := []struct {
		name      string
		item      func(*testing.T, uuid.UUID) *domain.Item
		action    app.StatusAction
		eventType events.Type
	}{
		{
			name: "publish emits item.published", item: newDraftItem,
			action: app.ActionPublish, eventType: events.TypeItemPublished,
		},
		{
			name: "sell emits item.sold", item: newPublishedItem,
			action: app.ActionSell, eventType: events.TypeItemSold,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			bus := events.NewBus(nil)
			received := make([]events.Event, 0, 1)
			bus.Subscribe(tt.eventType, func(_ context.Context, e events.Event) error {
				received = append(received, e)

				return nil
			})

			repo := &stubRepository{item: tt.item(t, ownerID)}
			handler := newHandler(repo, bus)

			_, err := handler.Handle(context.Background(), app.ChangeStatusCommand{
				ItemID: uuid.New(),
				Actor:  auth.Actor{ID: ownerID, Role: auth.RoleUser},
				Action: tt.action,
			})
			require.NoError(t, err)

			require.Len(t, received, 1)
			require.Equal(t, ownerID, received[0].UserID)
			require.Equal(t, 2, received[0].Payload.PhotoCount)
		})
	}
}

func TestChangeStatusHandler_ArchiveEmitsNothing(t *testing.T) {
	t.Parallel()

	ownerID := uuid.New()
	bus := events.NewBus(nil)

	calls := 0
	for _, eventType := range []events.Type{events.TypeItemPublished, events.TypeItemSold, events.TypeItemUpdated} {
		bus.Subscribe(eventType, func(context.Context, events.Event) error {
			calls++

			return nil
		})
	}

	repo := &stubRepository{item: newPublishedItem(t, ownerID)}
	handler := newHandler(repo, bus)

	_, err := handler.Handle(context.Background(), app.ChangeStatusCommand{
		ItemID: uuid.New(),
		Actor:  auth.Actor{ID: ownerID, Role: auth.RoleUser},
		Action: app.ActionArchive,
	})
	require.NoError(t, err)
	require.Equal(t, 0, calls)
}

package app_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/item/app"
	"github.com/avito-hack/backend/internal/module/item/domain"
	"github.com/avito-hack/backend/internal/shared/auth"
	"github.com/avito-hack/backend/internal/shared/domainerr"
	"github.com/avito-hack/backend/internal/shared/events"
)

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

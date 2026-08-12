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

func TestChangeStatusHandler_ModeratorMayOnlyArchive(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		item    func(*testing.T, uuid.UUID) *domain.Item
		action  app.StatusAction
		wantErr error
	}{
		{name: "archive is allowed", item: newPublishedItem, action: app.ActionArchive},
		{
			name: "submit is denied", item: newDraftItem,
			action: app.ActionSubmit, wantErr: domainerr.ErrForbidden,
		},
		{
			name: "restore is denied", item: newArchivedItem,
			action: app.ActionRestore, wantErr: domainerr.ErrForbidden,
		},
		{
			name: "sell is denied", item: newPublishedItem,
			action: app.ActionSell, wantErr: domainerr.ErrForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			for _, role := range []auth.Role{auth.RoleModerator, auth.RoleAdmin} {
				repo := &stubRepository{item: tt.item(t, uuid.New())}
				handler := newHandler(repo, events.NopPublisher{})

				_, err := handler.Handle(context.Background(), app.ChangeStatusCommand{
					ItemID: uuid.New(),
					Actor:  auth.Actor{ID: uuid.New(), Role: role},
					Action: tt.action,
				})

				if tt.wantErr != nil {
					require.ErrorIs(t, err, tt.wantErr)
					require.Nil(t, repo.saved)

					continue
				}

				require.NoError(t, err)
				require.NotNil(t, repo.saved)
			}
		})
	}
}

func TestChangeStatusHandler_OwnerKeepsFullControl(t *testing.T) {
	t.Parallel()

	ownerID := uuid.New()
	repo := &stubRepository{item: newArchivedItem(t, ownerID)}
	handler := newHandler(repo, events.NopPublisher{})

	result, err := handler.Handle(context.Background(), app.ChangeStatusCommand{
		ItemID: uuid.New(),
		Actor:  auth.Actor{ID: ownerID, Role: auth.RoleUser},
		Action: app.ActionRestore,
	})

	require.NoError(t, err)
	require.Equal(t, domain.StatusDraft, result.Status())
}

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
	"github.com/avito-hack/backend/internal/shared/events"
)

func collectEvents(t *testing.T, types ...events.Type) (*events.Bus, *[]events.Event) {
	t.Helper()

	bus := events.NewBus(nil)
	collected := make([]events.Event, 0)

	for _, eventType := range types {
		bus.Subscribe(eventType, func(_ context.Context, e events.Event) error {
			collected = append(collected, e)

			return nil
		})
	}

	return bus, &collected
}

func TestUpdateItemHandlerAppliesChanges(t *testing.T) {
	t.Parallel()

	ownerID := uuid.New()
	item := itemWithStatus(t, ownerID, domain.StatusDraft)
	repo := &stubRepository{item: item}
	bus, collected := collectEvents(t, events.TypeItemUpdated)

	handler := app.NewUpdateItemHandler(repo, &stubPhotos{count: 2}, passthroughTx{}, fakeClock{}, bus)

	title := "Updated title"
	price := int64(777)
	attributes := map[string]string{"color": "blue"}

	updated, err := handler.Handle(t.Context(), app.UpdateItemCommand{
		ItemID: item.ID(), ActorID: ownerID,
		Title: &title, PriceKopeks: &price, Attributes: &attributes,
	})

	require.NoError(t, err)
	assert.Equal(t, "Updated title", updated.Title())
	assert.Equal(t, int64(777), updated.Price().Kopeks())
	require.NotNil(t, repo.saved)

	require.Len(t, *collected, 1)
	event := (*collected)[0]
	assert.Equal(t, events.TypeItemUpdated, event.Type)
	assert.Equal(t, ownerID, event.UserID)
	assert.Equal(t, item.ID(), event.SubjectID)
	assert.Equal(t, "Updated title", event.Payload.Title)
	assert.Equal(t, int64(777), event.Payload.PriceKopeks)
	assert.Equal(t, 2, event.Payload.PhotoCount)

	color, ok := event.Payload.Attribute("color")
	require.True(t, ok)
	assert.Equal(t, "blue", color)
}

func TestUpdateItemHandlerRejectsNonOwner(t *testing.T) {
	t.Parallel()

	item := itemWithStatus(t, uuid.New(), domain.StatusDraft)
	repo := &stubRepository{item: item}
	bus, collected := collectEvents(t, events.TypeItemUpdated)

	handler := app.NewUpdateItemHandler(repo, &stubPhotos{}, passthroughTx{}, fakeClock{}, bus)

	title := "Hijacked"
	_, err := handler.Handle(t.Context(), app.UpdateItemCommand{
		ItemID: item.ID(), ActorID: uuid.New(), Title: &title,
	})

	require.ErrorIs(t, err, domainerr.ErrForbidden)
	assert.Nil(t, repo.saved)
	assert.Empty(t, *collected)
}

func TestUpdateItemHandlerValidatesPrice(t *testing.T) {
	t.Parallel()

	ownerID := uuid.New()
	item := itemWithStatus(t, ownerID, domain.StatusDraft)
	repo := &stubRepository{item: item}

	handler := app.NewUpdateItemHandler(repo, &stubPhotos{}, passthroughTx{}, fakeClock{}, events.NopPublisher{})

	price := int64(-1)
	_, err := handler.Handle(t.Context(), app.UpdateItemCommand{
		ItemID: item.ID(), ActorID: ownerID, PriceKopeks: &price,
	})

	var invalid *domainerr.InvalidError
	require.ErrorAs(t, err, &invalid)
	assert.Equal(t, "price", invalid.Field)
	assert.Nil(t, repo.saved)
}

func TestUpdateItemHandlerRejectsSoldItem(t *testing.T) {
	t.Parallel()

	ownerID := uuid.New()
	item := itemWithStatus(t, ownerID, domain.StatusSold)
	repo := &stubRepository{item: item}
	bus, collected := collectEvents(t, events.TypeItemUpdated)

	handler := app.NewUpdateItemHandler(repo, &stubPhotos{}, passthroughTx{}, fakeClock{}, bus)

	title := "Too late"
	_, err := handler.Handle(t.Context(), app.UpdateItemCommand{
		ItemID: item.ID(), ActorID: ownerID, Title: &title,
	})

	require.ErrorIs(t, err, domainerr.ErrConflict)
	assert.Empty(t, *collected)
}

func TestUpdateItemHandlerPropagatesLoadFailure(t *testing.T) {
	t.Parallel()

	sentinel := errors.New("row is locked")
	handler := app.NewUpdateItemHandler(
		&failingRepository{err: sentinel}, &stubPhotos{}, passthroughTx{}, fakeClock{}, events.NopPublisher{})

	_, err := handler.Handle(t.Context(), app.UpdateItemCommand{ItemID: uuid.New(), ActorID: uuid.New()})

	require.ErrorIs(t, err, sentinel)
}

func TestChangeStatusFullFlow(t *testing.T) {
	t.Parallel()

	ownerID := uuid.New()

	tests := []struct {
		name      string
		status    domain.Status
		action    app.StatusAction
		wantState domain.Status
	}{
		{
			name: "submit for moderation", status: domain.StatusDraft,
			action: app.ActionSubmit, wantState: domain.StatusModeration,
		},
		{
			name: "archive published", status: domain.StatusPublished,
			action: app.ActionArchive, wantState: domain.StatusArchived,
		},
		{
			name: "restore archived", status: domain.StatusArchived,
			action: app.ActionRestore, wantState: domain.StatusDraft,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			item := itemWithStatus(t, ownerID, tc.status)
			repo := &stubRepository{item: item}
			handler := app.NewChangeStatusHandler(
				repo, &stubPhotos{}, passthroughTx{}, fakeClock{}, events.NopPublisher{})

			updated, err := handler.Handle(t.Context(), app.ChangeStatusCommand{
				ItemID: item.ID(), Actor: actorOf(ownerID), Action: tc.action,
			})

			require.NoError(t, err)
			assert.Equal(t, tc.wantState, updated.Status())
		})
	}
}

func TestChangeStatusRejectsUnknownAction(t *testing.T) {
	t.Parallel()

	ownerID := uuid.New()
	item := itemWithStatus(t, ownerID, domain.StatusDraft)
	repo := &stubRepository{item: item}

	handler := app.NewChangeStatusHandler(repo, &stubPhotos{}, passthroughTx{}, fakeClock{}, events.NopPublisher{})

	_, err := handler.Handle(t.Context(), app.ChangeStatusCommand{
		ItemID: item.ID(), Actor: actorOf(ownerID), Action: app.StatusAction("delete"),
	})

	var invalid *domainerr.InvalidError
	require.ErrorAs(t, err, &invalid)
	assert.Equal(t, "action", invalid.Field)
	assert.Nil(t, repo.saved)
}

func TestChangeStatusPropagatesLoadFailure(t *testing.T) {
	t.Parallel()

	sentinel := errors.New("deadlock detected")
	handler := app.NewChangeStatusHandler(
		&failingRepository{err: sentinel}, &stubPhotos{}, passthroughTx{}, fakeClock{}, events.NopPublisher{})

	_, err := handler.Handle(t.Context(), app.ChangeStatusCommand{
		ItemID: uuid.New(), Actor: actorOf(uuid.New()), Action: app.ActionPublish,
	})

	require.ErrorIs(t, err, sentinel)
}

func TestChangeStatusToleratesPhotoCountFailure(t *testing.T) {
	t.Parallel()

	ownerID := uuid.New()
	item := itemWithStatus(t, ownerID, domain.StatusDraft)
	bus, collected := collectEvents(t, events.TypeItemPublished)

	handler := app.NewChangeStatusHandler(
		&stubRepository{item: item}, &countingPhotos{countErr: errors.New("boom")},
		passthroughTx{}, fakeClock{}, bus)

	_, err := handler.Handle(t.Context(), app.ChangeStatusCommand{
		ItemID: item.ID(), Actor: actorOf(ownerID), Action: app.ActionPublish,
	})

	require.NoError(t, err)
	require.Len(t, *collected, 1)
	assert.Zero(t, (*collected)[0].Payload.PhotoCount)
}

func actorOf(id uuid.UUID) auth.Actor {
	return auth.Actor{ID: id, Role: auth.RoleUser}
}

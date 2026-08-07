package app_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/favorite/app"
	"github.com/avito-hack/backend/internal/module/favorite/domain"
	"github.com/avito-hack/backend/internal/shared/events"
)

var fixedTime = time.Date(2026, time.January, 1, 12, 0, 0, 0, time.UTC)

type fakeClock struct{}

func (fakeClock) Now() time.Time { return fixedTime }

type passthroughTx struct{}

func (passthroughTx) WithTx(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

type stubRepository struct {
	stored map[string]struct{}
}

func newStubRepository() *stubRepository {
	return &stubRepository{stored: make(map[string]struct{})}
}

func key(userID, itemID uuid.UUID) string { return userID.String() + "|" + itemID.String() }

func (s *stubRepository) Add(_ context.Context, f *domain.Favorite) (bool, error) {
	k := key(f.UserID(), f.ItemID())
	if _, ok := s.stored[k]; ok {
		return false, nil
	}
	s.stored[k] = struct{}{}

	return true, nil
}

func (s *stubRepository) Remove(_ context.Context, userID, itemID uuid.UUID) error {
	k := key(userID, itemID)
	if _, ok := s.stored[k]; !ok {
		return domain.ErrFavoriteNotFound
	}
	delete(s.stored, k)

	return nil
}

type stubItems struct{ exists bool }

func (s stubItems) Exists(context.Context, uuid.UUID) (bool, error) { return s.exists, nil }

func TestAddFavorite_IsIdempotent(t *testing.T) {
	t.Parallel()

	bus := events.NewBus(nil)
	emitted := 0
	bus.Subscribe(events.TypeFavoriteAdded, func(context.Context, events.Event) error {
		emitted++

		return nil
	})

	repo := newStubRepository()
	handler := app.NewAddFavoriteHandler(repo, stubItems{exists: true}, passthroughTx{}, fakeClock{}, bus)

	cmd := app.AddFavoriteCommand{UserID: uuid.New(), ItemID: uuid.New()}

	require.NoError(t, handler.Handle(context.Background(), cmd))
	require.NoError(t, handler.Handle(context.Background(), cmd))

	require.Len(t, repo.stored, 1)
	require.Equal(t, 1, emitted)
}

func TestAddFavorite_UnknownItem(t *testing.T) {
	t.Parallel()

	handler := app.NewAddFavoriteHandler(
		newStubRepository(), stubItems{exists: false}, passthroughTx{}, fakeClock{}, events.NopPublisher{})

	err := handler.Handle(context.Background(), app.AddFavoriteCommand{UserID: uuid.New(), ItemID: uuid.New()})

	require.ErrorIs(t, err, domain.ErrItemNotFound)
}

func TestRemoveFavorite(t *testing.T) {
	t.Parallel()

	repo := newStubRepository()
	userID, itemID := uuid.New(), uuid.New()

	add := app.NewAddFavoriteHandler(repo, stubItems{exists: true}, passthroughTx{}, fakeClock{}, events.NopPublisher{})
	require.NoError(t, add.Handle(context.Background(), app.AddFavoriteCommand{UserID: userID, ItemID: itemID}))

	remove := app.NewRemoveFavoriteHandler(repo, passthroughTx{})
	require.NoError(t, remove.Handle(context.Background(), app.RemoveFavoriteCommand{UserID: userID, ItemID: itemID}))

	err := remove.Handle(context.Background(), app.RemoveFavoriteCommand{UserID: userID, ItemID: itemID})
	require.ErrorIs(t, err, domain.ErrFavoriteNotFound)
}

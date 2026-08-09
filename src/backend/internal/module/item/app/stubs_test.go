package app_test

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/item/app"
	"github.com/avito-hack/backend/internal/module/item/domain"
	"github.com/avito-hack/backend/internal/shared/auth"
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

type stubStorage struct {
	saved       int
	removed     int
	lastType    string
	saveErr     error
	returnedURL string
}

func (s *stubStorage) Save(
	_ context.Context, _ uuid.UUID, content io.Reader, contentType string,
) (app.StoredFile, error) {
	if s.saveErr != nil {
		return app.StoredFile{}, s.saveErr
	}

	body, err := io.ReadAll(content)
	if err != nil {
		return app.StoredFile{}, err
	}

	s.saved++
	s.lastType = contentType

	return app.StoredFile{Name: "photo.jpg", ContentType: contentType, Size: int64(len(body))}, nil
}

func (s *stubStorage) Delete(_ context.Context, _ uuid.UUID, _ string) error {
	s.removed++

	return nil
}

func (s *stubStorage) URL(_ uuid.UUID, name string) string {
	if s.returnedURL != "" {
		return s.returnedURL
	}

	return "/media/" + name
}

type countingPhotos struct {
	stubPhotos
	added    []*domain.Photo
	deleted  int
	countErr error
	addErr   error
}

func (c *countingPhotos) Add(_ context.Context, photo *domain.Photo) error {
	if c.addErr != nil {
		return c.addErr
	}

	c.added = append(c.added, photo)

	return nil
}

func (c *countingPhotos) CountByItemID(_ context.Context, _ uuid.UUID) (int, error) {
	if c.countErr != nil {
		return 0, c.countErr
	}

	return c.count, nil
}

func (c *countingPhotos) ByItemID(_ context.Context, itemID uuid.UUID) ([]*domain.Photo, error) {
	if c.countErr != nil {
		return nil, c.countErr
	}

	result := make([]*domain.Photo, 0, len(c.added))
	for _, photo := range c.added {
		if photo.ItemID() == itemID {
			result = append(result, photo)
		}
	}

	return result, nil
}

func (c *countingPhotos) DeleteByID(_ context.Context, _, _ uuid.UUID) (string, error) {
	c.deleted++

	return "/media/photo.jpg", nil
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

func newArchivedItem(t *testing.T, ownerID uuid.UUID) *domain.Item {
	t.Helper()

	item := newPublishedItem(t, ownerID)
	require.NoError(t, item.Archive(fixedTime))

	return item
}

func newHandler(repo *stubRepository, bus events.Publisher) *app.ChangeStatusHandler {
	return app.NewChangeStatusHandler(repo, &stubPhotos{count: 2}, passthroughTx{}, fakeClock{}, bus)
}

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

func actorOf(id uuid.UUID) auth.Actor {
	return auth.Actor{ID: id, Role: auth.RoleUser}
}

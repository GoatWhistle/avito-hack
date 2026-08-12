package api_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/item/api"
	"github.com/avito-hack/backend/internal/module/item/app"
	"github.com/avito-hack/backend/internal/module/item/domain"
	"github.com/avito-hack/backend/internal/shared/auth"
	"github.com/avito-hack/backend/internal/shared/events"
	"github.com/avito-hack/backend/internal/shared/validate"
	"github.com/avito-hack/backend/internal/shared/vo"
)

var fixedTime = time.Date(2026, time.February, 2, 12, 0, 0, 0, time.UTC)

type fakeClock struct{}

func (fakeClock) Now() time.Time { return fixedTime }

type passthroughTx struct{}

func (passthroughTx) WithTx(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

type memoryItems struct {
	items map[uuid.UUID]*domain.Item
}

func newMemoryItems() *memoryItems {
	return &memoryItems{items: map[uuid.UUID]*domain.Item{}}
}

func (m *memoryItems) Save(_ context.Context, item *domain.Item) error {
	m.items[item.ID()] = item

	return nil
}

func (m *memoryItems) ByID(_ context.Context, id uuid.UUID) (*domain.Item, error) {
	item, ok := m.items[id]
	if !ok {
		return nil, domain.ErrItemNotFound
	}

	return item, nil
}

func (m *memoryItems) ByIDForUpdate(ctx context.Context, id uuid.UUID) (*domain.Item, error) {
	return m.ByID(ctx, id)
}

func (m *memoryItems) ByDisplayID(_ context.Context, displayID string) (*domain.Item, error) {
	for _, item := range m.items {
		if item.DisplayID() == displayID {
			return item, nil
		}
	}

	return nil, domain.ErrItemNotFound
}

type memoryPhotos struct {
	photos map[uuid.UUID][]*domain.Photo
}

func newMemoryPhotos() *memoryPhotos {
	return &memoryPhotos{photos: map[uuid.UUID][]*domain.Photo{}}
}

func (m *memoryPhotos) Add(_ context.Context, photo *domain.Photo) error {
	m.photos[photo.ItemID()] = append(m.photos[photo.ItemID()], photo)

	return nil
}

func (m *memoryPhotos) ByItemID(_ context.Context, itemID uuid.UUID) ([]*domain.Photo, error) {
	return m.photos[itemID], nil
}

func (m *memoryPhotos) CountByItemID(_ context.Context, itemID uuid.UUID) (int, error) {
	return len(m.photos[itemID]), nil
}

func (m *memoryPhotos) DeleteByDisplayID(_ context.Context, itemID uuid.UUID, displayID string) (string, error) {
	kept := make([]*domain.Photo, 0, len(m.photos[itemID]))
	url := ""

	for _, photo := range m.photos[itemID] {
		if photo.DisplayID() == displayID {
			url = photo.URL()

			continue
		}
		kept = append(kept, photo)
	}
	m.photos[itemID] = kept

	if url == "" {
		return "", domain.ErrPhotoNotFound
	}

	return url, nil
}

type memoryModerationLog struct {
	entries []domain.ModerationLogEntry
}

func (m *memoryModerationLog) Add(_ context.Context, entry domain.ModerationLogEntry) error {
	m.entries = append(m.entries, entry)

	return nil
}

func (m *memoryModerationLog) LatestByItemID(
	_ context.Context, itemID uuid.UUID,
) (*domain.ModerationLogEntry, error) {
	for i := len(m.entries) - 1; i >= 0; i-- {
		if m.entries[i].ItemID == itemID {
			return &m.entries[i], nil
		}
	}

	return nil, nil //nolint:nilnil // absence of a moderation entry is a valid, expected state
}

type memoryReadModel struct {
	rows []app.ListItem
}

func (m *memoryReadModel) List(_ context.Context, f app.ListFilter) ([]app.ListItem, error) {
	filtered := make([]app.ListItem, 0, len(m.rows))
	for _, row := range m.rows {
		if !row.Status.IsPublic() && row.OwnerID != f.ViewerID {
			continue
		}
		if f.Status != "" && row.Status != f.Status {
			continue
		}
		if f.OwnerID != uuid.Nil && row.OwnerID != f.OwnerID {
			continue
		}
		filtered = append(filtered, row)
	}

	if len(filtered) > f.Limit {
		return filtered[:f.Limit], nil
	}

	return filtered, nil
}

type nopStorage struct{}

func (nopStorage) Save(context.Context, uuid.UUID, io.Reader, string) (app.StoredFile, error) {
	return app.StoredFile{Name: "photo.jpg", ContentType: "image/jpeg"}, nil
}

func (nopStorage) Delete(context.Context, uuid.UUID, string) error { return nil }

func (nopStorage) URL(itemDisplayID, name string) string { return "/uploads/" + itemDisplayID + "/" + name }

type fixture struct {
	router http.Handler
	items  *memoryItems
	photos *memoryPhotos
	read   *memoryReadModel
	actor  *auth.Actor
}

func newFixture(t *testing.T, actor *auth.Actor) *fixture {
	t.Helper()

	items := newMemoryItems()
	photos := newMemoryPhotos()
	read := &memoryReadModel{}
	moderationLog := &memoryModerationLog{}
	bus := events.NopPublisher{}

	authenticate := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if actor == nil {
				w.WriteHeader(http.StatusUnauthorized)

				return
			}

			next.ServeHTTP(w, r.WithContext(auth.WithActor(r.Context(), *actor)))
		})
	}

	optional := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if actor == nil {
				next.ServeHTTP(w, r)

				return
			}

			next.ServeHTTP(w, r.WithContext(auth.WithActor(r.Context(), *actor)))
		})
	}

	handlers := api.NewHandlers(api.Deps{
		Items:         items,
		CreateItem:    app.NewCreateItemHandler(items, passthroughTx{}, fakeClock{}),
		UpdateItem:    app.NewUpdateItemHandler(items, photos, passthroughTx{}, fakeClock{}, bus, nil),
		ChangeStatus:  app.NewChangeStatusHandler(items, photos, passthroughTx{}, fakeClock{}, bus, nil),
		GetItem:       app.NewGetItemHandler(items, moderationLog, nil),
		ListItems:     app.NewListItemsHandler(read, nil),
		AddPhoto:      app.NewAddPhotoHandler(items, photos, nopStorage{}, passthroughTx{}, fakeClock{}, nil),
		ListPhotos:    app.NewListPhotosHandler(photos),
		DeletePhoto:   app.NewDeletePhotoHandler(items, photos, nopStorage{}, passthroughTx{}, fakeClock{}, nil),
		Validator:     validate.New(),
		Authenticate:  authenticate,
		OptionalAuth:  optional,
		MaxBodyBytes:  1 << 20,
		MaxPhotoBytes: 1 << 20,
	})

	router := chi.NewRouter()
	handlers.RegisterRoutes(router)

	return &fixture{router: router, items: items, photos: photos, read: read, actor: actor}
}

func (f *fixture) do(t *testing.T, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	f.router.ServeHTTP(rec, req)

	return rec
}

func (f *fixture) seedItem(t *testing.T, ownerID uuid.UUID, status domain.Status) *domain.Item {
	t.Helper()

	item := domain.RestoreItem(domain.RestoreItemParams{
		ID: uuid.New(), DisplayID: domain.NewDisplayID(), OwnerID: ownerID, Title: "Bicycle", Description: "fast one",
		Price: vo.MustMoney(150000), Status: status, Attributes: domain.NewAttributes(nil),
		CreatedAt: fixedTime, UpdatedAt: fixedTime,
	})
	require.NoError(t, f.items.Save(context.Background(), item))

	return item
}

func userActor() auth.Actor {
	return auth.Actor{ID: uuid.New(), Role: auth.RoleUser}
}

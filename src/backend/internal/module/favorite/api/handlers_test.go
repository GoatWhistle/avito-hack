package api_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/favorite/api"
	"github.com/avito-hack/backend/internal/module/favorite/app"
	"github.com/avito-hack/backend/internal/module/favorite/domain"
	"github.com/avito-hack/backend/internal/shared/auth"
	"github.com/avito-hack/backend/internal/shared/events"
)

var fixedTime = time.Date(2026, time.March, 3, 9, 0, 0, 0, time.UTC)

type fixedClock struct{}

func (fixedClock) Now() time.Time { return fixedTime }

type passthroughTx struct{}

func (passthroughTx) WithTx(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

type memoryFavorites struct {
	stored map[string]struct{}
}

func newMemoryFavorites() *memoryFavorites {
	return &memoryFavorites{stored: map[string]struct{}{}}
}

func key(userID, itemID uuid.UUID) string { return userID.String() + "|" + itemID.String() }

func (m *memoryFavorites) Add(_ context.Context, f *domain.Favorite) (bool, error) {
	k := key(f.UserID(), f.ItemID())
	if _, ok := m.stored[k]; ok {
		return false, nil
	}
	m.stored[k] = struct{}{}

	return true, nil
}

func (m *memoryFavorites) Remove(_ context.Context, userID, itemID uuid.UUID) error {
	k := key(userID, itemID)
	if _, ok := m.stored[k]; !ok {
		return domain.ErrFavoriteNotFound
	}
	delete(m.stored, k)

	return nil
}

type stubItems struct{ exists bool }

func (s stubItems) Exists(context.Context, uuid.UUID) (bool, error) { return s.exists, nil }

type stubRead struct {
	rows []app.FavoriteItem
	err  error
}

func (s stubRead) List(_ context.Context, f app.ListFilter) ([]app.FavoriteItem, error) {
	if s.err != nil {
		return nil, s.err
	}

	if len(s.rows) > f.Limit {
		return s.rows[:f.Limit], nil
	}

	return s.rows, nil
}

type routerOptions struct {
	actor      *auth.Actor
	itemExists bool
	read       stubRead
}

func newRouter(t *testing.T, opts routerOptions) http.Handler {
	t.Helper()

	favorites := newMemoryFavorites()
	handlers := api.NewHandlers(api.Deps{
		AddFavorite: app.NewAddFavoriteHandler(
			favorites, stubItems{exists: opts.itemExists}, passthroughTx{}, fixedClock{}, events.NopPublisher{}),
		RemoveFavorite: app.NewRemoveFavoriteHandler(favorites, passthroughTx{}),
		ListFavorites:  app.NewListFavoritesHandler(opts.read),
		Authenticate: func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if opts.actor == nil {
					next.ServeHTTP(w, r)

					return
				}

				next.ServeHTTP(w, r.WithContext(auth.WithActor(r.Context(), *opts.actor)))
			})
		},
	})

	router := chi.NewRouter()
	handlers.RegisterRoutes(router)

	return router
}

func do(t *testing.T, router http.Handler, method, path string) *httptest.ResponseRecorder {
	t.Helper()

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(method, path, http.NoBody))

	return rec
}

func actor() auth.Actor {
	return auth.Actor{ID: uuid.New(), Role: auth.RoleUser}
}

func TestAddFavoriteEndpoint(t *testing.T) {
	t.Parallel()

	current := actor()
	router := newRouter(t, routerOptions{actor: &current, itemExists: true})
	path := "/items/" + uuid.NewString() + "/favorite"

	require.Equal(t, http.StatusNoContent, do(t, router, http.MethodPost, path).Code)
	assert.Equal(t, http.StatusNoContent, do(t, router, http.MethodPost, path).Code)
}

func TestAddFavoriteEndpointErrors(t *testing.T) {
	t.Parallel()

	current := actor()

	tests := []struct {
		name       string
		opts       routerOptions
		path       string
		wantStatus int
	}{
		{
			name:       "unauthenticated",
			opts:       routerOptions{itemExists: true},
			path:       "/items/" + uuid.NewString() + "/favorite",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "malformed item id",
			opts:       routerOptions{actor: &current, itemExists: true},
			path:       "/items/not-a-uuid/favorite",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "unknown item",
			opts:       routerOptions{actor: &current},
			path:       "/items/" + uuid.NewString() + "/favorite",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.wantStatus, do(t, newRouter(t, tc.opts), http.MethodPost, tc.path).Code)
		})
	}
}

func TestRemoveFavoriteEndpoint(t *testing.T) {
	t.Parallel()

	current := actor()
	router := newRouter(t, routerOptions{actor: &current, itemExists: true})
	path := "/items/" + uuid.NewString() + "/favorite"

	require.Equal(t, http.StatusNoContent, do(t, router, http.MethodPost, path).Code)
	assert.Equal(t, http.StatusNoContent, do(t, router, http.MethodDelete, path).Code)
	assert.Equal(t, http.StatusNotFound, do(t, router, http.MethodDelete, path).Code)
}

func TestRemoveFavoriteEndpointRequiresAuth(t *testing.T) {
	t.Parallel()

	router := newRouter(t, routerOptions{itemExists: true})
	path := "/items/" + uuid.NewString() + "/favorite"

	assert.Equal(t, http.StatusUnauthorized, do(t, router, http.MethodDelete, path).Code)
}

package api_test

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/weeklylottery/api"
)

func TestRoutes(t *testing.T) {
	t.Parallel()

	router := chi.NewRouter()
	handlers := api.NewHandlers(api.Deps{
		Authenticate: func(next http.Handler) http.Handler { return next },
	})
	handlers.RegisterRoutes(router)

	routes := make(map[string]string)
	err := chi.Walk(router, func(method, route string, _ http.Handler, _ ...func(http.Handler) http.Handler) error {
		routes[route] = method
		return nil
	})
	require.NoError(t, err)

	assert.Equal(t, http.MethodGet, routes["/weekly-lottery/prizes"])
	assert.Equal(t, http.MethodGet, routes["/weekly-lottery/state"])
	assert.Equal(t, http.MethodPost, routes["/weekly-lottery/runs"])
	assert.Equal(t, http.MethodPost, routes["/weekly-lottery/runs/{rid}/slots/{slot}/reveal"])
}

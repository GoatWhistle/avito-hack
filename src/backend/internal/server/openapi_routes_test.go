package server_test

import (
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/avito-hack/backend/internal/config"
	favoriteapi "github.com/avito-hack/backend/internal/module/favorite/api"
	itemapi "github.com/avito-hack/backend/internal/module/item/api"
	petapi "github.com/avito-hack/backend/internal/module/pet/api"
	raccoonapi "github.com/avito-hack/backend/internal/module/raccoon/api"
	userapi "github.com/avito-hack/backend/internal/module/user/api"
	"github.com/avito-hack/backend/internal/server"
)

const (
	specRelPath   = "../../../../docs/openapi.yaml"
	uploadURLPath = "/uploads"
	uploadDirPath = "./uploads"
)

var undocumentedPaths = map[string]struct{}{
	"/metrics": {},
}

type openapiSpec struct {
	Paths map[string]yaml.Node `yaml:"paths"`
}

func newTestRouter(t *testing.T) http.Handler {
	t.Helper()

	passthrough := func(next http.Handler) http.Handler { return next }

	return server.NewRouter(server.RouterDeps{
		Config: config.Config{
			UploadURL: uploadURLPath,
			UploadDir: uploadDirPath,
		},
		Modules: []server.ModuleRegistrar{
			userapi.NewHandlers(userapi.Deps{Authenticate: passthrough}),
			itemapi.NewHandlers(itemapi.Deps{
				Authenticate: passthrough,
				OptionalAuth: passthrough,
			}),
			favoriteapi.NewHandlers(favoriteapi.Deps{Authenticate: passthrough}),
			raccoonapi.NewHandlers(raccoonapi.Deps{Authenticate: passthrough}),
			petapi.NewLeaderboardHandlers(petapi.LeaderboardDeps{Authenticate: passthrough}),
			petapi.NewPetHandlers(petapi.PetDeps{Authenticate: passthrough}),
			petapi.NewRewardHandlers(petapi.RewardDeps{Authenticate: passthrough}),
			petapi.NewQuestHandlers(petapi.QuestDeps{Authenticate: passthrough}),
			petapi.NewSummaryHandlers(petapi.SummaryDeps{Authenticate: passthrough}),
		},
		WebSocket: http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}),
	})
}

func registeredRoutes(t *testing.T, handler http.Handler) map[string]struct{} {
	t.Helper()

	router, ok := handler.(chi.Routes)
	require.True(t, ok, "router must expose chi.Routes for introspection")

	routes := map[string]struct{}{}

	walk := func(_, route string, _ http.Handler, _ ...func(http.Handler) http.Handler) error {
		routes[normalizeRoute(route)] = struct{}{}

		return nil
	}

	require.NoError(t, chi.Walk(router, walk))

	return routes
}

func normalizeRoute(route string) string {
	if route != "/" {
		route = strings.TrimSuffix(route, "/")
	}

	if route == "" {
		route = "/"
	}

	return route
}

func specPaths(t *testing.T) map[string]struct{} {
	t.Helper()

	raw, err := os.ReadFile(filepath.Clean(specRelPath))
	require.NoError(t, err, "docs/openapi.yaml must exist next to the backend module")

	var spec openapiSpec
	require.NoError(t, yaml.Unmarshal(raw, &spec))
	require.NotEmpty(t, spec.Paths, "spec must declare at least one path")

	paths := map[string]struct{}{}
	for path := range spec.Paths {
		paths[normalizeSpecPath(path)] = struct{}{}
	}

	return paths
}

var wildcardPrefixes = []string{uploadURLPath}

func normalizeSpecPath(path string) string {
	return collapseWildcardPrefix(placeholdersToStar(path))
}

func normalizeChiPath(route string) string {
	return collapseWildcardPrefix(placeholdersToStar(route))
}

func placeholdersToStar(path string) string {
	segments := strings.Split(path, "/")
	for i, segment := range segments {
		if strings.HasPrefix(segment, "{") && strings.HasSuffix(segment, "}") {
			segments[i] = "*"
		}
	}

	return strings.Join(segments, "/")
}

func collapseWildcardPrefix(path string) string {
	for _, prefix := range wildcardPrefixes {
		if path == prefix || strings.HasPrefix(path, prefix+"/") {
			return prefix + "/*"
		}
	}

	return path
}

func sorted(set map[string]struct{}) []string {
	values := make([]string, 0, len(set))
	for value := range set {
		values = append(values, value)
	}

	sort.Strings(values)

	return values
}

func TestOpenAPICoversEveryRegisteredRoute(t *testing.T) {
	t.Parallel()

	documented := specPaths(t)

	var missing []string

	for route := range registeredRoutes(t, newTestRouter(t)) {
		if _, skip := undocumentedPaths[route]; skip {
			continue
		}

		if _, ok := documented[normalizeChiPath(route)]; !ok {
			missing = append(missing, route)
		}
	}

	sort.Strings(missing)

	assert.Emptyf(t, missing, "routes registered in chi but absent from docs/openapi.yaml: %v", missing)
}

func TestOpenAPIDeclaresNoUnknownRoute(t *testing.T) {
	t.Parallel()

	registered := map[string]struct{}{}
	for route := range registeredRoutes(t, newTestRouter(t)) {
		registered[normalizeChiPath(route)] = struct{}{}
	}

	var extra []string

	for path := range specPaths(t) {
		if _, ok := registered[path]; !ok {
			extra = append(extra, path)
		}
	}

	sort.Strings(extra)

	assert.Emptyf(t, extra,
		"paths documented in docs/openapi.yaml but not registered in chi: %v\nregistered routes: %v",
		extra, sorted(registered))
}

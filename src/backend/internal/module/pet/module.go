package pet

import (
	"net/http"
	"net/url"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/avito-hack/backend/internal/module/pet/api"
	"github.com/avito-hack/backend/internal/module/pet/app"
	"github.com/avito-hack/backend/internal/module/pet/infra"
	"github.com/avito-hack/backend/internal/shared/auth"
)

type Options struct {
	Pool   *pgxpool.Pool
	Redis  *redis.Client
	Tx     app.TxManager
	Clock  app.Clock
	Tokens interface {
		Parse(string) (auth.Actor, error)
	}
	AllowedOrigins []string
}

type Module struct {
	WebSocket http.Handler
}

func New(opts Options) *Module {
	repository := infra.NewPgRepository(opts.Pool)
	cache := infra.NewRedisCache(opts.Redis)
	service := app.NewService(repository, cache, opts.Tx, opts.Clock)

	return &Module{WebSocket: api.NewWebSocketHandler(service, opts.Tokens, originPatterns(opts.AllowedOrigins))}
}

func originPatterns(origins []string) []string {
	patterns := make([]string, 0, len(origins))
	for _, origin := range origins {
		parsed, err := url.Parse(origin)
		if err == nil && parsed.Host != "" {
			patterns = append(patterns, parsed.Host)
		}
	}

	return patterns
}

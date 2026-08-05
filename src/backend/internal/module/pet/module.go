package pet

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/avito-hack/backend/internal/module/pet/api"
	"github.com/avito-hack/backend/internal/module/pet/app"
	"github.com/avito-hack/backend/internal/module/pet/domain"
	"github.com/avito-hack/backend/internal/module/pet/infra"
	"github.com/avito-hack/backend/internal/shared/auth"
	"github.com/avito-hack/backend/internal/shared/events"
	"github.com/avito-hack/backend/internal/shared/ws"
)

type Options struct {
	Pool   *pgxpool.Pool
	Redis  *redis.Client
	Tx     app.TxManager
	Clock  app.Clock
	Tokens interface {
		Parse(string) (auth.Actor, error)
	}
	AllowedOrigins   []string
	RewardHMACSecret string
	Bus              events.Subscriber
	Authenticate     func(http.Handler) http.Handler
}

type Module struct {
	WebSocket   http.Handler
	Service     *app.Service
	Rewards     *app.RewardService
	Signer      *domain.RewardSigner
	Hub         *ws.Hub
	Leaderboard *api.LeaderboardHandlers
	Pet         *api.PetHandlers
	RewardAPI   *api.RewardHandlers
}

func New(opts Options) (*Module, error) {
	signer, err := domain.NewRewardSigner(opts.RewardHMACSecret)
	if err != nil {
		return nil, fmt.Errorf("build reward signer: %w", err)
	}

	repository := infra.NewPgRepository(opts.Pool)
	journal := infra.NewPgXPEventRepository(opts.Pool)
	cache := infra.NewRedisCache(opts.Redis)

	hub := ws.NewHub()
	notifier := api.NewNotifier(hub)

	service := app.NewService(repository, journal, cache, opts.Tx, opts.Clock).WithHatchNotifier(notifier)

	rewards := app.NewRewardService(app.RewardServiceDeps{
		Rewards:  infra.NewPgRewardRepository(opts.Pool),
		Pets:     repository,
		Signer:   signer,
		Tx:       opts.Tx,
		Clock:    opts.Clock,
		Notifier: notifier,
	})

	app.NewSubscriber(service, notifier, opts.Clock).WithRewards(rewards).Register(opts.Bus)

	return &Module{
		WebSocket: api.NewWebSocketHandler(service, opts.Tokens, hub, originPatterns(opts.AllowedOrigins)),
		Service:   service,
		Rewards:   rewards,
		Signer:    signer,
		Hub:       hub,
		Leaderboard: api.NewLeaderboardHandlers(api.LeaderboardDeps{
			Leaderboard:  app.NewLeaderboardHandler(infra.NewPgLeaderboard(opts.Pool)),
			Authenticate: opts.Authenticate,
		}),
		Pet: api.NewPetHandlers(api.PetDeps{
			Service:      service,
			Rewards:      rewards,
			Authenticate: opts.Authenticate,
		}),
		RewardAPI: api.NewRewardHandlers(api.RewardDeps{
			Rewards:      rewards,
			Authenticate: opts.Authenticate,
		}),
	}, nil
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

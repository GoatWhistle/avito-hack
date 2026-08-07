package pet

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"

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
	Logger           *slog.Logger
	FlushInterval    time.Duration
	FlushBatchSize   int64
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
	SummaryAPI  *api.SummaryHandlers
	Async       *app.AsyncHandler
	Flusher     *app.FlushHotStateHandler

	FlushInterval time.Duration
}

func New(opts Options) (*Module, error) {
	signer, err := domain.NewRewardSigner(opts.RewardHMACSecret)
	if err != nil {
		return nil, fmt.Errorf("build reward signer: %w", err)
	}
	if opts.Logger == nil {
		opts.Logger = slog.Default()
	}

	repository := infra.NewPgRepository(opts.Pool)
	journal := infra.NewPgXPEventRepository(opts.Pool)
	cache := infra.NewRedisCache(opts.Redis)
	hot := infra.NewRedisHotStateStore(opts.Redis)

	hub := ws.NewHub()
	notifier := api.NewNotifier(hub)

	service := app.NewService(repository, journal, cache, opts.Tx, opts.Clock).
		WithHatchNotifier(notifier).
		WithHotState(hot)

	rewards := app.NewRewardService(app.RewardServiceDeps{
		Rewards:  infra.NewPgRewardRepository(opts.Pool),
		Pets:     repository,
		Signer:   signer,
		Tx:       opts.Tx,
		Clock:    opts.Clock,
		Notifier: notifier,
	})

	subscriber := app.NewSubscriber(service, notifier, opts.Clock).WithRewards(rewards)
	subscriber.Register(opts.Bus)

	leaderboard := infra.NewPgLeaderboard(opts.Pool)
	summaries := newSummaryService(opts, repository, leaderboard)

	module := &Module{
		WebSocket: api.NewWebSocketHandler(service, opts.Tokens, hub, originPatterns(opts.AllowedOrigins)),
		Service:   service,
		Rewards:   rewards,
		Signer:    signer,
		Hub:       hub,
		Async: app.NewAsyncHandler(
			subscriber, infra.NewPgEventDeduplicator(opts.Pool), opts.Tx, opts.Logger,
		),
		Flusher: app.NewFlushHotStateHandler(
			repository, cache, hot, opts.Tx, opts.FlushBatchSize, opts.Logger,
		),
		FlushInterval: opts.FlushInterval,
	}
	module.attachHandlers(opts, service, rewards, leaderboard, summaries)

	return module, nil
}

func newSummaryService(
	opts Options, repository *infra.PgRepository, leaderboard *infra.PgLeaderboard,
) *app.SummaryService {
	return app.NewSummaryService(app.SummaryServiceDeps{
		Summaries: infra.NewPgSummaryRepository(opts.Pool),
		Collector: app.NewFactCollector(app.FactCollectorDeps{
			Pets:     repository,
			Activity: infra.NewPgDayActivity(opts.Pool),
			Rank:     leaderboard,
		}),
		Tx:    opts.Tx,
		Clock: opts.Clock,
	})
}

func (m *Module) attachHandlers(
	opts Options,
	service *app.Service,
	rewards *app.RewardService,
	leaderboard *infra.PgLeaderboard,
	summaries *app.SummaryService,
) {
	m.Leaderboard = api.NewLeaderboardHandlers(api.LeaderboardDeps{
		Leaderboard:  app.NewLeaderboardHandler(leaderboard),
		Authenticate: opts.Authenticate,
	})
	m.SummaryAPI = api.NewSummaryHandlers(api.SummaryDeps{
		Summaries:    summaries,
		Authenticate: opts.Authenticate,
	})
	m.Pet = api.NewPetHandlers(api.PetDeps{
		Service:      service,
		Rewards:      rewards,
		Authenticate: opts.Authenticate,
	})
	m.RewardAPI = api.NewRewardHandlers(api.RewardDeps{
		Rewards:      rewards,
		Authenticate: opts.Authenticate,
	})
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

package games

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/avito-hack/backend/internal/module/games/api"
	"github.com/avito-hack/backend/internal/module/games/app"
	"github.com/avito-hack/backend/internal/module/games/domain"
	"github.com/avito-hack/backend/internal/module/games/domain/bukovki"
	"github.com/avito-hack/backend/internal/module/games/domain/moreless"
	"github.com/avito-hack/backend/internal/module/games/infra"
)

type Options struct {
	Pool         *pgxpool.Pool
	Tx           app.TxManager
	Clock        app.Clock
	PhotoSecret  string
	Validator    interface{ Struct(dst any) error }
	Authenticate func(http.Handler) http.Handler
	OptionalAuth func(http.Handler) http.Handler
	MaxBodyBytes int64
}

type Module struct {
	Handlers *api.Handlers
	Registry *domain.Registry
}

func New(opts Options) *Module {
	rounds := infra.NewPgRoundRepository(opts.Pool)
	progress := infra.NewPgProgressRepository(opts.Pool)

	words := infra.NewPgWordPool(opts.Pool)

	signer := moreless.NewPhotoSigner(opts.PhotoSecret)

	registry := domain.NewRegistry(
		moreless.New(infra.NewPgItemPool(opts.Pool), signer),
		bukovki.New(words, words),
	)

	handlers := api.NewHandlers(api.Deps{
		ListGames:    app.NewListGamesHandler(registry, progress),
		GetState:     app.NewGetStateHandler(registry, rounds, progress),
		StartRound:   app.NewStartRoundHandler(registry, rounds, opts.Tx, opts.Clock),
		Guess:        app.NewGuessHandler(registry, rounds, progress, opts.Tx, opts.Clock),
		ClaimReward:  app.NewClaimRewardHandler(progress, opts.Tx, opts.Clock),
		HiddenPhoto:  app.NewHiddenPhotoHandler(rounds, signer),
		Clock:        opts.Clock,
		Validator:    opts.Validator,
		Authenticate: opts.Authenticate,
		OptionalAuth: opts.OptionalAuth,
		MaxBodyBytes: opts.MaxBodyBytes,
	})

	return &Module{Handlers: handlers, Registry: registry}
}

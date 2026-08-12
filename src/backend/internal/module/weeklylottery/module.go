package weeklylottery

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/avito-hack/backend/internal/module/weeklylottery/api"
	"github.com/avito-hack/backend/internal/module/weeklylottery/app"
	"github.com/avito-hack/backend/internal/module/weeklylottery/domain"
	"github.com/avito-hack/backend/internal/module/weeklylottery/infra"
)

type Options struct {
	Pool         *pgxpool.Pool
	Tx           app.TxManager
	Clock        app.Clock
	Signer       app.CodeSigner
	Authenticate func(http.Handler) http.Handler
}

type Module struct {
	Handlers *api.Handlers
}

func New(opts Options) *Module {
	repository := infra.NewPgRepository(opts.Pool)
	generator := domain.NewGenerator(domain.CryptoRandom{})

	return &Module{Handlers: api.NewHandlers(api.Deps{
		GetState:     app.NewGetStateHandler(repository, opts.Clock),
		ListPrizes:   app.NewListPrizesHandler(),
		Start:        app.NewStartHandler(repository, opts.Tx, opts.Clock, generator),
		Reveal:       app.NewRevealHandler(repository, opts.Tx, opts.Clock, opts.Signer),
		Authenticate: opts.Authenticate,
	})}
}

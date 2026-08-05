package user

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/avito-hack/backend/internal/module/user/api"
	"github.com/avito-hack/backend/internal/module/user/app"
	"github.com/avito-hack/backend/internal/module/user/domain"
	"github.com/avito-hack/backend/internal/module/user/infra"
)

type Options struct {
	Pool         *pgxpool.Pool
	Tx           app.TxManager
	Clock        app.Clock
	Tokens       app.TokenIssuer
	Validator    interface{ Struct(dst any) error }
	Authenticate func(http.Handler) http.Handler
	MaxBodyBytes int64
}

type Module struct {
	Handlers   *api.Handlers
	Repository domain.Repository
}

func New(opts Options) *Module {
	repo := infra.NewPgRepository(opts.Pool)

	handlers := api.NewHandlers(api.Deps{
		Register:      app.NewRegisterUserHandler(repo, opts.Tx, opts.Clock),
		Login:         app.NewLoginUserHandler(repo, opts.Tokens),
		GetProfile:    app.NewGetProfileHandler(repo),
		UpdateProfile: app.NewUpdateProfileHandler(repo, opts.Tx, opts.Clock),
		Validator:     opts.Validator,
		Authenticate:  opts.Authenticate,
		MaxBodyBytes:  opts.MaxBodyBytes,
	})

	return &Module{Handlers: handlers, Repository: repo}
}

package item

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/avito-hack/backend/internal/module/item/api"
	"github.com/avito-hack/backend/internal/module/item/app"
	"github.com/avito-hack/backend/internal/module/item/infra"
)

type Options struct {
	Pool         *pgxpool.Pool
	Tx           app.TxManager
	Clock        app.Clock
	Validator    interface{ Struct(dst any) error }
	Authenticate func(http.Handler) http.Handler
	OptionalAuth func(http.Handler) http.Handler
	MaxBodyBytes int64
}

type Module struct {
	Handlers *api.Handlers
}

func New(opts Options) *Module {
	repo := infra.NewPgRepository(opts.Pool)
	read := infra.NewPgReadModel(opts.Pool)
	owners := infra.NewPgOwnerProvider(opts.Pool)

	handlers := api.NewHandlers(api.Deps{
		CreateItem:   app.NewCreateItemHandler(repo, opts.Tx, opts.Clock),
		UpdateItem:   app.NewUpdateItemHandler(repo, opts.Tx, opts.Clock),
		ChangeStatus: app.NewChangeStatusHandler(repo, opts.Tx, opts.Clock),
		GetItem:      app.NewGetItemHandler(repo),
		ListItems:    app.NewListItemsHandler(read, owners),
		Validator:    opts.Validator,
		Authenticate: opts.Authenticate,
		OptionalAuth: opts.OptionalAuth,
		MaxBodyBytes: opts.MaxBodyBytes,
	})

	return &Module{Handlers: handlers}
}

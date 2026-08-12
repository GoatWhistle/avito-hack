package favorite

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	itemdomain "github.com/avito-hack/backend/internal/module/item/domain"

	"github.com/avito-hack/backend/internal/module/favorite/api"
	"github.com/avito-hack/backend/internal/module/favorite/app"
	"github.com/avito-hack/backend/internal/module/favorite/infra"
	"github.com/avito-hack/backend/internal/shared/events"
)

type Options struct {
	Pool         *pgxpool.Pool
	Tx           app.TxManager
	Clock        app.Clock
	Bus          events.Publisher
	Items        itemdomain.Repository
	Authenticate func(http.Handler) http.Handler
}

type Module struct {
	Handlers *api.Handlers
}

func New(opts Options) *Module {
	repo := infra.NewPgRepository(opts.Pool)
	read := infra.NewPgReadModel(opts.Pool)
	items := infra.NewPgItemChecker(opts.Pool)

	bus := opts.Bus
	if bus == nil {
		bus = events.NopPublisher{}
	}

	handlers := api.NewHandlers(api.Deps{
		Items:          opts.Items,
		AddFavorite:    app.NewAddFavoriteHandler(repo, items, opts.Tx, opts.Clock, bus),
		RemoveFavorite: app.NewRemoveFavoriteHandler(repo, opts.Tx),
		ListFavorites:  app.NewListFavoritesHandler(read),
		Authenticate:   opts.Authenticate,
	})

	return &Module{Handlers: handlers}
}

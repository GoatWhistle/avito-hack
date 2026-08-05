package item

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/avito-hack/backend/internal/module/item/api"
	"github.com/avito-hack/backend/internal/module/item/app"
	"github.com/avito-hack/backend/internal/module/item/infra"
	"github.com/avito-hack/backend/internal/shared/events"
)

type Options struct {
	Pool          *pgxpool.Pool
	Tx            app.TxManager
	Clock         app.Clock
	Bus           events.Publisher
	Validator     interface{ Struct(dst any) error }
	Authenticate  func(http.Handler) http.Handler
	OptionalAuth  func(http.Handler) http.Handler
	MaxBodyBytes  int64
	MaxPhotoBytes int64
	UploadDir     string
	UploadURL     string
}

type Module struct {
	Handlers *api.Handlers
}

func New(opts Options) *Module {
	repo := infra.NewPgRepository(opts.Pool)
	read := infra.NewPgReadModel(opts.Pool)
	owners := infra.NewPgOwnerProvider(opts.Pool)
	photos := infra.NewPgPhotoRepository(opts.Pool)
	storage := infra.NewLocalPhotoStorage(opts.UploadDir, opts.UploadURL)

	bus := opts.Bus
	if bus == nil {
		bus = events.NopPublisher{}
	}

	handlers := api.NewHandlers(api.Deps{
		CreateItem:    app.NewCreateItemHandler(repo, opts.Tx, opts.Clock),
		UpdateItem:    app.NewUpdateItemHandler(repo, photos, opts.Tx, opts.Clock, bus),
		ChangeStatus:  app.NewChangeStatusHandler(repo, photos, opts.Tx, opts.Clock, bus),
		GetItem:       app.NewGetItemHandler(repo),
		ListItems:     app.NewListItemsHandler(read, owners),
		AddPhoto:      app.NewAddPhotoHandler(repo, photos, storage, opts.Tx, opts.Clock),
		ListPhotos:    app.NewListPhotosHandler(photos),
		DeletePhoto:   app.NewDeletePhotoHandler(repo, photos, storage, opts.Tx),
		Validator:     opts.Validator,
		Authenticate:  opts.Authenticate,
		OptionalAuth:  opts.OptionalAuth,
		MaxBodyBytes:  opts.MaxBodyBytes,
		MaxPhotoBytes: opts.MaxPhotoBytes,
	})

	return &Module{Handlers: handlers}
}

package item

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/avito-hack/backend/internal/module/item/api"
	"github.com/avito-hack/backend/internal/module/item/app"
	"github.com/avito-hack/backend/internal/module/item/domain"
	"github.com/avito-hack/backend/internal/module/item/infra"
	"github.com/avito-hack/backend/internal/shared/events"
)

type Options struct {
	Pool             *pgxpool.Pool
	Tx               app.TxManager
	Clock            app.Clock
	Bus              events.Publisher
	Validator        interface{ Struct(dst any) error }
	Authenticate     func(http.Handler) http.Handler
	OptionalAuth     func(http.Handler) http.Handler
	MaxBodyBytes     int64
	MaxPhotoBytes    int64
	UploadDir        string
	UploadURL        string
	OpenRouterAPIKey string
	OpenRouterModel  string
}

type Module struct {
	Handlers   *api.Handlers
	Repository domain.Repository
}

func New(opts Options) *Module {
	repo := infra.NewPgRepository(opts.Pool)
	read := infra.NewPgReadModel(opts.Pool)
	owners := infra.NewPgOwnerProvider(opts.Pool)
	photos := infra.NewPgPhotoRepository(opts.Pool)
	storage := infra.NewLocalPhotoStorage(opts.UploadDir, opts.UploadURL)
	moderationLog := infra.NewPgModerationLogRepository(opts.Pool)
	photoBytes := infra.NewLocalPhotoBytesLoader(opts.UploadDir, opts.UploadURL)
	moderationProvider := infra.NewOpenRouterModerationProvider(opts.OpenRouterAPIKey, opts.OpenRouterModel)

	bus := opts.Bus
	if bus == nil {
		bus = events.NopPublisher{}
	}

	moderate := app.NewModerateItemHandler(
		repo, photos, moderationProvider, moderationLog, photoBytes, opts.Tx, opts.Clock, bus,
	)

	handlers := api.NewHandlers(api.Deps{
		Items:         repo,
		Owners:        owners,
		CreateItem:    app.NewCreateItemHandler(repo, opts.Tx, opts.Clock),
		UpdateItem:    app.NewUpdateItemHandler(repo, photos, opts.Tx, opts.Clock, bus, moderate),
		ChangeStatus:  app.NewChangeStatusHandler(repo, photos, opts.Tx, opts.Clock, bus, moderate),
		GetItem:       app.NewGetItemHandler(repo, moderationLog, owners),
		ViewItem:      app.NewViewItemHandler(repo, opts.Clock, bus),
		ListItems:     app.NewListItemsHandler(read, owners),
		AddPhoto:      app.NewAddPhotoHandler(repo, photos, storage, opts.Tx, opts.Clock, moderate),
		ListPhotos:    app.NewListPhotosHandler(photos),
		DeletePhoto:   app.NewDeletePhotoHandler(repo, photos, storage, opts.Tx, opts.Clock, moderate),
		Validator:     opts.Validator,
		Authenticate:  opts.Authenticate,
		OptionalAuth:  opts.OptionalAuth,
		MaxBodyBytes:  opts.MaxBodyBytes,
		MaxPhotoBytes: opts.MaxPhotoBytes,
	})

	return &Module{Handlers: handlers, Repository: repo}
}

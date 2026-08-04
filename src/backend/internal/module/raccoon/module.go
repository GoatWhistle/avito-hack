package raccoon

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/avito-hack/backend/internal/module/raccoon/api"
	"github.com/avito-hack/backend/internal/module/raccoon/app"
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
	stubReader := app.NewStubRaccoonStateReader()
	stubPromo := app.NewStubPromocodeGeneratorAdapter()
	stubNotifs := app.NewStubNotificationAdapter()
	stubCore := app.NewStubCoreRaccoonService()

	getProfile := app.NewGetRaccoonProfileUseCase(stubReader)
	claimReward := app.NewClaimRewardUseCase(stubPromo, stubNotifs)

	handlers := api.NewHandlers(api.Deps{
		GetProfile:   getProfile,
		ClaimReward:  claimReward,
		Core:         stubCore,
		Validator:    opts.Validator,
		Authenticate: opts.Authenticate,
		OptionalAuth: opts.OptionalAuth,
		MaxBodyBytes: opts.MaxBodyBytes,
	})

	return &Module{Handlers: handlers}
}

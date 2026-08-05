package raccoon

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	petapp "github.com/avito-hack/backend/internal/module/pet/app"
	petinfra "github.com/avito-hack/backend/internal/module/pet/infra"
	"github.com/avito-hack/backend/internal/module/raccoon/api"
	"github.com/avito-hack/backend/internal/module/raccoon/app"
	"github.com/avito-hack/backend/internal/module/raccoon/infra"
)

type Options struct {
	Pool         *pgxpool.Pool
	Pets         *petapp.Service
	Rewards      *petapp.RewardService
	Validator    interface{ Struct(dst any) error }
	Authenticate func(http.Handler) http.Handler
	MaxBodyBytes int64
}

type Module struct {
	Handlers *api.Handlers
}

func New(opts Options) *Module {
	badges := infra.NewBadgeAdapter(petinfra.NewPgBadgeRepository(opts.Pool))

	return &Module{Handlers: api.NewHandlers(api.Deps{
		GetProfile:   app.NewGetRaccoonProfileUseCase(infra.NewPetAdapter(opts.Pets), badges),
		ClaimReward:  app.NewClaimRewardUseCase(infra.NewRewardAdapter(opts.Rewards)),
		Validator:    opts.Validator,
		Authenticate: opts.Authenticate,
		MaxBodyBytes: opts.MaxBodyBytes,
	})}
}

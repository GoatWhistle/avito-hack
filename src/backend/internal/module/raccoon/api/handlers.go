package api

import (
	"net/http"

	"github.com/avito-hack/backend/internal/module/raccoon/app"
	"github.com/avito-hack/backend/internal/shared/httpx"
)

type Deps struct {
	GetProfile   *app.GetRaccoonProfileUseCase
	ClaimReward  *app.ClaimRewardUseCase
	Validator    httpx.Validator
	Authenticate func(http.Handler) http.Handler
	MaxBodyBytes int64
}

type Handlers struct {
	deps    Deps
	decoder httpx.Decoder
}

func NewHandlers(deps Deps) *Handlers {
	return &Handlers{deps: deps, decoder: httpx.NewDecoder(deps.Validator, deps.MaxBodyBytes)}
}

func (h *Handlers) decode(w http.ResponseWriter, r *http.Request, dst any) bool {
	return h.decoder.Decode(w, r, dst)
}

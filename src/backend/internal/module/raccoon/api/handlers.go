package api

import (
	"net/http"

	"github.com/avito-hack/backend/internal/module/raccoon/app"
	"github.com/avito-hack/backend/internal/shared/apierr"
	"github.com/avito-hack/backend/internal/shared/httpx"
)

type validator interface {
	Struct(dst any) error
}

type Deps struct {
	GetProfile   *app.GetRaccoonProfileUseCase
	ClaimReward  *app.ClaimRewardUseCase
	Core         app.CoreRaccoonService
	Validator    validator
	Authenticate func(http.Handler) http.Handler
	OptionalAuth func(http.Handler) http.Handler
	MaxBodyBytes int64
}

type Handlers struct {
	deps Deps
}

func NewHandlers(deps Deps) *Handlers {
	return &Handlers{deps: deps}
}

func (h *Handlers) decode(w http.ResponseWriter, r *http.Request, dst any) bool {
	if err := httpx.DecodeJSON(w, r, h.deps.MaxBodyBytes, dst); err != nil {
		apierr.WriteBadRequest(w, r, err.Error())
		return false
	}

	if err := h.deps.Validator.Struct(dst); err != nil {
		apierr.Write(w, r, err)
		return false
	}

	return true
}

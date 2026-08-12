package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/avito-hack/backend/internal/module/games/app"
	"github.com/avito-hack/backend/internal/module/games/domain"
	"github.com/avito-hack/backend/internal/shared/apierr"
	"github.com/avito-hack/backend/internal/shared/httpx"
	"github.com/avito-hack/backend/internal/shared/publicid"
)

type Deps struct {
	ListGames    *app.ListGamesHandler
	GetState     *app.GetStateHandler
	StartRound   *app.StartRoundHandler
	Guess        *app.GuessHandler
	ClaimReward  *app.ClaimRewardHandler
	Clock        app.Clock
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

func (h *Handlers) today() domain.Day {
	if h.deps.Clock == nil {
		return domain.Day{}
	}

	return domain.DayOf(h.deps.Clock.Now())
}

func (h *Handlers) resolveRoundID(w http.ResponseWriter, r *http.Request, name string) (string, bool) {
	roundID := chi.URLParam(r, name)

	if !publicid.IsValid(roundID) {
		apierr.Write(w, r, domain.ErrRoundNotFound)

		return "", false
	}

	return roundID, true
}

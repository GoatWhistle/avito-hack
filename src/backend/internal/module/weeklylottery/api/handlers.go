package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/avito-hack/backend/internal/module/weeklylottery/app"
	"github.com/avito-hack/backend/internal/module/weeklylottery/domain"
	"github.com/avito-hack/backend/internal/shared/apierr"
	"github.com/avito-hack/backend/internal/shared/auth"
	"github.com/avito-hack/backend/internal/shared/httpx"
	"github.com/avito-hack/backend/internal/shared/publicid"
)

type Deps struct {
	GetState     *app.GetStateHandler
	ListPrizes   *app.ListPrizesHandler
	Start        *app.StartHandler
	Reveal       *app.RevealHandler
	Authenticate func(http.Handler) http.Handler
}

type Handlers struct {
	deps Deps
}

func NewHandlers(deps Deps) *Handlers {
	return &Handlers{deps: deps}
}

func (h *Handlers) RegisterRoutes(r chi.Router) {
	r.Route("/weekly-lottery", func(r chi.Router) {
		r.Use(h.deps.Authenticate)
		r.Get("/prizes", h.ListPrizes)
		r.Get("/state", h.GetState)
		r.Post("/runs", h.Start)
		r.Post("/runs/{rid}/slots/{slot}/reveal", h.Reveal)
	})
}

// ListPrizes возвращает публичные условия призов без внутренних вероятностей.
// @Id listWeeklyLotteryPrizes
// @Summary Призы еженедельной лотереи
// @Tags WeeklyLottery
// @Produce json
// @Success 200 {array} PrizeCatalogResponse
// @Failure 401 {object} apierr.ErrorEnvelope
// @Security bearerAuth
// @Router /api/v1/weekly-lottery/prizes [get]
func (h *Handlers) ListPrizes(w http.ResponseWriter, r *http.Request) {
	if _, err := auth.ActorFrom(r.Context()); err != nil {
		apierr.Write(w, r, err)
		return
	}

	httpx.OK(w, toPrizeCatalogResponse(h.deps.ListPrizes.Handle()))
}

// GetState возвращает состояние еженедельной лотереи.
// @Id getWeeklyLotteryState
// @Summary Состояние еженедельной лотереи
// @Tags WeeklyLottery
// @Produce json
// @Success 200 {object} StateResponse
// @Failure 401 {object} apierr.ErrorEnvelope
// @Failure 500 {object} apierr.ErrorEnvelope
// @Security bearerAuth
// @Router /api/v1/weekly-lottery/state [get]
func (h *Handlers) GetState(w http.ResponseWriter, r *http.Request) {
	actor, err := auth.ActorFrom(r.Context())
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	view, err := h.deps.GetState.Handle(r.Context(), actor.ID)
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	httpx.OK(w, toStateResponse(view))
}

// Start создаёт или возвращает розыгрыш текущей недели.
// @Id startWeeklyLottery
// @Summary Начать еженедельную лотерею
// @Tags WeeklyLottery
// @Produce json
// @Success 200 {object} RunResponse "Уже существующий розыгрыш"
// @Success 201 {object} RunResponse "Новый розыгрыш"
// @Failure 401 {object} apierr.ErrorEnvelope
// @Failure 500 {object} apierr.ErrorEnvelope
// @Security bearerAuth
// @Router /api/v1/weekly-lottery/runs [post]
func (h *Handlers) Start(w http.ResponseWriter, r *http.Request) {
	actor, err := auth.ActorFrom(r.Context())
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	view, created, err := h.deps.Start.Handle(r.Context(), actor.ID)
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	response := toRunResponse(view)
	if created {
		httpx.Created(w, response)
		return
	}
	httpx.OK(w, response)
}

// Reveal открывает одну ячейку, не раскрывая остальные.
// @Id revealWeeklyLotterySlot
// @Summary Открыть ячейку еженедельной лотереи
// @Tags WeeklyLottery
// @Produce json
// @Param rid path string true "Публичный идентификатор розыгрыша"
// @Param slot path int true "Индекс ячейки от 0 до 8"
// @Success 200 {object} RevealResponse
// @Failure 400 {object} apierr.ErrorEnvelope
// @Failure 401 {object} apierr.ErrorEnvelope
// @Failure 404 {object} apierr.ErrorEnvelope
// @Failure 409 {object} apierr.ErrorEnvelope
// @Failure 500 {object} apierr.ErrorEnvelope
// @Security bearerAuth
// @Router /api/v1/weekly-lottery/runs/{rid}/slots/{slot}/reveal [post]
func (h *Handlers) Reveal(w http.ResponseWriter, r *http.Request) {
	actor, err := auth.ActorFrom(r.Context())
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	runID := chi.URLParam(r, "rid")
	if !publicid.IsValid(runID) {
		apierr.Write(w, r, domain.ErrRunNotFound)
		return
	}
	slot, err := strconv.Atoi(chi.URLParam(r, "slot"))
	if err != nil {
		apierr.Write(w, r, domain.ErrInvalidSlot())
		return
	}

	result, err := h.deps.Reveal.Handle(r.Context(), app.RevealCommand{
		UserID: actor.ID, RunID: runID, Slot: slot,
	})
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	httpx.OK(w, RevealResponse{
		Index: result.Index, Symbol: string(result.Symbol), Run: toRunResponse(result.Run),
	})
}

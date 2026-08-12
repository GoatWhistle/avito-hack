package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/avito-hack/backend/internal/module/games/app"
	"github.com/avito-hack/backend/internal/shared/apierr"
	"github.com/avito-hack/backend/internal/shared/auth"
	"github.com/avito-hack/backend/internal/shared/httpx"
)

// @Id listGames
// @Summary Список мини-игр
// @Description Все зарегистрированные мини-игры с прогрессом текущего пользователя:
// @Description выполнен ли дневной норматив и на какой неделе серия. Новая игра
// @Description появляется здесь автоматически после регистрации в реестре.
// @Tags Games
// @Produce json
// @Success 200 {array} GameResponse "Мини-игры"
// @Failure 401 {object} apierr.ErrorEnvelope
// @Failure 500 {object} apierr.ErrorEnvelope
// @Security bearerAuth
// @Router /api/v1/games [get]
func (h *Handlers) List(w http.ResponseWriter, r *http.Request) {
	actor, err := auth.ActorFrom(r.Context())
	if err != nil {
		apierr.Write(w, r, err)

		return
	}

	views, err := h.deps.ListGames.Handle(r.Context(), app.ListGamesQuery{
		UserID:    actor.ID,
		ClientDay: h.today(),
	})
	if err != nil {
		apierr.Write(w, r, err)

		return
	}

	httpx.OK(w, toGameResponses(views))
}

// @Id getGameState
// @Summary Состояние мини-игры
// @Description Недельная серия, дневной прогресс и активный раунд, если он есть.
// @Description Активный раунд отдаётся вместе с текущим вопросом, поэтому клиент
// @Description может продолжить игру после перезагрузки страницы. Скрытая цена
// @Description в вопрос никогда не попадает.
// @Tags Games
// @Produce json
// @Param slug path string true "Идентификатор мини-игры" example(moreless)
// @Success 200 {object} GameStateResponse "Состояние игры"
// @Failure 401 {object} apierr.ErrorEnvelope
// @Failure 404 {object} apierr.ErrorEnvelope
// @Failure 500 {object} apierr.ErrorEnvelope
// @Security bearerAuth
// @Router /api/v1/games/{slug}/state [get]
func (h *Handlers) GetState(w http.ResponseWriter, r *http.Request) {
	actor, err := auth.ActorFrom(r.Context())
	if err != nil {
		apierr.Write(w, r, err)

		return
	}

	view, err := h.deps.GetState.Handle(r.Context(), app.GetStateQuery{
		UserID:    actor.ID,
		GameSlug:  chi.URLParam(r, "slug"),
		ClientDay: h.today(),
	})
	if err != nil {
		apierr.Write(w, r, err)

		return
	}

	httpx.OK(w, toStateResponse(view))
}

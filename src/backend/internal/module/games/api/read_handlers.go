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
// @Description Все зарегистрированные мини-игры с прогрессом текущего пользователя.
// @Description Недельная серия одна на все мини-игры и лежит в корне ответа, а у каждой
// @Description игры остаётся только признак `daily_done` — играли ли в неё сегодня.
// @Description Новая игра появляется здесь автоматически после регистрации в реестре.
// @Tags Games
// @Produce json
// @Success 200 {object} GameListResponse "Мини-игры"
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

	view, err := h.deps.ListGames.Handle(r.Context(), app.ListGamesQuery{
		UserID:    actor.ID,
		ClientDay: h.today(),
	})
	if err != nil {
		apierr.Write(w, r, err)

		return
	}

	httpx.OK(w, toGameListResponse(view))
}

// @Id getGameState
// @Summary Состояние мини-игры
// @Description Общая недельная серия, дневной прогресс этой игры и активный раунд, если он есть.
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

// @Id getGameHiddenPhoto
// @Summary Фото скрытой стороны раунда
// @Description Отдаёт фотографию правого объявления по одноразовому подписанному токену
// @Description из вопроса. Публичный идентификатор скрытого объявления клиенту не
// @Description раскрывается, поэтому его цену нельзя подсмотреть через каталог.
// @Tags Games
// @Produce json
// @Param token path string true "Подписанный токен фото из поля prompt.right.photo_token"
// @Success 302 "Редирект на файл фотографии"
// @Failure 404 {object} apierr.ErrorEnvelope
// @Security bearerAuth
// @Security []
// @Router /api/v1/games/photo/{token} [get]
func (h *Handlers) HiddenPhoto(w http.ResponseWriter, r *http.Request) {
	actor, _ := auth.ActorFrom(r.Context()) //nolint:errcheck // the signed token carries the round binding

	url, err := h.deps.HiddenPhoto.Handle(r.Context(), app.HiddenPhotoQuery{
		UserID: actor.ID,
		Token:  chi.URLParam(r, "token"),
	})
	if err != nil {
		apierr.Write(w, r, err)

		return
	}

	w.Header().Set("Cache-Control", "private, no-store")
	http.Redirect(w, r, url, http.StatusFound)
}

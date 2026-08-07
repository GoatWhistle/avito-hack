package api

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/module/favorite/app"
	"github.com/avito-hack/backend/internal/shared/apierr"
	"github.com/avito-hack/backend/internal/shared/auth"
	"github.com/avito-hack/backend/internal/shared/httpx"
)

type Deps struct {
	AddFavorite    *app.AddFavoriteHandler
	RemoveFavorite *app.RemoveFavoriteHandler
	ListFavorites  *app.ListFavoritesHandler
	Authenticate   func(http.Handler) http.Handler
}

type Handlers struct {
	deps Deps
}

func NewHandlers(deps Deps) *Handlers {
	return &Handlers{deps: deps}
}

// @Id addFavorite
// @Summary Добавить в избранное
// @Description Идемпотентна: повторное добавление того же объявления не создаёт дубль
// @Description и возвращает те же 204 (ограничение уникальности пары пользователь-объявление
// @Description обрабатывается как no-op). Тело ответа пустое.
// @Description
// @Description Побочный эффект: первое добавление порождает доменное событие,
// @Description за которое питомцу асинхронно начисляется опыт. Повторное добавление
// @Description опыт не приносит.
// @Tags Favorites
// @Produce json
// @Param id path string true "Идентификатор объявления" format(uuid)
// @Success 204 "Добавлено в избранное"
// @Failure 400 {object} apierr.ErrorEnvelope
// @Failure 401 {object} apierr.ErrorEnvelope
// @Failure 404 {object} apierr.ErrorEnvelope
// @Failure 500 {object} apierr.ErrorEnvelope
// @Security bearerAuth
// @Router /api/v1/items/{id}/favorite [post]
func (h *Handlers) Add(w http.ResponseWriter, r *http.Request) {
	actor, itemID, ok := h.actorAndItem(w, r)
	if !ok {
		return
	}

	if err := h.deps.AddFavorite.Handle(r.Context(), app.AddFavoriteCommand{
		UserID: actor.ID,
		ItemID: itemID,
	}); err != nil {
		apierr.Write(w, r, err)
		return
	}

	httpx.NoContent(w)
}

// @Id removeFavorite
// @Summary Убрать из избранного
// @Description Идемпотентна: удаление отсутствующей записи также даёт 204.
// @Description Тело ответа пустое. Опыт за удаление не начисляется и не отнимается.
// @Tags Favorites
// @Produce json
// @Param id path string true "Идентификатор объявления" format(uuid)
// @Success 204 "Убрано из избранного"
// @Failure 400 {object} apierr.ErrorEnvelope
// @Failure 401 {object} apierr.ErrorEnvelope
// @Failure 500 {object} apierr.ErrorEnvelope
// @Security bearerAuth
// @Router /api/v1/items/{id}/favorite [delete]
func (h *Handlers) Remove(w http.ResponseWriter, r *http.Request) {
	actor, itemID, ok := h.actorAndItem(w, r)
	if !ok {
		return
	}

	if err := h.deps.RemoveFavorite.Handle(r.Context(), app.RemoveFavoriteCommand{
		UserID: actor.ID,
		ItemID: itemID,
	}); err != nil {
		apierr.Write(w, r, err)
		return
	}

	httpx.NoContent(w)
}

// @Id listFavorites
// @Summary Список избранного
// @Description Избранные объявления владельца токена, свежие сверху.
// @Description Пагинация — временной курсор по времени добавления в избранное.
// @Description Поле `photo_url` содержит первую фотографию объявления и опускается,
// @Description если фотографий нет.
// @Tags Favorites
// @Produce json
// @Param limit query int false "Размер страницы. По умолчанию 20, максимум 100." default(20) maximum(100)
// @Param cursor query string false "Временной курсор из `next_cursor` предыдущей страницы."
// @Success 200 {object} FavoriteListResponse "Страница избранного"
// @Failure 400 {object} apierr.ErrorEnvelope
// @Failure 401 {object} apierr.ErrorEnvelope
// @Failure 500 {object} apierr.ErrorEnvelope
// @Security bearerAuth
// @Router /api/v1/favorites [get]
func (h *Handlers) List(w http.ResponseWriter, r *http.Request) {
	actor, err := auth.ActorFrom(r.Context())
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	page, err := httpx.PageFromRequest(r)
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	result, err := h.deps.ListFavorites.Handle(r.Context(), app.ListFavoritesQuery{
		UserID: actor.ID,
		Cursor: page.Cursor,
		Limit:  page.Limit,
	})
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	httpx.OK(w, httpx.NewListResponse(toFavoriteListResponse(result.Items), result.NextCursor))
}

func (h *Handlers) actorAndItem(w http.ResponseWriter, r *http.Request) (auth.Actor, uuid.UUID, bool) {
	actor, err := auth.ActorFrom(r.Context())
	if err != nil {
		apierr.Write(w, r, err)

		return auth.Actor{}, uuid.Nil, false
	}

	itemID, err := httpx.UUIDParam(r, "id")
	if err != nil {
		apierr.Write(w, r, err)

		return auth.Actor{}, uuid.Nil, false
	}

	return actor, itemID, true
}

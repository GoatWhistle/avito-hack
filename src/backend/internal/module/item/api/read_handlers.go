package api

import (
	"net/http"

	"github.com/avito-hack/backend/internal/module/item/app"
	"github.com/avito-hack/backend/internal/module/item/domain"
	"github.com/avito-hack/backend/internal/shared/apierr"
	"github.com/avito-hack/backend/internal/shared/auth"
	"github.com/avito-hack/backend/internal/shared/domainerr"
	"github.com/avito-hack/backend/internal/shared/httpx"
)

// @Id getItem
// @Summary Объявление по идентификатору
// @Description Анонимный доступ разрешён. Чужой черновик или архив отдаётся как 404,
// @Description а не 403, чтобы не раскрывать существование записи.
// @Tags Items
// @Produce json
// @Param id path string true "Идентификатор объявления" format(uuid)
// @Success 200 {object} itemResponse "Объявление"
// @Failure 400 {object} apierr.ErrorEnvelope
// @Failure 404 {object} apierr.ErrorEnvelope
// @Failure 500 {object} apierr.ErrorEnvelope
// @Security bearerAuth
// @Security []
// @Router /api/v1/items/{id} [get]
func (h *Handlers) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.UUIDParam(r, "id")
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	actor, _ := auth.ActorFrom(r.Context()) //nolint:errcheck // anonymous access is allowed, the actor is optional

	item, err := h.deps.GetItem.Handle(r.Context(), app.GetItemQuery{ItemID: id, Actor: actor})
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	httpx.OK(w, toItemResponse(item))
}

// @Id listItems
// @Summary Публичный список объявлений
// @Description Анонимный доступ разрешён (`OptionalAuth`): токен читается, если передан,
// @Description но не требуется. Анонимно и для чужого пользователя видны только
// @Description объявления в статусах `published` и `sold`; владелец дополнительно видит
// @Description собственные черновики.
// @Description
// @Description Пагинация — временной курсор. Побочных эффектов нет.
// @Tags Items
// @Produce json
// @Param limit query int false "Размер страницы. По умолчанию 20, максимум 100." default(20) maximum(100)
// @Param cursor query string false "Временной курсор из `next_cursor` предыдущей страницы."
// @Param status query string false "Фильтр по статусу." Enums(draft, moderation, published, sold, archived)
// @Param owner_id query string false "Фильтр по владельцу. Невалидный UUID даёт 400." format(uuid)
// @Param search query string false "Полнотекстовый поиск по заголовку."
// @Success 200 {object} ItemListResponse "Страница объявлений"
// @Failure 400 {object} apierr.ErrorEnvelope
// @Failure 500 {object} apierr.ErrorEnvelope
// @Security bearerAuth
// @Security []
// @Router /api/v1/items [get]
func (h *Handlers) List(w http.ResponseWriter, r *http.Request) {
	page, err := httpx.PageFromRequest(r)
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	status, err := parseStatusQuery(r)
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	ownerID, _, err := httpx.OptionalUUIDQuery(r, "owner_id")
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	result, err := h.deps.ListItems.Handle(r.Context(), app.ListItemsQuery{
		Status:  status,
		OwnerID: ownerID,
		Search:  r.URL.Query().Get("search"),
		Cursor:  page.Cursor,
		Limit:   page.Limit,
	})
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	httpx.OK(w, httpx.NewListResponse(toListResponse(result.Items), result.NextCursor))
}

// @Id listMyItems
// @Summary Мои объявления
// @Description Список объявлений владельца токена во всех статусах, включая черновики
// @Description и архив. Фильтр `owner_id` недоступен — владелец подставляется из токена.
// @Tags Items
// @Produce json
// @Param limit query int false "Размер страницы. По умолчанию 20, максимум 100." default(20) maximum(100)
// @Param cursor query string false "Временной курсор из `next_cursor` предыдущей страницы."
// @Param status query string false "Фильтр по статусу." Enums(draft, moderation, published, sold, archived)
// @Success 200 {object} ItemListResponse "Страница объявлений"
// @Failure 400 {object} apierr.ErrorEnvelope
// @Failure 401 {object} apierr.ErrorEnvelope
// @Failure 500 {object} apierr.ErrorEnvelope
// @Security bearerAuth
// @Router /api/v1/items/mine [get]
func (h *Handlers) ListMine(w http.ResponseWriter, r *http.Request) {
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

	status, err := parseStatusQuery(r)
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	result, err := h.deps.ListItems.Handle(r.Context(), app.ListItemsQuery{
		Status:  status,
		OwnerID: actor.ID,
		Cursor:  page.Cursor,
		Limit:   page.Limit,
	})
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	httpx.OK(w, httpx.NewListResponse(toListResponse(result.Items), result.NextCursor))
}

func parseStatusQuery(r *http.Request) (domain.Status, error) {
	raw := r.URL.Query().Get("status")
	if raw == "" {
		return "", nil
	}

	status := domain.Status(raw)
	if !status.Valid() {
		return "", domainerr.NewInvalid("status", "allowed values: draft moderation published archived")
	}

	return status, nil
}

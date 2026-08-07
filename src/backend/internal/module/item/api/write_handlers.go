package api

import (
	"net/http"

	"github.com/avito-hack/backend/internal/module/item/app"
	"github.com/avito-hack/backend/internal/shared/apierr"
	"github.com/avito-hack/backend/internal/shared/auth"
	"github.com/avito-hack/backend/internal/shared/httpx"
)

// @Id createItem
// @Summary Создание объявления
// @Description Создаёт объявление в статусе `draft` от имени владельца токена.
// @Description
// @Description Не идемпотентна — каждый вызов создаёт новую запись.
// @Description
// @Description Побочный эффект: публикуется доменное событие, на которое подписан
// @Description игровой модуль. Опыт за создание объявления начисляется асинхронно
// @Description (через Kafka либо через in-process шину, если Kafka выключена),
// @Description поэтому в ответе изменение XP не отражается — его нужно получать
// @Description из WebSocket-события `xp.gained` или повторным `GET /api/v1/pet`.
// @Tags Items
// @Accept json
// @Produce json
// @Param request body createItemRequest true "Данные объявления"
// @Success 201 {object} itemResponse "Объявление создано"
// @Failure 400 {object} apierr.ErrorEnvelope
// @Failure 401 {object} apierr.ErrorEnvelope
// @Failure 500 {object} apierr.ErrorEnvelope
// @Security bearerAuth
// @Router /api/v1/items [post]
func (h *Handlers) Create(w http.ResponseWriter, r *http.Request) {
	actor, err := auth.ActorFrom(r.Context())
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	var req createItemRequest
	if !h.decode(w, r, &req) {
		return
	}

	item, err := h.deps.CreateItem.Handle(r.Context(), app.CreateItemCommand{
		OwnerID:     actor.ID,
		Title:       req.Title,
		Description: req.Description,
		PriceKopeks: req.PriceKopeks,
		Attributes:  req.Attributes,
	})
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	httpx.Created(w, toItemResponse(item))
}

// @Id updateItem
// @Summary Редактирование объявления
// @Description Частичное обновление: все поля тела необязательны, применяются только
// @Description переданные. Доступно только владельцу — чужое объявление даёт 403.
// @Description
// @Description Идемпотентна: повтор с тем же телом приводит к тому же состоянию.
// @Description Терминальные статусы (`sold`) редактировать нельзя — 409.
// @Tags Items
// @Accept json
// @Produce json
// @Param id path string true "Идентификатор объявления" format(uuid)
// @Param request body updateItemRequest true "Изменяемые поля"
// @Success 200 {object} itemResponse "Объявление обновлено"
// @Failure 400 {object} apierr.ErrorEnvelope
// @Failure 401 {object} apierr.ErrorEnvelope
// @Failure 403 {object} apierr.ErrorEnvelope
// @Failure 404 {object} apierr.ErrorEnvelope
// @Failure 409 {object} apierr.ErrorEnvelope
// @Failure 500 {object} apierr.ErrorEnvelope
// @Security bearerAuth
// @Router /api/v1/items/{id} [patch]
func (h *Handlers) Update(w http.ResponseWriter, r *http.Request) {
	actor, err := auth.ActorFrom(r.Context())
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	id, err := httpx.UUIDParam(r, "id")
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	var req updateItemRequest
	if !h.decode(w, r, &req) {
		return
	}

	item, err := h.deps.UpdateItem.Handle(r.Context(), app.UpdateItemCommand{
		ItemID:      id,
		ActorID:     actor.ID,
		Title:       req.Title,
		Description: req.Description,
		PriceKopeks: req.PriceKopeks,
		Attributes:  req.Attributes,
	})
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	httpx.OK(w, toItemResponse(item))
}

// @Id changeItemStatus
// @Summary Смена статуса объявления
// @Description Переводит объявление по конечному автомату статусов. Действие задаётся
// @Description полем `action`, а не целевым статусом.
// @Description
// @Description Разрешённые переходы:
// @Description
// @Description | Из | В |
// @Description |---|---|
// @Description | `draft` | `published`, `moderation`, `archived` |
// @Description | `moderation` | `published`, `draft`, `archived` |
// @Description | `published` | `sold`, `archived` |
// @Description | `sold` | — (терминальный) |
// @Description | `archived` | `draft` |
// @Description
// @Description Соответствие действий: `submit` → `moderation`, `publish` → `published`,
// @Description `sell` → `sold`, `archive` → `archived`, `restore` → `draft`.
// @Description
// @Description Недопустимый переход даёт 409. Доступно только владельцу.
// @Description
// @Description Побочный эффект: публикация (`publish`) и продажа (`sell`) порождают
// @Description доменные события, за которые игровой модуль асинхронно начисляет опыт
// @Description питомцу.
// @Tags Items
// @Accept json
// @Produce json
// @Param id path string true "Идентификатор объявления" format(uuid)
// @Param request body changeStatusRequest true "Действие"
// @Success 200 {object} itemResponse "Статус изменён"
// @Failure 400 {object} apierr.ErrorEnvelope
// @Failure 401 {object} apierr.ErrorEnvelope
// @Failure 403 {object} apierr.ErrorEnvelope
// @Failure 404 {object} apierr.ErrorEnvelope
// @Failure 409 {object} apierr.ErrorEnvelope "Переход не разрешён из текущего статуса"
// @Failure 500 {object} apierr.ErrorEnvelope
// @Security bearerAuth
// @Router /api/v1/items/{id}/status [post]
func (h *Handlers) ChangeStatus(w http.ResponseWriter, r *http.Request) {
	actor, err := auth.ActorFrom(r.Context())
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	id, err := httpx.UUIDParam(r, "id")
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	var req changeStatusRequest
	if !h.decode(w, r, &req) {
		return
	}

	item, err := h.deps.ChangeStatus.Handle(r.Context(), app.ChangeStatusCommand{
		ItemID: id,
		Actor:  actor,
		Action: app.StatusAction(req.Action),
	})
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	httpx.OK(w, toItemResponse(item))
}

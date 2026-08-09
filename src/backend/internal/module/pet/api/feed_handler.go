package api

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/module/pet/domain"
	"github.com/avito-hack/backend/internal/shared/apierr"
	"github.com/avito-hack/backend/internal/shared/httpx"
)

func (h *PetHandlers) notifyUpdated(userID uuid.UUID, pet *domain.Pet) {
	if h.deps.Notifier == nil || pet == nil {
		return
	}

	h.deps.Notifier.PetUpdated(userID, pet)
}

// @Id feedPet
// @Summary Покормить питомца
// @Description Повышает `satiety`. Тело запроса не требуется.
// @Description
// @Description Не идемпотентна, но эффект ограничен: сытость упирается в потолок 100,
// @Description а частота ограничена кулдауном в 5 часов. Повторный вызов внутри кулдауна
// @Description возвращает 200 с неизменённой сытостью, а не ошибку — клиенту не нужно
// @Description отличать «покормили» от «уже сыт».
// @Description
// @Description `feed_available_at` в ответе — момент, когда кормление снова станет
// @Description доступным (RFC3339). Значение `null` означает «доступно сейчас»;
// @Description после успешного кормления там будет время на 5 часов вперёд,
// @Description по нему клиент строит обратный отсчёт.
// @Description
// @Description Опыт за кормление не начисляется: это забота о питомце, а не целевое
// @Description действие на площадке. Влияние косвенное — сытость ниже 30 включает
// @Description штрафной множитель опыта, и кормление его снимает.
// @Description
// @Description Побочный эффект: всем активным WebSocket-соединениям пользователя
// @Description рассылается `pet.updated`.
// @Tags Pet
// @Produce json
// @Success 200 {object} petPayload "Обновлённое состояние питомца"
// @Failure 401 {object} apierr.ErrorEnvelope
// @Failure 404 {object} apierr.ErrorEnvelope
// @Failure 500 {object} apierr.ErrorEnvelope
// @Security bearerAuth
// @Router /api/v1/pet/actions/feed [post]
func (h *PetHandlers) Feed(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorOf(w, r)
	if !ok {
		return
	}

	view, err := h.deps.Service.FeedView(r.Context(), actor.ID)
	if err != nil {
		apierr.Write(w, r, mapPetError(err))

		return
	}

	h.notifyUpdated(actor.ID, view.Pet)

	httpx.OK(w, toStatePayload(view))
}

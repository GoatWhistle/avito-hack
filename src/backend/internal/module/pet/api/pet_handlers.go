package api

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/module/pet/app"
	"github.com/avito-hack/backend/internal/module/pet/domain"
	"github.com/avito-hack/backend/internal/shared/apierr"
	"github.com/avito-hack/backend/internal/shared/auth"
	"github.com/avito-hack/backend/internal/shared/httpx"
)

type petStateService interface {
	StateView(ctx context.Context, userID uuid.UUID) (app.StateView, error)
	FeedView(ctx context.Context, userID uuid.UUID) (app.StateView, error)
	Stroke(ctx context.Context, userID uuid.UUID) (*domain.Pet, error)
	CheckIn(ctx context.Context, userID uuid.UUID) (app.ActionResult, error)
	Progress(ctx context.Context, userID uuid.UUID) (app.ProgressView, error)
}

type rewardGranter interface {
	GrantEligible(ctx context.Context, userID uuid.UUID) ([]domain.Reward, error)
}

type petUpdateNotifier interface {
	PetUpdated(userID uuid.UUID, pet *domain.Pet)
}

type PetDeps struct {
	Service      petStateService
	Rewards      rewardGranter
	Notifier     petUpdateNotifier
	Authenticate func(http.Handler) http.Handler
}

type PetHandlers struct {
	deps PetDeps
}

func NewPetHandlers(deps PetDeps) *PetHandlers {
	return &PetHandlers{deps: deps}
}

// @Id getPet
// @Summary Состояние питомца
// @Description Текущее состояние питомца владельца токена. Питомец создаётся
// @Description автоматически при регистрации, поэтому 404 здесь не ожидается.
// @Description
// @Description Побочный эффект чтения: перед выдачей применяется ленивое затухание
// @Description характеристик (`satiety`, `happiness`, `energy`) пропорционально времени,
// @Description прошедшему с `last_decay_time`. Поэтому два последовательных чтения
// @Description без действий пользователя могут вернуть разные значения характеристик
// @Description и производное поле `state`.
// @Description
// @Description Второй побочный эффект — автоматический чек-ин. Если за текущие сутки
// @Description по московскому времени чек-ина ещё не было, он засчитывается прямо здесь:
// @Description начисляется опыт, продлевается серия, выдаются награды. Тогда
// @Description `checkin_applied` равно `true`, а объект `checkin` содержит начисленный
// @Description опыт, уровень и параметры серии — этого достаточно, чтобы показать тост.
// @Description Повторные чтения в тот же день ничего не меняют и возвращают
// @Description `checkin_applied: false` без объекта `checkin`. Операция идемпотентна
// @Description в пределах суток: параллельные запросы сериализуются блокировкой строки
// @Description питомца, поэтому опыт начисляется ровно один раз.
// @Description Явный `POST /api/v1/checkin` продолжает работать и остаётся
// @Description единственным способом получить полный ответ `CheckInResult`.
// @Description
// @Description `feed_available_at` — момент, когда снова можно кормить (RFC3339).
// @Description `null` означает, что кормление доступно прямо сейчас.
// @Tags Pet
// @Produce json
// @Success 200 {object} petPayload "Состояние питомца"
// @Failure 401 {object} apierr.ErrorEnvelope
// @Failure 404 {object} apierr.ErrorEnvelope
// @Failure 500 {object} apierr.ErrorEnvelope
// @Security bearerAuth
// @Router /api/v1/pet [get]
func (h *PetHandlers) Get(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorOf(w, r)
	if !ok {
		return
	}

	view, err := h.deps.Service.StateView(r.Context(), actor.ID)
	if err != nil {
		apierr.Write(w, r, err)

		return
	}

	if view.CheckInApplied {
		h.grantEligible(r, actor.ID)
	}

	httpx.OK(w, toStatePayload(view))
}

// @Id strokePet
// @Summary Погладить питомца
// @Description Повышает `happiness`. Тело запроса не требуется.
// @Description
// @Description Не идемпотентна, но эффект ограничен: приросты счастья упираются
// @Description в потолок 100, а частота ограничена внутренним лимитом действий —
// @Description превышение даёт 409 (`action limit reached`).
// @Description
// @Description Побочный эффект: всем активным WebSocket-соединениям пользователя
// @Description рассылается `pet.updated`.
// @Tags Pet
// @Produce json
// @Success 200 {object} petPayload "Обновлённое состояние питомца"
// @Failure 401 {object} apierr.ErrorEnvelope
// @Failure 404 {object} apierr.ErrorEnvelope
// @Failure 409 {object} apierr.ErrorEnvelope "Лимит действий исчерпан"
// @Failure 500 {object} apierr.ErrorEnvelope
// @Security bearerAuth
// @Router /api/v1/pet/actions/stroke [post]
func (h *PetHandlers) Stroke(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorOf(w, r)
	if !ok {
		return
	}

	pet, err := h.deps.Service.Stroke(r.Context(), actor.ID)
	if err != nil {
		apierr.Write(w, r, err)

		return
	}

	httpx.OK(w, toPetPayload(pet))
}

// @Id checkIn
// @Summary Ежедневный чек-ин
// @Description Отмечает вход за текущий день, начисляет опыт и продлевает стрик.
// @Description Тело запроса не требуется.
// @Description
// @Description Идемпотентна в пределах суток: повторный чек-ин в тот же день
// @Description даёт 409 (`action has already been performed today`). Именно этот
// @Description конфликт — нормальный, ожидаемый ответ для клиента, который
// @Description не знает, был ли уже чек-ин.
// @Description
// @Description Эндпоинт сохранён для обратной совместимости. Клиенту он больше
// @Description не нужен: `GET /api/v1/pet` засчитывает чек-ин сам и сообщает об этом
// @Description полями `checkin_applied` и `checkin`. Оба пути делят одно состояние,
// @Description поэтому после автоматического чек-ина этот вызов вернёт 409.
// @Description
// @Description Логика стрика:
// @Description
// @Description - Вход на следующий день после предыдущего — стрик продолжается
// @Description (`continued: true`).
// @Description - Пропуск дня при наличии заморозки — она тратится
// @Description (`freeze_used: true`), стрик сохраняется. Заморозок не больше 3.
// @Description - Пропуск дня без заморозки — стрик обнуляется (`reset: true`).
// @Description - На рубежных значениях стрика начисляется бонус
// @Description (`milestone_bonus` > 0, `milestone_reached` — достигнутый рубеж).
// @Description
// @Description Побочные эффекты: начисление опыта, возможный рост уровня,
// @Description автоматическая выдача наград, для которых выполнились условия
// @Description (их идентификаторы — в `unlocked_rewards`). По WebSocket уходят
// @Description `xp.gained`, при росте уровня `level.up`, при выдаче наград
// @Description `reward.granted`, а также `streak.updated`.
// @Tags Pet
// @Produce json
// @Success 200 {object} checkInPayload "Чек-ин засчитан"
// @Failure 401 {object} apierr.ErrorEnvelope
// @Failure 404 {object} apierr.ErrorEnvelope
// @Failure 409 {object} apierr.ErrorEnvelope "Чек-ин за сегодня уже выполнен"
// @Failure 500 {object} apierr.ErrorEnvelope
// @Security bearerAuth
// @Router /api/v1/checkin [post]
func (h *PetHandlers) CheckIn(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorOf(w, r)
	if !ok {
		return
	}

	result, err := h.deps.Service.CheckIn(r.Context(), actor.ID)
	if err != nil {
		apierr.Write(w, r, mapPetError(err))

		return
	}

	h.grantEligible(r, actor.ID)

	httpx.OK(w, toCheckInPayload(result))
}

// @Id getProgress
// @Summary Прогресс уровня и стрика
// @Description Компактная выжимка для полосы прогресса: уровень, опыт, сколько
// @Description осталось до следующего уровня, стадия и стрик. Подмножество данных
// @Description `GET /api/v1/pet` — отдельный запрос нужен, чтобы не тянуть
// @Description характеристики питомца ради шапки интерфейса.
// @Description
// @Description На максимальном уровне `is_max_level` равно `true`, а `xp_to_next_level`
// @Description равно нулю.
// @Tags Pet
// @Produce json
// @Success 200 {object} progressPayload "Прогресс"
// @Failure 401 {object} apierr.ErrorEnvelope
// @Failure 404 {object} apierr.ErrorEnvelope
// @Failure 500 {object} apierr.ErrorEnvelope
// @Security bearerAuth
// @Router /api/v1/progress [get]
func (h *PetHandlers) Progress(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorOf(w, r)
	if !ok {
		return
	}

	progress, err := h.deps.Service.Progress(r.Context(), actor.ID)
	if err != nil {
		apierr.Write(w, r, err)

		return
	}

	httpx.OK(w, toProgressPayload(progress))
}

func (h *PetHandlers) grantEligible(r *http.Request, userID uuid.UUID) {
	if h.deps.Rewards == nil {
		return
	}

	//nolint:errcheck // reward granting is best-effort, the check-in itself already succeeded
	_, _ = h.deps.Rewards.GrantEligible(r.Context(), userID)
}

func actorOf(w http.ResponseWriter, r *http.Request) (auth.Actor, bool) {
	actor, err := auth.ActorFrom(r.Context())
	if err != nil {
		apierr.Write(w, r, err)

		return auth.Actor{}, false
	}

	return actor, true
}

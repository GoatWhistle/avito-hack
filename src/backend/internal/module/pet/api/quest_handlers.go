package api

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/module/pet/app"
	"github.com/avito-hack/backend/internal/shared/apierr"
	"github.com/avito-hack/backend/internal/shared/httpx"
)

type questService interface {
	Today(ctx context.Context, userID uuid.UUID) ([]app.QuestView, error)
	ClaimCompleted(ctx context.Context, userID uuid.UUID) ([]app.QuestView, error)
}

type QuestDeps struct {
	Quests       questService
	Authenticate func(http.Handler) http.Handler
}

type QuestHandlers struct {
	deps QuestDeps
}

func NewQuestHandlers(deps QuestDeps) *QuestHandlers {
	return &QuestHandlers{deps: deps}
}

// @Id listDailyQuests
// @Summary Ежедневные задания
// @Description Набор заданий на текущий день. Состав заданий детерминирован: он
// @Description зависит от пользователя и календарной даты (МСК), поэтому в течение
// @Description суток список не меняется между запросами.
// @Description
// @Description Прогресс не хранится отдельным счётчиком, а вычисляется из журнала
// @Description начислений опыта за текущие сутки, поэтому он не может разойтись
// @Description с реальными действиями. `completed` — задание выполнено,
// @Description `claimed` — награда за него уже начислена.
// @Tags Quests
// @Produce json
// @Success 200 {object} QuestListResponse "Задания на сегодня"
// @Failure 401 {object} apierr.ErrorEnvelope
// @Failure 500 {object} apierr.ErrorEnvelope
// @Security bearerAuth
// @Router /api/v1/quests [get]
func (h *QuestHandlers) Today(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorOf(w, r)
	if !ok {
		return
	}

	views, err := h.deps.Quests.Today(r.Context(), actor.ID)
	if err != nil {
		apierr.Write(w, r, mapPetError(err))

		return
	}

	httpx.OK(w, httpx.NewListResponse(toQuestItems(views), ""))
}

// @Id claimDailyQuests
// @Summary Забрать награды за выполненные задания
// @Description Начисляет опыт за все выполненные, но ещё не оплаченные задания дня.
// @Description Идемпотентна: награда за каждое задание выдаётся не более одного раза
// @Description в сутки, повторный вызов ничего не начисляет и возвращает то же
// @Description состояние списка.
// @Tags Quests
// @Produce json
// @Success 200 {object} QuestListResponse "Состояние заданий после начисления"
// @Failure 401 {object} apierr.ErrorEnvelope
// @Failure 500 {object} apierr.ErrorEnvelope
// @Security bearerAuth
// @Router /api/v1/quests/claim [post]
func (h *QuestHandlers) Claim(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorOf(w, r)
	if !ok {
		return
	}

	views, err := h.deps.Quests.ClaimCompleted(r.Context(), actor.ID)
	if err != nil {
		apierr.Write(w, r, mapPetError(err))

		return
	}

	httpx.OK(w, httpx.NewListResponse(toQuestItems(views), ""))
}

func toQuestItems(views []app.QuestView) []questItem {
	items := make([]questItem, 0, len(views))
	for _, view := range views {
		items = append(items, questItem{
			ID:        view.Quest.ID,
			Action:    string(view.Quest.Action),
			Target:    view.Quest.Target,
			Reward:    view.Quest.Reward,
			Current:   view.Current,
			Completed: view.Completed,
			Claimed:   view.Claimed,
		})
	}

	return items
}

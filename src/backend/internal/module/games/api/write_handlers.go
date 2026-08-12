package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/avito-hack/backend/internal/module/games/app"
	"github.com/avito-hack/backend/internal/shared/apierr"
	"github.com/avito-hack/backend/internal/shared/auth"
	"github.com/avito-hack/backend/internal/shared/httpx"
)

// @Id startGameRound
// @Summary Начать раунд мини-игры
// @Description Создаёт серверный раунд и отдаёт первый вопрос. Незавершённый предыдущий
// @Description раунд той же игры автоматически засчитывается как проигранный, поэтому
// @Description активный раунд всегда ровно один. Ответы хранятся только на сервере.
// @Tags Games
// @Produce json
// @Param slug path string true "Идентификатор мини-игры" example(moreless)
// @Success 201 {object} StartRoundResponse "Раунд начат"
// @Failure 401 {object} apierr.ErrorEnvelope
// @Failure 404 {object} apierr.ErrorEnvelope
// @Failure 500 {object} apierr.ErrorEnvelope
// @Security bearerAuth
// @Router /api/v1/games/{slug}/rounds [post]
func (h *Handlers) StartRound(w http.ResponseWriter, r *http.Request) {
	actor, err := auth.ActorFrom(r.Context())
	if err != nil {
		apierr.Write(w, r, err)

		return
	}

	result, err := h.deps.StartRound.Handle(r.Context(), app.StartRoundCommand{
		UserID:   actor.ID,
		GameSlug: chi.URLParam(r, "slug"),
	})
	if err != nil {
		apierr.Write(w, r, err)

		return
	}

	httpx.Created(w, StartRoundResponse{
		RoundID:      result.RoundID,
		Streak:       result.Streak,
		TargetStreak: result.TargetStreak,
		Prompt:       result.Prompt,
	})
}

// @Id guessGameRound
// @Summary Сделать ход в раунде
// @Description Применяет ход игрока. Формат поля `move` задаёт конкретная игра:
// @Description для `moreless` это `{"choice":"higher"}` или `{"choice":"lower"}`.
// @Description
// @Description Чужой раунд отдаётся как 404, а не 403. Повторный ход в уже завершённом
// @Description раунде даёт 409, поэтому переигровка невозможна. Ответ содержит `reveal`
// @Description с раскрытой правдой только за уже сыгранный ход.
// @Tags Games
// @Accept json
// @Produce json
// @Param slug path string true "Идентификатор мини-игры" example(moreless)
// @Param rid path string true "Публичный идентификатор раунда (display_id)"
// @Param request body GuessRequest true "Ход игрока"
// @Success 200 {object} GuessResponse "Ход учтён"
// @Failure 400 {object} apierr.ErrorEnvelope
// @Failure 401 {object} apierr.ErrorEnvelope
// @Failure 404 {object} apierr.ErrorEnvelope
// @Failure 409 {object} apierr.ErrorEnvelope
// @Failure 500 {object} apierr.ErrorEnvelope
// @Security bearerAuth
// @Router /api/v1/games/{slug}/rounds/{rid}/guess [post]
func (h *Handlers) Guess(w http.ResponseWriter, r *http.Request) {
	actor, err := auth.ActorFrom(r.Context())
	if err != nil {
		apierr.Write(w, r, err)

		return
	}

	roundID, ok := h.resolveRoundID(w, r, "rid")
	if !ok {
		return
	}

	var request GuessRequest
	if !h.decode(w, r, &request) {
		return
	}

	result, err := h.deps.Guess.Handle(r.Context(), app.GuessCommand{
		UserID:    actor.ID,
		GameSlug:  chi.URLParam(r, "slug"),
		RoundID:   roundID,
		Move:      request.Move,
		ClientDay: h.today(),
	})
	if err != nil {
		apierr.Write(w, r, err)

		return
	}

	httpx.OK(w, toGuessResponse(result))
}

// @Id claimGameReward
// @Summary Забрать награду за недельную серию
// @Description Выдаёт промокод, когда недельная серия достигла семи дней подряд.
// @Description После выдачи счётчик серии обнуляется и следующий цикл начинается заново,
// @Description поэтому повторный вызов до новой серии даёт 409.
// @Tags Games
// @Produce json
// @Param slug path string true "Идентификатор мини-игры" example(moreless)
// @Success 200 {object} ClaimGameRewardResponse "Промокод выдан"
// @Failure 401 {object} apierr.ErrorEnvelope
// @Failure 404 {object} apierr.ErrorEnvelope
// @Failure 409 {object} apierr.ErrorEnvelope
// @Failure 500 {object} apierr.ErrorEnvelope
// @Security bearerAuth
// @Router /api/v1/games/{slug}/reward/claim [post]
func (h *Handlers) ClaimReward(w http.ResponseWriter, r *http.Request) {
	actor, err := auth.ActorFrom(r.Context())
	if err != nil {
		apierr.Write(w, r, err)

		return
	}

	result, err := h.deps.ClaimReward.Handle(r.Context(), app.ClaimRewardCommand{
		UserID:   actor.ID,
		GameSlug: chi.URLParam(r, "slug"),
	})
	if err != nil {
		apierr.Write(w, r, err)

		return
	}

	httpx.OK(w, ClaimGameRewardResponse{Code: result.Code})
}

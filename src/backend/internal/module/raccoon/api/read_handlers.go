package api

import (
	"net/http"
	"time"

	petdomain "github.com/avito-hack/backend/internal/module/pet/domain"
	"github.com/avito-hack/backend/internal/module/raccoon/app"
	"github.com/avito-hack/backend/internal/shared/apierr"
	"github.com/avito-hack/backend/internal/shared/auth"
	"github.com/avito-hack/backend/internal/shared/httpx"
)

// @Id getRaccoonProfile
// @Summary Профиль енота
// @Description Альтернативное представление того же питомца, ориентированное на
// @Description экран профиля: уровень, опыт, остаток до следующего уровня, стрик,
// @Description стадия, состояние и заработанные бейджи.
// @Description
// @Description Данные берутся из того же сервиса, что и `GET /api/v1/pet`, но поля
// @Description названы иначе — `xp_to_next_level` вместо `next_level_xp`,
// @Description `current_streak` вместо `streak_days`.
// @Tags Raccoon
// @Produce json
// @Success 200 {object} RaccoonProfileResponse "Профиль енота"
// @Failure 401 {object} apierr.ErrorEnvelope
// @Failure 404 {object} apierr.ErrorEnvelope
// @Failure 500 {object} apierr.ErrorEnvelope
// @Security bearerAuth
// @Router /api/v1/raccoon/profile [get]
func (h *Handlers) GetProfile(w http.ResponseWriter, r *http.Request) {
	actor, err := auth.ActorFrom(r.Context())
	if err != nil {
		apierr.Write(w, r, err)

		return
	}

	profile, err := h.deps.GetProfile.Execute(r.Context(), actor.ID)
	if err != nil {
		apierr.Write(w, r, err)

		return
	}

	httpx.OK(w, RaccoonProfileResponse{
		ID:            petdomain.PublicToken(profile.ID),
		UserID:        petdomain.PublicToken(profile.UserID),
		Name:          profile.Name,
		Level:         profile.Level,
		XP:            profile.XP,
		XPToNextLevel: profile.XPToNextLevel,
		CurrentStreak: profile.CurrentStreak,
		Stage:         profile.Stage,
		State:         profile.State,
		Badges:        toBadgeResponses(profile.Badges),
	})
}

// @Id listBadges
// @Summary Бейджи пользователя
// @Description Плоский массив бейджей — без обёртки `{items}` и без пагинации.
// @Description Поле `earned_at` пусто (и опущено) у ещё не полученных бейджей;
// @Description формат времени — RFC 3339 в UTC.
// @Tags Raccoon
// @Produce json
// @Success 200 {array} BadgeResponse "Массив бейджей"
// @Failure 401 {object} apierr.ErrorEnvelope
// @Failure 500 {object} apierr.ErrorEnvelope
// @Security bearerAuth
// @Router /api/v1/badges [get]
func (h *Handlers) GetBadges(w http.ResponseWriter, r *http.Request) {
	actor, err := auth.ActorFrom(r.Context())
	if err != nil {
		apierr.Write(w, r, err)

		return
	}

	badges, err := h.deps.GetProfile.ListBadges(r.Context(), actor.ID)
	if err != nil {
		apierr.Write(w, r, err)

		return
	}

	httpx.OK(w, toBadgeResponses(badges))
}

// @Id claimReward
// @Summary Получить промокод награды
// @Description Вариант активации награды из модуля енота: идентификатор награды
// @Description передаётся в теле, а не в пути, и в ответе поле называется
// @Description `promocode`, а не `code`.
// @Description
// @Description Функционально дублирует `POST /api/v1/rewards/{id}/activate`
// @Description и работает поверх того же сервиса наград. Для нового кода
// @Description предпочтительна версия с идентификатором в пути.
// @Description
// @Description Не идемпотентна: повторный вызов для уже активированной награды
// @Description даёт 409.
// @Tags Raccoon
// @Accept json
// @Produce json
// @Param request body ClaimRewardRequest true "Идентификатор награды"
// @Success 200 {object} ClaimRewardResponse "Промокод выдан"
// @Failure 400 {object} apierr.ErrorEnvelope
// @Failure 401 {object} apierr.ErrorEnvelope
// @Failure 404 {object} apierr.ErrorEnvelope
// @Failure 409 {object} apierr.ErrorEnvelope
// @Failure 500 {object} apierr.ErrorEnvelope
// @Security bearerAuth
// @Router /api/v1/rewards/claim [post]
func (h *Handlers) ClaimReward(w http.ResponseWriter, r *http.Request) {
	actor, err := auth.ActorFrom(r.Context())
	if err != nil {
		apierr.Write(w, r, err)

		return
	}

	var req ClaimRewardRequest
	if !h.decode(w, r, &req) {
		return
	}

	result, err := h.deps.ClaimReward.Execute(r.Context(), actor.ID, req.RewardID)
	if err != nil {
		apierr.Write(w, r, err)

		return
	}

	httpx.OK(w, ClaimRewardResponse{RewardID: result.RewardID, Promocode: result.Promocode})
}

func toBadgeResponses(badges []app.BadgeView) []BadgeResponse {
	responses := make([]BadgeResponse, 0, len(badges))
	for _, badge := range badges {
		earnedAt := ""
		if badge.EarnedAt != nil {
			earnedAt = badge.EarnedAt.UTC().Format(time.RFC3339)
		}
		responses = append(responses, BadgeResponse{
			ID:          badge.ID,
			Name:        badge.Name,
			Description: badge.Description,
			Icon:        badge.IconURL,
			EarnedAt:    earnedAt,
			Current:     badge.Current,
			Target:      badge.Target,
		})
	}

	return responses
}

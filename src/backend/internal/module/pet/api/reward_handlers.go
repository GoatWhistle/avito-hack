package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/module/pet/app"
	"github.com/avito-hack/backend/internal/module/pet/domain"
	"github.com/avito-hack/backend/internal/shared/apierr"
	"github.com/avito-hack/backend/internal/shared/domainerr"
	"github.com/avito-hack/backend/internal/shared/httpx"
)

type rewardService interface {
	Catalog(ctx context.Context, userID uuid.UUID) ([]app.RewardCatalogEntry, error)
	Mine(ctx context.Context, userID uuid.UUID) ([]app.GrantedReward, error)
	Activate(ctx context.Context, userID uuid.UUID, rewardID string) (string, error)
}

type RewardDeps struct {
	Rewards      rewardService
	Authenticate func(http.Handler) http.Handler
}

type RewardHandlers struct {
	deps RewardDeps
}

func NewRewardHandlers(deps RewardDeps) *RewardHandlers {
	return &RewardHandlers{deps: deps}
}

// @Id listRewardCatalog
// @Summary Каталог наград
// @Description Полный каталог наград с отметкой прогресса текущего пользователя:
// @Description `unlocked` — условие выполнено, `claimed` — награда уже выдана,
// @Description `progress_current`/`progress_target` — прогресс к условию.
// @Description
// @Description Ответ обёрнут в `{items, next_cursor}` ради единообразия, но
// @Description пагинации нет: каталог отдаётся целиком, `next_cursor` всегда пуст
// @Description и потому отсутствует в JSON.
// @Tags Rewards
// @Produce json
// @Success 200 {object} RewardCatalogResponse "Каталог наград"
// @Failure 401 {object} apierr.ErrorEnvelope
// @Failure 500 {object} apierr.ErrorEnvelope
// @Security bearerAuth
// @Router /api/v1/rewards [get]
func (h *RewardHandlers) Catalog(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorOf(w, r)
	if !ok {
		return
	}

	entries, err := h.deps.Rewards.Catalog(r.Context(), actor.ID)
	if err != nil {
		apierr.Write(w, r, mapPetError(err))

		return
	}

	items := make([]rewardCatalogItem, 0, len(entries))
	for i := range entries {
		items = append(items, toCatalogItem(entries[i]))
	}

	httpx.OK(w, httpx.NewListResponse(items, ""))
}

// @Id listMyRewards
// @Summary Мои награды
// @Description Награды, выданные текущему пользователю. Поле `code` — промокод —
// @Description заполняется только после активации; до неё оно опущено.
// @Description `next_cursor` всегда пуст, пагинации нет.
// @Tags Rewards
// @Produce json
// @Success 200 {object} UserRewardListResponse "Выданные награды"
// @Failure 401 {object} apierr.ErrorEnvelope
// @Failure 500 {object} apierr.ErrorEnvelope
// @Security bearerAuth
// @Router /api/v1/rewards/my [get]
func (h *RewardHandlers) Mine(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorOf(w, r)
	if !ok {
		return
	}

	granted, err := h.deps.Rewards.Mine(r.Context(), actor.ID)
	if err != nil {
		apierr.Write(w, r, mapPetError(err))

		return
	}

	items := make([]myRewardItem, 0, len(granted))
	for i := range granted {
		items = append(items, toMyRewardItem(granted[i]))
	}

	httpx.OK(w, httpx.NewListResponse(items, ""))
}

// @Id activateReward
// @Summary Активация награды
// @Description Превращает выданную награду в промокод. Код детерминированно
// @Description выводится через HMAC от идентификаторов пользователя и награды,
// @Description поэтому одна и та же награда всегда даёт один и тот же код.
// @Description
// @Description Не идемпотентна на уровне ответа: первая активация даёт 200,
// @Description повторная — 409 (`reward has already been activated`).
// @Description
// @Description Ошибки:
// @Description
// @Description - награда не выдана пользователю — 400, `field: reward_id`;
// @Description - условие награды не выполнено — 400, `field: reward_id`;
// @Description - награда уже активирована — 409;
// @Description - награда просрочена — 409;
// @Description - награда неактивируемая по своей природе — 409.
// @Tags Rewards
// @Produce json
// @Param id path string true "Строковый идентификатор награды (например `streak_5`), НЕ UUID." minlength(1)
// @Success 200 {object} activateRewardResponse "Награда активирована"
// @Failure 400 {object} apierr.ErrorEnvelope "Награда не выдана или условие не выполнено"
// @Failure 401 {object} apierr.ErrorEnvelope
// @Failure 404 {object} apierr.ErrorEnvelope
// @Failure 409 {object} apierr.ErrorEnvelope "Награда уже активирована, просрочена или неактивируема"
// @Failure 500 {object} apierr.ErrorEnvelope
// @Security bearerAuth
// @Router /api/v1/rewards/{id}/activate [post]
func (h *RewardHandlers) Activate(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorOf(w, r)
	if !ok {
		return
	}

	rewardID := chi.URLParam(r, "id")
	if rewardID == "" {
		apierr.Write(w, r, domainerr.NewInvalid("id", "reward id is required"))

		return
	}

	code, err := h.deps.Rewards.Activate(r.Context(), actor.ID, rewardID)
	if err != nil {
		apierr.Write(w, r, mapPetError(err))

		return
	}

	httpx.OK(w, activateRewardResponse{
		RewardID: rewardID,
		Code:     code,
		Status:   string(domain.RewardActivated),
	})
}

func toCatalogItem(entry app.RewardCatalogEntry) rewardCatalogItem {
	return rewardCatalogItem{
		ID:             entry.Reward.ID(),
		Title:          entry.Reward.Title(),
		Description:    entry.Reward.Description(),
		Kind:           string(entry.Reward.Kind()),
		ConditionType:  string(entry.Reward.ConditionType()),
		ConditionValue: entry.Reward.ConditionValue(),
		Unlocked:       entry.Unlocked,
		Claimed:        entry.Claimed,
		Status:         string(entry.Status),
		Current:        entry.Current,
		Target:         entry.Target,
	}
}

func toMyRewardItem(item app.GrantedReward) myRewardItem {
	return myRewardItem{
		RewardID:    item.Grant.RewardID(),
		Title:       item.Reward.Title(),
		Description: item.Reward.Description(),
		Kind:        string(item.Reward.Kind()),
		Status:      string(item.Grant.Status()),
		Code:        item.Grant.Code(),
		GrantedAt:   item.Grant.GrantedAt(),
		ActivatedAt: item.Grant.ActivatedAt(),
		ExpiresAt:   item.Grant.ExpiresAt(),
	}
}

func mapPetError(err error) error {
	switch {
	case errors.Is(err, domain.ErrRewardNotGranted):
		return domainerr.NewInvalid("reward_id", "reward has not been granted to you")
	case errors.Is(err, domain.ErrRewardAlreadyActivated):
		return domainerr.NewConflict("reward has already been activated")
	case errors.Is(err, domain.ErrRewardNotActivatable):
		return domainerr.NewConflict("reward cannot be activated")
	case errors.Is(err, domain.ErrRewardExpired):
		return domainerr.NewConflict("reward has expired")
	case errors.Is(err, domain.ErrRewardLocked):
		return domainerr.NewInvalid("reward_id", "reward condition is not met")
	case errors.Is(err, domain.ErrDuplicateAction):
		return domainerr.NewConflict("action has already been performed today")
	case errors.Is(err, domain.ErrLimitReached):
		return domainerr.NewConflict("action limit reached")
	case errors.Is(err, domain.ErrConditionNotMet):
		return domainerr.NewInvalid("action", "action condition is not met")
	case errors.Is(err, domain.ErrInvalidAction):
		return domainerr.NewInvalid("action", "invalid action data")
	default:
		return err
	}
}

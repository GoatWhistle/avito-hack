package api

import (
	"context"
	"errors"
	"net/http"
	"time"

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

type rewardCatalogItem struct {
	ID             string `json:"id"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	Kind           string `json:"kind"`
	ConditionType  string `json:"condition_type"`
	ConditionValue int    `json:"condition_value"`
	Unlocked       bool   `json:"unlocked"`
	Claimed        bool   `json:"claimed"`
	Status         string `json:"status,omitempty"`
	Current        int    `json:"progress_current"`
	Target         int    `json:"progress_target"`
}

type myRewardItem struct {
	RewardID    string     `json:"reward_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Kind        string     `json:"kind"`
	Status      string     `json:"status"`
	Code        string     `json:"code,omitempty"`
	GrantedAt   time.Time  `json:"granted_at"`
	ActivatedAt *time.Time `json:"activated_at,omitempty"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
}

type activateRewardResponse struct {
	RewardID string `json:"reward_id"`
	Code     string `json:"code"`
	Status   string `json:"status"`
}

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
	for _, entry := range entries {
		items = append(items, toCatalogItem(entry))
	}

	httpx.OK(w, httpx.NewListResponse(items, ""))
}

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
	for _, item := range granted {
		items = append(items, toMyRewardItem(item))
	}

	httpx.OK(w, httpx.NewListResponse(items, ""))
}

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

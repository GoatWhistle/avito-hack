package api

import (
	"net/http"

	"github.com/avito-hack/backend/internal/module/raccoon/app"
	"github.com/avito-hack/backend/internal/shared/apierr"
	"github.com/avito-hack/backend/internal/shared/auth"
	"github.com/avito-hack/backend/internal/shared/httpx"
)

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
		ID:            profile.ID.String(),
		UserID:        profile.UserID.String(),
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
			earnedAt = badge.EarnedAt.UTC().Format(badgeTimeLayout)
		}
		responses = append(responses, BadgeResponse{
			ID:          badge.ID,
			Name:        badge.Name,
			Description: badge.Description,
			Icon:        badge.IconURL,
			EarnedAt:    earnedAt,
		})
	}

	return responses
}

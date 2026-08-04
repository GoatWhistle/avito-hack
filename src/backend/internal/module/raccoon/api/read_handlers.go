package api

import (
	"net/http"

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

	badges := make([]BadgeResponse, 0, len(profile.Badges))
	for _, b := range profile.Badges {
		earnedAt := ""
		if b.EarnedAt != nil {
			earnedAt = b.EarnedAt.Format("2006-01-02T15:04:05Z")
		}
		badges = append(badges, BadgeResponse{
			ID:          b.ID,
			Name:        b.Name,
			Description: b.Description,
			Icon:        b.IconURL,
			EarnedAt:    earnedAt,
		})
	}

	res := RaccoonProfileResponse{
		ID:            profile.ID.String(),
		UserID:        profile.UserID.String(),
		Name:          profile.Name,
		Level:         profile.Level,
		XP:            profile.XP,
		XPToNextLevel: profile.XPToNextLevel,
		CurrentStreak: profile.CurrentStreak,
		Badges:        badges,
	}

	httpx.OK(w, res)
}

func (h *Handlers) GetBadges(w http.ResponseWriter, r *http.Request) {
	actor, err := auth.ActorFrom(r.Context())
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	badgeViews, err := h.deps.Core.EvaluateBadges(r.Context(), actor.ID)
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	badges := make([]BadgeResponse, 0, len(badgeViews))
	for _, b := range badgeViews {
		earnedAt := ""
		if b.EarnedAt != nil {
			earnedAt = b.EarnedAt.Format("2006-01-02T15:04:05Z")
		}
		badges = append(badges, BadgeResponse{
			ID:          b.ID,
			Name:        b.Name,
			Description: b.Description,
			Icon:        b.IconURL,
			EarnedAt:    earnedAt,
		})
	}

	httpx.OK(w, badges)
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

	httpx.OK(w, ClaimRewardResponse{
		RewardID:  result.RewardID,
		Promocode: result.Promocode,
	})
}

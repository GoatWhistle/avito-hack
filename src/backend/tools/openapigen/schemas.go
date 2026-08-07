//go:build tools

package main

import (
	"strings"
)

const schemaRefPrefix = "#/components/schemas/"

var schemaNames = map[string]string{
	"apierr.ErrorEnvelope":         "ErrorEnvelope",
	"apierr.ErrorBody":             "ErrorBody",
	"api.registerRequest":          "RegisterRequest",
	"api.loginRequest":             "LoginRequest",
	"api.updateProfileRequest":     "UpdateProfileRequest",
	"api.userResponse":             "User",
	"api.sessionResponse":          "Session",
	"api.createItemRequest":        "CreateItemRequest",
	"api.updateItemRequest":        "UpdateItemRequest",
	"api.changeStatusRequest":      "ChangeStatusRequest",
	"api.itemResponse":             "Item",
	"api.itemListItemResponse":     "ItemListEntry",
	"api.ItemListResponse":         "ItemListResponse",
	"api.photoResponse":            "ItemPhoto",
	"api.favoriteItemResponse":     "Favorite",
	"api.FavoriteListResponse":     "FavoriteListResponse",
	"api.petPayload":               "Pet",
	"api.progressPayload":          "Progress",
	"api.streakPayload":            "Streak",
	"api.checkInPayload":           "CheckInResult",
	"api.rewardCatalogItem":        "Reward",
	"api.RewardCatalogResponse":    "RewardCatalogResponse",
	"api.myRewardItem":             "UserReward",
	"api.UserRewardListResponse":   "UserRewardListResponse",
	"api.activateRewardResponse":   "ActivateRewardResponse",
	"api.leaderboardEntryResponse": "LeaderboardEntry",
	"api.leaderboardResponse":      "LeaderboardResponse",
	"api.BadgeResponse":            "Badge",
	"api.RaccoonProfileResponse":   "RaccoonProfile",
	"api.ClaimRewardRequest":       "ClaimRewardRequest",
	"api.ClaimRewardResponse":      "ClaimRewardResponse",
	"domain.Stage":                 "PetStage",
	"domain.State":                 "PetState",
}

func renameSchemas(spec map[string]any) {
	components, ok := spec["components"].(map[string]any)
	if !ok {
		return
	}

	schemas, ok := components["schemas"].(map[string]any)
	if !ok {
		return
	}

	renamed := make(map[string]any, len(schemas))

	for name, schema := range schemas {
		renamed[schemaName(name)] = schema
	}

	components["schemas"] = renamed

	walk(spec, func(node map[string]any) {
		ref, ok := node["$ref"].(string)
		if !ok || !strings.HasPrefix(ref, schemaRefPrefix) {
			return
		}

		node["$ref"] = schemaRefPrefix + schemaName(strings.TrimPrefix(ref, schemaRefPrefix))
	})
}

func schemaName(generated string) string {
	if mapped, ok := schemaNames[generated]; ok {
		return mapped
	}

	return generated
}

package api

import (
	"time"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/module/favorite/app"
)

type favoriteItemResponse struct {
	ItemID      uuid.UUID `json:"item_id"`
	OwnerID     uuid.UUID `json:"owner_id"`
	Title       string    `json:"title"`
	PriceKopeks int64     `json:"price"`
	Status      string    `json:"status"`
	PhotoURL    string    `json:"photo_url,omitempty"`
	AddedAt     time.Time `json:"added_at"`
}

func toFavoriteListResponse(items []app.FavoriteItem) []favoriteItemResponse {
	result := make([]favoriteItemResponse, 0, len(items))

	for _, item := range items {
		result = append(result, favoriteItemResponse{
			ItemID:      item.ItemID,
			OwnerID:     item.OwnerID,
			Title:       item.Title,
			PriceKopeks: item.PriceKopeks,
			Status:      item.Status,
			PhotoURL:    item.PhotoURL,
			AddedAt:     item.CreatedAt,
		})
	}

	return result
}

package api

import (
	"time"

	"github.com/avito-hack/backend/internal/module/favorite/app"
)

type FavoriteListResponse struct {
	Items      []favoriteItemResponse `json:"items"`
	NextCursor string                 `json:"next_cursor,omitempty"`
}

type favoriteItemResponse struct {
	ItemID      string    `json:"item_id"`
	OwnerID     string    `json:"owner_id"`
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
			ItemID:      item.ItemDisplayID,
			OwnerID:     item.OwnerDisplayID,
			Title:       item.Title,
			PriceKopeks: item.PriceKopeks,
			Status:      item.Status,
			PhotoURL:    item.PhotoURL,
			AddedAt:     item.CreatedAt,
		})
	}

	return result
}

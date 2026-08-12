package api

import (
	"time"

	"github.com/avito-hack/backend/internal/module/item/app"
	"github.com/avito-hack/backend/internal/module/item/domain"
)

type createItemRequest struct {
	Title       string            `json:"title"       validate:"required,min=3,max=200"`
	Description string            `json:"description" validate:"max=5000"`
	PriceKopeks int64             `json:"price"       validate:"gte=0"`
	Attributes  map[string]string `json:"attributes"  validate:"omitempty,max=50,dive,keys,max=100,endkeys,max=500"`
}

type updateItemRequest struct {
	Title       *string            `json:"title"       validate:"omitempty,min=3,max=200"`
	Description *string            `json:"description" validate:"omitempty,max=5000"`
	PriceKopeks *int64             `json:"price"       validate:"omitempty,gte=0"`
	Attributes  *map[string]string `json:"attributes"  validate:"omitempty,dive,keys,max=100,endkeys,max=500"`
}

type changeStatusRequest struct {
	Action string `json:"action" validate:"required,oneof=submit sell archive restore"`
}

type photoResponse struct {
	ID        string    `json:"id"`
	ItemID    string    `json:"item_id"`
	URL       string    `json:"url"`
	Position  int       `json:"position"`
	CreatedAt time.Time `json:"created_at"`
}

type itemResponse struct {
	ID               string            `json:"id"`
	OwnerID          string            `json:"owner_id"`
	Title            string            `json:"title"`
	Description      string            `json:"description"`
	PriceKopeks      int64             `json:"price"`
	Status           string            `json:"status"`
	Attributes       map[string]string `json:"attributes"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
	ModerationReason string            `json:"moderation_reason,omitempty"`
	IsSeed           bool              `json:"is_seed"`
	AIVerified       bool              `json:"ai_verified"`
}

type ItemListResponse struct {
	Items      []itemListItemResponse `json:"items"`
	NextCursor string                 `json:"next_cursor,omitempty"`
}

type itemListItemResponse struct {
	ID          string    `json:"id"`
	OwnerID     string    `json:"owner_id"`
	OwnerName   string    `json:"owner_name,omitempty"`
	Title       string    `json:"title"`
	PriceKopeks int64     `json:"price"`
	Status      string    `json:"status"`
	Category    string    `json:"category,omitempty"`
	Condition   string    `json:"condition,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	IsSeed      bool      `json:"is_seed"`
	AIVerified  bool      `json:"ai_verified"`
}

func toItemResponse(i *domain.Item, ownerDisplayID string) itemResponse {
	return itemResponse{
		ID:          i.DisplayID(),
		OwnerID:     ownerDisplayID,
		Title:       i.Title(),
		Description: i.Description(),
		PriceKopeks: i.Price().Kopeks(),
		Status:      i.Status().String(),
		Attributes:  map[string]string(i.Attributes()),
		CreatedAt:   i.CreatedAt(),
		UpdatedAt:   i.UpdatedAt(),
		IsSeed:      i.IsSeed(),
		AIVerified:  i.AIVerified(),
	}
}

func toItemViewResponse(v app.ItemView) itemResponse {
	resp := toItemResponse(v.Item, v.OwnerDisplayID)
	resp.ModerationReason = v.ModerationReason

	return resp
}

func toPhotoResponse(p *domain.Photo, itemDisplayID string) photoResponse {
	return photoResponse{
		ID:        p.DisplayID(),
		ItemID:    itemDisplayID,
		URL:       p.URL(),
		Position:  p.Position(),
		CreatedAt: p.CreatedAt(),
	}
}

func toPhotoListResponse(photos []*domain.Photo, itemDisplayID string) []photoResponse {
	result := make([]photoResponse, 0, len(photos))
	for _, photo := range photos {
		result = append(result, toPhotoResponse(photo, itemDisplayID))
	}

	return result
}

func toListResponse(items []app.ListItem) []itemListItemResponse {
	result := make([]itemListItemResponse, 0, len(items))

	for _, item := range items {
		result = append(result, itemListItemResponse{
			ID:          item.DisplayID,
			OwnerID:     item.OwnerDisplayID,
			OwnerName:   item.OwnerName,
			Title:       item.Title,
			PriceKopeks: item.PriceKopeks,
			Status:      item.Status.String(),
			Category:    item.Category,
			Condition:   item.Condition,
			CreatedAt:   item.CreatedAt,
			IsSeed:      item.IsSeed,
			AIVerified:  item.AIVerified,
		})
	}

	return result
}

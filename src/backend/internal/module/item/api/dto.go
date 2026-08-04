package api

import (
	"time"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/module/item/app"
	"github.com/avito-hack/backend/internal/module/item/domain"
)

type createItemRequest struct {
	Title       string            `json:"title"       validate:"required,min=3,max=200"`
	Description string            `json:"description" validate:"max=5000"`
	PriceKopeks int64             `json:"price"       validate:"gte=0"`
	Attributes  map[string]string `json:"attributes"  validate:"omitempty,max=50"`
}

type updateItemRequest struct {
	Title       *string            `json:"title"       validate:"omitempty,min=3,max=200"`
	Description *string            `json:"description" validate:"omitempty,max=5000"`
	PriceKopeks *int64             `json:"price"       validate:"omitempty,gte=0"`
	Attributes  *map[string]string `json:"attributes"`
}

type changeStatusRequest struct {
	Action string `json:"action" validate:"required,oneof=submit publish archive restore"`
}

type itemResponse struct {
	ID          uuid.UUID         `json:"id"`
	OwnerID     uuid.UUID         `json:"owner_id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	PriceKopeks int64             `json:"price"`
	Status      string            `json:"status"`
	Attributes  map[string]string `json:"attributes"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

type itemListItemResponse struct {
	ID          uuid.UUID `json:"id"`
	OwnerID     uuid.UUID `json:"owner_id"`
	OwnerName   string    `json:"owner_name,omitempty"`
	Title       string    `json:"title"`
	PriceKopeks int64     `json:"price"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

func toItemResponse(i *domain.Item) itemResponse {
	return itemResponse{
		ID:          i.ID(),
		OwnerID:     i.OwnerID(),
		Title:       i.Title(),
		Description: i.Description(),
		PriceKopeks: i.Price().Kopeks(),
		Status:      i.Status().String(),
		Attributes:  map[string]string(i.Attributes()),
		CreatedAt:   i.CreatedAt(),
		UpdatedAt:   i.UpdatedAt(),
	}
}

func toListResponse(items []app.ListItem) []itemListItemResponse {
	result := make([]itemListItemResponse, 0, len(items))

	for _, item := range items {
		result = append(result, itemListItemResponse{
			ID:          item.ID,
			OwnerID:     item.OwnerID,
			OwnerName:   item.OwnerName,
			Title:       item.Title,
			PriceKopeks: item.PriceKopeks,
			Status:      item.Status.String(),
			CreatedAt:   item.CreatedAt,
		})
	}

	return result
}

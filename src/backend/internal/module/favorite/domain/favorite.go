package domain

import (
	"time"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/shared/domainerr"
)

type Favorite struct {
	userID    uuid.UUID
	itemID    uuid.UUID
	createdAt time.Time
}

func New(userID, itemID uuid.UUID, now time.Time) (*Favorite, error) {
	if userID == uuid.Nil {
		return nil, domainerr.NewInvalid("user_id", "field is required")
	}

	if itemID == uuid.Nil {
		return nil, domainerr.NewInvalid("item_id", "field is required")
	}

	return &Favorite{userID: userID, itemID: itemID, createdAt: now}, nil
}

func (f *Favorite) UserID() uuid.UUID    { return f.userID }
func (f *Favorite) ItemID() uuid.UUID    { return f.itemID }
func (f *Favorite) CreatedAt() time.Time { return f.createdAt }

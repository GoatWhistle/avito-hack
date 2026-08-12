package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/shared/domainerr"
)

const MaxPhotosPerItem = 10

type Photo struct {
	id        uuid.UUID
	displayID string
	itemID    uuid.UUID
	url       string
	position  int
	createdAt time.Time
}

type NewPhotoParams struct {
	ItemID   uuid.UUID
	URL      string
	Position int
	Now      time.Time
}

func NewPhoto(p NewPhotoParams) (*Photo, error) {
	url := strings.TrimSpace(p.URL)
	if url == "" {
		return nil, domainerr.NewInvalid("url", "photo url is required")
	}

	if p.ItemID == uuid.Nil {
		return nil, domainerr.NewInvalid("item_id", "photo must belong to an item")
	}

	if p.Position < 0 || p.Position >= MaxPhotosPerItem {
		return nil, domainerr.NewInvalid("position", "photo position is out of range")
	}

	return &Photo{
		id:        uuid.New(),
		displayID: NewDisplayID(),
		itemID:    p.ItemID,
		url:       url,
		position:  p.Position,
		createdAt: p.Now,
	}, nil
}

func RestorePhoto(id uuid.UUID, displayID string, itemID uuid.UUID, url string, position int, createdAt time.Time) *Photo {
	return &Photo{id: id, displayID: displayID, itemID: itemID, url: url, position: position, createdAt: createdAt}
}

func (p *Photo) ID() uuid.UUID        { return p.id }
func (p *Photo) DisplayID() string    { return p.displayID }
func (p *Photo) ItemID() uuid.UUID    { return p.itemID }
func (p *Photo) URL() string          { return p.url }
func (p *Photo) Position() int        { return p.position }
func (p *Photo) CreatedAt() time.Time { return p.createdAt }

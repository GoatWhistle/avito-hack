package domain

import (
	"time"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/shared/vo"
)

type Item struct {
	id          uuid.UUID
	ownerID     uuid.UUID
	title       string
	description string
	price       vo.Money
	status      Status
	attributes  Attributes
	createdAt   time.Time
	updatedAt   time.Time
}

type NewItemParams struct {
	OwnerID     uuid.UUID
	Title       string
	Description string
	Price       vo.Money
	Attributes  Attributes
	Now         time.Time
}

func NewItem(p NewItemParams) (*Item, error) {
	title, err := normalizeTitle(p.Title)
	if err != nil {
		return nil, err
	}

	description, err := normalizeDescription(p.Description)
	if err != nil {
		return nil, err
	}

	if p.OwnerID == uuid.Nil {
		return nil, errInvalidOwner()
	}

	return &Item{
		id:          uuid.New(),
		ownerID:     p.OwnerID,
		title:       title,
		description: description,
		price:       p.Price,
		status:      StatusDraft,
		attributes:  p.Attributes.Clone(),
		createdAt:   p.Now,
		updatedAt:   p.Now,
	}, nil
}

type RestoreItemParams struct {
	ID          uuid.UUID
	OwnerID     uuid.UUID
	Title       string
	Description string
	Price       vo.Money
	Status      Status
	Attributes  Attributes
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func RestoreItem(p RestoreItemParams) *Item {
	return &Item{
		id:          p.ID,
		ownerID:     p.OwnerID,
		title:       p.Title,
		description: p.Description,
		price:       p.Price,
		status:      p.Status,
		attributes:  p.Attributes,
		createdAt:   p.CreatedAt,
		updatedAt:   p.UpdatedAt,
	}
}

func (i *Item) ID() uuid.UUID          { return i.id }
func (i *Item) OwnerID() uuid.UUID     { return i.ownerID }
func (i *Item) Title() string          { return i.title }
func (i *Item) Description() string    { return i.description }
func (i *Item) Price() vo.Money        { return i.price }
func (i *Item) Status() Status         { return i.status }
func (i *Item) Attributes() Attributes { return i.attributes.Clone() }
func (i *Item) CreatedAt() time.Time   { return i.createdAt }
func (i *Item) UpdatedAt() time.Time   { return i.updatedAt }

func (i *Item) IsOwnedBy(actorID uuid.UUID) bool {
	return i.ownerID == actorID
}

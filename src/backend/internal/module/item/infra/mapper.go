package infra

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/avito-hack/backend/internal/module/item/domain"
	"github.com/avito-hack/backend/internal/shared/vo"
)

type itemRow struct {
	id          uuid.UUID
	displayID   string
	ownerID     uuid.UUID
	title       string
	description string
	priceKopeks int64
	status      string
	attributes  []byte
	createdAt   time.Time
	updatedAt   time.Time
	isSeed      bool
	aiVerified  bool
}

func scanItem(row pgx.Row) (*domain.Item, error) {
	var r itemRow

	err := row.Scan(
		&r.id, &r.displayID, &r.ownerID, &r.title, &r.description,
		&r.priceKopeks, &r.status, &r.attributes,
		&r.createdAt, &r.updatedAt, &r.isSeed, &r.aiVerified,
	)
	if err != nil {
		return nil, err
	}

	return toDomain(r)
}

func toDomain(r itemRow) (*domain.Item, error) {
	price, err := vo.NewMoney(r.priceKopeks)
	if err != nil {
		return nil, fmt.Errorf("restore price for item %s: %w", r.id, err)
	}

	attributes, err := decodeAttributes(r.attributes)
	if err != nil {
		return nil, fmt.Errorf("restore attributes for item %s: %w", r.id, err)
	}

	return domain.RestoreItem(domain.RestoreItemParams{
		ID:          r.id,
		DisplayID:   r.displayID,
		OwnerID:     r.ownerID,
		Title:       r.title,
		Description: r.description,
		Price:       price,
		Status:      domain.Status(r.status),
		Attributes:  attributes,
		CreatedAt:   r.createdAt,
		UpdatedAt:   r.updatedAt,
		IsSeed:      r.isSeed,
		AIVerified:  r.aiVerified,
	}), nil
}

func scanPhoto(row pgx.Row) (*domain.Photo, error) {
	var (
		id, itemID uuid.UUID
		displayID  string
		url        string
		position   int
		createdAt  time.Time
	)

	if err := row.Scan(&id, &displayID, &itemID, &url, &position, &createdAt); err != nil {
		return nil, fmt.Errorf("scan photo: %w", err)
	}

	return domain.RestorePhoto(id, displayID, itemID, url, position, createdAt), nil
}

func statusFrom(raw string) domain.Status {
	status := domain.Status(raw)
	if !status.Valid() {
		return domain.StatusDraft
	}

	return status
}

func decodeAttributes(raw []byte) (domain.Attributes, error) {
	if len(raw) == 0 {
		return domain.NewAttributes(nil), nil
	}

	values := map[string]string{}
	if err := json.Unmarshal(raw, &values); err != nil {
		return nil, fmt.Errorf("decode attributes: %w", err)
	}

	return domain.NewAttributes(values), nil
}

func encodeAttributes(attributes domain.Attributes) ([]byte, error) {
	raw, err := json.Marshal(map[string]string(attributes))
	if err != nil {
		return nil, fmt.Errorf("encode attributes: %w", err)
	}

	return raw, nil
}

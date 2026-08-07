//go:build kafka

package kafka

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/shared/events"
)

const messageVersion = 1

var eventNamespace = uuid.MustParse("6f9619ff-8b86-d011-b42d-00c04fc964ff")

type payload struct {
	Title       string            `json:"title,omitempty"`
	Description string            `json:"description,omitempty"`
	PriceKopeks int64             `json:"price_kopeks,omitempty"`
	PhotoCount  int               `json:"photo_count,omitempty"`
	HasVideo    bool              `json:"has_video,omitempty"`
	Attributes  map[string]string `json:"attributes,omitempty"`
}

type message struct {
	Version    int       `json:"version"`
	ID         uuid.UUID `json:"id"`
	Type       string    `json:"type"`
	UserID     uuid.UUID `json:"user_id"`
	SubjectID  uuid.UUID `json:"subject_id,omitempty"`
	OccurredAt time.Time `json:"occurred_at"`
	Payload    payload   `json:"payload,omitempty"`
}

func toMessage(e events.Event) message {
	return message{
		Version: messageVersion, ID: EventID(e), Type: string(e.Type), UserID: e.UserID,
		SubjectID: e.SubjectID, OccurredAt: e.OccurredAt,
		Payload: payload{
			Title: e.Payload.Title, Description: e.Payload.Description,
			PriceKopeks: e.Payload.PriceKopeks, PhotoCount: e.Payload.PhotoCount,
			HasVideo: e.Payload.HasVideo, Attributes: e.Payload.Attributes,
		},
	}
}

func toEvent(m message) events.Event {
	return events.Event{
		Type: events.Type(m.Type), UserID: m.UserID, SubjectID: m.SubjectID,
		OccurredAt: m.OccurredAt,
		Payload: events.Payload{
			Title: m.Payload.Title, Description: m.Payload.Description,
			PriceKopeks: m.Payload.PriceKopeks, PhotoCount: m.Payload.PhotoCount,
			HasVideo: m.Payload.HasVideo, Attributes: m.Payload.Attributes,
		},
	}
}

func EventID(e events.Event) uuid.UUID {
	seed := fmt.Sprintf("%s|%s|%s|%d",
		e.Type, e.UserID, e.SubjectID, e.OccurredAt.UTC().UnixNano())

	return uuid.NewHash(sha256.New(), eventNamespace, []byte(seed), 5)
}

func encodeMessage(m message) ([]byte, error) {
	raw, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("encode kafka message: %w", err)
	}

	return raw, nil
}

func decodeMessage(raw []byte) (message, error) {
	var m message
	if err := json.Unmarshal(raw, &m); err != nil {
		return message{}, fmt.Errorf("decode kafka message: %w", err)
	}
	if err := validateMessage(m); err != nil {
		return message{}, err
	}

	return m, nil
}

func validateMessage(m message) error {
	if m.Version != messageVersion {
		return fmt.Errorf("%w: unsupported version %d", ErrInvalidMessage, m.Version)
	}
	if m.ID == uuid.Nil || m.UserID == uuid.Nil || m.Type == "" || m.OccurredAt.IsZero() {
		return ErrInvalidMessage
	}

	return nil
}

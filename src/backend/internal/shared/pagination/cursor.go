package pagination

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	DefaultLimit = 20
	MaxLimit     = 100
)

var ErrInvalidCursor = errors.New("invalid cursor")

type Cursor struct {
	CreatedAt time.Time
	ID        uuid.UUID
}

func (c Cursor) IsZero() bool {
	return c.ID == uuid.Nil
}

func (c Cursor) Encode() string {
	if c.IsZero() {
		return ""
	}

	raw := fmt.Sprintf("%s|%s", c.CreatedAt.UTC().Format(time.RFC3339Nano), c.ID)

	return base64.RawURLEncoding.EncodeToString([]byte(raw))
}

func DecodeCursor(raw string) (Cursor, error) {
	if raw == "" {
		return Cursor{}, nil
	}

	decoded, err := base64.RawURLEncoding.DecodeString(raw)
	if err != nil {
		return Cursor{}, ErrInvalidCursor
	}

	parts := strings.SplitN(string(decoded), "|", 2)
	if len(parts) != 2 {
		return Cursor{}, ErrInvalidCursor
	}

	createdAt, err := time.Parse(time.RFC3339Nano, parts[0])
	if err != nil {
		return Cursor{}, ErrInvalidCursor
	}

	id, err := uuid.Parse(parts[1])
	if err != nil || id == uuid.Nil {
		return Cursor{}, ErrInvalidCursor
	}

	return Cursor{CreatedAt: createdAt, ID: id}, nil
}

func NormalizeLimit(limit int) int {
	if limit <= 0 {
		return DefaultLimit
	}
	if limit > MaxLimit {
		return MaxLimit
	}

	return limit
}

package domain

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/avito-hack/backend/internal/shared/domainerr"
)

const (
	minTitleLen       = 3
	maxTitleLen       = 200
	maxDescriptionLen = 5000
)

var ErrItemNotFound = fmt.Errorf("item not found: %w", domainerr.ErrNotFound)

func errTransition(from, to Status) error {
	return domainerr.NewConflict(fmt.Sprintf("cannot change status from %s to %s", from, to))
}

func errInvalidOwner() error {
	return domainerr.NewInvalid("owner_id", "field is required")
}

func normalizeTitle(raw string) (string, error) {
	title := strings.TrimSpace(raw)

	length := utf8.RuneCountInString(title)
	if length < minTitleLen {
		return "", domainerr.NewInvalid("title", fmt.Sprintf("value is shorter than minimum: %d", minTitleLen))
	}
	if length > maxTitleLen {
		return "", domainerr.NewInvalid("title", fmt.Sprintf("value is longer than maximum: %d", maxTitleLen))
	}

	return title, nil
}

func normalizeDescription(raw string) (string, error) {
	description := strings.TrimSpace(raw)

	if utf8.RuneCountInString(description) > maxDescriptionLen {
		return "", domainerr.NewInvalid("description",
			fmt.Sprintf("value is longer than maximum: %d", maxDescriptionLen))
	}

	return description, nil
}

package vo

import (
	"errors"
	"net/mail"
	"strings"
)

var ErrInvalidEmail = errors.New("email: invalid address")

const maxEmailLen = 254

type Email struct {
	value string
}

func NewEmail(raw string) (Email, error) {
	normalized := strings.ToLower(strings.TrimSpace(raw))

	if normalized == "" || len(normalized) > maxEmailLen {
		return Email{}, ErrInvalidEmail
	}

	if _, err := mail.ParseAddress(normalized); err != nil {
		return Email{}, ErrInvalidEmail
	}

	return Email{value: normalized}, nil
}

func (e Email) String() string {
	return e.value
}

func (e Email) IsZero() bool {
	return e.value == ""
}

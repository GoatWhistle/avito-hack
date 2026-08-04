package password

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const (
	cost      = 12
	MinLength = 8
	MaxLength = 72
)

var (
	ErrMismatch  = errors.New("password mismatch")
	ErrTooShort  = errors.New("password is too short")
	ErrTooLong   = errors.New("password is too long")
	ErrEmptyHash = errors.New("password hash is empty")
)

type Hash struct {
	value string
}

func NewHash(plain string) (Hash, error) {
	if len(plain) < MinLength {
		return Hash{}, ErrTooShort
	}
	if len(plain) > MaxLength {
		return Hash{}, ErrTooLong
	}

	raw, err := bcrypt.GenerateFromPassword([]byte(plain), cost)
	if err != nil {
		return Hash{}, fmt.Errorf("hash password: %w", err)
	}

	return Hash{value: string(raw)}, nil
}

func RestoreHash(stored string) (Hash, error) {
	if stored == "" {
		return Hash{}, ErrEmptyHash
	}

	return Hash{value: stored}, nil
}

func (h Hash) Compare(plain string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(h.value), []byte(plain)); err != nil {
		return ErrMismatch
	}

	return nil
}

func (h Hash) String() string {
	return h.value
}

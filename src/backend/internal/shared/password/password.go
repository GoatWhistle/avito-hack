package password

import (
	"errors"
	"fmt"
	"sync"

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

	raw, err := bcrypt.GenerateFromPassword(truncate(plain), cost)
	if err != nil {
		return Hash{}, fmt.Errorf("hash password: %w", err)
	}

	return Hash{value: string(raw)}, nil
}

func truncate(plain string) []byte {
	raw := []byte(plain)
	if len(raw) > MaxLength {
		return raw[:MaxLength]
	}

	return raw
}

func RestoreHash(stored string) (Hash, error) {
	if stored == "" {
		return Hash{}, ErrEmptyHash
	}

	return Hash{value: stored}, nil
}

func (h Hash) Compare(plain string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(h.value), truncate(plain)); err != nil {
		return ErrMismatch
	}

	return nil
}

func (h Hash) String() string {
	return h.value
}

var decoyHash = sync.OnceValue(func() []byte {
	raw, err := bcrypt.GenerateFromPassword([]byte("decoy-password-for-constant-time-login"), cost)
	if err != nil {
		return nil
	}

	return raw
})

func CompareDecoy(plain string) {
	hash := decoyHash()
	if hash == nil {
		return
	}

	if err := bcrypt.CompareHashAndPassword(hash, truncate(plain)); err != nil {
		return
	}
}

package domain

import (
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/shared/auth"
	"github.com/avito-hack/backend/internal/shared/domainerr"
	"github.com/avito-hack/backend/internal/shared/password"
	"github.com/avito-hack/backend/internal/shared/vo"
)

const (
	minDisplayNameLen = 2
	maxDisplayNameLen = 100
)

type User struct {
	id           uuid.UUID
	email        vo.Email
	passwordHash password.Hash
	displayName  string
	role         auth.Role
	createdAt    time.Time
	updatedAt    time.Time
}

type NewUserParams struct {
	Email        vo.Email
	PasswordHash password.Hash
	DisplayName  string
	Role         auth.Role
	Now          time.Time
}

func NewUser(p NewUserParams) (*User, error) {
	name, err := normalizeDisplayName(p.DisplayName)
	if err != nil {
		return nil, err
	}

	if p.Email.IsZero() {
		return nil, domainerr.NewInvalid("email", "field is required")
	}

	role := p.Role
	if role == "" {
		role = auth.RoleUser
	}
	if !role.Valid() {
		return nil, domainerr.NewInvalid("role", "unknown role")
	}

	return &User{
		id:           uuid.New(),
		email:        p.Email,
		passwordHash: p.PasswordHash,
		displayName:  name,
		role:         role,
		createdAt:    p.Now,
		updatedAt:    p.Now,
	}, nil
}

type RestoreUserParams struct {
	ID           uuid.UUID
	Email        vo.Email
	PasswordHash password.Hash
	DisplayName  string
	Role         auth.Role
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func RestoreUser(p RestoreUserParams) *User {
	return &User{
		id:           p.ID,
		email:        p.Email,
		passwordHash: p.PasswordHash,
		displayName:  p.DisplayName,
		role:         p.Role,
		createdAt:    p.CreatedAt,
		updatedAt:    p.UpdatedAt,
	}
}

func (u *User) ID() uuid.UUID                   { return u.id }
func (u *User) Email() vo.Email                 { return u.email }
func (u *User) PasswordHash() password.Hash     { return u.passwordHash }
func (u *User) DisplayName() string             { return u.displayName }
func (u *User) Role() auth.Role                 { return u.role }
func (u *User) CreatedAt() time.Time            { return u.createdAt }
func (u *User) UpdatedAt() time.Time            { return u.updatedAt }
func (u *User) Actor() auth.Actor               { return auth.Actor{ID: u.id, Role: u.role} }
func (u *User) Authenticate(plain string) error { return u.passwordHash.Compare(plain) }

func (u *User) Rename(displayName string, now time.Time) error {
	name, err := normalizeDisplayName(displayName)
	if err != nil {
		return err
	}

	u.displayName = name
	u.updatedAt = now

	return nil
}

func (u *User) ChangePassword(hash password.Hash, now time.Time) {
	u.passwordHash = hash
	u.updatedAt = now
}

func normalizeDisplayName(raw string) (string, error) {
	name := strings.TrimSpace(raw)

	length := utf8.RuneCountInString(name)
	if length < minDisplayNameLen {
		return "", domainerr.NewInvalid("display_name", "value is shorter than minimum: 2")
	}
	if length > maxDisplayNameLen {
		return "", domainerr.NewInvalid("display_name", "value is longer than maximum: 100")
	}

	return name, nil
}

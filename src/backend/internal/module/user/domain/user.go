package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/shared/auth"
	"github.com/avito-hack/backend/internal/shared/domainerr"
	"github.com/avito-hack/backend/internal/shared/password"
	"github.com/avito-hack/backend/internal/shared/vo"
)

type User struct {
	id           uuid.UUID
	email        vo.Email
	passwordHash password.Hash
	fullName     string
	role         auth.Role
	createdAt    time.Time
	updatedAt    time.Time
}

type NewUserParams struct {
	Email        vo.Email
	PasswordHash password.Hash
	FullName     string
	Role         auth.Role
	Now          time.Time
}

func NewUser(p NewUserParams) (*User, error) {
	name, err := normalizeFullName(p.FullName)
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
		fullName:     name,
		role:         role,
		createdAt:    p.Now,
		updatedAt:    p.Now,
	}, nil
}

type RestoreUserParams struct {
	ID           uuid.UUID
	Email        vo.Email
	PasswordHash password.Hash
	FullName     string
	Role         auth.Role
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func RestoreUser(p RestoreUserParams) *User {
	return &User{
		id:           p.ID,
		email:        p.Email,
		passwordHash: p.PasswordHash,
		fullName:     p.FullName,
		role:         p.Role,
		createdAt:    p.CreatedAt,
		updatedAt:    p.UpdatedAt,
	}
}

func (u *User) ID() uuid.UUID                   { return u.id }
func (u *User) Email() vo.Email                 { return u.email }
func (u *User) PasswordHash() password.Hash     { return u.passwordHash }
func (u *User) FullName() string                { return u.fullName }
func (u *User) Role() auth.Role                 { return u.role }
func (u *User) CreatedAt() time.Time            { return u.createdAt }
func (u *User) UpdatedAt() time.Time            { return u.updatedAt }
func (u *User) Actor() auth.Actor               { return auth.Actor{ID: u.id, Role: u.role} }
func (u *User) Authenticate(plain string) error { return u.passwordHash.Compare(plain) }

func (u *User) Rename(fullName string, now time.Time) error {
	name, err := normalizeFullName(fullName)
	if err != nil {
		return err
	}

	u.fullName = name
	u.updatedAt = now

	return nil
}

func normalizeFullName(raw string) (string, error) {
	name := strings.TrimSpace(raw)
	if name == "" {
		return "", domainerr.NewInvalid("full_name", "field is required")
	}

	return name, nil
}

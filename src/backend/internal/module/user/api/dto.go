package api

import (
	"time"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/module/user/domain"
)

type registerRequest struct {
	Email    string `json:"email"     validate:"required,email"`
	Password string `json:"password"  validate:"required,min=8,max=72" minLength:"8" maxLength:"72"`
	FullName string `json:"full_name" validate:"required"`
}

type loginRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type updateProfileRequest struct {
	FullName string `json:"full_name" validate:"required"`
}

type userResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type sessionResponse struct {
	Token string       `json:"token"`
	User  userResponse `json:"user"`
}

func toUserResponse(u *domain.User) userResponse {
	return userResponse{
		ID:        u.ID(),
		Email:     u.Email().String(),
		FullName:  u.FullName(),
		Role:      string(u.Role()),
		CreatedAt: u.CreatedAt(),
	}
}

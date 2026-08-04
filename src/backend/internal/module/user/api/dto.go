package api

import (
	"time"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/module/user/domain"
)

type registerRequest struct {
	Email       string `json:"email"        validate:"required,email,max=254"`
	Password    string `json:"password"     validate:"required,min=8,max=72"`
	DisplayName string `json:"display_name" validate:"required,min=2,max=100"`
}

type loginRequest struct {
	Email    string `json:"email"    validate:"required,email,max=254"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

type updateProfileRequest struct {
	DisplayName string `json:"display_name" validate:"required,min=2,max=100"`
}

type userResponse struct {
	ID          uuid.UUID `json:"id"`
	Email       string    `json:"email"`
	DisplayName string    `json:"display_name"`
	Role        string    `json:"role"`
	CreatedAt   time.Time `json:"created_at"`
}

type sessionResponse struct {
	Token     string       `json:"token"`
	ExpiresAt time.Time    `json:"expires_at"`
	User      userResponse `json:"user"`
}

func toUserResponse(u *domain.User) userResponse {
	return userResponse{
		ID:          u.ID(),
		Email:       u.Email().String(),
		DisplayName: u.DisplayName(),
		Role:        string(u.Role()),
		CreatedAt:   u.CreatedAt(),
	}
}

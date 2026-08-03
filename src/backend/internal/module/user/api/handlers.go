package api

import (
	"net/http"

	"github.com/avito-hack/backend/internal/module/user/app"
	"github.com/avito-hack/backend/internal/shared/apierr"
	"github.com/avito-hack/backend/internal/shared/auth"
	"github.com/avito-hack/backend/internal/shared/httpx"
)

type validator interface {
	Struct(dst any) error
}

type Deps struct {
	Register      *app.RegisterUserHandler
	Login         *app.LoginUserHandler
	GetProfile    *app.GetProfileHandler
	UpdateProfile *app.UpdateProfileHandler
	Validator     validator
	Authenticate  func(http.Handler) http.Handler
	MaxBodyBytes  int64
}

type Handlers struct {
	deps Deps
}

func NewHandlers(deps Deps) *Handlers {
	return &Handlers{deps: deps}
}

func (h *Handlers) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if !h.decode(w, r, &req) {
		return
	}

	result, err := h.deps.Register.Handle(r.Context(), app.RegisterUserCommand{
		Email:       req.Email,
		Password:    req.Password,
		DisplayName: req.DisplayName,
	})
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	httpx.Created(w, toUserResponse(result.User))
}

func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if !h.decode(w, r, &req) {
		return
	}

	result, err := h.deps.Login.Handle(r.Context(), app.LoginUserCommand{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	httpx.OK(w, sessionResponse{
		Token:     result.Token,
		ExpiresAt: result.ExpiresAt,
		User:      toUserResponse(result.User),
	})
}

func (h *Handlers) Me(w http.ResponseWriter, r *http.Request) {
	actor, err := auth.ActorFrom(r.Context())
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	user, err := h.deps.GetProfile.Handle(r.Context(), app.GetProfileQuery{UserID: actor.ID})
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	httpx.OK(w, toUserResponse(user))
}

func (h *Handlers) UpdateMe(w http.ResponseWriter, r *http.Request) {
	actor, err := auth.ActorFrom(r.Context())
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	var req updateProfileRequest
	if !h.decode(w, r, &req) {
		return
	}

	user, err := h.deps.UpdateProfile.Handle(r.Context(), app.UpdateProfileCommand{
		UserID:      actor.ID,
		DisplayName: req.DisplayName,
	})
	if err != nil {
		apierr.Write(w, r, err)
		return
	}

	httpx.OK(w, toUserResponse(user))
}

func (h *Handlers) decode(w http.ResponseWriter, r *http.Request, dst any) bool {
	if err := httpx.DecodeJSON(w, r, h.deps.MaxBodyBytes, dst); err != nil {
		apierr.WriteBadRequest(w, r, err.Error())
		return false
	}

	if err := h.deps.Validator.Struct(dst); err != nil {
		apierr.Write(w, r, err)
		return false
	}

	return true
}

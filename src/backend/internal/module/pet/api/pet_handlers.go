package api

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/module/pet/app"
	"github.com/avito-hack/backend/internal/module/pet/domain"
	"github.com/avito-hack/backend/internal/shared/apierr"
	"github.com/avito-hack/backend/internal/shared/auth"
	"github.com/avito-hack/backend/internal/shared/httpx"
)

type petStateService interface {
	State(ctx context.Context, userID uuid.UUID) (*domain.Pet, error)
	Stroke(ctx context.Context, userID uuid.UUID) (*domain.Pet, error)
	CheckIn(ctx context.Context, userID uuid.UUID) (app.ActionResult, error)
	Progress(ctx context.Context, userID uuid.UUID) (app.ProgressView, error)
}

type rewardGranter interface {
	GrantEligible(ctx context.Context, userID uuid.UUID) ([]domain.Reward, error)
}

type PetDeps struct {
	Service      petStateService
	Rewards      rewardGranter
	Authenticate func(http.Handler) http.Handler
}

type PetHandlers struct {
	deps PetDeps
}

func NewPetHandlers(deps PetDeps) *PetHandlers {
	return &PetHandlers{deps: deps}
}

func (h *PetHandlers) Get(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorOf(w, r)
	if !ok {
		return
	}

	pet, err := h.deps.Service.State(r.Context(), actor.ID)
	if err != nil {
		apierr.Write(w, r, err)

		return
	}

	httpx.OK(w, toPetPayload(pet))
}

func (h *PetHandlers) Stroke(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorOf(w, r)
	if !ok {
		return
	}

	pet, err := h.deps.Service.Stroke(r.Context(), actor.ID)
	if err != nil {
		apierr.Write(w, r, err)

		return
	}

	httpx.OK(w, toPetPayload(pet))
}

func (h *PetHandlers) CheckIn(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorOf(w, r)
	if !ok {
		return
	}

	result, err := h.deps.Service.CheckIn(r.Context(), actor.ID)
	if err != nil {
		apierr.Write(w, r, mapPetError(err))

		return
	}

	h.grantEligible(r, actor.ID)

	httpx.OK(w, toCheckInPayload(result))
}

func (h *PetHandlers) Progress(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorOf(w, r)
	if !ok {
		return
	}

	progress, err := h.deps.Service.Progress(r.Context(), actor.ID)
	if err != nil {
		apierr.Write(w, r, err)

		return
	}

	httpx.OK(w, toProgressPayload(progress))
}

func (h *PetHandlers) grantEligible(r *http.Request, userID uuid.UUID) {
	if h.deps.Rewards == nil {
		return
	}

	//nolint:errcheck // reward granting is best-effort, the check-in itself already succeeded
	_, _ = h.deps.Rewards.GrantEligible(r.Context(), userID)
}

func actorOf(w http.ResponseWriter, r *http.Request) (auth.Actor, bool) {
	actor, err := auth.ActorFrom(r.Context())
	if err != nil {
		apierr.Write(w, r, err)

		return auth.Actor{}, false
	}

	return actor, true
}

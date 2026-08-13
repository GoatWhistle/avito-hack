package app

import (
	"context"

	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/module/games/domain"
	"github.com/avito-hack/backend/internal/module/games/domain/moreless"
)

type HiddenPhotoQuery struct {
	UserID uuid.UUID
	Token  string
}

type HiddenPhotoHandler struct {
	rounds RoundRepository
	signer moreless.PhotoSigner
}

func NewHiddenPhotoHandler(rounds RoundRepository, signer moreless.PhotoSigner) *HiddenPhotoHandler {
	return &HiddenPhotoHandler{rounds: rounds, signer: signer}
}

func (h *HiddenPhotoHandler) Handle(ctx context.Context, q HiddenPhotoQuery) (string, error) {
	roundID, displayID, ok := h.signer.Verify(q.Token)
	if !ok {
		return "", domain.ErrRoundNotFound
	}

	round, err := h.rounds.ByDisplayID(ctx, roundID)
	if err != nil {
		return "", err
	}

	if round.GameSlug() != moreless.Slug {
		return "", domain.ErrRoundNotFound
	}

	if q.UserID != uuid.Nil && round.UserID() != q.UserID {
		return "", domain.ErrRoundNotFound
	}

	url, ok := moreless.HiddenPhotoURL(round.Payload(), displayID)
	if !ok {
		return "", domain.ErrRoundNotFound
	}

	return url, nil
}

package app

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type GetRaccoonProfileUseCase struct {
	reader RaccoonStateReader
}

func NewGetRaccoonProfileUseCase(reader RaccoonStateReader) *GetRaccoonProfileUseCase {
	return &GetRaccoonProfileUseCase{reader: reader}
}

func (uc *GetRaccoonProfileUseCase) Execute(ctx context.Context, userID uuid.UUID) (*RaccoonProfileView, error) {
	if userID == uuid.Nil {
		return nil, fmt.Errorf("invalid user id")
	}

	profile, err := uc.reader.GetRaccoonProfile(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get raccoon profile: %w", err)
	}

	return profile, nil
}

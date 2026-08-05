package domain

import (
	"fmt"

	"github.com/avito-hack/backend/internal/shared/domainerr"
)

var (
	ErrFavoriteNotFound = fmt.Errorf("favorite not found: %w", domainerr.ErrNotFound)
	ErrItemNotFound     = fmt.Errorf("item not found: %w", domainerr.ErrNotFound)
)

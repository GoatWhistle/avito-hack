package bukovki

import (
	"fmt"

	"github.com/avito-hack/backend/internal/shared/domainerr"
)

var ErrNoWordsInPool = fmt.Errorf("bukovki word pool is empty: %w", domainerr.ErrNotFound)

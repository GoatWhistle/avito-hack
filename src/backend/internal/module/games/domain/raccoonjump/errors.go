package raccoonjump

import (
	"github.com/avito-hack/backend/internal/shared/domainerr"
)

var ErrImplausibleScore = domainerr.NewInvalid("score", "reported score is not achievable for this round")

package raccoonjump

import (
	"github.com/avito-hack/backend/internal/shared/domainerr"
)

var ErrImplausibleScore = domainerr.NewInvalid("score", "reported score is not achievable for this round")

var ErrRoundExpired = domainerr.NewInvalid("round", "round has expired and can no longer be submitted")

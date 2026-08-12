package domain

import (
	"strings"

	"github.com/avito-hack/backend/internal/shared/publicid"
)

const (
	promoPrefix = "PROMO-"
	promoLength = 8
)

func NewPromoCode() string {
	return promoPrefix + strings.ToUpper(publicid.New()[:promoLength])
}

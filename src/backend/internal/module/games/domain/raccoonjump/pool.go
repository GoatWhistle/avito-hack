package raccoonjump

import (
	"context"
	"time"
)

type Listing struct {
	DisplayID   string
	Title       string
	PhotoURL    string
	PriceKopeks int64
}

type ListingPool interface {
	RandomListings(ctx context.Context, limit int) ([]Listing, error)
}

type Clock interface {
	Now() time.Time
}

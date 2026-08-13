package moreless_test

import (
	"context"
	"encoding/json"
	"math"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/games/domain"
	"github.com/avito-hack/backend/internal/module/games/domain/moreless"
)

const hiddenPrice = int64(987654321)

type nearCall struct {
	seen      []string
	reference int64
	minRatio  float64
	maxRatio  float64
}

type fakePool struct {
	items      []moreless.Item
	cursor     int
	excludes   [][]string
	nearCalls  []nearCall
	nearEmpty  bool
	ignoreBand bool
	err        error
}

func (p *fakePool) Random(_ context.Context, exclude []string) (moreless.Item, error) {
	if p.err != nil {
		return moreless.Item{}, p.err
	}

	p.excludes = append(p.excludes, append([]string(nil), exclude...))

	blocked := make(map[string]struct{}, len(exclude))
	for _, id := range exclude {
		blocked[id] = struct{}{}
	}

	for p.cursor < len(p.items) {
		item := p.items[p.cursor]
		p.cursor++

		if _, skip := blocked[item.DisplayID]; !skip {
			return item, nil
		}
	}

	return moreless.Item{}, domain.ErrNoItemsInPool
}

func (p *fakePool) RandomNear(
	_ context.Context,
	seen []string,
	referenceKopeks int64,
	minRatio, maxRatio float64,
) (moreless.Item, error) {
	if p.err != nil {
		return moreless.Item{}, p.err
	}

	p.nearCalls = append(p.nearCalls, nearCall{
		seen:      append([]string(nil), seen...),
		reference: referenceKopeks,
		minRatio:  minRatio,
		maxRatio:  maxRatio,
	})

	if p.nearEmpty {
		return moreless.Item{}, domain.ErrNoItemsInPool
	}

	if p.ignoreBand {
		minRatio, maxRatio = 0, math.Inf(1)
	}

	lo := float64(referenceKopeks) * minRatio
	hi := float64(referenceKopeks) * maxRatio

	blocked := make(map[string]struct{}, len(seen))
	for _, id := range seen {
		blocked[id] = struct{}{}
	}

	for cursor := p.cursor; cursor < len(p.items); cursor++ {
		item := p.items[cursor]

		if _, skip := blocked[item.DisplayID]; skip {
			continue
		}

		price := float64(item.PriceKopeks)
		if price < lo || price > hi {
			continue
		}

		p.cursor = cursor + 1

		return item, nil
	}

	return moreless.Item{}, domain.ErrNoItemsInPool
}

func poolOf(prices ...int64) *fakePool {
	items := make([]moreless.Item, 0, len(prices))

	for i, price := range prices {
		suffix := strconv.Itoa(i)
		items = append(items, moreless.Item{
			DisplayID:   "item00000" + strings.Repeat("0", 3-len(suffix)) + suffix,
			Title:       "item " + suffix,
			PhotoURL:    "/uploads/item" + suffix + "/photo.jpg",
			PriceKopeks: price,
		})
	}

	return &fakePool{items: items, ignoreBand: true}
}

func bandedPoolOf(prices ...int64) *fakePool {
	pool := poolOf(prices...)
	pool.ignoreBand = false

	return pool
}

const testPhotoSecret = "moreless-test-photo-secret"

func testSigner() moreless.PhotoSigner {
	return moreless.NewPhotoSigner(testPhotoSecret)
}

func newGame(pool moreless.ItemPool) *moreless.Game {
	return moreless.New(pool, testSigner())
}

func newRound() *domain.Round {
	return domain.NewRound(uuid.New(), moreless.Slug, time.Date(2026, time.March, 14, 12, 0, 0, 0, time.UTC))
}

func decodePrompt(t *testing.T, raw json.RawMessage) map[string]any {
	t.Helper()

	var decoded map[string]any
	require.NoError(t, json.Unmarshal(raw, &decoded))

	return decoded
}

func keysOf(m map[string]any) []string {
	out := make([]string, 0, len(m))
	for key := range m {
		out = append(out, key)
	}

	return out
}

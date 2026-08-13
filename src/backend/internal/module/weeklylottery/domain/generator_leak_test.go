package domain_test

import (
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/weeklylottery/domain"
)

const fuzzIterations = 20000

func boardFingerprint(board [domain.BoardSize]domain.Symbol) string {
	counts := symbolCounts(board)
	shape := make([]int, 0, len(counts))
	for _, count := range counts {
		shape = append(shape, count)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(shape)))

	parts := make([]string, 0, len(shape))
	for _, count := range shape {
		parts = append(parts, string(rune('0'+count)))
	}

	return strings.Join(parts, ",")
}

func TestGeneratorNeverContradictsStoredPrize(t *testing.T) {
	t.Parallel()

	generator := domain.NewGenerator(domain.CryptoRandom{})
	for range fuzzIterations {
		board, prize, err := generator.Generate()
		require.NoError(t, err)

		prizeID := ""
		if prize != nil {
			prizeID = prize.ID
		}
		require.NoErrorf(t, domain.ValidateBoard(board, prizeID),
			"generated board %v contradicts prize %q", board, prizeID)
	}
}

func TestGeneratorDoesNotLeakOutcomeThroughSymbolSet(t *testing.T) {
	t.Parallel()

	generator := domain.NewGenerator(domain.CryptoRandom{})
	lossSets := make(map[string]int)
	winSets := make(map[string]int)

	for range fuzzIterations {
		board, prize, err := generator.Generate()
		require.NoError(t, err)

		distinct := make([]string, 0, len(symbolCounts(board)))
		for symbol := range symbolCounts(board) {
			distinct = append(distinct, string(symbol))
		}
		sort.Strings(distinct)
		key := strings.Join(distinct, ",")

		if prize == nil {
			lossSets[key]++
		} else {
			winSets[key]++
		}
	}

	require.NotEmpty(t, lossSets)
	require.NotEmpty(t, winSets)

	shared := 0
	for key := range winSets {
		if lossSets[key] > 0 {
			shared++
		}
	}
	require.NotZerof(t, shared,
		"no symbol set is shared between winning and losing boards: the distinct symbol set reveals the outcome")
}

func TestGeneratorLosingBoardCanContainEverySymbol(t *testing.T) {
	t.Parallel()

	generator := domain.NewGenerator(domain.CryptoRandom{})
	seen := make(map[domain.Symbol]bool)

	for range fuzzIterations {
		board, prize, err := generator.Generate()
		require.NoError(t, err)
		if prize != nil {
			continue
		}
		for _, symbol := range board {
			seen[symbol] = true
		}
	}

	for _, symbol := range []domain.Symbol{
		domain.SymbolBicycle, domain.SymbolSmartphone, domain.SymbolSofa,
		domain.SymbolSneakers, domain.SymbolDelivery, domain.SymbolPromotion,
	} {
		require.Truef(t, seen[symbol],
			"symbol %q never appears on a losing board: its presence proves a win", symbol)
	}
}

func TestGeneratorMultisetShapeDoesNotRevealOutcome(t *testing.T) {
	t.Parallel()

	generator := domain.NewGenerator(domain.CryptoRandom{})
	lossShapes := make(map[string]int)
	winShapes := make(map[string]int)

	for range fuzzIterations {
		board, prize, err := generator.Generate()
		require.NoError(t, err)
		if prize == nil {
			lossShapes[boardFingerprint(board)]++
		} else {
			winShapes[boardFingerprint(board)]++
		}
	}

	require.NotEmpty(t, lossShapes)
	require.NotEmpty(t, winShapes)
	require.Greaterf(t, len(lossShapes), 1,
		"losing boards always have the same multiset shape %v, which reveals the outcome", lossShapes)
}

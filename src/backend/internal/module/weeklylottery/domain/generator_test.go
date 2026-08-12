package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/weeklylottery/domain"
)

type sequenceRandom struct {
	values []int
	index  int
}

func (r *sequenceRandom) Intn(limit int) (int, error) {
	value := r.values[r.index%len(r.values)] % limit
	r.index++

	return value, nil
}

func TestGeneratorCreatesLosingBoardWithoutTriples(t *testing.T) {
	t.Parallel()

	generator := domain.NewGenerator(&sequenceRandom{values: []int{0, 0, 0, 0, 0, 0, 0, 0, 0}})
	board, prize, err := generator.Generate()

	require.NoError(t, err)
	assert.Nil(t, prize)
	assert.Len(t, board, domain.BoardSize)
	assertNoUnexpectedTriples(t, board)
}

func TestGeneratorCreatesExactlyTheSelectedPrizeTriple(t *testing.T) {
	t.Parallel()

	// 50 is the first winning ticket after the loss weight and selects the bicycle prize.
	generator := domain.NewGenerator(&sequenceRandom{values: []int{50, 0, 0, 0, 0, 0, 0, 0, 0}})
	board, prize, err := generator.Generate()

	require.NoError(t, err)
	require.NotNil(t, prize)
	assert.Equal(t, domain.SymbolBicycle, prize.Symbol)

	counts := symbolCounts(board)
	assert.Equal(t, 3, counts[domain.SymbolBicycle])
	for symbol, count := range counts {
		if symbol != domain.SymbolBicycle {
			assert.LessOrEqual(t, count, 2)
		}
	}
}

func TestGeneratorCreatesOneTripleForEveryPrize(t *testing.T) {
	t.Parallel()

	draws := []int{50, 65, 75, 85, 92, 97}
	prizes := domain.Prizes()
	require.Len(t, prizes, len(draws))

	for index, prize := range prizes {
		prize := prize
		t.Run(prize.ID, func(t *testing.T) {
			t.Parallel()
			random := &sequenceRandom{values: []int{draws[index], 0, 0, 0, 0, 0, 0, 0, 0}}
			board, selected, err := domain.NewGenerator(random).Generate()

			require.NoError(t, err)
			require.NotNil(t, selected)
			assert.Equal(t, prize.ID, selected.ID)
			require.NoError(t, domain.ValidateBoard(board, selected.ID))
			assert.Equal(t, 3, symbolCounts(board)[selected.Symbol])
		})
	}
}

func TestProductPrizesUseApplicationCategories(t *testing.T) {
	t.Parallel()

	applicationCategories := map[string]struct{}{
		"electronics": {}, "appliances": {}, "furniture": {}, "clothes": {},
		"kids": {}, "sport": {}, "hobby": {}, "music": {}, "books": {},
		"auto": {}, "realty": {}, "beauty": {}, "animals": {}, "other": {},
	}

	for _, prize := range domain.Prizes() {
		if prize.ScopeType != domain.ScopeTypeCategory {
			continue
		}

		_, exists := applicationCategories[prize.ScopeValue]
		assert.Truef(t, exists, "prize %s uses unknown item category %q", prize.ID, prize.ScopeValue)
	}
}

func assertNoUnexpectedTriples(t *testing.T, board [domain.BoardSize]domain.Symbol) {
	t.Helper()
	for _, count := range symbolCounts(board) {
		assert.LessOrEqual(t, count, 2)
	}
}

func symbolCounts(board [domain.BoardSize]domain.Symbol) map[domain.Symbol]int {
	counts := make(map[domain.Symbol]int)
	for _, symbol := range board {
		counts[symbol]++
	}

	return counts
}

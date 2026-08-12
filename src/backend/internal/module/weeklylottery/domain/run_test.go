package domain_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/weeklylottery/domain"
	"github.com/avito-hack/backend/internal/shared/domainerr"
)

func TestRunWinsWhenThirdMatchingSlotIsOpened(t *testing.T) {
	t.Parallel()

	run := newWinningRun()
	for _, slot := range []int{0, 1} {
		_, err := run.Reveal(slot, run.CreatedAt)
		require.NoError(t, err)
		assert.Equal(t, domain.StateActive, run.State)
	}

	symbol, err := run.Reveal(2, run.CreatedAt)
	require.NoError(t, err)
	assert.Equal(t, domain.SymbolBicycle, symbol)
	assert.Equal(t, domain.StateWon, run.State)
}

func TestRunRejectsRepeatedAndInvalidSlots(t *testing.T) {
	t.Parallel()

	run := newWinningRun()
	_, err := run.Reveal(0, run.CreatedAt)
	require.NoError(t, err)

	_, err = run.Reveal(0, run.CreatedAt)
	require.ErrorIs(t, err, domainerr.ErrConflict)
	_, err = run.Reveal(9, run.CreatedAt)
	require.Error(t, err)
}

func TestRunLosesOnlyAfterAllNineNonMatchingSlots(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, time.August, 12, 12, 0, 0, 0, time.UTC)
	board := [domain.BoardSize]domain.Symbol{
		domain.SymbolBicycle, domain.SymbolBicycle,
		domain.SymbolSmartphone, domain.SymbolSmartphone,
		domain.SymbolSofa, domain.SymbolSofa,
		domain.SymbolSneakers, domain.SymbolSneakers,
		domain.SymbolDelivery,
	}
	run := domain.NewRun(uuid.New(), domain.WeekAt(now), board, nil, now)

	for slot := 0; slot < domain.BoardSize-1; slot++ {
		_, err := run.Reveal(slot, now)
		require.NoError(t, err)
		assert.Equal(t, domain.StateActive, run.State)
	}
	_, err := run.Reveal(domain.BoardSize-1, now)
	require.NoError(t, err)
	assert.Equal(t, domain.StateLost, run.State)

	_, err = run.Reveal(0, now)
	require.ErrorIs(t, err, domainerr.ErrConflict)
}

func TestWeekStartsOnMoscowMonday(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, time.August, 16, 23, 30, 0, 0, time.UTC) // Monday 02:30 in Moscow.
	week := domain.WeekAt(now)

	assert.Equal(t, time.Monday, week.Start.Weekday())
	assert.Equal(t, 7*24*time.Hour, week.End.Sub(week.Start))
	_, offset := week.Start.Zone()
	assert.Equal(t, 3*60*60, offset)
	assert.Equal(t, time.Date(2026, time.August, 17, 0, 0, 0, 0, time.UTC), week.Key)
}

func TestValidateRunRejectsCorruptOpenedSlots(t *testing.T) {
	t.Parallel()

	run := newWinningRun()
	require.NoError(t, domain.ValidateRun(run))

	run.Opened = []int{0, 0}
	require.Error(t, domain.ValidateRun(run))
	run.Opened = []int{domain.BoardSize}
	require.Error(t, domain.ValidateRun(run))
}

func newWinningRun() domain.Run {
	board := [domain.BoardSize]domain.Symbol{
		domain.SymbolBicycle, domain.SymbolBicycle, domain.SymbolBicycle,
		domain.SymbolSofa, domain.SymbolSofa, domain.SymbolSneakers,
		domain.SymbolSneakers, domain.SymbolDelivery, domain.SymbolPromotion,
	}
	prize, _ := domain.PrizeByID("weekly_bicycle_5")
	now := time.Date(2026, time.August, 12, 12, 0, 0, 0, time.UTC)

	return domain.NewRun(uuid.New(), domain.WeekAt(now), board, &prize, now)
}

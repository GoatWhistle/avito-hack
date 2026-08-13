package infra_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/weeklylottery/domain"
	"github.com/avito-hack/backend/internal/module/weeklylottery/infra"
	"github.com/avito-hack/backend/internal/shared/postgres"
	"github.com/avito-hack/backend/internal/shared/postgres/pgtest"
)

func TestByDisplayIDUsesRowLock(t *testing.T) {
	t.Parallel()

	run := lotteryRun()
	tx := &pgtest.Tx{RowResults: []pgtest.Row{{Values: runRow(run)}}}

	loaded, err := infra.NewPgRepository(nil).ByDisplayID(
		postgres.ContextWithTx(context.Background(), tx), run.DisplayID, true,
	)

	require.NoError(t, err)
	assert.Equal(t, run.DisplayID, loaded.DisplayID)
	assert.Equal(t, run.Board, loaded.Board)
	require.Len(t, tx.QueryRowCalls, 1)
	assert.Contains(t, strings.ToUpper(tx.QueryRowCalls[0].SQL), "FOR UPDATE")
}

func TestCreateReliesOnWeeklyUniqueConflict(t *testing.T) {
	t.Parallel()

	run := lotteryRun()
	tx := &pgtest.Tx{ExecTags: []pgconn.CommandTag{pgconn.NewCommandTag("INSERT 0 0")}}

	created, err := infra.NewPgRepository(nil).Create(
		postgres.ContextWithTx(context.Background(), tx), run,
	)

	require.NoError(t, err)
	assert.False(t, created)
	require.Len(t, tx.ExecCalls, 1)
	assert.Contains(t, tx.ExecCalls[0].SQL, "ON CONFLICT (user_id, week_start) DO NOTHING")
}

func lotteryRun() domain.Run {
	now := time.Date(2026, time.August, 12, 12, 0, 0, 0, time.UTC)
	board := [domain.BoardSize]domain.Symbol{
		domain.SymbolBicycle, domain.SymbolBicycle, domain.SymbolBicycle,
		domain.SymbolSofa, domain.SymbolSofa, domain.SymbolSneakers,
		domain.SymbolSneakers, domain.SymbolDelivery, domain.SymbolPromotion,
	}
	prize, _ := domain.PrizeByID("weekly_bicycle_5")

	return domain.NewRun(uuid.New(), domain.WeekAt(now), board, &prize, now)
}

func runRow(run domain.Run) []any {
	return []any{
		run.ID, run.DisplayID, run.UserID, run.WeekStart, string(run.State),
		[]byte(`[
			"bicycle","bicycle","bicycle","sofa","sofa",
			"sneakers","sneakers","delivery","promotion"
		]`),
		[]byte(`[]`), &run.PrizeID, (*string)(nil), (*time.Time)(nil), run.CreatedAt, run.UpdatedAt,
	}
}

package infra_test

import (
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/pet/domain"
	"github.com/avito-hack/backend/internal/module/pet/infra"
	"github.com/avito-hack/backend/internal/shared/domainerr"
	"github.com/avito-hack/backend/internal/shared/postgres/pgtest"
)

func rewardRow() []any {
	return []any{"promo10", "Promo 10%", "Discount", "promo", "level", 5}
}

func TestRewardCatalogReturnsRewards(t *testing.T) {
	t.Parallel()

	tx := &pgtest.Tx{QueryRows: []pgx.Rows{rowsOf(
		rewardRow(),
		[]any{"boost", "Boost", "Push up", "utility", "streak", 7},
	)}}

	rewards, err := infra.NewPgRewardRepository(nil).Catalog(ctxWith(tx))

	require.NoError(t, err)
	require.Len(t, rewards, 2)
	assert.Equal(t, "promo10", rewards[0].ID())
	assert.Equal(t, domain.RewardKindPromo, rewards[0].Kind())
	assert.Equal(t, domain.ConditionLevel, rewards[0].ConditionType())
	assert.Equal(t, 5, rewards[0].ConditionValue())
	assert.Equal(t, domain.ConditionStreak, rewards[1].ConditionType())
	requireSQL(t, tx.QueryCalls[0], "FROM rewards", "ORDER BY condition_type")
}

func TestRewardCatalogErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		tx   *pgtest.Tx
	}{
		{name: "query", tx: &pgtest.Tx{QueryErrs: []error{errDB}}},
		{
			name: "scan",
			tx: &pgtest.Tx{QueryRows: []pgx.Rows{
				&pgtest.Rows{Records: [][]any{rewardRow()}, ScanErr: errDB},
			}},
		},
		{name: "iterate", tx: &pgtest.Tx{QueryRows: []pgx.Rows{&pgtest.Rows{IterErr: errDB}}}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := infra.NewPgRewardRepository(nil).Catalog(ctxWith(tc.tx))

			require.ErrorIs(t, err, errDB)
		})
	}
}

func TestRewardCatalogRejectsInvalidRow(t *testing.T) {
	t.Parallel()

	tx := &pgtest.Tx{QueryRows: []pgx.Rows{rowsOf(
		[]any{"broken", "Broken", "", "nonsense", "level", 5},
	)}}

	_, err := infra.NewPgRewardRepository(nil).Catalog(ctxWith(tx))

	require.ErrorIs(t, err, domain.ErrInvalidReward)
}

func TestRewardByID(t *testing.T) {
	t.Parallel()

	tx := &pgtest.Tx{RowResults: []pgtest.Row{{Values: rewardRow()}}}

	reward, err := infra.NewPgRewardRepository(nil).ByID(ctxWith(tx), "promo10")

	require.NoError(t, err)
	assert.Equal(t, "Promo 10%", reward.Title())
	assert.Equal(t, "Discount", reward.Description())
	requireSQL(t, tx.QueryRowCalls[0], "FROM rewards WHERE id = $1")
	assert.Equal(t, []any{"promo10"}, tx.QueryRowCalls[0].Args)
}

func TestRewardByIDErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		row  pgtest.Row
		want error
	}{
		{name: "not found", row: pgtest.Row{Err: pgx.ErrNoRows}, want: domainerr.ErrNotFound},
		{name: "db error", row: pgtest.Row{Err: errDB}, want: errDB},
		{
			name: "invalid reward",
			row:  pgtest.Row{Values: []any{"", "T", "D", "promo", "level", 1}},
			want: domain.ErrInvalidReward,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tx := &pgtest.Tx{RowResults: []pgtest.Row{tc.row}}

			reward, err := infra.NewPgRewardRepository(nil).ByID(ctxWith(tx), "promo10")

			require.ErrorIs(t, err, tc.want)
			assert.Equal(t, domain.Reward{}, reward)
		})
	}
}

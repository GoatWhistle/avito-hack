package infra_test

import (
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/pet/infra"
	"github.com/avito-hack/backend/internal/shared/postgres/pgtest"
)

func TestBadgeCatalogReturnsBadges(t *testing.T) {
	t.Parallel()

	rows := rowsOf(
		[]any{"first_step", "First Step", "Do it", "https://cdn/1.png"},
		[]any{"pro", "Pro", "Keep going", "https://cdn/2.png"},
	)
	tx := &pgtest.Tx{QueryRows: []pgx.Rows{rows}}

	badges, err := infra.NewPgBadgeRepository(nil).Catalog(ctxWith(tx))

	require.NoError(t, err)
	require.Len(t, badges, 2)
	assert.Equal(t, "first_step", badges[0].ID())
	assert.Equal(t, "Pro", badges[1].Name())
	assert.Equal(t, "https://cdn/1.png", badges[0].IconURL())
	requireSQL(t, tx.QueryCalls[0], "FROM badges", "ORDER BY id")
	assert.True(t, rows.Closed)
}

func TestBadgeCatalogReturnsEmptySlice(t *testing.T) {
	t.Parallel()

	tx := &pgtest.Tx{QueryRows: []pgx.Rows{rowsOf()}}

	badges, err := infra.NewPgBadgeRepository(nil).Catalog(ctxWith(tx))

	require.NoError(t, err)
	assert.Empty(t, badges)
}

func TestBadgeCatalogErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		tx   *pgtest.Tx
	}{
		{name: "query", tx: &pgtest.Tx{QueryErrs: []error{errDB}}},
		{
			name: "scan",
			tx: &pgtest.Tx{QueryRows: []pgx.Rows{
				&pgtest.Rows{Records: [][]any{{"a", "b", "c", "d"}}, ScanErr: errDB},
			}},
		},
		{
			name: "iterate",
			tx:   &pgtest.Tx{QueryRows: []pgx.Rows{&pgtest.Rows{IterErr: errDB}}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := infra.NewPgBadgeRepository(nil).Catalog(ctxWith(tc.tx))

			require.ErrorIs(t, err, errDB)
		})
	}
}

func TestBadgeEarnedByReturnsBadges(t *testing.T) {
	t.Parallel()

	tx := &pgtest.Tx{QueryRows: []pgx.Rows{rowsOf(
		[]any{"pro", "Pro", "Keep going", "https://cdn/2.png", fixedTime},
	)}}

	earned, err := infra.NewPgBadgeRepository(nil).EarnedBy(ctxWith(tx), userA)

	require.NoError(t, err)
	require.Len(t, earned, 1)
	assert.Equal(t, "pro", earned[0].ID())
	assert.Equal(t, userA, earned[0].UserID())
	assert.Equal(t, fixedTime, earned[0].EarnedAt())
	requireSQL(t, tx.QueryCalls[0], "FROM user_badges ub", "ORDER BY ub.earned_at DESC")
	assert.Equal(t, []any{userA}, tx.QueryCalls[0].Args)
}

func TestBadgeEarnedByErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		tx   *pgtest.Tx
	}{
		{name: "query", tx: &pgtest.Tx{QueryErrs: []error{errDB}}},
		{
			name: "scan",
			tx: &pgtest.Tx{QueryRows: []pgx.Rows{
				&pgtest.Rows{Records: [][]any{{"a", "b", "c", "d", fixedTime}}, ScanErr: errDB},
			}},
		},
		{name: "iterate", tx: &pgtest.Tx{QueryRows: []pgx.Rows{&pgtest.Rows{IterErr: errDB}}}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := infra.NewPgBadgeRepository(nil).EarnedBy(ctxWith(tc.tx), userA)

			require.ErrorIs(t, err, errDB)
		})
	}
}

func TestBadgeAward(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		tx      *pgtest.Tx
		want    bool
		wantErr error
	}{
		{name: "awarded", tx: &pgtest.Tx{ExecTags: []pgconn.CommandTag{pgtest.Affected(1)}}, want: true},
		{name: "conflict", tx: &pgtest.Tx{ExecTags: []pgconn.CommandTag{pgtest.Affected(0)}}, want: false},
		{name: "error", tx: &pgtest.Tx{ExecErrs: []error{errDB}}, wantErr: errDB},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := infra.NewPgBadgeRepository(nil).Award(ctxWith(tc.tx), userA, "pro", fixedTime)

			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
			requireSQL(t, tc.tx.ExecCalls[0], "INSERT INTO user_badges", "ON CONFLICT")
			assert.Equal(t, userA, tc.tx.ExecCalls[0].Args[1])
			assert.Equal(t, "pro", tc.tx.ExecCalls[0].Args[2])
			assert.Equal(t, fixedTime, tc.tx.ExecCalls[0].Args[3])
		})
	}
}

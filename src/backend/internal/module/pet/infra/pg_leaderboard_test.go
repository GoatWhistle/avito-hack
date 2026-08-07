package infra_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/pet/app"
	"github.com/avito-hack/backend/internal/module/pet/infra"
	"github.com/avito-hack/backend/internal/shared/pagination"
	"github.com/avito-hack/backend/internal/shared/postgres/pgtest"
)

func TestLeaderboardPageWithoutCursorUsesOffset(t *testing.T) {
	t.Parallel()

	tx := &pgtest.Tx{QueryRows: []pgx.Rows{rowsOf(
		[]any{userA, "Alice", 7, 420, 3},
		[]any{userB, "Bob", 6, 300, 1},
	)}}

	entries, err := infra.NewPgLeaderboard(nil).
		Page(ctxWith(tx), app.LeaderboardFilter{Limit: 10, Offset: 20})

	require.NoError(t, err)
	require.Len(t, entries, 2)
	assert.Equal(t, 21, entries[0].Rank)
	assert.Equal(t, 22, entries[1].Rank)
	assert.Equal(t, "Alice", entries[0].Name)
	assert.Equal(t, 420, entries[0].XP)
	requireSQL(t, tx.QueryCalls[0], "FROM pets p", "LIMIT $1 OFFSET $2")
	assert.NotContains(t, tx.QueryCalls[0].SQL, "WHERE (p.level <")
	assert.Equal(t, []any{10, 20}, tx.QueryCalls[0].Args)
	assert.Empty(t, tx.QueryRowCalls)
}

func TestLeaderboardPageWithCursorUsesKeyset(t *testing.T) {
	t.Parallel()

	cursor := pagination.LeaderboardCursor{Level: 7, XP: 420, UserID: userA}
	tx := &pgtest.Tx{
		RowResults: []pgtest.Row{{Values: []any{4}}},
		QueryRows:  []pgx.Rows{rowsOf([]any{userB, "Bob", 6, 300, 1})},
	}

	entries, err := infra.NewPgLeaderboard(nil).
		Page(ctxWith(tx), app.LeaderboardFilter{Cursor: cursor, Limit: 10})

	require.NoError(t, err)
	require.Len(t, entries, 1)
	assert.Equal(t, 5, entries[0].Rank)
	requireSQL(t, tx.QueryRowCalls[0], "count(*)", "p.user_id <= $3")
	assert.Equal(t, []any{7, 420, userA}, tx.QueryRowCalls[0].Args)
	requireSQL(t, tx.QueryCalls[0], "WHERE (p.level < $1)", "LIMIT $4 OFFSET $5")
	assert.Equal(t, []any{7, 420, userA, 10, 0}, tx.QueryCalls[0].Args)
}

func TestLeaderboardPageCursorRankErrorPropagates(t *testing.T) {
	t.Parallel()

	tx := &pgtest.Tx{RowResults: []pgtest.Row{{Err: errDB}}}

	_, err := infra.NewPgLeaderboard(nil).Page(ctxWith(tx), app.LeaderboardFilter{
		Cursor: pagination.LeaderboardCursor{Level: 1, XP: 2, UserID: userA},
		Limit:  10,
	})

	require.ErrorIs(t, err, errDB)
	assert.Empty(t, tx.QueryCalls)
}

func TestLeaderboardPageErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		tx   *pgtest.Tx
	}{
		{name: "query", tx: &pgtest.Tx{QueryErrs: []error{errDB}}},
		{
			name: "scan",
			tx: &pgtest.Tx{QueryRows: []pgx.Rows{
				&pgtest.Rows{Records: [][]any{{userA, "A", 1, 2, 3}}, ScanErr: errDB},
			}},
		},
		{name: "iterate", tx: &pgtest.Tx{QueryRows: []pgx.Rows{&pgtest.Rows{IterErr: errDB}}}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := infra.NewPgLeaderboard(nil).
				Page(ctxWith(tc.tx), app.LeaderboardFilter{Limit: 10})

			require.ErrorIs(t, err, errDB)
		})
	}
}

func TestLeaderboardRankOf(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		userID    uuid.UUID
		row       *pgtest.Row
		wantRank  int
		wantFound bool
		wantErr   error
	}{
		{name: "nil user short circuits", userID: uuid.Nil},
		{
			name:      "found",
			userID:    userA,
			row:       &pgtest.Row{Values: []any{12}},
			wantRank:  12,
			wantFound: true,
		},
		{name: "no rows", userID: userA, row: &pgtest.Row{Err: pgx.ErrNoRows}},
		{name: "error", userID: userA, row: &pgtest.Row{Err: errDB}, wantErr: errDB},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tx := &pgtest.Tx{}
			if tc.row != nil {
				tx.RowResults = []pgtest.Row{*tc.row}
			}

			rank, found, err := infra.NewPgLeaderboard(nil).RankOf(ctxWith(tx), tc.userID)

			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.wantRank, rank)
			assert.Equal(t, tc.wantFound, found)
		})
	}
}

func TestLeaderboardRankOfNilUserSkipsQuery(t *testing.T) {
	t.Parallel()

	tx := &pgtest.Tx{}

	_, _, err := infra.NewPgLeaderboard(nil).RankOf(ctxWith(tx), uuid.Nil)

	require.NoError(t, err)
	assert.Empty(t, tx.QueryRowCalls)
}

func TestLeaderboardRankOfQueriesWithCTE(t *testing.T) {
	t.Parallel()

	tx := &pgtest.Tx{RowResults: []pgtest.Row{{Values: []any{3}}}}

	_, _, err := infra.NewPgLeaderboard(nil).RankOf(ctxWith(tx), userA)

	require.NoError(t, err)
	requireSQL(t, tx.QueryRowCalls[0], "WITH me AS", "CROSS JOIN me")
	assert.Equal(t, []any{userA}, tx.QueryRowCalls[0].Args)
}

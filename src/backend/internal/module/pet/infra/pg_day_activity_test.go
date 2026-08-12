package infra_test

import (
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/pet/domain"
	"github.com/avito-hack/backend/internal/module/pet/infra"
	"github.com/avito-hack/backend/internal/shared/postgres/pgtest"
)

func TestDayActivityXPByActionAggregates(t *testing.T) {
	t.Parallel()

	tx := &pgtest.Tx{QueryRows: []pgx.Rows{rowsOf(
		[]any{"daily_checkin", 2, 40},
		[]any{"favorite", 5, 25},
	)}}

	got, err := infra.NewPgDayActivity(nil).XPByAction(ctxWith(tx), userA, dayStart, dayEnd)

	require.NoError(t, err)
	require.Len(t, got, 2)
	assert.Equal(t, domain.ActionDailyCheckIn, got[0].Action)
	assert.Equal(t, 2, got[0].Count)
	assert.Equal(t, 40, got[0].Amount)
	assert.Equal(t, domain.ActionFavorite, got[1].Action)
	requireSQL(t, tx.QueryCalls[0], "FROM xp_events", "GROUP BY action")
	assert.Equal(t, []any{userA, dayStart, dayEnd}, tx.QueryCalls[0].Args)
}

func TestDayActivityXPByActionErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		tx   *pgtest.Tx
	}{
		{name: "query", tx: &pgtest.Tx{QueryErrs: []error{errDB}}},
		{
			name: "scan",
			tx: &pgtest.Tx{QueryRows: []pgx.Rows{
				&pgtest.Rows{Records: [][]any{{"favorite", 1, 5}}, ScanErr: errDB},
			}},
		},
		{name: "iterate", tx: &pgtest.Tx{QueryRows: []pgx.Rows{&pgtest.Rows{IterErr: errDB}}}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := infra.NewPgDayActivity(nil).XPByAction(ctxWith(tc.tx), userA, dayStart, dayEnd)

			require.ErrorIs(t, err, errDB)
		})
	}
}

func TestDayActivityXPTotalBefore(t *testing.T) {
	t.Parallel()

	tx := &pgtest.Tx{RowResults: []pgtest.Row{{Values: []any{321}}}}

	total, err := infra.NewPgDayActivity(nil).XPTotalBefore(ctxWith(tx), userA, dayStart)

	require.NoError(t, err)
	assert.Equal(t, 321, total)
	requireSQL(t, tx.QueryRowCalls[0], "sum(amount)", "created_at < $2")
	assert.Equal(t, []any{userA, dayStart}, tx.QueryRowCalls[0].Args)
}

func TestDayActivityXPTotalBeforeWrapsError(t *testing.T) {
	t.Parallel()

	tx := &pgtest.Tx{RowResults: []pgtest.Row{{Err: errDB}}}

	total, err := infra.NewPgDayActivity(nil).XPTotalBefore(ctxWith(tx), userA, dayStart)

	require.ErrorIs(t, err, errDB)
	assert.Zero(t, total)
}

func TestDayActivityRewardsGranted(t *testing.T) {
	t.Parallel()

	tx := &pgtest.Tx{QueryRows: []pgx.Rows{rowsOf([]any{"Promo 10%"}, []any{"Boost"})}}

	titles, err := infra.NewPgDayActivity(nil).RewardsGranted(ctxWith(tx), userA, dayStart, dayEnd)

	require.NoError(t, err)
	assert.Equal(t, []string{"Promo 10%", "Boost"}, titles)
	requireSQL(t, tx.QueryCalls[0], "FROM user_rewards ur", "JOIN rewards r")
}

func TestDayActivityBadgesEarned(t *testing.T) {
	t.Parallel()

	tx := &pgtest.Tx{QueryRows: []pgx.Rows{rowsOf([]any{"Raccoon Friend"})}}

	titles, err := infra.NewPgDayActivity(nil).BadgesEarned(ctxWith(tx), userA, dayStart, dayEnd)

	require.NoError(t, err)
	assert.Equal(t, []string{"Raccoon Friend"}, titles)
	requireSQL(t, tx.QueryCalls[0], "FROM user_badges ub", "JOIN badges b")
}

func TestDayActivityTitlesErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		tx   *pgtest.Tx
	}{
		{name: "query", tx: &pgtest.Tx{QueryErrs: []error{errDB}}},
		{
			name: "scan",
			tx: &pgtest.Tx{QueryRows: []pgx.Rows{
				&pgtest.Rows{Records: [][]any{{"x"}}, ScanErr: errDB},
			}},
		},
		{name: "iterate", tx: &pgtest.Tx{QueryRows: []pgx.Rows{&pgtest.Rows{IterErr: errDB}}}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := infra.NewPgDayActivity(nil).RewardsGranted(ctxWith(tc.tx), userA, dayStart, dayEnd)

			require.ErrorIs(t, err, errDB)
		})
	}
}

func TestDayActivityListingIssues(t *testing.T) {
	t.Parallel()

	itemDisplayID := "abc123def456"

	tx := &pgtest.Tx{QueryRows: []pgx.Rows{rowsOf(
		[]any{itemDisplayID, "Bike", "no_photo", 3},
		[]any{"def456abc123", "Chair", "stale", 30},
	)}}

	issues, err := infra.NewPgDayActivity(nil).ListingIssues(ctxWith(tx), userA, 5)

	require.NoError(t, err)
	require.Len(t, issues, 2)
	assert.Equal(t, itemDisplayID, issues[0].ItemID)
	assert.Equal(t, domain.IssueNoPhoto, issues[0].Kind)
	assert.Equal(t, 3, issues[0].StaleDays)
	assert.Equal(t, domain.IssueStale, issues[1].Kind)
	requireSQL(t, tx.QueryCalls[0], "FROM items i", "LIMIT $2")
	assert.Equal(t, []any{userA, 5, 14}, tx.QueryCalls[0].Args)
}

func TestDayActivityListingIssuesErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		tx   *pgtest.Tx
	}{
		{name: "query", tx: &pgtest.Tx{QueryErrs: []error{errDB}}},
		{
			name: "scan",
			tx: &pgtest.Tx{QueryRows: []pgx.Rows{
				&pgtest.Rows{Records: [][]any{{itemA, "t", "stale", 1}}, ScanErr: errDB},
			}},
		},
		{name: "iterate", tx: &pgtest.Tx{QueryRows: []pgx.Rows{&pgtest.Rows{IterErr: errDB}}}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := infra.NewPgDayActivity(nil).ListingIssues(ctxWith(tc.tx), userA, 5)

			require.ErrorIs(t, err, errDB)
		})
	}
}

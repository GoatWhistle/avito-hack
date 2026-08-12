package infra_test

import (
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/pet/app"
	"github.com/avito-hack/backend/internal/module/pet/domain"
	"github.com/avito-hack/backend/internal/module/pet/infra"
	"github.com/avito-hack/backend/internal/shared/postgres/pgtest"
)

func summaryFilter() app.SummaryHistoryFilter {
	return app.SummaryHistoryFilter{UserID: userA, Limit: 10}
}

func newSummary(t *testing.T, advice *domain.Advice) *domain.DailySummary {
	t.Helper()

	facts := domain.DayFacts{Date: dayStart, PetName: "Enot", TotalXP: 120}

	summary, err := domain.NewDailySummary(
		userA, facts, "Great day", advice, domain.SummarySourceTemplate, fixedTime)
	require.NoError(t, err)

	return summary
}

func TestSummaryListByUserWithoutCursor(t *testing.T) {
	t.Parallel()

	tx := &pgtest.Tx{QueryRows: []pgx.Rows{rowsOf(summaryRow(factsJSON, adviceJSON))}}

	summaries, err := infra.NewPgSummaryRepository(nil).ListByUser(ctxWith(tx), summaryFilter())

	require.NoError(t, err)
	require.Len(t, summaries, 1)
	assert.Equal(t, 120, summaries[0].Facts().TotalXP)
	requireSQL(t, tx.QueryCalls[0], "WHERE user_id = $1", "LIMIT $2")
	assert.NotContains(t, tx.QueryCalls[0].SQL, "AND date <")
	assert.Equal(t, []any{userA, 10}, tx.QueryCalls[0].Args)
}

func TestSummaryListByUserWithBeforeCursor(t *testing.T) {
	t.Parallel()

	tx := &pgtest.Tx{QueryRows: []pgx.Rows{rowsOf()}}
	filter := summaryFilter()
	filter.Before = dayEnd

	summaries, err := infra.NewPgSummaryRepository(nil).ListByUser(ctxWith(tx), filter)

	require.NoError(t, err)
	assert.Empty(t, summaries)
	requireSQL(t, tx.QueryCalls[0], "AND date < $2", "LIMIT $3")
	assert.Equal(t, []any{userA, dayEnd, 10}, tx.QueryCalls[0].Args)
}

func TestSummaryInsert(t *testing.T) {
	t.Parallel()

	adviceItemID := "abc123def456"
	advice := &domain.Advice{Text: "Add a photo", ItemID: &adviceItemID, Action: domain.AdviceAddPhoto}

	tests := []struct {
		name    string
		tx      *pgtest.Tx
		want    bool
		wantErr error
	}{
		{name: "inserted", tx: &pgtest.Tx{ExecTags: []pgconn.CommandTag{pgtest.Affected(1)}}, want: true},
		{name: "conflict", tx: &pgtest.Tx{ExecTags: []pgconn.CommandTag{pgtest.Affected(0)}}, want: false},
		{name: "error", tx: &pgtest.Tx{ExecErrs: []error{errDB}}, wantErr: errDB},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			summary := newSummary(t, advice)

			got, err := infra.NewPgSummaryRepository(nil).Insert(ctxWith(tc.tx), summary)

			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
			requireSQL(t, tc.tx.ExecCalls[0], "INSERT INTO daily_summaries", "ON CONFLICT")

			args := tc.tx.ExecCalls[0].Args
			require.Len(t, args, 8)
			assert.Equal(t, userA, args[1])
			assert.Contains(t, string(args[3].([]byte)), `"total_xp":120`)
			assert.Equal(t, "Great day", args[4])
			assert.Contains(t, string(args[5].([]byte)), `"add_photo"`)
			assert.Equal(t, "template", args[6])
			assert.Equal(t, fixedTime, args[7])
		})
	}
}

func TestSummaryInsertPassesNilAdvice(t *testing.T) {
	t.Parallel()

	tx := &pgtest.Tx{ExecTags: []pgconn.CommandTag{pgtest.Affected(1)}}

	got, err := infra.NewPgSummaryRepository(nil).Insert(ctxWith(tx), newSummary(t, nil))

	require.NoError(t, err)
	assert.True(t, got)
	assert.Nil(t, tx.ExecCalls[0].Args[5])
}

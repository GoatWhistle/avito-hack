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

const factsJSON = `{"total_xp":120,"pet_name":"Enot"}`

const adviceJSON = `{"text":"Add a photo","action":"add_photo"}`

func summaryRow(facts, advice string) []any {
	var rawAdvice []byte
	if advice != "" {
		rawAdvice = []byte(advice)
	}

	return []any{itemA, userA, dayStart, []byte(facts), "Great day", rawAdvice, "template", fixedTime}
}

func TestSummaryByUserAndDateDecodesRow(t *testing.T) {
	t.Parallel()

	tx := &pgtest.Tx{RowResults: []pgtest.Row{{Values: summaryRow(factsJSON, adviceJSON)}}}

	summary, err := infra.NewPgSummaryRepository(nil).ByUserAndDate(ctxWith(tx), userA, dayStart)

	require.NoError(t, err)
	require.NotNil(t, summary)
	assert.Equal(t, itemA, summary.ID())
	assert.Equal(t, userA, summary.UserID())
	assert.Equal(t, dayStart, summary.Date())
	assert.Equal(t, 120, summary.Facts().TotalXP)
	assert.Equal(t, "Enot", summary.Facts().PetName)
	assert.Equal(t, "Great day", summary.Message())
	assert.Equal(t, domain.SummarySourceTemplate, summary.GeneratedBy())
	require.NotNil(t, summary.Advice())
	assert.Equal(t, domain.AdviceAddPhoto, summary.Advice().Action)
	requireSQL(t, tx.QueryRowCalls[0], "FROM daily_summaries WHERE user_id = $1 AND date = $2")
	assert.Equal(t, []any{userA, dayStart}, tx.QueryRowCalls[0].Args)
}

func TestSummaryByUserAndDateLeavesAdviceNilWhenEmpty(t *testing.T) {
	t.Parallel()

	tx := &pgtest.Tx{RowResults: []pgtest.Row{{Values: summaryRow(factsJSON, "")}}}

	summary, err := infra.NewPgSummaryRepository(nil).ByUserAndDate(ctxWith(tx), userA, dayStart)

	require.NoError(t, err)
	assert.Nil(t, summary.Advice())
}

func TestSummaryByUserAndDateErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		row  pgtest.Row
		want error
	}{
		{name: "not found", row: pgtest.Row{Err: pgx.ErrNoRows}, want: domainerr.ErrNotFound},
		{name: "db error", row: pgtest.Row{Err: errDB}, want: errDB},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tx := &pgtest.Tx{RowResults: []pgtest.Row{tc.row}}

			summary, err := infra.NewPgSummaryRepository(nil).
				ByUserAndDate(ctxWith(tx), userA, dayStart)

			require.ErrorIs(t, err, tc.want)
			assert.Nil(t, summary)
		})
	}
}

func TestSummaryByUserAndDateRejectsMalformedJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		row  []any
	}{
		{name: "facts", row: summaryRow(`{"total_xp":`, adviceJSON)},
		{name: "advice", row: summaryRow(factsJSON, `{"text":`)},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tx := &pgtest.Tx{RowResults: []pgtest.Row{{Values: tc.row}}}

			_, err := infra.NewPgSummaryRepository(nil).ByUserAndDate(ctxWith(tx), userA, dayStart)

			require.Error(t, err)
			assert.Contains(t, err.Error(), "decode summary")
		})
	}
}

func TestSummaryListByUserErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		tx   *pgtest.Tx
	}{
		{name: "query", tx: &pgtest.Tx{QueryErrs: []error{errDB}}},
		{
			name: "scan",
			tx: &pgtest.Tx{QueryRows: []pgx.Rows{
				&pgtest.Rows{Records: [][]any{summaryRow(factsJSON, "")}, ScanErr: errDB},
			}},
		},
		{name: "iterate", tx: &pgtest.Tx{QueryRows: []pgx.Rows{&pgtest.Rows{IterErr: errDB}}}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := infra.NewPgSummaryRepository(nil).
				ListByUser(ctxWith(tc.tx), summaryFilter())

			require.ErrorIs(t, err, errDB)
		})
	}
}

package infra_test

import (
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/pet/domain"
	"github.com/avito-hack/backend/internal/module/pet/infra"
	"github.com/avito-hack/backend/internal/shared/postgres/pgtest"
)

func TestXPEventAppendInserts(t *testing.T) {
	t.Parallel()

	tx := &pgtest.Tx{}
	event := domain.NewXPEvent(userA, domain.ActionFavorite, &itemA, 5, fixedTime)

	require.NoError(t, infra.NewPgXPEventRepository(nil).Append(ctxWith(tx), event))
	requireSQL(t, tx.ExecCalls[0], "INSERT INTO xp_events")
	assert.Equal(t,
		[]any{event.ID, userA, "favorite", &itemA, 5, fixedTime},
		tx.ExecCalls[0].Args)
}

func TestXPEventAppendMapsUniqueViolation(t *testing.T) {
	t.Parallel()

	tx := &pgtest.Tx{ExecErrs: []error{&pgconn.PgError{Code: "23505"}}}
	event := domain.NewXPEvent(userA, domain.ActionFavorite, nil, 5, fixedTime)

	err := infra.NewPgXPEventRepository(nil).Append(ctxWith(tx), event)

	require.ErrorIs(t, err, domain.ErrDuplicateAction)
}

func TestXPEventAppendWrapsOtherPgErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
	}{
		{name: "other pg code", err: &pgconn.PgError{Code: "23503"}},
		{name: "plain error", err: errDB},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tx := &pgtest.Tx{ExecErrs: []error{tc.err}}
			event := domain.NewXPEvent(userA, domain.ActionFavorite, nil, 5, fixedTime)

			err := infra.NewPgXPEventRepository(nil).Append(ctxWith(tx), event)

			require.ErrorIs(t, err, tc.err)
			assert.NotErrorIs(t, err, domain.ErrDuplicateAction)
		})
	}
}

func TestXPEventCountSince(t *testing.T) {
	t.Parallel()

	tx := &pgtest.Tx{RowResults: []pgtest.Row{{Values: []any{3}}}}

	count, err := infra.NewPgXPEventRepository(nil).
		CountSince(ctxWith(tx), userA, domain.ActionDailyCheckIn, dayStart)

	require.NoError(t, err)
	assert.Equal(t, 3, count)
	requireSQL(t, tx.QueryRowCalls[0], "count(*) FROM xp_events", "created_at >= $3")
	assert.Equal(t, []any{userA, "daily_checkin", dayStart}, tx.QueryRowCalls[0].Args)
}

func TestXPEventCountSinceWrapsError(t *testing.T) {
	t.Parallel()

	tx := &pgtest.Tx{RowResults: []pgtest.Row{{Err: errDB}}}

	count, err := infra.NewPgXPEventRepository(nil).
		CountSince(ctxWith(tx), userA, domain.ActionDailyCheckIn, dayStart)

	require.ErrorIs(t, err, errDB)
	assert.Zero(t, count)
}

func TestXPEventExistsBySubject(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		row     pgtest.Row
		want    bool
		wantErr error
	}{
		{name: "exists", row: pgtest.Row{Values: []any{true}}, want: true},
		{name: "missing", row: pgtest.Row{Values: []any{false}}, want: false},
		{name: "error", row: pgtest.Row{Err: errDB}, wantErr: errDB},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tx := &pgtest.Tx{RowResults: []pgtest.Row{tc.row}}

			got, err := infra.NewPgXPEventRepository(nil).
				ExistsBySubject(ctxWith(tx), userA, domain.ActionItemSold, itemA)

			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
				assert.False(t, got)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
			requireSQL(t, tx.QueryRowCalls[0], "SELECT EXISTS", "subject_id = $3")
			assert.Equal(t, []any{userA, "item_sold", itemA}, tx.QueryRowCalls[0].Args)
		})
	}
}

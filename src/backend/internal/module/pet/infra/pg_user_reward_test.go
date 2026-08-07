package infra_test

import (
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/pet/domain"
	"github.com/avito-hack/backend/internal/module/pet/infra"
	"github.com/avito-hack/backend/internal/shared/domainerr"
	"github.com/avito-hack/backend/internal/shared/postgres/pgtest"
)

func userRewardRow(code *string, activatedAt, expiresAt *time.Time) []any {
	return []any{itemA, userA, "promo10", "granted", code, fixedTime, activatedAt, expiresAt}
}

func TestRewardListByUserMapsNullables(t *testing.T) {
	t.Parallel()

	code := "ABC123"
	activated := fixedTime.Add(time.Hour)
	tx := &pgtest.Tx{QueryRows: []pgx.Rows{rowsOf(
		userRewardRow(&code, &activated, nil),
		userRewardRow(nil, nil, nil),
	)}}

	granted, err := infra.NewPgRewardRepository(nil).ListByUser(ctxWith(tx), userA)

	require.NoError(t, err)
	require.Len(t, granted, 2)
	assert.Equal(t, "ABC123", granted[0].Code())
	require.NotNil(t, granted[0].ActivatedAt())
	assert.Equal(t, activated, *granted[0].ActivatedAt())
	assert.Nil(t, granted[0].ExpiresAt())
	assert.Empty(t, granted[1].Code())
	assert.Nil(t, granted[1].ActivatedAt())
	assert.Equal(t, domain.RewardGranted, granted[1].Status())
	requireSQL(t, tx.QueryCalls[0], "FROM user_rewards WHERE user_id = $1", "ORDER BY granted_at DESC")
}

func TestRewardListByUserErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		tx   *pgtest.Tx
	}{
		{name: "query", tx: &pgtest.Tx{QueryErrs: []error{errDB}}},
		{
			name: "scan",
			tx: &pgtest.Tx{QueryRows: []pgx.Rows{
				&pgtest.Rows{Records: [][]any{userRewardRow(nil, nil, nil)}, ScanErr: errDB},
			}},
		},
		{name: "iterate", tx: &pgtest.Tx{QueryRows: []pgx.Rows{&pgtest.Rows{IterErr: errDB}}}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := infra.NewPgRewardRepository(nil).ListByUser(ctxWith(tc.tx), userA)

			require.ErrorIs(t, err, errDB)
		})
	}
}

func TestRewardByUserAndReward(t *testing.T) {
	t.Parallel()

	expires := fixedTime.Add(48 * time.Hour)
	tx := &pgtest.Tx{RowResults: []pgtest.Row{{Values: userRewardRow(nil, nil, &expires)}}}

	granted, err := infra.NewPgRewardRepository(nil).ByUserAndReward(ctxWith(tx), userA, "promo10")

	require.NoError(t, err)
	assert.Equal(t, userA, granted.UserID())
	assert.Equal(t, "promo10", granted.RewardID())
	require.NotNil(t, granted.ExpiresAt())
	assert.Equal(t, expires, *granted.ExpiresAt())
	assert.Equal(t, []any{userA, "promo10"}, tx.QueryRowCalls[0].Args)
}

func TestRewardByUserAndRewardErrors(t *testing.T) {
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

			granted, err := infra.NewPgRewardRepository(nil).
				ByUserAndReward(ctxWith(tx), userA, "promo10")

			require.ErrorIs(t, err, tc.want)
			assert.Equal(t, domain.UserReward{}, granted)
		})
	}
}

func TestRewardGrant(t *testing.T) {
	t.Parallel()

	granted, err := domain.GrantReward(userA, "promo10", fixedTime)
	require.NoError(t, err)

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

			got, err := infra.NewPgRewardRepository(nil).Grant(ctxWith(tc.tx), granted)

			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
			requireSQL(t, tc.tx.ExecCalls[0], "INSERT INTO user_rewards", "ON CONFLICT")
			assert.Equal(t,
				[]any{granted.ID(), userA, "promo10", "granted", fixedTime},
				tc.tx.ExecCalls[0].Args)
		})
	}
}

func TestRewardActivate(t *testing.T) {
	t.Parallel()

	granted, err := domain.GrantReward(userA, "promo10", fixedTime)
	require.NoError(t, err)
	require.NoError(t, granted.Activate("CODE99", fixedTime.Add(time.Hour)))

	tests := []struct {
		name    string
		tx      *pgtest.Tx
		want    bool
		wantErr error
	}{
		{name: "updated", tx: &pgtest.Tx{ExecTags: []pgconn.CommandTag{pgtest.Affected(1)}}, want: true},
		{name: "no rows", tx: &pgtest.Tx{ExecTags: []pgconn.CommandTag{pgtest.Affected(0)}}, want: false},
		{name: "error", tx: &pgtest.Tx{ExecErrs: []error{errDB}}, wantErr: errDB},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := infra.NewPgRewardRepository(nil).Activate(ctxWith(tc.tx), granted)

			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
			requireSQL(t, tc.tx.ExecCalls[0], "UPDATE user_rewards", "status = 'granted'")
			assert.Equal(t, "activated", tc.tx.ExecCalls[0].Args[0])
			assert.Equal(t, "CODE99", tc.tx.ExecCalls[0].Args[1])
			assert.Equal(t, userA, tc.tx.ExecCalls[0].Args[3])
			assert.Equal(t, "promo10", tc.tx.ExecCalls[0].Args[4])
		})
	}
}

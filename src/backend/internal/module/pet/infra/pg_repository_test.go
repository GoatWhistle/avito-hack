package infra_test

import (
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/pet/domain"
	"github.com/avito-hack/backend/internal/module/pet/infra"
	"github.com/avito-hack/backend/internal/shared/domainerr"
	"github.com/avito-hack/backend/internal/shared/postgres/pgtest"
)

func petRow(t *testing.T) []any {
	t.Helper()

	checkIn := dayStart
	hatched := fixedTime.Add(-24 * time.Hour)

	return []any{
		itemA, userA, "Enot", "adult", 7, 420, 500, 60, 80, 90,
		4, 1, &checkIn, &hatched, fixedTime, fixedTime, int64(12),
	}
}

func TestPgRepositoryByUserIDRestoresPet(t *testing.T) {
	t.Parallel()

	tx := &pgtest.Tx{RowResults: []pgtest.Row{{Values: petRow(t)}}}

	pet, err := infra.NewPgRepository(nil).ByUserID(ctxWith(tx), userA)

	require.NoError(t, err)
	require.NotNil(t, pet)
	assert.Equal(t, itemA, pet.ID())
	assert.Equal(t, userA, pet.UserID())
	assert.Equal(t, "Enot", pet.Name())
	assert.Equal(t, domain.Stage("adult"), pet.Stage())
	assert.Equal(t, 7, pet.Level())
	assert.Equal(t, 420, pet.XP())
	assert.Equal(t, int64(12), pet.InteractionVersion())
	assert.True(t, pet.IsHatched())
	requireSQL(t, tx.QueryRowCalls[0], "FROM pets WHERE user_id = $1")
	assert.NotContains(t, tx.QueryRowCalls[0].SQL, "FOR UPDATE")
}

func TestPgRepositoryByUserIDHandlesNullTimestamps(t *testing.T) {
	t.Parallel()

	values := petRow(t)
	values[12] = nil
	values[13] = nil
	tx := &pgtest.Tx{RowResults: []pgtest.Row{{Values: values}}}

	pet, err := infra.NewPgRepository(nil).ByUserID(ctxWith(tx), userA)

	require.NoError(t, err)
	assert.Nil(t, pet.LastCheckInDate())
	assert.Nil(t, pet.HatchedAt())
	assert.False(t, pet.IsHatched())
}

func TestPgRepositoryByUserIDNotFound(t *testing.T) {
	t.Parallel()

	tx := &pgtest.Tx{RowResults: []pgtest.Row{{Err: pgx.ErrNoRows}}}

	pet, err := infra.NewPgRepository(nil).ByUserID(ctxWith(tx), userA)

	require.ErrorIs(t, err, domainerr.ErrNotFound)
	assert.Nil(t, pet)
}

func TestPgRepositoryByUserIDWrapsError(t *testing.T) {
	t.Parallel()

	tx := &pgtest.Tx{RowResults: []pgtest.Row{{Err: errDB}}}

	_, err := infra.NewPgRepository(nil).ByUserID(ctxWith(tx), userA)

	require.ErrorIs(t, err, errDB)
}

func TestPgRepositoryByUserIDForUpdateLocksThenSelects(t *testing.T) {
	t.Parallel()

	tx := &pgtest.Tx{RowResults: []pgtest.Row{{Values: petRow(t)}}}

	pet, err := infra.NewPgRepository(nil).ByUserIDForUpdate(ctxWith(tx), userA)

	require.NoError(t, err)
	require.NotNil(t, pet)
	require.Len(t, tx.ExecCalls, 1)
	requireSQL(t, tx.ExecCalls[0], "pg_advisory_xact_lock")
	assert.Equal(t, []any{userA}, tx.ExecCalls[0].Args)
	requireSQL(t, tx.QueryRowCalls[0], "FOR UPDATE")
}

func TestPgRepositoryByUserIDForUpdateLockErrorStopsQuery(t *testing.T) {
	t.Parallel()

	tx := &pgtest.Tx{ExecErrs: []error{errDB}}

	_, err := infra.NewPgRepository(nil).ByUserIDForUpdate(ctxWith(tx), userA)

	require.ErrorIs(t, err, errDB)
	assert.Empty(t, tx.QueryRowCalls)
}

func TestPgRepositorySaveUpsertsAllColumns(t *testing.T) {
	t.Parallel()

	tx := &pgtest.Tx{}
	pet := domain.New(userA, fixedTime)

	require.NoError(t, infra.NewPgRepository(nil).Save(ctxWith(tx), pet))
	require.Len(t, tx.ExecCalls, 1)
	requireSQL(t, tx.ExecCalls[0], "INSERT INTO pets", "ON CONFLICT (user_id) DO UPDATE")
	require.Len(t, tx.ExecCalls[0].Args, 17)
	assert.Equal(t, pet.ID(), tx.ExecCalls[0].Args[0])
	assert.Equal(t, userA, tx.ExecCalls[0].Args[1])
	assert.Equal(t, pet.Level(), tx.ExecCalls[0].Args[4])
	assert.Equal(t, pet.InteractionVersion(), tx.ExecCalls[0].Args[16])
}

func TestPgRepositorySaveWrapsError(t *testing.T) {
	t.Parallel()

	tx := &pgtest.Tx{ExecErrs: []error{errDB}}

	err := infra.NewPgRepository(nil).Save(ctxWith(tx), domain.New(userA, fixedTime))

	require.ErrorIs(t, err, errDB)
}

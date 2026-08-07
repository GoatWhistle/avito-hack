package infra_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/pet/app"
	"github.com/avito-hack/backend/internal/module/pet/infra"
	"github.com/avito-hack/backend/internal/shared/postgres/pgtest"
)

func hotStates() []app.HotState {
	return []app.HotState{
		{UserID: userA, Happiness: 80, Satiety: 60, Version: 3, UpdatedAt: fixedTime},
		{UserID: userB, Happiness: 50, Satiety: 40, Version: 9, UpdatedAt: fixedTime},
	}
}

func TestPgRepositorySaveHotStatesQueuesEveryState(t *testing.T) {
	t.Parallel()

	tx := &pgtest.BatchTx{Results: &pgtest.BatchResults{}}

	err := infra.NewPgRepository(nil).SaveHotStates(ctxWith(tx), hotStates())

	require.NoError(t, err)
	require.Len(t, tx.BatchCalls, 1)
	assert.Equal(t, 2, tx.BatchCalls[0].Len())
	assert.True(t, tx.Results.Closed)
}

func TestPgRepositorySaveHotStatesEmptyBatch(t *testing.T) {
	t.Parallel()

	tx := &pgtest.BatchTx{}

	err := infra.NewPgRepository(nil).SaveHotStates(ctxWith(tx), nil)

	require.NoError(t, err)
	require.Len(t, tx.BatchCalls, 1)
	assert.Zero(t, tx.BatchCalls[0].Len())
}

func TestPgRepositorySaveHotStatesWrapsExecError(t *testing.T) {
	t.Parallel()

	tx := &pgtest.BatchTx{Results: &pgtest.BatchResults{ExecErrs: []error{nil, errDB}}}

	err := infra.NewPgRepository(nil).SaveHotStates(ctxWith(tx), hotStates())

	require.ErrorIs(t, err, errDB)
	assert.True(t, tx.Results.Closed)
}

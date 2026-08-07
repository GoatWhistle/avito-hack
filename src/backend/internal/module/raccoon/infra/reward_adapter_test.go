package infra_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	petapp "github.com/avito-hack/backend/internal/module/pet/app"
	"github.com/avito-hack/backend/internal/module/raccoon/infra"
)

type stubTx struct {
	err    error
	called bool
}

func (s *stubTx) WithTx(ctx context.Context, fn func(context.Context) error) error {
	s.called = true
	if s.err != nil {
		return s.err
	}

	return fn(ctx)
}

type stubSigner struct{}

func (stubSigner) Issue(uuid.UUID, string) (string, error) { return "code", nil }

func (stubSigner) Verify(uuid.UUID, string, string) error { return nil }

func TestRewardAdapterActivateRejectsEmptyRewardID(t *testing.T) {
	t.Parallel()

	tx := &stubTx{}
	adapter := infra.NewRewardAdapter(petapp.NewRewardService(petapp.RewardServiceDeps{
		Signer: stubSigner{}, Tx: tx,
	}))

	code, err := adapter.Activate(context.Background(), uuid.New(), "")
	require.Error(t, err)
	assert.Empty(t, code)
	assert.False(t, tx.called)
}

func TestRewardAdapterActivateRequiresSigner(t *testing.T) {
	t.Parallel()

	tx := &stubTx{}
	adapter := infra.NewRewardAdapter(petapp.NewRewardService(petapp.RewardServiceDeps{Tx: tx}))

	code, err := adapter.Activate(context.Background(), uuid.New(), "reward-1")
	require.EqualError(t, err, "reward signer is not configured")
	assert.Empty(t, code)
	assert.False(t, tx.called)
}

func TestRewardAdapterActivatePropagatesTxError(t *testing.T) {
	t.Parallel()

	tx := &stubTx{err: errAdapter}
	adapter := infra.NewRewardAdapter(petapp.NewRewardService(petapp.RewardServiceDeps{
		Signer: stubSigner{}, Tx: tx,
	}))

	code, err := adapter.Activate(context.Background(), uuid.New(), "reward-1")
	require.ErrorIs(t, err, errAdapter)
	assert.Empty(t, code)
	assert.True(t, tx.called)
}

func TestNewRewardAdapterNotNil(t *testing.T) {
	t.Parallel()

	assert.NotNil(t, infra.NewRewardAdapter(petapp.NewRewardService(petapp.RewardServiceDeps{})))
}

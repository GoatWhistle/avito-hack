package vo_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/shared/vo"
)

func TestNewMoney(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		kopeks  int64
		wantErr error
		want    int64
	}{
		{name: "zero", kopeks: 0, want: 0},
		{name: "positive", kopeks: 12345, want: 12345},
		{name: "max int64", kopeks: 1<<62 - 1, want: 1<<62 - 1},
		{name: "negative", kopeks: -1, wantErr: vo.ErrNegativeMoney},
		{name: "large negative", kopeks: -999999, wantErr: vo.ErrNegativeMoney},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := vo.NewMoney(tt.kopeks)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				assert.True(t, got.IsZero())

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got.Kopeks())
		})
	}
}

func TestMustMoney(t *testing.T) {
	t.Parallel()

	assert.Equal(t, int64(500), vo.MustMoney(500).Kopeks())
	assert.Panics(t, func() { vo.MustMoney(-1) })
}

func TestZeroMoney(t *testing.T) {
	t.Parallel()

	zero := vo.ZeroMoney()

	assert.True(t, zero.IsZero())
	assert.Equal(t, int64(0), zero.Kopeks())
	assert.False(t, vo.MustMoney(1).IsZero())
}

func TestMoneyString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		kopeks int64
		want   string
	}{
		{name: "zero", kopeks: 0, want: "0.00"},
		{name: "kopeks only", kopeks: 5, want: "0.05"},
		{name: "two digit kopeks", kopeks: 42, want: "0.42"},
		{name: "whole rubles", kopeks: 10000, want: "100.00"},
		{name: "mixed", kopeks: 123456, want: "1234.56"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, vo.MustMoney(tt.kopeks).String())
		})
	}
}

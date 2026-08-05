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

func TestMoneyAdd(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		a    int64
		b    int64
		want int64
	}{
		{name: "zero plus zero", a: 0, b: 0, want: 0},
		{name: "zero plus value", a: 0, b: 250, want: 250},
		{name: "two values", a: 199, b: 801, want: 1000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sum := vo.MustMoney(tt.a).Add(vo.MustMoney(tt.b))

			assert.Equal(t, tt.want, sum.Kopeks())
		})
	}
}

func TestMoneyAddIsImmutable(t *testing.T) {
	t.Parallel()

	base := vo.MustMoney(100)
	_ = base.Add(vo.MustMoney(50))

	assert.Equal(t, int64(100), base.Kopeks())
}

func TestMoneyEqual(t *testing.T) {
	t.Parallel()

	assert.True(t, vo.MustMoney(100).Equal(vo.MustMoney(100)))
	assert.False(t, vo.MustMoney(100).Equal(vo.MustMoney(101)))
	assert.True(t, vo.ZeroMoney().Equal(vo.MustMoney(0)))
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

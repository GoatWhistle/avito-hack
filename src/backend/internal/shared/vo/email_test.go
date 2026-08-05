package vo_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/shared/vo"
)

func TestNewEmailValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		raw  string
		want string
	}{
		{name: "simple", raw: "user@example.com", want: "user@example.com"},
		{name: "uppercase normalized", raw: "User@Example.COM", want: "user@example.com"},
		{name: "surrounding spaces trimmed", raw: "  user@example.com  ", want: "user@example.com"},
		{name: "plus tag", raw: "user+tag@example.com", want: "user+tag@example.com"},
		{name: "subdomain", raw: "u@mail.example.co.uk", want: "u@mail.example.co.uk"},
		{name: "dots in local part", raw: "first.last@example.com", want: "first.last@example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := vo.NewEmail(tt.raw)

			require.NoError(t, err)
			assert.Equal(t, tt.want, got.String())
			assert.False(t, got.IsZero())
		})
	}
}

func TestNewEmailInvalid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		raw  string
	}{
		{name: "empty", raw: ""},
		{name: "only spaces", raw: "   "},
		{name: "no at sign", raw: "userexample.com"},
		{name: "no domain", raw: "user@"},
		{name: "no local part", raw: "@example.com"},
		{name: "double at", raw: "user@@example.com"},
		{name: "spaces inside", raw: "user name@example.com"},
		{name: "too long", raw: strings.Repeat("a", 250) + "@example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := vo.NewEmail(tt.raw)

			require.ErrorIs(t, err, vo.ErrInvalidEmail)
			assert.True(t, got.IsZero())
			assert.Empty(t, got.String())
		})
	}
}

func TestEmailZeroValue(t *testing.T) {
	t.Parallel()

	var email vo.Email

	assert.True(t, email.IsZero())
	assert.Empty(t, email.String())
}

func TestNewEmailAtMaxLength(t *testing.T) {
	t.Parallel()

	local := strings.Repeat("a", 254-len("@example.com"))

	got, err := vo.NewEmail(local + "@example.com")

	require.NoError(t, err)
	assert.Len(t, got.String(), 254)
}

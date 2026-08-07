package password_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/shared/password"
)

func TestNewHash_RejectsEmptyPlain(t *testing.T) {
	t.Parallel()

	_, err := password.NewHash("")

	require.ErrorIs(t, err, password.ErrTooShort)
}

func TestNewHash_RejectsPlainBeyondMaxLength(t *testing.T) {
	t.Parallel()

	_, err := password.NewHash(strings.Repeat("a", password.MaxLength+1))

	require.ErrorIs(t, err, password.ErrTooLong)
}

func TestNewHash_DistinctSaltsForSamePlain(t *testing.T) {
	t.Parallel()

	first, err := password.NewHash("correct horse battery")
	require.NoError(t, err)

	second, err := password.NewHash("correct horse battery")
	require.NoError(t, err)

	assert.NotEqual(t, first.String(), second.String())
	assert.NoError(t, first.Compare("correct horse battery"))
	assert.NoError(t, second.Compare("correct horse battery"))
}

func TestNewHash_HandlesMultibytePlain(t *testing.T) {
	t.Parallel()

	plain := "пароль-с-кириллицей"
	hash, err := password.NewHash(plain)
	require.NoError(t, err)

	assert.NoError(t, hash.Compare(plain))
	require.ErrorIs(t, hash.Compare(plain+"x"), password.ErrMismatch)
}

func TestRestoreHash_RejectsEmptyStored(t *testing.T) {
	t.Parallel()

	_, err := password.RestoreHash("")

	require.ErrorIs(t, err, password.ErrEmptyHash)
}

func TestRestoreHash_RoundTripsThroughString(t *testing.T) {
	t.Parallel()

	original, err := password.NewHash("round-trip-me")
	require.NoError(t, err)

	restored, err := password.RestoreHash(original.String())
	require.NoError(t, err)

	assert.Equal(t, original.String(), restored.String())
	assert.NoError(t, restored.Compare("round-trip-me"))
}

func TestCompare_RejectsMalformedStoredHash(t *testing.T) {
	t.Parallel()

	restored, err := password.RestoreHash("not-a-bcrypt-hash")
	require.NoError(t, err)

	require.ErrorIs(t, restored.Compare("anything"), password.ErrMismatch)
}

func TestCompare_ZeroValueHashAlwaysMismatches(t *testing.T) {
	t.Parallel()

	var zero password.Hash

	assert.Empty(t, zero.String())
	require.ErrorIs(t, zero.Compare(""), password.ErrMismatch)
	require.ErrorIs(t, zero.Compare("secret"), password.ErrMismatch)
}

func TestCompare_MismatchCases(t *testing.T) {
	t.Parallel()

	hash, err := password.NewHash("s3cret-value")
	require.NoError(t, err)

	cases := []struct {
		name  string
		plain string
	}{
		{name: "empty", plain: ""},
		{name: "prefix", plain: "s3cret"},
		{name: "case shifted", plain: "S3CRET-VALUE"},
		{name: "trailing space", plain: "s3cret-value "},
		{name: "leading space", plain: " s3cret-value"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			require.ErrorIs(t, hash.Compare(tc.plain), password.ErrMismatch)
		})
	}
}

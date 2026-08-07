package password_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/shared/password"
)

func TestNewHashLengthBoundaries(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		plain   string
		wantErr error
	}{
		{name: "empty", plain: "", wantErr: password.ErrTooShort},
		{name: "single character", plain: "a", wantErr: password.ErrTooShort},
		{name: "below minimum", plain: strings.Repeat("a", password.MinLength-1), wantErr: password.ErrTooShort},
		{name: "exactly minimum", plain: strings.Repeat("a", password.MinLength)},
		{name: "exactly maximum", plain: strings.Repeat("a", password.MaxLength)},
		{name: "above maximum", plain: strings.Repeat("a", password.MaxLength+10), wantErr: password.ErrTooLong},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			hash, err := password.NewHash(tt.plain)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				assert.Empty(t, hash.String())

				return
			}

			require.NoError(t, err)
			assert.NotEmpty(t, hash.String())
			assert.NotEqual(t, tt.plain, hash.String())
		})
	}
}

func TestHashCompare(t *testing.T) {
	t.Parallel()

	const plain = "correct-horse-battery"

	hash, err := password.NewHash(plain)
	require.NoError(t, err)

	require.NoError(t, hash.Compare(plain))

	tests := []struct {
		name  string
		plain string
	}{
		{name: "wrong password", plain: "wrong-horse-battery"},
		{name: "empty", plain: ""},
		{name: "case differs", plain: "Correct-Horse-Battery"},
		{name: "prefix only", plain: "correct-horse"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			require.ErrorIs(t, hash.Compare(tt.plain), password.ErrMismatch)
		})
	}
}

func TestNewHashIsSalted(t *testing.T) {
	t.Parallel()

	const plain = "same-password-twice"

	first, err := password.NewHash(plain)
	require.NoError(t, err)

	second, err := password.NewHash(plain)
	require.NoError(t, err)

	assert.NotEqual(t, first.String(), second.String())
	require.NoError(t, first.Compare(plain))
	require.NoError(t, second.Compare(plain))
}

func TestRestoreHash(t *testing.T) {
	t.Parallel()

	const plain = "restore-me-please"

	original, err := password.NewHash(plain)
	require.NoError(t, err)

	restored, err := password.RestoreHash(original.String())
	require.NoError(t, err)

	assert.Equal(t, original.String(), restored.String())
	require.NoError(t, restored.Compare(plain))
}

func TestRestoreHashEmpty(t *testing.T) {
	t.Parallel()

	hash, err := password.RestoreHash("")

	require.ErrorIs(t, err, password.ErrEmptyHash)
	assert.Empty(t, hash.String())
}

func TestRestoredGarbageHashFailsCompare(t *testing.T) {
	t.Parallel()

	hash, err := password.RestoreHash("not-a-bcrypt-hash")
	require.NoError(t, err)

	require.ErrorIs(t, hash.Compare("anything"), password.ErrMismatch)
}

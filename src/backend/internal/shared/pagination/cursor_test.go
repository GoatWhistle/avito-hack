package pagination_test

import (
	"encoding/base64"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/shared/pagination"
)

func TestCursorRoundTrip(t *testing.T) {
	t.Parallel()

	original := pagination.Cursor{
		CreatedAt: time.Date(2026, time.February, 3, 14, 25, 36, 123456789, time.UTC),
		ID:        uuid.New(),
	}

	decoded, err := pagination.DecodeCursor(original.Encode())

	require.NoError(t, err)
	assert.True(t, original.CreatedAt.Equal(decoded.CreatedAt))
	assert.Equal(t, original.ID, decoded.ID)
	assert.False(t, decoded.IsZero())
}

func TestCursorNormalizesTimezoneToUTC(t *testing.T) {
	t.Parallel()

	zone := time.FixedZone("MSK", 3*60*60)
	original := pagination.Cursor{
		CreatedAt: time.Date(2026, time.February, 3, 14, 0, 0, 0, zone),
		ID:        uuid.New(),
	}

	decoded, err := pagination.DecodeCursor(original.Encode())

	require.NoError(t, err)
	assert.Equal(t, time.UTC, decoded.CreatedAt.Location())
	assert.True(t, original.CreatedAt.Equal(decoded.CreatedAt))
}

func TestZeroCursorEncodesToEmptyString(t *testing.T) {
	t.Parallel()

	zero := pagination.Cursor{}

	assert.True(t, zero.IsZero())
	assert.Empty(t, zero.Encode())

	withTimeOnly := pagination.Cursor{CreatedAt: time.Now()}
	assert.True(t, withTimeOnly.IsZero())
	assert.Empty(t, withTimeOnly.Encode())
}

func TestDecodeCursorFailures(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		raw     string
		wantErr bool
	}{
		{name: "empty string yields zero cursor", raw: ""},
		{name: "not base64", raw: "!!!not-base64!!!", wantErr: true},
		{
			name:    "missing separator",
			raw:     base64.RawURLEncoding.EncodeToString([]byte("2026-02-03T14:00:00Z")),
			wantErr: true,
		},
		{
			name:    "invalid timestamp",
			raw:     base64.RawURLEncoding.EncodeToString([]byte("yesterday|" + uuid.NewString())),
			wantErr: true,
		},
		{
			name:    "invalid uuid",
			raw:     base64.RawURLEncoding.EncodeToString([]byte("2026-02-03T14:00:00Z|not-a-uuid")),
			wantErr: true,
		},
		{
			name:    "nil uuid is rejected",
			raw:     base64.RawURLEncoding.EncodeToString([]byte("2026-02-03T14:00:00Z|" + uuid.Nil.String())),
			wantErr: true,
		},
		{
			name:    "separator without payload",
			raw:     base64.RawURLEncoding.EncodeToString([]byte("|")),
			wantErr: true,
		},
		{
			name:    "padded base64 is rejected",
			raw:     "MjAyNi0wMi0wM1QxNDowMDowMFo=",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			cursor, err := pagination.DecodeCursor(tc.raw)

			if tc.wantErr {
				require.ErrorIs(t, err, pagination.ErrInvalidCursor)
				assert.True(t, cursor.IsZero())

				return
			}

			require.NoError(t, err)
			assert.True(t, cursor.IsZero())
		})
	}
}

func TestNormalizeLimit(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		limit int
		want  int
	}{
		{name: "zero falls back to default", limit: 0, want: pagination.DefaultLimit},
		{name: "negative falls back to default", limit: -42, want: pagination.DefaultLimit},
		{name: "one stays one", limit: 1, want: 1},
		{name: "value within range is kept", limit: 37, want: 37},
		{name: "max is kept", limit: pagination.MaxLimit, want: pagination.MaxLimit},
		{name: "above max is clamped", limit: pagination.MaxLimit + 1, want: pagination.MaxLimit},
		{name: "huge value is clamped", limit: 1 << 20, want: pagination.MaxLimit},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.want, pagination.NormalizeLimit(tc.limit))
		})
	}
}

package infra

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/module/item/domain"
	"github.com/avito-hack/backend/internal/shared/postgres/pgtest"
	"github.com/avito-hack/backend/internal/shared/vo"
)

var mapperTime = time.Date(2023, 9, 20, 14, 0, 0, 0, time.UTC)

func validItemRow(id uuid.UUID) itemRow {
	return itemRow{
		id:          id,
		ownerID:     uuid.New(),
		title:       "Bike",
		description: "A fast bike",
		priceKopeks: 150000,
		status:      string(domain.StatusPublished),
		attributes:  []byte(`{"color":"red"}`),
		createdAt:   mapperTime,
		updatedAt:   mapperTime,
	}
}

func TestItemToDomainRestoresItem(t *testing.T) {
	t.Parallel()

	id := uuid.New()
	row := validItemRow(id)

	item, err := toDomain(row)

	require.NoError(t, err)
	assert.Equal(t, id, item.ID())
	assert.Equal(t, row.ownerID, item.OwnerID())
	assert.Equal(t, "Bike", item.Title())
	assert.Equal(t, "A fast bike", item.Description())
	assert.Equal(t, int64(150000), item.Price().Kopeks())
	assert.Equal(t, domain.StatusPublished, item.Status())
	assert.Equal(t, domain.Attributes{"color": "red"}, item.Attributes())
	assert.Equal(t, mapperTime, item.CreatedAt())
}

func TestItemToDomainRejectsNegativePrice(t *testing.T) {
	t.Parallel()

	row := validItemRow(uuid.New())
	row.priceKopeks = -1

	_, err := toDomain(row)

	require.ErrorIs(t, err, vo.ErrNegativeMoney)
	assert.Contains(t, err.Error(), "restore price for item")
}

func TestItemToDomainRejectsMalformedAttributes(t *testing.T) {
	t.Parallel()

	row := validItemRow(uuid.New())
	row.attributes = []byte(`{"color":`)

	_, err := toDomain(row)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "restore attributes for item")
}

func TestItemToDomainAcceptsEmptyAttributes(t *testing.T) {
	t.Parallel()

	row := validItemRow(uuid.New())
	row.attributes = nil

	item, err := toDomain(row)

	require.NoError(t, err)
	assert.Equal(t, 0, item.Attributes().Len())
}

func TestScanItemReadsAllColumns(t *testing.T) {
	t.Parallel()

	id, ownerID := uuid.New(), uuid.New()
	row := pgtest.Row{Values: []any{
		id, ownerID, "Chair", "Wooden", int64(500),
		string(domain.StatusDraft), []byte(`{}`), mapperTime, mapperTime,
	}}

	item, err := scanItem(row)

	require.NoError(t, err)
	assert.Equal(t, id, item.ID())
	assert.Equal(t, ownerID, item.OwnerID())
	assert.Equal(t, domain.StatusDraft, item.Status())
}

func TestScanItemPropagatesErrors(t *testing.T) {
	t.Parallel()

	sentinel := errors.New("scan failed")

	_, err := scanItem(pgtest.Row{Err: sentinel})

	require.ErrorIs(t, err, sentinel)
}

func TestScanPhotoReadsRow(t *testing.T) {
	t.Parallel()

	id, itemID := uuid.New(), uuid.New()
	row := pgtest.Row{Values: []any{id, itemID, "http://cdn/p.jpg", 2, mapperTime}}

	photo, err := scanPhoto(row)

	require.NoError(t, err)
	assert.Equal(t, id, photo.ID())
	assert.Equal(t, itemID, photo.ItemID())
	assert.Equal(t, "http://cdn/p.jpg", photo.URL())
	assert.Equal(t, 2, photo.Position())
	assert.Equal(t, mapperTime, photo.CreatedAt())
}

func TestScanPhotoWrapsScanError(t *testing.T) {
	t.Parallel()

	sentinel := errors.New("boom")

	_, err := scanPhoto(pgtest.Row{Err: sentinel})

	require.ErrorIs(t, err, sentinel)
	assert.Contains(t, err.Error(), "scan photo")
}

func TestStatusFrom(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		raw  string
		want domain.Status
	}{
		{name: "published", raw: "published", want: domain.StatusPublished},
		{name: "archived", raw: "archived", want: domain.StatusArchived},
		{name: "unknown", raw: "wat", want: domain.StatusDraft},
		{name: "empty", raw: "", want: domain.StatusDraft},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.want, statusFrom(tc.raw))
		})
	}
}

func TestDecodeAttributes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		raw     []byte
		want    domain.Attributes
		wantErr bool
	}{
		{name: "nil", raw: nil, want: domain.Attributes{}},
		{name: "empty slice", raw: []byte{}, want: domain.Attributes{}},
		{name: "empty object", raw: []byte(`{}`), want: domain.Attributes{}},
		{name: "values", raw: []byte(`{"a":"b"}`), want: domain.Attributes{"a": "b"}},
		{name: "malformed", raw: []byte(`{`), wantErr: true},
		{name: "wrong shape", raw: []byte(`{"a":1}`), wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := decodeAttributes(tc.raw)

			if tc.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "decode attributes")

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestEncodeAttributes(t *testing.T) {
	t.Parallel()

	raw, err := encodeAttributes(domain.Attributes{"size": "L"})

	require.NoError(t, err)
	assert.JSONEq(t, `{"size":"L"}`, string(raw))

	nilAttrs, err := encodeAttributes(nil)

	require.NoError(t, err)
	assert.JSONEq(t, `null`, string(nilAttrs))

	empty, err := encodeAttributes(domain.Attributes{})

	require.NoError(t, err)
	assert.JSONEq(t, `{}`, string(empty))
}

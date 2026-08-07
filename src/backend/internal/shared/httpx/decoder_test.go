package httpx_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-hack/backend/internal/shared/httpx"
)

const bigBody = 1 << 20

type decodePayload struct {
	Name string `json:"name"`
}

type stubValidator struct{ err error }

func (s stubValidator) Struct(any) error { return s.err }

func postJSON(body string) (*httptest.ResponseRecorder, *http.Request) {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))

	return httptest.NewRecorder(), req
}

func TestDecoderAcceptsValidPayload(t *testing.T) {
	t.Parallel()

	rec, req := postJSON(`{"name":"raccoon"}`)
	decoder := httpx.NewDecoder(stubValidator{}, bigBody)

	var dst decodePayload
	ok := decoder.Decode(rec, req, &dst)

	require.True(t, ok)
	assert.Equal(t, "raccoon", dst.Name)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestDecoderRejectsMalformedJSON(t *testing.T) {
	t.Parallel()

	rec, req := postJSON(`{"name":`)
	decoder := httpx.NewDecoder(stubValidator{}, bigBody)

	var dst decodePayload
	ok := decoder.Decode(rec, req, &dst)

	require.False(t, ok)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestDecoderRejectsUnknownFields(t *testing.T) {
	t.Parallel()

	rec, req := postJSON(`{"name":"a","extra":1}`)
	decoder := httpx.NewDecoder(stubValidator{}, bigBody)

	var dst decodePayload
	ok := decoder.Decode(rec, req, &dst)

	require.False(t, ok)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestDecoderReportsValidationFailure(t *testing.T) {
	t.Parallel()

	rec, req := postJSON(`{"name":"a"}`)
	decoder := httpx.NewDecoder(stubValidator{err: errors.New("invalid")}, bigBody)

	var dst decodePayload
	ok := decoder.Decode(rec, req, &dst)

	require.False(t, ok)
	assert.GreaterOrEqual(t, rec.Code, http.StatusBadRequest)
}

func TestDecoderWithoutValidatorSkipsValidation(t *testing.T) {
	t.Parallel()

	rec, req := postJSON(`{"name":"a"}`)
	decoder := httpx.NewDecoder(nil, bigBody)

	var dst decodePayload
	ok := decoder.Decode(rec, req, &dst)

	require.True(t, ok)
	assert.Equal(t, "a", dst.Name)
}

func TestDecoderEnforcesBodyLimit(t *testing.T) {
	t.Parallel()

	rec, req := postJSON(`{"name":"` + strings.Repeat("x", 100) + `"}`)
	decoder := httpx.NewDecoder(stubValidator{}, 8)

	var dst decodePayload
	ok := decoder.Decode(rec, req, &dst)

	require.False(t, ok)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

package httpx

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/avito-hack/backend/internal/shared/domainerr"
)

var ErrMalformedBody = errors.New("malformed request body")

func DecodeJSON(w http.ResponseWriter, r *http.Request, maxBytes int64, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxBytes)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		return fmt.Errorf("%w: %s", ErrMalformedBody, err.Error())
	}

	return nil
}

func UUIDParam(r *http.Request, name string) (uuid.UUID, error) {
	raw := chi.URLParam(r, name)

	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, domainerr.NewInvalid(name, "invalid identifier")
	}

	return id, nil
}

func OptionalUUIDQuery(r *http.Request, name string) (uuid.UUID, bool, error) {
	raw := r.URL.Query().Get(name)
	if raw == "" {
		return uuid.Nil, false, nil
	}

	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, false, domainerr.NewInvalid(name, "invalid identifier")
	}

	return id, true, nil
}

func IntQuery(r *http.Request, name string, fallback int) (int, error) {
	raw := r.URL.Query().Get(name)
	if raw == "" {
		return fallback, nil
	}

	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0, domainerr.NewInvalid(name, "integer expected")
	}

	return value, nil
}

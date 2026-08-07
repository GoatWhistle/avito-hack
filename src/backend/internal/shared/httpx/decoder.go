package httpx

import (
	"net/http"

	"github.com/avito-hack/backend/internal/shared/apierr"
)

type Validator interface {
	Struct(dst any) error
}

type Decoder struct {
	validator    Validator
	maxBodyBytes int64
}

func NewDecoder(v Validator, maxBodyBytes int64) Decoder {
	return Decoder{validator: v, maxBodyBytes: maxBodyBytes}
}

func (d Decoder) Decode(w http.ResponseWriter, r *http.Request, dst any) bool {
	if err := DecodeJSON(w, r, d.maxBodyBytes, dst); err != nil {
		apierr.WriteBadRequest(w, r, err.Error())

		return false
	}

	if d.validator == nil {
		return true
	}

	if err := d.validator.Struct(dst); err != nil {
		apierr.Write(w, r, err)

		return false
	}

	return true
}

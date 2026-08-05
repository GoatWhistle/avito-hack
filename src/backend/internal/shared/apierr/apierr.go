package apierr

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/avito-hack/backend/internal/shared/domainerr"
)

const (
	CodeNotFound     = "not_found"
	CodeForbidden    = "forbidden"
	CodeUnauthorized = "unauthorized"
	CodeConflict     = "conflict"
	CodeValidation   = "validation_error"
	CodeBadRequest   = "bad_request"
	CodeInternal     = "internal_error"
)

type errorBody struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	Field     string `json:"field,omitempty"`
	RequestID string `json:"request_id,omitempty"`
}

type errorEnvelope struct {
	Error errorBody `json:"error"`
}

func Write(w http.ResponseWriter, r *http.Request, err error) {
	status, body := classify(r, err)

	if status == http.StatusInternalServerError {
		slog.ErrorContext(r.Context(), "unhandled error",
			slog.Any("error", err),
			slog.String("request_id", body.RequestID),
		)
	}

	respond(w, status, errorEnvelope{Error: body})
}

func classify(r *http.Request, err error) (int, errorBody) {
	body := errorBody{RequestID: middleware.GetReqID(r.Context())}

	var (
		invalid  *domainerr.InvalidError
		conflict *domainerr.ConflictError
	)

	switch {
	case errors.As(err, &invalid):
		body.Code, body.Message, body.Field = CodeValidation, invalid.Message, invalid.Field
		return http.StatusBadRequest, body

	case errors.As(err, &conflict):
		body.Code, body.Message = CodeConflict, conflict.Message
		return http.StatusConflict, body

	case errors.Is(err, domainerr.ErrNotFound):
		body.Code, body.Message = CodeNotFound, err.Error()
		return http.StatusNotFound, body

	case errors.Is(err, domainerr.ErrForbidden):
		body.Code, body.Message = CodeForbidden, err.Error()
		return http.StatusForbidden, body

	case errors.Is(err, domainerr.ErrUnauthorized):
		body.Code, body.Message = CodeUnauthorized, err.Error()
		return http.StatusUnauthorized, body

	case errors.Is(err, domainerr.ErrConflict):
		body.Code, body.Message = CodeConflict, err.Error()
		return http.StatusConflict, body

	default:
		body.Code, body.Message = CodeInternal, "internal server error"
		return http.StatusInternalServerError, body
	}
}

func WriteBadRequest(w http.ResponseWriter, r *http.Request, message string) {
	respond(w, http.StatusBadRequest, errorEnvelope{Error: errorBody{
		Code:      CodeBadRequest,
		Message:   message,
		RequestID: middleware.GetReqID(r.Context()),
	}})
}

func respond(w http.ResponseWriter, status int, payload errorEnvelope) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		slog.Error("encode error response", slog.Any("error", err))
	}
}

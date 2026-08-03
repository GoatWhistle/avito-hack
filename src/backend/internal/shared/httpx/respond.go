package httpx

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func JSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if payload == nil {
		return
	}

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		slog.Error("encode response", slog.Any("error", err))
	}
}

func OK(w http.ResponseWriter, payload any) {
	JSON(w, http.StatusOK, payload)
}

func Created(w http.ResponseWriter, payload any) {
	JSON(w, http.StatusCreated, payload)
}

func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

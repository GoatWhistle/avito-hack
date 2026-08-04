package main

import (
	"errors"
	"log/slog"
	"time"

	"github.com/avito-hack/backend/internal/shared/logger"
)

func main() {
	log := logger.New(logger.Options{Level: slog.LevelDebug, Format: "pretty", Color: "always"})
	log.Info("server started", slog.String("addr", ":8080"), slog.String("version", "dev"))
	log.Info("http request", slog.String("method", "GET"), slog.String("path", "/healthz"),
		slog.Int("status", 200), slog.Int("bytes", 15), slog.Duration("duration", 34137*time.Nanosecond),
		slog.String("request_id", "2f6c52301a10/gZihcvlq5L-000264"), slog.String("remote_ip", "127.0.0.1:58788"))
	log.Info("http request", slog.String("method", "POST"), slog.String("path", "/api/v1/items"),
		slog.Int("status", 201), slog.Int("bytes", 2048), slog.Duration("duration", 137*time.Millisecond),
		slog.String("request_id", "abc/xyz-000265"), slog.String("remote_ip", "10.0.0.5:4411"))
	log.Warn("http request", slog.String("method", "DELETE"), slog.String("path", "/api/v1/items/42"),
		slog.Int("status", 404), slog.Int("bytes", 88), slog.Duration("duration", 720*time.Millisecond),
		slog.String("request_id", "abc/xyz-000266"), slog.String("remote_ip", "10.0.0.5:4412"))
	log.Error("connect database", slog.Any("error", errors.New("dial tcp: connection refused")))
	log.Debug("cache miss", slog.String("key", "user:42"))
}

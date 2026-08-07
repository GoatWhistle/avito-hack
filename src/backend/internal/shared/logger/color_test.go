package logger

import (
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLevelLabel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		level slog.Level
		want  string
	}{
		{name: "debug", level: slog.LevelDebug, want: "DEBUG"},
		{name: "below debug", level: slog.LevelDebug - 1, want: "DEBUG"},
		{name: "info", level: slog.LevelInfo, want: "INFO"},
		{name: "warn", level: slog.LevelWarn, want: "WARNING"},
		{name: "error", level: slog.LevelError, want: "ERROR"},
		{name: "above error", level: slog.LevelError + 4, want: "ERROR"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, levelLabel(tt.level))
		})
	}
}

func TestLevelColor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		level slog.Level
		want  string
	}{
		{name: "debug is magenta", level: slog.LevelDebug, want: ansiMagenta},
		{name: "info is green", level: slog.LevelInfo, want: ansiGreen},
		{name: "warn is yellow", level: slog.LevelWarn, want: ansiYellow},
		{name: "error is red", level: slog.LevelError, want: ansiRed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, levelColor(tt.level))
		})
	}
}

func TestMethodColor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		method string
		want   string
	}{
		{name: "get", method: "GET", want: ansiBlue},
		{name: "post", method: "POST", want: ansiGreen},
		{name: "put", method: "PUT", want: ansiYellow},
		{name: "patch", method: "PATCH", want: ansiYellow},
		{name: "delete", method: "DELETE", want: ansiRed},
		{name: "options falls back", method: "OPTIONS", want: ansiMagenta},
		{name: "lowercase is not matched", method: "get", want: ansiMagenta},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, methodColor(tt.method))
		})
	}
}

func TestStatusColor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		status int
		want   string
	}{
		{name: "ok is green", status: 200, want: ansiGreen},
		{name: "just below redirect", status: 299, want: ansiGreen},
		{name: "redirect is cyan", status: 301, want: ansiCyan},
		{name: "client error is yellow", status: 404, want: ansiYellow},
		{name: "server error is red", status: 500, want: ansiRed},
		{name: "above server error", status: 503, want: ansiRed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, statusColor(tt.status))
		})
	}
}

func TestDurationColor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		elapsed time.Duration
		want    string
	}{
		{name: "fast is dim", elapsed: 10 * time.Millisecond, want: ansiDim},
		{name: "medium is yellow", elapsed: 200 * time.Millisecond, want: ansiYellow},
		{name: "exactly medium", elapsed: mediumRequest, want: ansiYellow},
		{name: "slow is red", elapsed: time.Second, want: ansiRed},
		{name: "exactly slow", elapsed: slowRequest, want: ansiRed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, durationColor(tt.elapsed))
		})
	}
}

package logger

import (
	"bytes"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func accessOutput(t *testing.T, color bool, attrs ...slog.Attr) string {
	t.Helper()

	buf := &bytes.Buffer{}
	h := NewPrettyHandler(buf, PrettyOptions{Color: color})

	handle(t, h, fixedRecord(slog.LevelInfo, accessLogMessage, attrs...))

	return buf.String()
}

func TestWriteAccessLogRendersFullLine(t *testing.T) {
	t.Parallel()

	out := accessOutput(t, false,
		slog.String("method", "GET"),
		slog.String("path", "/api/v1/items"),
		slog.Int("status", 200),
		slog.Duration("duration", 12*time.Millisecond),
		slog.Int("bytes", 2048),
		slog.String("request_id", "host/abcdef123456"),
		slog.String("remote_ip", "10.0.0.1:5555"),
	)

	assert.Contains(t, out, "GET   ")
	assert.Contains(t, out, "/api/v1/items")
	assert.Contains(t, out, "200")
	assert.Contains(t, out, "12.0ms")
	assert.Contains(t, out, "2.0KB")
	assert.Contains(t, out, "#ef123456")
	assert.Contains(t, out, "from 10.0.0.1:5555")
}

func TestWriteAccessLogOmitsOptionalAttrs(t *testing.T) {
	t.Parallel()

	out := accessOutput(t, false,
		slog.String("method", "POST"),
		slog.String("path", "/api/v1/items"),
		slog.Int("status", 201),
	)

	assert.Contains(t, out, "POST")
	assert.Contains(t, out, "201")
	assert.NotContains(t, out, "#")
	assert.NotContains(t, out, "from ")
	assert.NotContains(t, out, "B ")
}

func TestWriteAccessLogSkipsEmptyRequestIDAndIP(t *testing.T) {
	t.Parallel()

	out := accessOutput(t, false,
		slog.String("method", "GET"),
		slog.String("path", "/health"),
		slog.Int("status", 200),
		slog.String("request_id", ""),
		slog.String("remote_ip", ""),
	)

	assert.NotContains(t, out, "#")
	assert.NotContains(t, out, "from ")
}

func TestWriteAccessLogConsumesKnownKeysOnly(t *testing.T) {
	t.Parallel()

	out := accessOutput(t, false,
		slog.String("method", "GET"),
		slog.String("path", "/x"),
		slog.Int("status", 200),
		slog.String("tenant", "avito"),
	)

	assert.Contains(t, out, "tenant=avito")
	assert.NotContains(t, out, "method=")
	assert.NotContains(t, out, "path=")
	assert.NotContains(t, out, "status=")
}

func TestWriteAccessLogColorizesByStatusAndMethod(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		method string
		status int
		want   string
	}{
		{name: "server error is red", method: "GET", status: 500, want: ansiRed},
		{name: "client error is yellow", method: "GET", status: 404, want: ansiYellow},
		{name: "redirect is cyan", method: "GET", status: 302, want: ansiCyan},
		{name: "success is green", method: "GET", status: 200, want: ansiGreen},
		{name: "delete method is red", method: "DELETE", status: 204, want: ansiRed},
		{name: "unknown method is magenta", method: "TRACE", status: 200, want: ansiMagenta},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			out := accessOutput(t, true,
				slog.String("method", tc.method),
				slog.String("path", "/x"),
				slog.Int("status", tc.status),
			)

			assert.Contains(t, out, tc.want)
		})
	}
}

func TestWriteAccessLogColorsSlowDurations(t *testing.T) {
	t.Parallel()

	slow := accessOutput(t, true,
		slog.String("method", "GET"), slog.String("path", "/x"), slog.Int("status", 200),
		slog.Duration("duration", 1500*time.Millisecond),
	)
	medium := accessOutput(t, true,
		slog.String("method", "GET"), slog.String("path", "/x"), slog.Int("status", 200),
		slog.Duration("duration", 250*time.Millisecond),
	)
	fast := accessOutput(t, true,
		slog.String("method", "GET"), slog.String("path", "/x"), slog.Int("status", 200),
		slog.Duration("duration", 400*time.Microsecond),
	)

	assert.Contains(t, slow, ansiRed)
	assert.Contains(t, slow, "1.50s")
	assert.Contains(t, medium, ansiYellow)
	assert.Contains(t, medium, "250.0ms")
	assert.Contains(t, fast, "400µs")
}

func TestShortRequestID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		id   string
		want string
	}{
		{name: "chi style id keeps tail after slash", id: "hostname/aBcDeF/000001", want: "000001"},
		{name: "long tail is truncated to eight runes", id: "0123456789abcdef", want: "89abcdef"},
		{name: "short id is unchanged", id: "abc123", want: "abc123"},
		{name: "exactly eight is unchanged", id: "12345678", want: "12345678"},
		{name: "empty stays empty", id: "", want: ""},
		{name: "trailing slash keeps whole string", id: "host/", want: "host/"},
		{name: "leading slash keeps tail", id: "/abc", want: "abc"},
		{name: "multibyte tail is cut by runes", id: "аааааааааааа", want: "аааааааа"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.want, shortRequestID(tc.id))
		})
	}
}

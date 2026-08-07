package logger

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func fixedRecord(level slog.Level, msg string, attrs ...slog.Attr) slog.Record {
	rec := slog.NewRecord(time.Date(2024, time.March, 1, 12, 30, 45, 0, time.UTC), level, msg, 0)
	rec.AddAttrs(attrs...)

	return rec
}

func handle(t *testing.T, h slog.Handler, rec slog.Record) {
	t.Helper()

	require.NoError(t, h.Handle(context.Background(), rec))
}

func TestPrettyHandlerWritesHeader(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	h := NewPrettyHandler(buf, PrettyOptions{})

	handle(t, h, fixedRecord(slog.LevelInfo, "started"))

	out := buf.String()
	assert.Contains(t, out, "12:30:45.000")
	assert.Contains(t, out, "INFO")
	assert.Contains(t, out, "started")
	assert.True(t, strings.HasSuffix(out, "\n"))
}

func TestPrettyHandlerDefaultLevelIsInfo(t *testing.T) {
	t.Parallel()

	h := NewPrettyHandler(&bytes.Buffer{}, PrettyOptions{})

	assert.False(t, h.Enabled(context.Background(), slog.LevelDebug))
	assert.True(t, h.Enabled(context.Background(), slog.LevelInfo))
	assert.True(t, h.Enabled(context.Background(), slog.LevelError))
}

func TestPrettyHandlerRespectsConfiguredLevel(t *testing.T) {
	t.Parallel()

	h := NewPrettyHandler(&bytes.Buffer{}, PrettyOptions{Level: slog.LevelWarn})

	assert.False(t, h.Enabled(context.Background(), slog.LevelInfo))
	assert.True(t, h.Enabled(context.Background(), slog.LevelWarn))
}

func TestPrettyHandlerColorWrapsAnsi(t *testing.T) {
	t.Parallel()

	plain := &bytes.Buffer{}
	colored := &bytes.Buffer{}

	handle(t, NewPrettyHandler(plain, PrettyOptions{}), fixedRecord(slog.LevelError, "boom"))
	handle(t, NewPrettyHandler(colored, PrettyOptions{Color: true}), fixedRecord(slog.LevelError, "boom"))

	assert.NotContains(t, plain.String(), ansiReset)
	assert.Contains(t, colored.String(), ansiReset)
	assert.Contains(t, colored.String(), ansiRed)
}

func TestPrettyHandlerRendersAttrs(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	h := NewPrettyHandler(buf, PrettyOptions{})

	handle(t, h, fixedRecord(slog.LevelInfo, "msg", slog.String("user", "ivan"), slog.Int("count", 3)))

	out := buf.String()
	assert.Contains(t, out, "user=ivan")
	assert.Contains(t, out, "count=3")
}

func TestPrettyHandlerStackAttrIsIndentedOnNewLine(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	h := NewPrettyHandler(buf, PrettyOptions{})

	handle(t, h, fixedRecord(slog.LevelError, "crash", slog.String("stack", "line one\nline two")))

	assert.Contains(t, buf.String(), "stack=\n    line one\n    line two")
}

func TestPrettyHandlerStackAttrKeepsRealTraceReadable(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	h := NewPrettyHandler(buf, PrettyOptions{})

	trace := "goroutine 1 [running]:\nruntime/debug.Stack()\n\t/usr/local/go/src/runtime/debug/stack.go:26 +0x5e\n"

	handle(t, h, fixedRecord(slog.LevelError, "panic recovered", slog.String("stack", trace)))

	out := buf.String()
	assert.Contains(t, out, "stack=\n    goroutine 1 [running]:")
	assert.Contains(t, out, "\n    runtime/debug.Stack()")
	assert.NotContains(t, out, `\n`)
}

func TestIndentBlock(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "single line", input: "only", want: "    only"},
		{name: "two lines", input: "a\nb", want: "    a\n    b"},
		{name: "trailing newline is trimmed", input: "a\nb\n", want: "    a\n    b"},
		{name: "many trailing newlines are trimmed", input: "a\n\n\n", want: "    a"},
		{name: "blank interior line is still indented", input: "a\n\nb", want: "    a\n    \n    b"},
		{name: "empty string", input: "", want: "    "},
		{name: "already indented gains more indent", input: "\tx", want: "    \tx"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.want, indentBlock(tc.input))
		})
	}
}

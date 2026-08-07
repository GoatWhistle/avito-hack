package logger

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewReturnsUsableLogger(t *testing.T) {
	tests := []struct {
		name   string
		format string
	}{
		{name: "json format", format: formatJSON},
		{name: "pretty format", format: formatPretty},
		{name: "unknown format falls back to pretty", format: "yaml"},
		{name: "empty format falls back to pretty", format: ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("NO_COLOR", "1")

			log := New(Options{Level: slog.LevelWarn, Format: tc.format})

			require.NotNil(t, log)
			assert.False(t, log.Enabled(context.Background(), slog.LevelInfo))
			assert.True(t, log.Enabled(context.Background(), slog.LevelWarn))
		})
	}
}

func TestNewJSONHandlerTypeDiffersFromPretty(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	jsonLog := New(Options{Format: formatJSON})
	prettyLog := New(Options{Format: formatPretty})

	_, isJSON := jsonLog.Handler().(*slog.JSONHandler)
	_, prettyIsJSON := prettyLog.Handler().(*slog.JSONHandler)

	assert.True(t, isJSON)
	assert.False(t, prettyIsJSON)
}

func TestResolveFormat(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "json", input: "json", want: formatJSON},
		{name: "pretty", input: "pretty", want: formatPretty},
		{name: "text alias", input: "text", want: formatPretty},
		{name: "console alias", input: "console", want: formatPretty},
		{name: "unknown", input: "logfmt", want: formatPretty},
		{name: "empty", input: "", want: formatPretty},
		{name: "case sensitive JSON is not json", input: "JSON", want: formatPretty},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.want, resolveFormat(tc.input))
		})
	}
}

func TestResolveColorNoColorEnvWins(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	for _, mode := range []string{"always", "on", "true", "1", "auto", ""} {
		assert.False(t, resolveColor(mode), mode)
	}
}

func TestResolveColorDisabledModes(t *testing.T) {
	t.Setenv("NO_COLOR", "")

	for _, mode := range []string{"never", "off", "false", "0"} {
		assert.False(t, resolveColor(mode), mode)
	}
}

func TestResolveColorEnabledModes(t *testing.T) {
	t.Setenv("NO_COLOR", "")

	for _, mode := range []string{"always", "on", "true", "1"} {
		assert.True(t, resolveColor(mode), mode)
	}
}

func TestResolveColorAutoDependsOnStdout(t *testing.T) {
	t.Setenv("NO_COLOR", "")

	assert.Equal(t, isTerminal(os.Stdout), resolveColor("auto"))
	assert.Equal(t, isTerminal(os.Stdout), resolveColor("unrecognised"))
}

func TestIsTerminalOnRegularFileIsFalse(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "out.log")

	f, err := os.Create(path)
	require.NoError(t, err)

	t.Cleanup(func() { _ = f.Close() })

	assert.False(t, isTerminal(f))
}

func TestIsTerminalOnClosedFileIsFalse(t *testing.T) {
	t.Parallel()

	f, err := os.Create(filepath.Join(t.TempDir(), "closed.log"))
	require.NoError(t, err)
	require.NoError(t, f.Close())

	assert.False(t, isTerminal(f))
}

func TestIsTerminalOnPipeIsFalse(t *testing.T) {
	t.Parallel()

	reader, writer, err := os.Pipe()
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = reader.Close()
		_ = writer.Close()
	})

	assert.False(t, isTerminal(writer))
	assert.False(t, isTerminal(reader))
}

func TestIsTerminalOnCharDeviceIsTrue(t *testing.T) {
	t.Parallel()

	f, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0o600)
	require.NoError(t, err)

	t.Cleanup(func() { _ = f.Close() })

	assert.True(t, isTerminal(f))
}

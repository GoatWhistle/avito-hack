package logger

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrettyHandlerWithAttrs(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	h := NewPrettyHandler(buf, PrettyOptions{}).WithAttrs([]slog.Attr{slog.String("service", "api")})

	handle(t, h, fixedRecord(slog.LevelInfo, "msg"))

	assert.Contains(t, buf.String(), "service=api")
}

func TestPrettyHandlerWithAttrsEmptyReturnsSameHandler(t *testing.T) {
	t.Parallel()

	h := NewPrettyHandler(&bytes.Buffer{}, PrettyOptions{})

	assert.Same(t, h, h.WithAttrs(nil))
}

func TestPrettyHandlerWithGroupEmptyReturnsSameHandler(t *testing.T) {
	t.Parallel()

	h := NewPrettyHandler(&bytes.Buffer{}, PrettyOptions{})

	assert.Same(t, h, h.WithGroup(""))
}

func TestPrettyHandlerWithGroupQualifiesKeys(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	h := NewPrettyHandler(buf, PrettyOptions{}).WithGroup("http")

	handle(t, h, fixedRecord(slog.LevelInfo, "msg", slog.String("path", "/items")))

	assert.Contains(t, buf.String(), "http.path=/items")
}

func TestPrettyHandlerNestedGroupsJoinWithDot(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	h := NewPrettyHandler(buf, PrettyOptions{}).WithGroup("a").WithGroup("b")

	handle(t, h, fixedRecord(slog.LevelInfo, "msg", slog.String("k", "v")))

	assert.Contains(t, buf.String(), "a.b.k=v")
}

func TestPrettyHandlerWithAttrsDoesNotLeakToSibling(t *testing.T) {
	t.Parallel()

	base := NewPrettyHandler(&bytes.Buffer{}, PrettyOptions{})
	left := base.WithAttrs([]slog.Attr{slog.String("side", "left")})

	rightBuf := &bytes.Buffer{}
	right := NewPrettyHandler(rightBuf, PrettyOptions{}).WithAttrs([]slog.Attr{slog.String("side", "right")})

	handle(t, left, fixedRecord(slog.LevelInfo, "msg"))
	handle(t, right, fixedRecord(slog.LevelInfo, "msg"))

	assert.Contains(t, rightBuf.String(), "side=right")
	assert.NotContains(t, rightBuf.String(), "side=left")
}

func TestPrettyHandlerRecordAttrOverridesHandlerAttr(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	h := NewPrettyHandler(buf, PrettyOptions{}).WithAttrs([]slog.Attr{slog.String("env", "dev")})

	handle(t, h, fixedRecord(slog.LevelInfo, "msg", slog.String("env", "prod")))

	out := buf.String()
	assert.Contains(t, out, "env=prod")
	assert.NotContains(t, out, "env=dev")
}

func TestPrettyHandlerZeroTimeUsesNow(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	h := NewPrettyHandler(buf, PrettyOptions{})

	require.NoError(t, h.Handle(context.Background(), slog.NewRecord(time.Time{}, slog.LevelInfo, "msg", 0)))

	assert.NotEmpty(t, buf.String())
}

func TestPrettyHandlerConcurrentWritesAreSerialized(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	h := NewPrettyHandler(buf, PrettyOptions{})

	var wg sync.WaitGroup

	for range 50 {
		wg.Add(1)

		go func() {
			defer wg.Done()

			_ = h.Handle(context.Background(), fixedRecord(slog.LevelInfo, "concurrent"))
		}()
	}

	wg.Wait()

	assert.Equal(t, 50, strings.Count(buf.String(), "concurrent"))
}

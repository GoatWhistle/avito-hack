package logger

import (
	"context"
	"log/slog"
	"os"
)

const (
	formatJSON   = "json"
	formatPretty = "pretty"
)

type ctxKey struct{}

type Options struct {
	Level  slog.Level
	Format string
	Color  string
}

func New(opts Options) *slog.Logger {
	if resolveFormat(opts.Format) == formatJSON {
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: opts.Level}))
	}

	return slog.New(NewPrettyHandler(os.Stdout, PrettyOptions{
		Level: opts.Level,
		Color: resolveColor(opts.Color),
	}))
}

func resolveFormat(format string) string {
	switch format {
	case formatJSON:
		return formatJSON
	case formatPretty, "text", "console":
		return formatPretty
	default:
		return formatPretty
	}
}

func resolveColor(mode string) bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}

	switch mode {
	case "never", "off", "false", "0":
		return false
	default:
		return true
	}
}

func ToContext(ctx context.Context, log *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, log)
}

func FromContext(ctx context.Context) *slog.Logger {
	if log, ok := ctx.Value(ctxKey{}).(*slog.Logger); ok {
		return log
	}

	return slog.Default()
}

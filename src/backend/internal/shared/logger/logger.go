package logger

import (
	"log/slog"
	"os"
)

const (
	formatJSON   = "json"
	formatPretty = "pretty"
)

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
	case "always", "on", "true", "1":
		return true
	default:
		return isTerminal(os.Stdout)
	}
}

func isTerminal(f *os.File) bool {
	info, err := f.Stat()
	if err != nil {
		return false
	}

	return info.Mode()&os.ModeCharDevice != 0
}

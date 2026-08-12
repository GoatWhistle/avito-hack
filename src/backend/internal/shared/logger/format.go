package logger

import (
	"log/slog"
	"math"
	"strconv"
	"strings"
	"time"
)

const (
	unitKiB = 1024
	unitMiB = 1024 * 1024
)

func intValue(value slog.Value) int64 {
	kind := value.Kind()

	if kind == slog.KindInt64 {
		return value.Int64()
	}

	if kind == slog.KindUint64 {
		unsigned := value.Uint64()
		if unsigned > math.MaxInt64 {
			return math.MaxInt64
		}

		return int64(unsigned)
	}

	if kind == slog.KindFloat64 {
		return int64(value.Float64())
	}

	if kind == slog.KindDuration {
		return int64(value.Duration())
	}

	parsed, err := strconv.ParseInt(value.String(), 10, 64)
	if err != nil {
		return 0
	}

	return parsed
}

func durationValue(value slog.Value) time.Duration {
	if value.Kind() == slog.KindDuration {
		return value.Duration()
	}

	return time.Duration(intValue(value))
}

func renderValue(value slog.Value) string {
	if value.Kind() == slog.KindDuration {
		return formatDuration(value.Duration())
	}

	text := value.String()
	if needsQuoting(text) {
		return strconv.Quote(text)
	}

	return text
}

func needsQuoting(text string) bool {
	if strings.ContainsAny(text, " \t\"") {
		return true
	}

	for _, r := range text {
		if r < 0x20 || r == 0x7f {
			return true
		}
	}

	return false
}

func formatDuration(elapsed time.Duration) string {
	switch {
	case elapsed >= time.Second:
		return strconv.FormatFloat(elapsed.Seconds(), 'f', 2, 64) + "s"
	case elapsed >= time.Millisecond:
		return strconv.FormatFloat(float64(elapsed.Microseconds())/1000, 'f', 1, 64) + "ms"
	case elapsed >= time.Microsecond:
		return strconv.FormatInt(elapsed.Microseconds(), 10) + "µs"
	default:
		return strconv.FormatInt(elapsed.Nanoseconds(), 10) + "ns"
	}
}

func formatBytes(size int64) string {
	switch {
	case size >= unitMiB:
		return strconv.FormatFloat(float64(size)/unitMiB, 'f', 1, 64) + "MB"
	case size >= unitKiB:
		return strconv.FormatFloat(float64(size)/unitKiB, 'f', 1, 64) + "KB"
	default:
		return strconv.FormatInt(size, 10) + "B"
	}
}

func indentBlock(text string) string {
	lines := strings.Split(strings.TrimRight(text, "\n"), "\n")
	for i, line := range lines {
		lines[i] = "    " + line
	}

	return strings.Join(lines, "\n")
}

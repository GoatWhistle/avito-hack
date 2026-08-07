package logger

import (
	"log/slog"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIntValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value slog.Value
		want  int64
	}{
		{name: "int64", value: slog.Int64Value(42), want: 42},
		{name: "negative int64", value: slog.Int64Value(-7), want: -7},
		{name: "uint64 in range", value: slog.Uint64Value(99), want: 99},
		{name: "uint64 overflow clamps", value: slog.Uint64Value(math.MaxUint64), want: math.MaxInt64},
		{name: "float truncates", value: slog.Float64Value(3.9), want: 3},
		{name: "duration as nanos", value: slog.DurationValue(2 * time.Second), want: int64(2 * time.Second)},
		{name: "numeric string parses", value: slog.StringValue("123"), want: 123},
		{name: "non numeric string is zero", value: slog.StringValue("abc"), want: 0},
		{name: "bool is zero", value: slog.BoolValue(true), want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, intValue(tt.value))
		})
	}
}

func TestDurationValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value slog.Value
		want  time.Duration
	}{
		{name: "duration kind", value: slog.DurationValue(time.Minute), want: time.Minute},
		{name: "int64 nanos", value: slog.Int64Value(int64(time.Second)), want: time.Second},
		{name: "unparsable is zero", value: slog.StringValue("nope"), want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, durationValue(tt.value))
		})
	}
}

func TestRenderValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value slog.Value
		want  string
	}{
		{name: "plain string untouched", value: slog.StringValue("hello"), want: "hello"},
		{name: "string with space quoted", value: slog.StringValue("a b"), want: `"a b"`},
		{name: "string with tab quoted", value: slog.StringValue("a\tb"), want: `"a\tb"`},
		{name: "string with quote quoted", value: slog.StringValue(`a"b`), want: `"a\"b"`},
		{name: "duration formatted", value: slog.DurationValue(1500 * time.Millisecond), want: "1.50s"},
		{name: "int rendered", value: slog.Int64Value(5), want: "5"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, renderValue(tt.value))
		})
	}
}

func TestFormatDuration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		elapsed time.Duration
		want    string
	}{
		{name: "seconds", elapsed: 2500 * time.Millisecond, want: "2.50s"},
		{name: "exactly one second", elapsed: time.Second, want: "1.00s"},
		{name: "milliseconds", elapsed: 250 * time.Millisecond, want: "250.0ms"},
		{name: "exactly one millisecond", elapsed: time.Millisecond, want: "1.0ms"},
		{name: "microseconds", elapsed: 500 * time.Microsecond, want: "500µs"},
		{name: "exactly one microsecond", elapsed: time.Microsecond, want: "1µs"},
		{name: "nanoseconds", elapsed: 300 * time.Nanosecond, want: "300ns"},
		{name: "zero", elapsed: 0, want: "0ns"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, formatDuration(tt.elapsed))
		})
	}
}

func TestFormatBytes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		size int64
		want string
	}{
		{name: "bytes", size: 512, want: "512B"},
		{name: "zero", size: 0, want: "0B"},
		{name: "exactly one kib", size: unitKiB, want: "1.0KB"},
		{name: "kilobytes", size: 2048, want: "2.0KB"},
		{name: "exactly one mib", size: unitMiB, want: "1.0MB"},
		{name: "megabytes", size: 5 * unitMiB, want: "5.0MB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, formatBytes(tt.size))
		})
	}
}

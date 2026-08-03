package logger

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"
)

const (
	accessLogMessage = "http request"
	methodPadding    = 6
	requestIDTail    = 8
)

func (h *prettyHandler) writeAccessLog(buf *strings.Builder, attrs map[string]slog.Value) {
	method := attrs["method"].String()
	status := intValue(attrs["status"])

	buf.WriteString(" ")
	buf.WriteString(h.paint(methodColor(method)+ansiBold, fmt.Sprintf("%-*s", methodPadding, method)))
	buf.WriteString(" ")
	buf.WriteString(h.paint(ansiCyan, attrs["path"].String()))
	buf.WriteString(" ")
	buf.WriteString(h.paint(statusColor(int(status))+ansiBold, strconv.FormatInt(status, 10)))

	if dur, ok := attrs["duration"]; ok {
		elapsed := durationValue(dur)
		buf.WriteString(" ")
		buf.WriteString(h.paint(durationColor(elapsed), formatDuration(elapsed)))
	}

	if size, ok := attrs["bytes"]; ok {
		buf.WriteString(" ")
		buf.WriteString(h.paint(ansiDim, formatBytes(intValue(size))))
	}

	if id, ok := attrs["request_id"]; ok && id.String() != "" {
		buf.WriteString(" ")
		buf.WriteString(h.paint(ansiDim, "#"+shortRequestID(id.String())))
	}

	if ip, ok := attrs["remote_ip"]; ok && ip.String() != "" {
		buf.WriteString(" ")
		buf.WriteString(h.paint(ansiDim, "from "+ip.String()))
	}
}

func shortRequestID(id string) string {
	if idx := strings.LastIndex(id, "/"); idx >= 0 && idx+1 < len(id) {
		id = id[idx+1:]
	}

	if runes := []rune(id); len(runes) > requestIDTail {
		return string(runes[len(runes)-requestIDTail:])
	}

	return id
}

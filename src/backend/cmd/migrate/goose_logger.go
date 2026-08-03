package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

type gooseLogger struct{}

func (gooseLogger) Printf(format string, v ...any) {
	message := strings.TrimSpace(fmt.Sprintf(format, v...))
	if message == "" {
		return
	}

	slog.Info(message)
}

func (gooseLogger) Fatalf(format string, v ...any) {
	slog.Error(strings.TrimSpace(fmt.Sprintf(format, v...)))
	os.Exit(1)
}

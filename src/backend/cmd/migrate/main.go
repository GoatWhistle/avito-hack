package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"github.com/avito-hack/backend/internal/shared/logger"
	"github.com/avito-hack/backend/migrations"
)

const dialect = "postgres"

func main() {
	slog.SetDefault(logger.New(logger.Options{
		Level:  slog.LevelInfo,
		Format: os.Getenv("LOG_FORMAT"),
		Color:  os.Getenv("LOG_COLOR"),
	}))

	if err := run(); err != nil {
		slog.Error("migration failed", slog.Any("error", err))
		os.Exit(1)
	}
}

func run() error {
	command := "up"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return errors.New("DATABASE_URL is required")
	}

	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			slog.Warn("close database", slog.Any("error", closeErr))
		}
	}()

	goose.SetBaseFS(migrations.FS)
	goose.SetLogger(gooseLogger{})

	if err := goose.SetDialect(dialect); err != nil {
		return fmt.Errorf("set dialect: %w", err)
	}

	if err := goose.RunContext(context.Background(), command, db, "."); err != nil {
		return fmt.Errorf("run %q: %w", command, err)
	}

	slog.Info("migration completed", slog.String("command", command))

	return nil
}

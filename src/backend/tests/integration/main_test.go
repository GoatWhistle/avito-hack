//go:build integration

package integration_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"github.com/avito-hack/backend/migrations"
)

const migrationDialect = "postgres"

var pool *pgxpool.Pool

func TestMain(m *testing.M) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		fmt.Fprintln(os.Stderr, "DATABASE_URL is not set, skipping integration tests")
		os.Exit(0)
	}

	if err := prepare(dsn); err != nil {
		fmt.Fprintf(os.Stderr, "integration setup failed: %v\n", err)
		os.Exit(1)
	}

	code := m.Run()

	pool.Close()
	os.Exit(code)
}

func prepare(dsn string) error {
	if err := migrateUp(dsn); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	created, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return fmt.Errorf("open pool: %w", err)
	}

	if err := created.Ping(ctx); err != nil {
		return fmt.Errorf("ping database: %w", err)
	}

	pool = created

	return nil
}

func migrateUp(dsn string) error {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer db.Close()

	goose.SetBaseFS(migrations.FS)
	goose.SetLogger(goose.NopLogger())

	if err := goose.SetDialect(migrationDialect); err != nil {
		return fmt.Errorf("set dialect: %w", err)
	}

	if err := goose.Up(db, "."); err != nil {
		return fmt.Errorf("apply migrations: %w", err)
	}

	return nil
}

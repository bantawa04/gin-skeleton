package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"gin/database/migrations"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		args = []string{"up"}
	}

	goose.SetBaseFS(migrations.FS)

	if err := goose.SetDialect("postgres"); err != nil {
		fatalf("set dialect: %v", err)
	}

	cmd := args[0]
	cmdArgs := args[1:]

	switch cmd {
	case "baseline":
		db := openDB()
		defer db.Close()
		if err := baseline(db); err != nil {
			fatalf("baseline: %v", err)
		}
		fmt.Println("Baseline complete; embedded migrations are marked as applied.")
	case "create":
		if len(cmdArgs) < 1 {
			fatalf("usage: migrate create <name> [sql|go]")
		}
		style := "sql"
		if len(cmdArgs) >= 2 {
			style = cmdArgs[1]
		}
		db := openDB()
		defer db.Close()
		if err := goose.Create(db, "database/migrations", cmdArgs[0], style); err != nil {
			fatalf("create: %v", err)
		}
	default:
		db := openDB()
		defer db.Close()
		if err := goose.RunContext(context.Background(), cmd, db, "."); err != nil {
			fatalf("goose %s: %v", cmd, err)
		}
	}
}

func baseline(db *sql.DB) error {
	ctx := context.Background()

	if _, err := goose.EnsureDBVersionContext(ctx, db); err != nil {
		return fmt.Errorf("ensure version table: %w", err)
	}

	ms, err := goose.CollectMigrations(".", 0, goose.MaxVersion)
	if err != nil {
		return fmt.Errorf("collect migrations: %w", err)
	}

	rows, err := db.QueryContext(ctx, `SELECT version_id FROM goose_db_version WHERE is_applied = true`)
	if err != nil {
		return fmt.Errorf("query existing versions: %w", err)
	}
	defer rows.Close()

	applied := make(map[int64]bool)
	for rows.Next() {
		var v int64
		if err := rows.Scan(&v); err != nil {
			return err
		}
		applied[v] = true
	}
	if err := rows.Err(); err != nil {
		return err
	}

	for _, m := range ms {
		if applied[m.Version] {
			fmt.Printf("skip v%d (already recorded)\n", m.Version)
			continue
		}
		if _, err := db.ExecContext(ctx, `INSERT INTO goose_db_version (version_id, is_applied) VALUES ($1, true)`, m.Version); err != nil {
			return fmt.Errorf("insert v%d: %w", m.Version, err)
		}
		fmt.Printf("mark v%d %s\n", m.Version, m.Source)
	}
	return nil
}

func openDB() *sql.DB {
	db, err := sql.Open("pgx", buildDSN())
	if err != nil {
		fatalf("open db: %v", err)
	}
	return db
}

func buildDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		requireEnv("DB_USER"),
		requireEnv("DB_PASSWORD"),
		envOr("DB_HOST", "localhost"),
		envOr("DB_PORT", "5432"),
		requireEnv("DB_NAME"),
		envOr("DB_SSL_MODE", "disable"),
	)
}

func requireEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		fatalf("required env var %s is not set", key)
	}
	return v
}

func envOr(key string, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "migrate: "+format+"\n", args...)
	os.Exit(1)
}

package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"nanamiku-blog/backend/biz/bootstrap"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

const migrationTableName = "schema_migrations"

func main() {
	_ = godotenv.Load()
	cfg := bootstrap.LoadConfig()
	ctx := context.Background()

	db, err := bootstrap.NewDBPool(ctx, cfg.DB)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer db.Close()

	direction := "up"
	if len(os.Args) > 1 {
		direction = os.Args[1]
	}

	migrationsDir := "sql/migrations"
	if len(os.Args) > 2 {
		migrationsDir = os.Args[2]
	}

	if err := ensureMigrationTable(ctx, db); err != nil {
		log.Fatalf("ensure migration table: %v", err)
	}

	applied, err := loadAppliedMigrations(ctx, db)
	if err != nil {
		log.Fatalf("load applied migrations: %v", err)
	}

	suffix := "." + direction + ".sql"

	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		log.Fatalf("read migrations dir: %v", err)
	}

	var files []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), suffix) {
			files = append(files, e.Name())
		}
	}

	if direction == "down" {
		sort.Sort(sort.Reverse(sort.StringSlice(files)))
	} else {
		sort.Strings(files)
	}

	for _, f := range files {
		path := filepath.Join(migrationsDir, f)
		sql, err := os.ReadFile(path)
		if err != nil {
			log.Fatalf("read %s: %v", f, err)
		}

		if direction == "up" {
			if _, ok := applied[f]; ok {
				fmt.Printf("skip (already applied): %s\n", f)
				continue
			}

			if err := applyUpMigration(ctx, db, f, string(sql)); err != nil {
				log.Fatalf("execute %s: %v", f, err)
			}
			applied[f] = struct{}{}
			continue
		}

		if _, ok := applied[f]; !ok {
			fmt.Printf("skip (not applied): %s\n", f)
			continue
		}
		if err := applyDownMigration(ctx, db, f, string(sql)); err != nil {
			log.Fatalf("execute %s: %v", f, err)
		}
		delete(applied, f)
	}

	fmt.Println("migrations complete")
}

func applyUpMigration(ctx context.Context, db *pgxpool.Pool, version, sql string) error {
	err := runInTx(ctx, db, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, sql); err != nil {
			return err
		}
		_, err := tx.Exec(ctx, "INSERT INTO "+migrationTableName+" (version) VALUES ($1)", version)
		return err
	})
	if err == nil {
		fmt.Printf("applied: %s\n", version)
		return nil
	}

	if !isDuplicateObjectErr(err) {
		return err
	}

	// 兼容老库：当数据库对象已存在但 migration 记录缺失时，自动补记录，避免重复失败。
	if err := markMigrationApplied(ctx, db, version); err != nil {
		return fmt.Errorf("mark legacy migration as applied: %w", err)
	}
	fmt.Printf("mark as applied (legacy schema detected): %s\n", version)
	return nil
}

func applyDownMigration(ctx context.Context, db *pgxpool.Pool, version, sql string) error {
	err := runInTx(ctx, db, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, sql); err != nil {
			return err
		}
		_, err := tx.Exec(ctx, "DELETE FROM "+migrationTableName+" WHERE version = $1", version)
		return err
	})
	if err != nil {
		return err
	}
	fmt.Printf("reverted: %s\n", version)
	return nil
}

func runInTx(ctx context.Context, db *pgxpool.Pool, fn func(tx pgx.Tx) error) (err error) {
	tx, err := db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	if err = fn(tx); err != nil {
		return err
	}
	err = tx.Commit(ctx)
	return err
}

func ensureMigrationTable(ctx context.Context, db *pgxpool.Pool) error {
	sql := `
CREATE TABLE IF NOT EXISTS schema_migrations (
    version text PRIMARY KEY,
    applied_at timestamptz NOT NULL DEFAULT now()
);`
	_, err := db.Exec(ctx, sql)
	return err
}

func loadAppliedMigrations(ctx context.Context, db *pgxpool.Pool) (map[string]struct{}, error) {
	rows, err := db.Query(ctx, "SELECT version FROM "+migrationTableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]struct{})
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		result[version] = struct{}{}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func markMigrationApplied(ctx context.Context, db *pgxpool.Pool, version string) error {
	_, err := db.Exec(ctx, "INSERT INTO "+migrationTableName+" (version) VALUES ($1) ON CONFLICT (version) DO NOTHING", version)
	return err
}

func isDuplicateObjectErr(err error) bool {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return false
	}
	return pgErr.Code == "42710" || pgErr.Code == "42P07"
}

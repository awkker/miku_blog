package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"nanamiku-blog/backend/biz/bootstrap"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

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
		if err := execSQL(ctx, db, string(sql)); err != nil {
			log.Fatalf("execute %s: %v", f, err)
		}
		fmt.Printf("applied: %s\n", f)
	}

	fmt.Println("migrations complete")
}

func execSQL(ctx context.Context, db *pgxpool.Pool, sql string) error {
	_, err := db.Exec(ctx, sql)
	return err
}

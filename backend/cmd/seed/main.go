package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"nanamiku-blog/backend/biz/bootstrap"
	"nanamiku-blog/backend/biz/service"
	"nanamiku-blog/backend/query"

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

	q := query.New(db)

	username := "admin"
	email := "admin@miku.blog"
	password := "admin123"

	if len(os.Args) > 1 {
		password = os.Args[1]
	}

	hash, err := service.HashPassword(password)
	if err != nil {
		log.Fatalf("hash password: %v", err)
	}

	row, err := q.CreateAdminUser(ctx, query.CreateAdminUserParams{
		Username:     username,
		Email:        email,
		PasswordHash: hash,
		Role:         "admin",
	})
	if err != nil {
		log.Fatalf("create admin: %v", err)
	}

	fmt.Printf("Admin user created: id=%s username=%s email=%s\n", row.ID, username, email)
	fmt.Printf("Password: %s (change this immediately)\n", password)
}

package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"nanamiku-blog/backend/biz/bootstrap"
	"nanamiku-blog/backend/biz/jobs"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	cfg := bootstrap.LoadConfig()
	ctx := context.Background()

	db, err := bootstrap.NewDBPool(ctx, cfg.DB)
	if err != nil {
		slog.Error("failed to connect database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	rdb, err := bootstrap.NewRedisClient(ctx, cfg.Redis)
	if err != nil {
		slog.Error("failed to connect redis", "error", err)
		os.Exit(1)
	}
	defer rdb.Close()

	h := server.Default(server.WithHostPorts("0.0.0.0:" + cfg.Server.Port))
	svcs := bootstrap.RegisterRoutes(h, db, rdb, cfg)

	jobs.StartHealthCheckJob(ctx, svcs.Friends, 1*time.Hour)

	slog.Info("server starting", "port", cfg.Server.Port)
	h.Spin()
}

package jobs

import (
	"context"
	"log/slog"
	"time"

	"nanamiku-blog/backend/biz/service"
)

func StartHealthCheckJob(ctx context.Context, friendsSvc *service.FriendsService, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				slog.Info("health check job stopped")
				return
			case <-ticker.C:
				slog.Info("running health checks")
				friendsSvc.RunHealthChecks(ctx)
			}
		}
	}()
	slog.Info("health check job started", "interval", interval.String())
}

package jobs

import (
	"context"
	"log/slog"
	"time"

	"nanamiku-blog/backend/biz/service"
)

func StartPublishSchedulerJob(ctx context.Context, postsSvc *service.PostsService, momentsSvc *service.MomentsService, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				slog.Info("publish scheduler stopped")
				return
			case <-ticker.C:
				if postsSvc != nil {
					count, err := postsSvc.PublishDueScheduled(ctx)
					if err != nil {
						slog.Error("publish scheduled posts failed", "error", err)
					} else if count > 0 {
						slog.Info("published scheduled posts", "count", count)
					}
				}
				if momentsSvc != nil {
					count, err := momentsSvc.PublishDueScheduled(ctx)
					if err != nil {
						slog.Error("publish scheduled moments failed", "error", err)
					} else if count > 0 {
						slog.Info("published scheduled moments", "count", count)
					}
				}
			}
		}
	}()
	slog.Info("publish scheduler started", "interval", interval.String())
}

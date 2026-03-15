package middleware

import (
	"context"
	"log/slog"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
)

func Logger() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		start := time.Now()
		c.Next(ctx)
		latency := time.Since(start)

		status := c.Response.StatusCode()
		level := slog.LevelInfo
		if status >= 500 {
			level = slog.LevelError
		} else if status >= 400 {
			level = slog.LevelWarn
		}

		slog.Log(ctx, level, "request",
			"method", string(c.Method()),
			"path", string(c.URI().Path()),
			"status", status,
			"latency", latency.String(),
			"ip", c.ClientIP(),
			"request_id", c.GetString("request_id"),
		)
	}
}

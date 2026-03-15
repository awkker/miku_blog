package middleware

import (
	"context"
	"log/slog"
	"runtime/debug"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func Recovery() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("panic recovered",
					"error", r,
					"stack", string(debug.Stack()),
					"path", string(c.URI().Path()),
					"method", string(c.Method()),
				)
				c.AbortWithStatusJSON(consts.StatusInternalServerError, map[string]interface{}{
					"code":    5000,
					"message": "internal server error",
				})
			}
		}()
		c.Next(ctx)
	}
}

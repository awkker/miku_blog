package middleware

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/google/uuid"
)

const RequestIDHeader = "X-Request-ID"

func RequestID() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		rid := string(c.GetHeader(RequestIDHeader))
		if rid == "" {
			rid = uuid.New().String()
		}
		c.Set("request_id", rid)
		c.Response.Header.Set(RequestIDHeader, rid)
		c.Next(ctx)
	}
}

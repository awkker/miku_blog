package middleware

import (
	"context"
	"fmt"
	"time"

	"nanamiku-blog/backend/biz/dto"
	"nanamiku-blog/backend/biz/errcode"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/redis/go-redis/v9"
)

var rateLimitScript = redis.NewScript(`
local key = KEYS[1]
local limit = tonumber(ARGV[1])
local window = tonumber(ARGV[2])

local current = redis.call("INCR", key)
if current == 1 then
    redis.call("EXPIRE", key, window)
end

if current > limit then
    return 0
end
return 1
`)

func RateLimit(rdb *redis.Client, keyPrefix string, limit int, window time.Duration) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		ip := c.ClientIP()
		key := fmt.Sprintf("rl:%s:%s", keyPrefix, ip)
		windowSec := int(window.Seconds())

		result, err := rateLimitScript.Run(ctx, rdb, []string{key}, limit, windowSec).Int()
		if err != nil {
			c.Next(ctx)
			return
		}

		if result == 0 {
			c.AbortWithStatusJSON(consts.StatusTooManyRequests,
				dto.Err(errcode.ErrRateLimited, "too many requests, please try again later"))
			return
		}

		c.Next(ctx)
	}
}

package public

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type HealthHandler struct {
	db  *pgxpool.Pool
	rdb *redis.Client
}

func NewHealthHandler(db *pgxpool.Pool, rdb *redis.Client) *HealthHandler {
	return &HealthHandler{db: db, rdb: rdb}
}

func (h *HealthHandler) Check(ctx context.Context, c *app.RequestContext) {
	dbOK := h.db.Ping(ctx) == nil
	redisOK := h.rdb.Ping(ctx).Err() == nil

	status := consts.StatusOK
	if !dbOK || !redisOK {
		status = consts.StatusServiceUnavailable
	}

	c.JSON(status, map[string]interface{}{
		"status":   status == consts.StatusOK,
		"postgres": dbOK,
		"redis":    redisOK,
	})
}

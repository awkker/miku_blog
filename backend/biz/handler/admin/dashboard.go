package admin

import (
	"context"

	"nanamiku-blog/backend/biz/dto"
	"nanamiku-blog/backend/biz/errcode"
	"nanamiku-blog/backend/biz/service"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type DashboardHandler struct {
	svc *service.DashboardService
}

func NewDashboardHandler(svc *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{svc: svc}
}

func (h *DashboardHandler) Stats(ctx context.Context, c *app.RequestContext) {
	stats, err := h.svc.GetStats(ctx)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "failed to get stats"))
		return
	}
	c.JSON(consts.StatusOK, dto.OK(stats))
}

func (h *DashboardHandler) ViewTrend(ctx context.Context, c *app.RequestContext) {
	days := queryInt(c, "days", 30)
	points, err := h.svc.GetViewTrend(ctx, days)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "failed to get view trend"))
		return
	}
	c.JSON(consts.StatusOK, dto.OK(points))
}

func (h *DashboardHandler) CommentTrend(ctx context.Context, c *app.RequestContext) {
	days := queryInt(c, "days", 30)
	points, err := h.svc.GetCommentTrend(ctx, days)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "failed to get comment trend"))
		return
	}
	c.JSON(consts.StatusOK, dto.OK(points))
}

func (h *DashboardHandler) LikeTrend(ctx context.Context, c *app.RequestContext) {
	days := queryInt(c, "days", 30)
	points, err := h.svc.GetLikeTrend(ctx, days)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "failed to get like trend"))
		return
	}
	c.JSON(consts.StatusOK, dto.OK(points))
}

func queryInt(c *app.RequestContext, key string, def int) int {
	v := c.DefaultQuery(key, "")
	if v == "" {
		return def
	}
	n := 0
	for _, ch := range v {
		if ch >= '0' && ch <= '9' {
			n = n*10 + int(ch-'0')
		} else {
			return def
		}
	}
	if n <= 0 {
		return def
	}
	return n
}

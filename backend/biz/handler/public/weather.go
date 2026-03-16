package public

import (
	"context"
	"log/slog"

	"nanamiku-blog/backend/biz/dto"
	"nanamiku-blog/backend/biz/errcode"
	"nanamiku-blog/backend/biz/service"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type WeatherHandler struct {
	svc *service.WeatherService
}

func NewWeatherHandler(svc *service.WeatherService) *WeatherHandler {
	return &WeatherHandler{svc: svc}
}

func (h *WeatherHandler) Current(ctx context.Context, c *app.RequestContext) {
	data, err := h.svc.GetCurrent(ctx)
	if err != nil {
		slog.Error("failed to get weather", "error", err)
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "failed to fetch weather"))
		return
	}
	c.JSON(consts.StatusOK, dto.OK(data))
}

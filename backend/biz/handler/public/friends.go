package public

import (
	"context"

	"nanamiku-blog/backend/biz/dto"
	"nanamiku-blog/backend/biz/errcode"
	"nanamiku-blog/backend/biz/service"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type FriendsHandler struct {
	svc *service.FriendsService
}

func NewFriendsHandler(svc *service.FriendsService) *FriendsHandler {
	return &FriendsHandler{svc: svc}
}

func (h *FriendsHandler) List(ctx context.Context, c *app.RequestContext) {
	items, err := h.svc.ListApproved(ctx)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "failed to list friend links"))
		return
	}
	c.JSON(consts.StatusOK, dto.OK(items))
}

package admin

import (
	"context"

	"nanamiku-blog/backend/biz/dto"
	"nanamiku-blog/backend/biz/errcode"
	"nanamiku-blog/backend/biz/service"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/google/uuid"
)

type MomentsAdminHandler struct {
	svc    *service.MomentsService
	modSvc *service.ModerationService
}

func NewMomentsAdminHandler(svc *service.MomentsService, modSvc *service.ModerationService) *MomentsAdminHandler {
	return &MomentsAdminHandler{svc: svc, modSvc: modSvc}
}

type updateMomentReq struct {
	AuthorName string   `json:"author_name"`
	Content    string   `json:"content"`
	ImageURLs  []string `json:"image_urls"`
}

func (h *MomentsAdminHandler) Update(ctx context.Context, c *app.RequestContext) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "invalid moment id"))
		return
	}

	var req updateMomentReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "invalid request"))
		return
	}
	if req.AuthorName == "" || req.Content == "" {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "author_name and content required"))
		return
	}
	if len(req.ImageURLs) > 4 {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "max 4 images allowed"))
		return
	}

	if err := h.svc.Update(ctx, id, req.AuthorName, req.Content, req.ImageURLs); err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "update moment failed"))
		return
	}
	_ = h.modSvc.LogAudit(ctx, getAdminID(c), "update", "moment", id.String(), nil, getClientIP(c))
	c.JSON(consts.StatusOK, dto.OK(nil))
}

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

type ModerationHandler struct {
	svc *service.ModerationService
}

func NewModerationHandler(svc *service.ModerationService) *ModerationHandler {
	return &ModerationHandler{svc: svc}
}

func (h *ModerationHandler) ListComments(ctx context.Context, c *app.RequestContext) {
	status := c.DefaultQuery("status", "")
	page := queryInt(c, "page", 1)
	size := queryInt(c, "size", 20)

	items, total, err := h.svc.ListComments(ctx, status, page, size)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "failed to list comments"))
		return
	}
	c.JSON(consts.StatusOK, dto.OKPaged(items, total, page, size))
}

func (h *ModerationHandler) ApproveComment(ctx context.Context, c *app.RequestContext) {
	commentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "invalid comment id"))
		return
	}
	adminID := getAdminID(c)
	if err := h.svc.ApproveComment(ctx, commentID, adminID); err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "approve failed"))
		return
	}
	_ = h.svc.LogAudit(ctx, adminID, "approve", "comment", commentID.String(), nil, getClientIP(c))
	c.JSON(consts.StatusOK, dto.OK(nil))
}

func (h *ModerationHandler) RejectComment(ctx context.Context, c *app.RequestContext) {
	commentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "invalid comment id"))
		return
	}
	adminID := getAdminID(c)
	if err := h.svc.RejectComment(ctx, commentID, adminID); err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "reject failed"))
		return
	}
	_ = h.svc.LogAudit(ctx, adminID, "reject", "comment", commentID.String(), nil, getClientIP(c))
	c.JSON(consts.StatusOK, dto.OK(nil))
}

func (h *ModerationHandler) DeleteComment(ctx context.Context, c *app.RequestContext) {
	commentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "invalid comment id"))
		return
	}
	if err := h.svc.DeleteComment(ctx, commentID); err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "delete failed"))
		return
	}
	_ = h.svc.LogAudit(ctx, getAdminID(c), "delete", "comment", commentID.String(), nil, getClientIP(c))
	c.JSON(consts.StatusOK, dto.OK(nil))
}

func (h *ModerationHandler) ListGuestbookMessages(ctx context.Context, c *app.RequestContext) {
	status := c.DefaultQuery("status", "")
	page := queryInt(c, "page", 1)
	size := queryInt(c, "size", 20)

	items, total, err := h.svc.ListGuestbookMessages(ctx, status, page, size)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "failed to list guestbook messages"))
		return
	}
	c.JSON(consts.StatusOK, dto.OKPaged(items, total, page, size))
}

func (h *ModerationHandler) ApproveGuestbookMessage(ctx context.Context, c *app.RequestContext) {
	messageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "invalid message id"))
		return
	}
	adminID := getAdminID(c)
	if err := h.svc.ApproveGuestbookMessage(ctx, messageID, adminID); err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "approve failed"))
		return
	}
	_ = h.svc.LogAudit(ctx, adminID, "approve", "guestbook", messageID.String(), nil, getClientIP(c))
	c.JSON(consts.StatusOK, dto.OK(nil))
}

func (h *ModerationHandler) RejectGuestbookMessage(ctx context.Context, c *app.RequestContext) {
	messageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "invalid message id"))
		return
	}
	adminID := getAdminID(c)
	if err := h.svc.RejectGuestbookMessage(ctx, messageID, adminID); err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "reject failed"))
		return
	}
	_ = h.svc.LogAudit(ctx, adminID, "reject", "guestbook", messageID.String(), nil, getClientIP(c))
	c.JSON(consts.StatusOK, dto.OK(nil))
}

func (h *ModerationHandler) DeleteGuestbookMessage(ctx context.Context, c *app.RequestContext) {
	messageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "invalid message id"))
		return
	}
	if err := h.svc.DeleteGuestbookMessage(ctx, messageID); err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "delete failed"))
		return
	}
	_ = h.svc.LogAudit(ctx, getAdminID(c), "delete", "guestbook", messageID.String(), nil, getClientIP(c))
	c.JSON(consts.StatusOK, dto.OK(nil))
}

func (h *ModerationHandler) ListAuditLogs(ctx context.Context, c *app.RequestContext) {
	page := queryInt(c, "page", 1)
	size := queryInt(c, "size", 20)

	items, err := h.svc.ListAuditLogs(ctx, page, size)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "failed to list audit logs"))
		return
	}
	c.JSON(consts.StatusOK, dto.OK(items))
}

func (h *ModerationHandler) RateLimitMetrics(ctx context.Context, c *app.RequestContext) {
	minutes := queryInt(c, "minutes", 60)
	data, err := h.svc.GetRateLimitMetrics(ctx, minutes)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "failed to get rate-limit metrics"))
		return
	}
	c.JSON(consts.StatusOK, dto.OK(data))
}

func getAdminID(c *app.RequestContext) uuid.UUID {
	val, exists := c.Get("admin_id")
	if !exists {
		return uuid.Nil
	}
	id, ok := val.(uuid.UUID)
	if !ok {
		return uuid.Nil
	}
	return id
}

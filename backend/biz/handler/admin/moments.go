package admin

import (
	"context"
	"time"

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

func (h *MomentsAdminHandler) List(ctx context.Context, c *app.RequestContext) {
	page := queryInt(c, "page", 1)
	size := queryInt(c, "size", 20)

	items, total, err := h.svc.ListAdmin(ctx, page, size)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "failed to list moments"))
		return
	}
	c.JSON(consts.StatusOK, dto.OKPaged(items, total, page, size))
}

type createMomentReq struct {
	AuthorName    string   `json:"author_name"`
	Content       string   `json:"content"`
	ImageURLs     []string `json:"image_urls"`
	PublishStatus string   `json:"publish_status"`
	ScheduledAt   string   `json:"scheduled_at"`
}

func (h *MomentsAdminHandler) Create(ctx context.Context, c *app.RequestContext) {
	var req createMomentReq
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

	publishStatus, scheduledAt, ok := parseMomentPublish(req.PublishStatus, req.ScheduledAt)
	if !ok {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "invalid publish_status or scheduled_at"))
		return
	}

	item, err := h.svc.Create(ctx, service.CreateMomentInput{
		AuthorName:    req.AuthorName,
		Content:       req.Content,
		ImageURLs:     req.ImageURLs,
		IPHash:        "",
		UAHash:        "",
		PublishStatus: publishStatus,
		ScheduledAt:   scheduledAt,
	})
	if err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "create moment failed"))
		return
	}
	_ = h.modSvc.LogAudit(ctx, getAdminID(c), "create", "moment", item.ID.String(), nil, getClientIP(c))
	c.JSON(consts.StatusCreated, dto.OK(item))
}

type updateMomentReq struct {
	AuthorName    string   `json:"author_name"`
	Content       string   `json:"content"`
	ImageURLs     []string `json:"image_urls"`
	PublishStatus string   `json:"publish_status"`
	ScheduledAt   string   `json:"scheduled_at"`
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

	publishStatus, scheduledAt, ok := parseMomentPublish(req.PublishStatus, req.ScheduledAt)
	if !ok {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "invalid publish_status or scheduled_at"))
		return
	}

	if err := h.svc.Update(ctx, id, req.AuthorName, req.Content, req.ImageURLs, publishStatus, scheduledAt); err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "update moment failed"))
		return
	}
	_ = h.modSvc.LogAudit(ctx, getAdminID(c), "update", "moment", id.String(), nil, getClientIP(c))
	c.JSON(consts.StatusOK, dto.OK(nil))
}

type scheduleMomentReq struct {
	ScheduledAt string `json:"scheduled_at"`
}

func (h *MomentsAdminHandler) Publish(ctx context.Context, c *app.RequestContext) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "invalid moment id"))
		return
	}
	if err := h.svc.Publish(ctx, id); err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "publish moment failed"))
		return
	}
	_ = h.modSvc.LogAudit(ctx, getAdminID(c), "publish", "moment", id.String(), nil, getClientIP(c))
	c.JSON(consts.StatusOK, dto.OK(nil))
}

func (h *MomentsAdminHandler) Unpublish(ctx context.Context, c *app.RequestContext) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "invalid moment id"))
		return
	}
	if err := h.svc.Unpublish(ctx, id); err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "unpublish moment failed"))
		return
	}
	_ = h.modSvc.LogAudit(ctx, getAdminID(c), "unpublish", "moment", id.String(), nil, getClientIP(c))
	c.JSON(consts.StatusOK, dto.OK(nil))
}

func (h *MomentsAdminHandler) Schedule(ctx context.Context, c *app.RequestContext) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "invalid moment id"))
		return
	}

	var req scheduleMomentReq
	if err := c.BindJSON(&req); err != nil || req.ScheduledAt == "" {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "scheduled_at required"))
		return
	}

	at, err := time.Parse(time.RFC3339, req.ScheduledAt)
	if err != nil {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "invalid scheduled_at format"))
		return
	}

	if err := h.svc.Schedule(ctx, id, at); err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "schedule moment failed"))
		return
	}
	_ = h.modSvc.LogAudit(ctx, getAdminID(c), "schedule", "moment", id.String(), nil, getClientIP(c))
	c.JSON(consts.StatusOK, dto.OK(nil))
}

func parseMomentPublish(status, scheduledAt string) (string, *time.Time, bool) {
	if status == "" {
		status = "published"
	}
	if status != "draft" && status != "published" && status != "scheduled" {
		return "", nil, false
	}
	if status != "scheduled" {
		return status, nil, true
	}
	if scheduledAt == "" {
		return "", nil, false
	}
	at, err := time.Parse(time.RFC3339, scheduledAt)
	if err != nil {
		return "", nil, false
	}
	return status, &at, true
}

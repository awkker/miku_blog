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

type PostsAdminHandler struct {
	svc *service.PostsService
}

func NewPostsAdminHandler(svc *service.PostsService) *PostsAdminHandler {
	return &PostsAdminHandler{svc: svc}
}

func (h *PostsAdminHandler) List(ctx context.Context, c *app.RequestContext) {
	page := queryInt(c, "page", 1)
	size := queryInt(c, "size", 20)

	items, total, err := h.svc.ListAdmin(ctx, page, size)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "failed to list posts"))
		return
	}
	c.JSON(consts.StatusOK, dto.OKPaged(items, total, page, size))
}

func (h *PostsAdminHandler) Get(ctx context.Context, c *app.RequestContext) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "invalid post id"))
		return
	}

	detail, err := h.svc.GetByID(ctx, id)
	if err != nil {
		c.JSON(consts.StatusNotFound, dto.Err(errcode.ErrNotFound, "post not found"))
		return
	}
	c.JSON(consts.StatusOK, dto.OK(detail))
}

type createPostReq struct {
	Slug            string   `json:"slug"`
	Title           string   `json:"title"`
	Excerpt         string   `json:"excerpt"`
	ContentMarkdown string   `json:"content_markdown"`
	HeroImageURL    string   `json:"hero_image_url"`
	Category        string   `json:"category"`
	Status          string   `json:"status"`
	Tags            []string `json:"tags"`
}

func (h *PostsAdminHandler) Create(ctx context.Context, c *app.RequestContext) {
	var req createPostReq
	if err := c.BindJSON(&req); err != nil || req.Title == "" || req.Slug == "" {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "title and slug required"))
		return
	}

	if req.Status == "" {
		req.Status = "draft"
	}

	adminID := getAdminID(c)
	id, err := h.svc.Create(ctx, service.CreatePostInput{
		Slug:            req.Slug,
		Title:           req.Title,
		Excerpt:         req.Excerpt,
		ContentMarkdown: req.ContentMarkdown,
		HeroImageURL:    req.HeroImageURL,
		Category:        req.Category,
		Status:          req.Status,
		Tags:            req.Tags,
		AdminID:         adminID,
	})
	if err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "create post failed"))
		return
	}
	c.JSON(consts.StatusCreated, dto.OK(map[string]interface{}{"id": id}))
}

type updatePostReq struct {
	Slug            string   `json:"slug"`
	Title           string   `json:"title"`
	Excerpt         string   `json:"excerpt"`
	ContentMarkdown string   `json:"content_markdown"`
	HeroImageURL    string   `json:"hero_image_url"`
	Category        string   `json:"category"`
	Tags            []string `json:"tags"`
}

func (h *PostsAdminHandler) Update(ctx context.Context, c *app.RequestContext) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "invalid post id"))
		return
	}

	var req updatePostReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "invalid request"))
		return
	}

	adminID := getAdminID(c)
	if err := h.svc.Update(ctx, service.UpdatePostInput{
		ID:              id,
		Slug:            req.Slug,
		Title:           req.Title,
		Excerpt:         req.Excerpt,
		ContentMarkdown: req.ContentMarkdown,
		HeroImageURL:    req.HeroImageURL,
		Category:        req.Category,
		Tags:            req.Tags,
		AdminID:         adminID,
	}); err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "update failed"))
		return
	}
	c.JSON(consts.StatusOK, dto.OK(nil))
}

func (h *PostsAdminHandler) Publish(ctx context.Context, c *app.RequestContext) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "invalid post id"))
		return
	}
	adminID := getAdminID(c)
	if err := h.svc.Publish(ctx, id, adminID); err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "publish failed"))
		return
	}
	c.JSON(consts.StatusOK, dto.OK(nil))
}

func (h *PostsAdminHandler) Unpublish(ctx context.Context, c *app.RequestContext) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "invalid post id"))
		return
	}
	adminID := getAdminID(c)
	if err := h.svc.Unpublish(ctx, id, adminID); err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "unpublish failed"))
		return
	}
	c.JSON(consts.StatusOK, dto.OK(nil))
}

type schedulePostReq struct {
	ScheduledAt string `json:"scheduled_at"`
}

func (h *PostsAdminHandler) Schedule(ctx context.Context, c *app.RequestContext) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "invalid post id"))
		return
	}

	var req schedulePostReq
	if err := c.BindJSON(&req); err != nil || req.ScheduledAt == "" {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "scheduled_at required"))
		return
	}

	at, err := time.Parse(time.RFC3339, req.ScheduledAt)
	if err != nil {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "invalid scheduled_at format"))
		return
	}

	adminID := getAdminID(c)
	if err := h.svc.Schedule(ctx, id, adminID, at); err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "schedule failed"))
		return
	}
	c.JSON(consts.StatusOK, dto.OK(nil))
}

func (h *PostsAdminHandler) Delete(ctx context.Context, c *app.RequestContext) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "invalid post id"))
		return
	}
	if err := h.svc.Delete(ctx, id); err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "delete failed"))
		return
	}
	c.JSON(consts.StatusOK, dto.OK(nil))
}

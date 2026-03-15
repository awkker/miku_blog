package public

import (
	"context"

	"nanamiku-blog/backend/biz/dto"
	"nanamiku-blog/backend/biz/errcode"
	"nanamiku-blog/backend/biz/service"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/google/uuid"
)

type PostsHandler struct {
	svc *service.PostsService
}

func NewPostsHandler(svc *service.PostsService) *PostsHandler {
	return &PostsHandler{svc: svc}
}

func (h *PostsHandler) List(ctx context.Context, c *app.RequestContext) {
	category := c.DefaultQuery("category", "")
	page := queryInt(c, "page", 1)
	size := queryInt(c, "size", 20)

	if category != "" {
		items, err := h.svc.ListByCategory(ctx, category, page, size)
		if err != nil {
			c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "failed to list posts"))
			return
		}
		c.JSON(consts.StatusOK, dto.OK(items))
		return
	}

	items, total, err := h.svc.ListPublished(ctx, page, size)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "failed to list posts"))
		return
	}
	c.JSON(consts.StatusOK, dto.OKPaged(items, total, page, size))
}

func (h *PostsHandler) Hot(ctx context.Context, c *app.RequestContext) {
	limit := queryInt(c, "limit", 10)
	items, err := h.svc.ListHot(ctx, limit)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "failed to list hot posts"))
		return
	}
	c.JSON(consts.StatusOK, dto.OK(items))
}

func (h *PostsHandler) Search(ctx context.Context, c *app.RequestContext) {
	q := c.DefaultQuery("q", "")
	if q == "" {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "query parameter 'q' required"))
		return
	}
	page := queryInt(c, "page", 1)
	size := queryInt(c, "size", 20)

	items, err := h.svc.Search(ctx, q, page, size)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "search failed"))
		return
	}
	c.JSON(consts.StatusOK, dto.OK(items))
}

func (h *PostsHandler) GetBySlug(ctx context.Context, c *app.RequestContext) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "slug required"))
		return
	}
	vid := getVisitorID(c)

	detail, err := h.svc.GetBySlug(ctx, slug, vid)
	if err != nil {
		c.JSON(consts.StatusNotFound, dto.Err(errcode.ErrNotFound, "post not found"))
		return
	}
	c.JSON(consts.StatusOK, dto.OK(detail))
}

func (h *PostsHandler) Like(ctx context.Context, c *app.RequestContext) {
	postID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "invalid post id"))
		return
	}
	vid := getVisitorID(c)
	if vid == uuid.Nil {
		c.JSON(consts.StatusForbidden, dto.Err(errcode.ErrForbidden, "visitor identity required"))
		return
	}

	liked, err := h.svc.ToggleLike(ctx, postID, vid)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "like failed"))
		return
	}
	c.JSON(consts.StatusOK, dto.OK(map[string]bool{"liked": liked}))
}

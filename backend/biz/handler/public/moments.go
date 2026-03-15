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

type MomentsHandler struct {
	svc *service.MomentsService
}

func NewMomentsHandler(svc *service.MomentsService) *MomentsHandler {
	return &MomentsHandler{svc: svc}
}

func (h *MomentsHandler) List(ctx context.Context, c *app.RequestContext) {
	page := queryInt(c, "page", 1)
	size := queryInt(c, "size", 20)
	vid := getVisitorID(c)

	items, total, err := h.svc.List(ctx, page, size, vid)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "failed to list moments"))
		return
	}
	c.JSON(consts.StatusOK, dto.OKPaged(items, total, page, size))
}

func (h *MomentsHandler) Latest(ctx context.Context, c *app.RequestContext) {
	limit := queryInt(c, "limit", 3)
	items, err := h.svc.Latest(ctx, limit)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "failed to list latest moments"))
		return
	}
	c.JSON(consts.StatusOK, dto.OK(items))
}

type createMomentReq struct {
	AuthorName string   `json:"author_name"`
	Content    string   `json:"content"`
	ImageURLs  []string `json:"image_urls"`
}

func (h *MomentsHandler) Create(ctx context.Context, c *app.RequestContext) {
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

	ipHash := hashStr(c.ClientIP())
	uaHash := hashStr(string(c.UserAgent()))

	item, err := h.svc.Create(ctx, req.AuthorName, req.Content, req.ImageURLs, ipHash, uaHash)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "failed to create moment"))
		return
	}
	c.JSON(consts.StatusCreated, dto.OK(item))
}

func (h *MomentsHandler) Like(ctx context.Context, c *app.RequestContext) {
	momentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "invalid moment id"))
		return
	}
	vid := getVisitorID(c)
	if vid == uuid.Nil {
		c.JSON(consts.StatusForbidden, dto.Err(errcode.ErrForbidden, "visitor identity required"))
		return
	}

	liked, err := h.svc.ToggleLike(ctx, momentID, vid)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "like failed"))
		return
	}
	c.JSON(consts.StatusOK, dto.OK(map[string]bool{"liked": liked}))
}

func (h *MomentsHandler) Repost(ctx context.Context, c *app.RequestContext) {
	momentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "invalid moment id"))
		return
	}
	vid := getVisitorID(c)
	if vid == uuid.Nil {
		c.JSON(consts.StatusForbidden, dto.Err(errcode.ErrForbidden, "visitor identity required"))
		return
	}

	reposted, err := h.svc.ToggleRepost(ctx, momentID, vid)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "repost failed"))
		return
	}
	c.JSON(consts.StatusOK, dto.OK(map[string]bool{"reposted": reposted}))
}

type createMomentCommentReq struct {
	AuthorName string `json:"author_name"`
	Content    string `json:"content"`
}

func (h *MomentsHandler) CreateComment(ctx context.Context, c *app.RequestContext) {
	momentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "invalid moment id"))
		return
	}

	var req createMomentCommentReq
	if err := c.BindJSON(&req); err != nil || req.AuthorName == "" || req.Content == "" {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "author_name and content required"))
		return
	}

	ipHash := hashStr(c.ClientIP())
	uaHash := hashStr(string(c.UserAgent()))

	item, err := h.svc.CreateComment(ctx, momentID, req.AuthorName, req.Content, ipHash, uaHash)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "failed to create comment"))
		return
	}
	c.JSON(consts.StatusCreated, dto.OK(item))
}

func (h *MomentsHandler) CommentLike(ctx context.Context, c *app.RequestContext) {
	commentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "invalid comment id"))
		return
	}
	vid := getVisitorID(c)
	if vid == uuid.Nil {
		c.JSON(consts.StatusForbidden, dto.Err(errcode.ErrForbidden, "visitor identity required"))
		return
	}

	liked, err := h.svc.ToggleCommentLike(ctx, commentID, vid)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "like failed"))
		return
	}
	c.JSON(consts.StatusOK, dto.OK(map[string]bool{"liked": liked}))
}

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

type PostCommentsHandler struct {
	svc *service.PostCommentsService
}

func NewPostCommentsHandler(svc *service.PostCommentsService) *PostCommentsHandler {
	return &PostCommentsHandler{svc: svc}
}

func (h *PostCommentsHandler) List(ctx context.Context, c *app.RequestContext) {
	postID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "invalid post id"))
		return
	}
	page := queryInt(c, "page", 1)
	size := queryInt(c, "size", 20)

	items, total, err := h.svc.ListComments(ctx, postID, page, size)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "failed to list comments"))
		return
	}
	c.JSON(consts.StatusOK, dto.OKPaged(items, total, page, size))
}

type createPostCommentReq struct {
	ParentID    *uuid.UUID `json:"parent_id"`
	AuthorName  string     `json:"author_name"`
	AuthorEmail string     `json:"author_email"`
	Content     string     `json:"content"`
}

func (h *PostCommentsHandler) Create(ctx context.Context, c *app.RequestContext) {
	postID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "invalid post id"))
		return
	}

	var req createPostCommentReq
	if err := c.BindJSON(&req); err != nil || req.AuthorName == "" || req.Content == "" {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "author_name and content required"))
		return
	}

	ipHash := hashStr(c.ClientIP())
	uaHash := hashStr(string(c.UserAgent()))
	vid := getVisitorID(c)

	item, err := h.svc.CreateComment(ctx, postID, req.ParentID, req.AuthorName, req.AuthorEmail, req.Content, ipHash, uaHash, vid)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "failed to create comment"))
		return
	}
	c.JSON(consts.StatusCreated, dto.OK(item))
}

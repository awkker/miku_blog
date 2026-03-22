package public

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

	"nanamiku-blog/backend/biz/dto"
	"nanamiku-blog/backend/biz/errcode"
	"nanamiku-blog/backend/biz/service"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/google/uuid"
)

type GuestbookHandler struct {
	svc    *service.GuestbookService
	modSvc *service.ModerationService
}

func NewGuestbookHandler(svc *service.GuestbookService, modSvc *service.ModerationService) *GuestbookHandler {
	return &GuestbookHandler{svc: svc, modSvc: modSvc}
}

func (h *GuestbookHandler) List(ctx context.Context, c *app.RequestContext) {
	sortBy := c.DefaultQuery("sort", "newest")
	page := queryInt(c, "page", 1)
	size := queryInt(c, "size", 20)
	visitorID := getVisitorID(c)

	items, total, err := h.svc.ListMessages(ctx, sortBy, page, size, visitorID)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "failed to list messages"))
		return
	}

	c.JSON(consts.StatusOK, dto.OKPaged(items, total, page, size))
}

type createMessageReq struct {
	ParentID      *uuid.UUID `json:"parent_id"`
	AuthorName    string     `json:"author_name"`
	AuthorWebsite string     `json:"author_website"`
	Content       string     `json:"content"`
}

func (h *GuestbookHandler) Create(ctx context.Context, c *app.RequestContext) {
	var req createMessageReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "invalid request"))
		return
	}
	if req.AuthorName == "" || req.Content == "" {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "author_name and content required"))
		return
	}

	if h.modSvc != nil {
		word, err := h.modSvc.FindSensitiveWord(ctx, req.AuthorName, req.Content)
		if err != nil {
			c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "sensitive-word check failed"))
			return
		}
		if word != "" {
			c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBlocked, "message contains blocked keyword"))
			return
		}
	}

	ipHash := hashStr(c.ClientIP())
	uaHash := hashStr(string(c.UserAgent()))

	msg, err := h.svc.CreateMessage(ctx, req.ParentID, req.AuthorName, req.AuthorWebsite, req.Content, ipHash, uaHash)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "failed to create message"))
		return
	}

	c.JSON(consts.StatusCreated, dto.OK(msg))
}

type voteReq struct {
	Vote string `json:"vote"`
}

func (h *GuestbookHandler) Vote(ctx context.Context, c *app.RequestContext) {
	msgIDStr := c.Param("id")
	msgID, err := uuid.Parse(msgIDStr)
	if err != nil {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "invalid message id"))
		return
	}

	visitorID := getVisitorID(c)
	if visitorID == uuid.Nil {
		c.JSON(consts.StatusForbidden, dto.Err(errcode.ErrForbidden, "visitor identity required"))
		return
	}

	var req voteReq
	_ = c.BindJSON(&req)

	if req.Vote != "" && req.Vote != "up" && req.Vote != "down" {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "vote must be 'up', 'down', or empty"))
		return
	}

	if err := h.svc.Vote(ctx, msgID, visitorID, req.Vote); err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "vote failed"))
		return
	}

	c.JSON(consts.StatusOK, dto.OK(nil))
}

func getVisitorID(c *app.RequestContext) uuid.UUID {
	val, exists := c.Get("visitor_id")
	if !exists {
		return uuid.Nil
	}
	id, ok := val.(uuid.UUID)
	if !ok {
		return uuid.Nil
	}
	return id
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

func hashStr(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}

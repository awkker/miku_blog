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

type FriendsAdminHandler struct {
	svc *service.ModerationService
}

func NewFriendsAdminHandler(svc *service.ModerationService) *FriendsAdminHandler {
	return &FriendsAdminHandler{svc: svc}
}

func (h *FriendsAdminHandler) List(ctx context.Context, c *app.RequestContext) {
	page := queryInt(c, "page", 1)
	size := queryInt(c, "size", 20)

	items, total, err := h.svc.ListAdminFriends(ctx, page, size)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "failed to list friends"))
		return
	}
	c.JSON(consts.StatusOK, dto.OKPaged(items, total, page, size))
}

type createFriendReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Domain      string `json:"domain"`
	AvatarURL   string `json:"avatar_url"`
	SortOrder   int32  `json:"sort_order"`
}

func (h *FriendsAdminHandler) Create(ctx context.Context, c *app.RequestContext) {
	var req createFriendReq
	if err := c.BindJSON(&req); err != nil || req.Name == "" || req.URL == "" {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "name and url required"))
		return
	}
	adminID := getAdminID(c)

	id, err := h.svc.CreateFriend(ctx, req.Name, req.Description, req.URL, req.Domain, req.AvatarURL, req.SortOrder, adminID)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "create failed"))
		return
	}
	_ = h.svc.LogAudit(ctx, adminID, "create", "friend_link", id.String(), map[string]string{"name": req.Name}, getClientIP(c))
	c.JSON(consts.StatusCreated, dto.OK(map[string]interface{}{"id": id}))
}

type updateFriendReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Domain      string `json:"domain"`
	AvatarURL   string `json:"avatar_url"`
	SortOrder   int32  `json:"sort_order"`
}

func (h *FriendsAdminHandler) Update(ctx context.Context, c *app.RequestContext) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "invalid id"))
		return
	}
	var req updateFriendReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "invalid request"))
		return
	}
	if err := h.svc.UpdateFriend(ctx, id, req.Name, req.Description, req.URL, req.Domain, req.AvatarURL, req.SortOrder); err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "update failed"))
		return
	}
	_ = h.svc.LogAudit(ctx, getAdminID(c), "update", "friend_link", id.String(), map[string]string{"name": req.Name}, getClientIP(c))
	c.JSON(consts.StatusOK, dto.OK(nil))
}

func (h *FriendsAdminHandler) Delete(ctx context.Context, c *app.RequestContext) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "invalid id"))
		return
	}
	if err := h.svc.DeleteFriend(ctx, id); err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "delete failed"))
		return
	}
	_ = h.svc.LogAudit(ctx, getAdminID(c), "delete", "friend_link", id.String(), nil, getClientIP(c))
	c.JSON(consts.StatusOK, dto.OK(nil))
}

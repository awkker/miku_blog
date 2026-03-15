package admin

import (
	"context"
	"errors"

	"nanamiku-blog/backend/biz/dto"
	"nanamiku-blog/backend/biz/errcode"
	"nanamiku-blog/backend/biz/service"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authSvc *service.AuthService
}

func NewAuthHandler(authSvc *service.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

type loginRequest struct {
	Username string `json:"username" vd:"len($)>0"`
	Password string `json:"password" vd:"len($)>0"`
}

func (h *AuthHandler) Login(ctx context.Context, c *app.RequestContext) {
	var req loginRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "invalid request body"))
		return
	}

	pair, err := h.authSvc.Login(ctx, req.Username, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			c.JSON(consts.StatusUnauthorized, dto.Err(errcode.ErrInvalidCredentials, "invalid username or password"))
			return
		}
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "login failed"))
		return
	}

	c.JSON(consts.StatusOK, dto.OK(pair))
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token" vd:"len($)>0"`
}

func (h *AuthHandler) Refresh(ctx context.Context, c *app.RequestContext) {
	var req refreshRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "refresh_token required"))
		return
	}

	pair, err := h.authSvc.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		if errors.Is(err, service.ErrTokenExpired) {
			c.JSON(consts.StatusUnauthorized, dto.Err(errcode.ErrTokenExpired, "refresh token expired"))
			return
		}
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "refresh failed"))
		return
	}

	c.JSON(consts.StatusOK, dto.OK(pair))
}

type logoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *AuthHandler) Logout(ctx context.Context, c *app.RequestContext) {
	var req logoutRequest
	_ = c.BindJSON(&req)

	_ = h.authSvc.Logout(ctx, req.RefreshToken)
	c.JSON(consts.StatusOK, dto.OK(nil))
}

func (h *AuthHandler) Me(ctx context.Context, c *app.RequestContext) {
	adminIDVal, exists := c.Get("admin_id")
	if !exists {
		c.JSON(consts.StatusUnauthorized, dto.Err(errcode.ErrUnauthorized, "unauthorized"))
		return
	}

	adminID, ok := adminIDVal.(uuid.UUID)
	if !ok {
		c.JSON(consts.StatusUnauthorized, dto.Err(errcode.ErrUnauthorized, "unauthorized"))
		return
	}

	admin, err := h.authSvc.GetAdminInfo(ctx, adminID)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "failed to get admin info"))
		return
	}

	c.JSON(consts.StatusOK, dto.OK(map[string]interface{}{
		"id":            admin.ID,
		"username":      admin.Username,
		"email":         admin.Email,
		"role":          admin.Role,
		"last_login_at": admin.LastLoginAt,
		"created_at":    admin.CreatedAt,
	}))
}

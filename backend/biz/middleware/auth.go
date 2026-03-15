package middleware

import (
	"context"
	"strings"

	"nanamiku-blog/backend/biz/dto"
	"nanamiku-blog/backend/biz/errcode"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/google/uuid"
)

type AdminClaims struct {
	AdminID  uuid.UUID
	Username string
	Role     string
}

type TokenValidator func(tokenStr string) (*AdminClaims, error)

func AdminAuth(validate TokenValidator) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		header := string(c.GetHeader("Authorization"))
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(consts.StatusUnauthorized, dto.Err(errcode.ErrUnauthorized, "missing or invalid authorization header"))
			return
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := validate(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(consts.StatusUnauthorized, dto.Err(errcode.ErrTokenInvalid, "invalid or expired token"))
			return
		}

		c.Set("admin_id", claims.AdminID)
		c.Set("admin_username", claims.Username)
		c.Set("admin_role", claims.Role)
		c.Next(ctx)
	}
}

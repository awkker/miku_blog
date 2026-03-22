package admin

import (
	"context"
	"fmt"
	"time"

	"nanamiku-blog/backend/biz/dto"
	"nanamiku-blog/backend/biz/errcode"
	"nanamiku-blog/backend/biz/service"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type BackupHandler struct {
	svc *service.BackupService
}

func NewBackupHandler(svc *service.BackupService) *BackupHandler {
	return &BackupHandler{svc: svc}
}

func (h *BackupHandler) Export(ctx context.Context, c *app.RequestContext) {
	format := c.DefaultQuery("format", "json")
	if format != "json" && format != "sql" {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "format must be json or sql"))
		return
	}

	stamp := time.Now().UTC().Format("20060102-150405")
	if format == "sql" {
		content, err := h.svc.BuildSQL(ctx)
		if err != nil {
			c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "failed to export sql backup"))
			return
		}
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"miku-backup-%s.sql\"", stamp))
		c.Data(consts.StatusOK, "application/sql; charset=utf-8", content)
		return
	}

	content, err := h.svc.BuildJSON(ctx)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "failed to export json backup"))
		return
	}
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"miku-backup-%s.json\"", stamp))
	c.Data(consts.StatusOK, "application/json; charset=utf-8", content)
}

package public

import (
	"context"
	"strings"
	"time"

	"nanamiku-blog/backend/biz/dto"
	"nanamiku-blog/backend/biz/errcode"
	"nanamiku-blog/backend/biz/service"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type AnalyticsHandler struct {
	svc *service.DashboardService
}

func NewAnalyticsHandler(svc *service.DashboardService) *AnalyticsHandler {
	return &AnalyticsHandler{svc: svc}
}

type collectAnalyticsReq struct {
	SessionKey string `json:"session_key"`
	Path       string `json:"path"`
	Title      string `json:"title"`
	Referrer   string `json:"referrer"`
	Timezone   string `json:"timezone"`
	Language   string `json:"language"`
	OccurredAt string `json:"occurred_at"`
}

func (h *AnalyticsHandler) Collect(ctx context.Context, c *app.RequestContext) {
	ua := strings.TrimSpace(string(c.UserAgent()))
	if isBotUserAgent(ua) {
		c.JSON(consts.StatusOK, dto.OK(nil))
		return
	}

	var req collectAnalyticsReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(consts.StatusBadRequest, dto.Err(errcode.ErrBadRequest, "invalid request"))
		return
	}

	occurredAt := time.Time{}
	if raw := strings.TrimSpace(req.OccurredAt); raw != "" {
		if t, err := time.Parse(time.RFC3339, raw); err == nil {
			occurredAt = t
		}
	}

	referrer := strings.TrimSpace(req.Referrer)
	if referrer == "" {
		referrer = string(c.GetHeader("Referer"))
	}

	err := h.svc.CollectPageview(ctx, service.AnalyticsCollectInput{
		VisitorID:   getVisitorID(c),
		SessionKey:  req.SessionKey,
		Path:        req.Path,
		Title:       req.Title,
		Referrer:    referrer,
		Timezone:    req.Timezone,
		Language:    req.Language,
		OccurredAt:  occurredAt,
		UserAgent:   ua,
		ClientIP:    getClientIP(c),
		CountryCode: pickHeader(c, "CF-IPCountry", "X-Vercel-IP-Country", "X-AppEngine-Country", "X-Country-Code"),
		Region:      pickHeader(c, "X-Vercel-IP-Country-Region", "X-AppEngine-Region", "CF-IPRegion"),
		City:        pickHeader(c, "X-Vercel-IP-City", "X-AppEngine-City", "CF-IPCity"),
	})
	if err != nil {
		c.JSON(consts.StatusInternalServerError, dto.Err(errcode.ErrInternal, "analytics collection failed"))
		return
	}

	c.JSON(consts.StatusOK, dto.OK(nil))
}

func pickHeader(c *app.RequestContext, keys ...string) string {
	for _, key := range keys {
		v := strings.TrimSpace(string(c.GetHeader(key)))
		if v != "" {
			return v
		}
	}
	return ""
}

func getClientIP(c *app.RequestContext) string {
	if ip := strings.TrimSpace(string(c.GetHeader("X-Real-IP"))); ip != "" {
		return ip
	}
	if ip := strings.TrimSpace(string(c.GetHeader("X-Forwarded-For"))); ip != "" {
		if idx := strings.Index(ip, ","); idx >= 0 {
			return strings.TrimSpace(ip[:idx])
		}
		return ip
	}
	return strings.TrimSpace(c.ClientIP())
}

func isBotUserAgent(ua string) bool {
	if ua == "" {
		return false
	}
	lower := strings.ToLower(ua)
	botKeywords := []string{"bot", "crawler", "spider", "headless", "pingdom", "curl", "wget"}
	for _, token := range botKeywords {
		if strings.Contains(lower, token) {
			return true
		}
	}
	return false
}

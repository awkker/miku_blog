package middleware

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log/slog"

	"nanamiku-blog/backend/query"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const visitorCookieName = "miku_vid"

func Visitor(db *pgxpool.Pool) app.HandlerFunc {
	q := query.New(db)

	return func(ctx context.Context, c *app.RequestContext) {
		var visitorID uuid.UUID

		cookieVal := string(c.Cookie(visitorCookieName))
		if cookieVal != "" {
			if id, err := uuid.Parse(cookieVal); err == nil {
				visitorID = id
				_ = q.TouchVisitor(ctx, id)
			}
		}

		if visitorID == uuid.Nil {
			ipHash := hashString(c.ClientIP())
			uaHash := hashString(string(c.UserAgent()))

			v, err := q.GetVisitorByIPUA(ctx, query.GetVisitorByIPUAParams{
				IpHash: ipHash,
				UaHash: uaHash,
			})
			if err != nil {
				if err == pgx.ErrNoRows {
					created, createErr := q.CreateVisitor(ctx, query.CreateVisitorParams{
						IpHash: ipHash,
						UaHash: uaHash,
					})
					if createErr != nil {
						slog.Error("create visitor failed", "error", createErr)
						c.Next(ctx)
						return
					}
					visitorID = created.ID
				} else {
					slog.Error("query visitor failed", "error", err)
					c.Next(ctx)
					return
				}
			} else {
				visitorID = v.ID
				_ = q.TouchVisitor(ctx, v.ID)
			}

			c.SetCookie(visitorCookieName, visitorID.String(), 365*24*3600, "/", "", protocol.CookieSameSiteLaxMode, false, true)
		}

		c.Set("visitor_id", visitorID)
		c.Next(ctx)
	}
}

func hashString(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}

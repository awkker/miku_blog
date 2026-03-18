package bootstrap

import (
	"time"

	"nanamiku-blog/backend/biz/handler/admin"
	"nanamiku-blog/backend/biz/handler/public"
	"nanamiku-blog/backend/biz/middleware"
	"nanamiku-blog/backend/biz/service"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Services struct {
	Auth       *service.AuthService
	Guestbook  *service.GuestbookService
	Moments    *service.MomentsService
	Friends    *service.FriendsService
	Dashboard  *service.DashboardService
	Moderation *service.ModerationService
}

func RegisterRoutes(h *server.Hertz, db *pgxpool.Pool, rdb *redis.Client, cfg *Config) *Services {
	h.Use(
		middleware.Recovery(),
		middleware.RequestID(),
		middleware.Logger(),
		middleware.CORS(cfg.CORS.Origins),
	)

	authSvc := service.NewAuthService(db, service.JWTConfig{
		Secret:     cfg.JWT.Secret,
		AccessTTL:  cfg.JWT.AccessTTL,
		RefreshTTL: cfg.JWT.RefreshTTL,
	})
	guestbookSvc := service.NewGuestbookService(db)
	momentsSvc := service.NewMomentsService(db)
	friendsSvc := service.NewFriendsService(db)
	dashboardSvc := service.NewDashboardService(db)
	moderationSvc := service.NewModerationService(db)
	postsSvc := service.NewPostsService(db)
	postCommentsSvc := service.NewPostCommentsService(db)
	weatherSvc := service.NewWeatherService(rdb, cfg.Weather.Location)

	tokenValidator := func(tokenStr string) (*middleware.AdminClaims, error) {
		claims, err := authSvc.ValidateAccessToken(tokenStr)
		if err != nil {
			return nil, err
		}
		return &middleware.AdminClaims{
			AdminID:  claims.AdminID,
			Username: claims.Username,
			Role:     claims.Role,
		}, nil
	}

	healthH := public.NewHealthHandler(db, rdb)
	authH := admin.NewAuthHandler(authSvc)
	guestbookH := public.NewGuestbookHandler(guestbookSvc)
	momentsH := public.NewMomentsHandler(momentsSvc)
	friendsH := public.NewFriendsHandler(friendsSvc)
	analyticsH := public.NewAnalyticsHandler(dashboardSvc)
	dashboardH := admin.NewDashboardHandler(dashboardSvc)
	moderationH := admin.NewModerationHandler(moderationSvc)
	friendsAdminH := admin.NewFriendsAdminHandler(moderationSvc)
	postsH := public.NewPostsHandler(postsSvc)
	postCommentsH := public.NewPostCommentsHandler(postCommentsSvc)
	postsAdminH := admin.NewPostsAdminHandler(postsSvc, moderationSvc)
	momentsAdminH := admin.NewMomentsAdminHandler(momentsSvc, moderationSvc)
	weatherH := public.NewWeatherHandler(weatherSvc)

	api := h.Group("/api/v1")
	api.Use(middleware.Visitor(db))
	{
		api.GET("/health", healthH.Check)
		api.GET("/weather", weatherH.Current)

		// Auth
		auth := api.Group("/auth")
		auth.Use(middleware.RateLimit(rdb, "login", 10, 1*time.Minute))
		{
			auth.POST("/login", authH.Login)
			auth.POST("/refresh", authH.Refresh)
			auth.POST("/logout", authH.Logout)
		}
		authed := api.Group("/auth")
		authed.Use(middleware.AdminAuth(tokenValidator))
		{
			authed.GET("/me", authH.Me)
		}

		// Guestbook
		gb := api.Group("/guestbook")
		{
			gb.GET("/messages", guestbookH.List)
			gb.POST("/messages", middleware.RateLimit(rdb, "gb:create", 5, 1*time.Minute), guestbookH.Create)
			gb.POST("/messages/:id/vote", middleware.RateLimit(rdb, "gb:vote", 30, 1*time.Minute), guestbookH.Vote)
		}

		// Moments
		mt := api.Group("/moments")
		{
			mt.GET("", momentsH.List)
			mt.GET("/latest", momentsH.Latest)
			mt.POST("", middleware.RateLimit(rdb, "mt:create", 3, 1*time.Minute), momentsH.Create)
			mt.POST("/:id/like", middleware.RateLimit(rdb, "mt:like", 30, 1*time.Minute), momentsH.Like)
			mt.POST("/:id/repost", middleware.RateLimit(rdb, "mt:repost", 10, 1*time.Minute), momentsH.Repost)
			mt.GET("/:id/comments", momentsH.ListComments)
			mt.POST("/:id/comments", middleware.RateLimit(rdb, "mt:comment", 5, 1*time.Minute), momentsH.CreateComment)
			mt.POST("/comments/:id/like", middleware.RateLimit(rdb, "mt:clike", 30, 1*time.Minute), momentsH.CommentLike)
		}

		// Friends
		api.GET("/friends", friendsH.List)
		api.POST("/analytics/collect", middleware.RateLimit(rdb, "analytics:collect", 240, 1*time.Minute), analyticsH.Collect)

		// Posts (public)
		posts := api.Group("/posts")
		{
			posts.GET("", postsH.List)
			posts.GET("/hot", postsH.Hot)
			posts.GET("/search", postsH.Search)
			posts.GET("/:slug", postsH.GetBySlug)
			posts.POST("/:id/like", middleware.RateLimit(rdb, "post:like", 30, 1*time.Minute), postsH.Like)
			posts.GET("/:id/comments", postCommentsH.List)
			posts.POST("/:id/comments", middleware.RateLimit(rdb, "post:comment", 5, 1*time.Minute), postCommentsH.Create)
		}

		// Admin routes (behind AdminAuth)
		adm := api.Group("/admin")
		adm.Use(middleware.AdminAuth(tokenValidator))
		{
			// Dashboard
			adm.GET("/dashboard/stats", dashboardH.Stats)
			adm.GET("/dashboard/trend/views", dashboardH.ViewTrend)
			adm.GET("/dashboard/trend/comments", dashboardH.CommentTrend)
			adm.GET("/dashboard/trend/likes", dashboardH.LikeTrend)
			adm.GET("/dashboard/analytics", dashboardH.Analytics)

			// Moderation - comments
			adm.GET("/comments", moderationH.ListComments)
			adm.POST("/comments/:id/approve", moderationH.ApproveComment)
			adm.POST("/comments/:id/reject", moderationH.RejectComment)
			adm.DELETE("/comments/:id", moderationH.DeleteComment)
			adm.GET("/guestbook/messages", moderationH.ListGuestbookMessages)
			adm.POST("/guestbook/messages/:id/approve", moderationH.ApproveGuestbookMessage)
			adm.POST("/guestbook/messages/:id/reject", moderationH.RejectGuestbookMessage)
			adm.DELETE("/guestbook/messages/:id", moderationH.DeleteGuestbookMessage)

			// Moderation - audit
			adm.GET("/audit-logs", moderationH.ListAuditLogs)

			// Friends CRUD
			adm.GET("/friends", friendsAdminH.List)
			adm.POST("/friends", friendsAdminH.Create)
			adm.PUT("/friends/:id", friendsAdminH.Update)
			adm.DELETE("/friends/:id", friendsAdminH.Delete)

			// Posts CRUD
			adm.GET("/posts", postsAdminH.List)
			adm.GET("/posts/:id", postsAdminH.Get)
			adm.POST("/posts", postsAdminH.Create)
			adm.PUT("/posts/:id", postsAdminH.Update)
			adm.POST("/posts/:id/publish", postsAdminH.Publish)
			adm.POST("/posts/:id/unpublish", postsAdminH.Unpublish)
			adm.POST("/posts/:id/schedule", postsAdminH.Schedule)
			adm.DELETE("/posts/:id", postsAdminH.Delete)

			// Moments CRUD
			adm.PUT("/moments/:id", momentsAdminH.Update)
		}
	}

	return &Services{
		Auth:       authSvc,
		Guestbook:  guestbookSvc,
		Moments:    momentsSvc,
		Friends:    friendsSvc,
		Dashboard:  dashboardSvc,
		Moderation: moderationSvc,
	}
}

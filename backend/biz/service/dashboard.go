package service

import (
	"context"
	"fmt"
	"time"

	"nanamiku-blog/backend/query"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DashboardService struct {
	q           *query.Queries
	db          *pgxpool.Pool
	geoResolver *GeoIPResolver
}

func NewDashboardService(db *pgxpool.Pool, geoResolver *GeoIPResolver) *DashboardService {
	return &DashboardService{q: query.New(db), db: db, geoResolver: geoResolver}
}

type DashboardStats struct {
	TotalPosts   int64 `json:"total_posts"`
	TotalLikes   int64 `json:"total_likes"`
	PendingCount int64 `json:"pending_comments"`
	FriendCount  int64 `json:"friend_count"`
	DraftCount   int64 `json:"draft_count"`
}

type TrendPoint struct {
	Day   string `json:"day"`
	Value int64  `json:"value"`
}

type ViewTrendPoint struct {
	Day string `json:"day"`
	PV  int64  `json:"pv"`
	UV  int64  `json:"uv"`
}

func (s *DashboardService) GetStats(ctx context.Context) (*DashboardStats, error) {
	totalPosts, err := s.q.GetTotalPostCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("total posts: %w", err)
	}
	totalLikes, err := s.q.GetTotalLikeCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("total likes: %w", err)
	}
	pendingComments, err := s.q.CountPendingComments(ctx)
	if err != nil {
		return nil, fmt.Errorf("pending comments: %w", err)
	}
	pendingGuestbook, err := s.q.CountAdminGuestbookMessages(ctx, query.NullModerationStatus{
		ModerationStatus: query.ModerationStatusPending,
		Valid:            true,
	})
	if err != nil {
		return nil, fmt.Errorf("pending guestbook messages: %w", err)
	}
	friendCount, err := s.q.CountApprovedFriendLinks(ctx)
	if err != nil {
		return nil, fmt.Errorf("friend count: %w", err)
	}
	draftCount, err := s.q.CountDraftPosts(ctx)
	if err != nil {
		return nil, fmt.Errorf("draft count: %w", err)
	}

	return &DashboardStats{
		TotalPosts:   totalPosts,
		TotalLikes:   totalLikes,
		PendingCount: pendingComments + pendingGuestbook,
		FriendCount:  friendCount,
		DraftCount:   draftCount,
	}, nil
}

func (s *DashboardService) GetViewTrend(ctx context.Context, days int) ([]ViewTrendPoint, error) {
	since := time.Now().AddDate(0, 0, -days)
	rows, err := s.q.GetDailyViewTrend(ctx, pgtype.Date{Time: since, Valid: true})
	if err != nil {
		return nil, err
	}
	points := make([]ViewTrendPoint, 0, len(rows))
	for _, r := range rows {
		points = append(points, ViewTrendPoint{
			Day: r.Day.Time.Format("2006-01-02"),
			PV:  r.Pv,
			UV:  r.Uv,
		})
	}
	return points, nil
}

func (s *DashboardService) GetCommentTrend(ctx context.Context, days int) ([]TrendPoint, error) {
	since := time.Now().AddDate(0, 0, -days)
	rows, err := s.q.GetDailyCommentTrend(ctx, since)
	if err != nil {
		return nil, err
	}
	points := make([]TrendPoint, 0, len(rows))
	for _, r := range rows {
		points = append(points, TrendPoint{
			Day:   r.Day.Time.Format("2006-01-02"),
			Value: r.Total,
		})
	}
	return points, nil
}

func (s *DashboardService) GetLikeTrend(ctx context.Context, days int) ([]TrendPoint, error) {
	since := time.Now().AddDate(0, 0, -days)
	rows, err := s.q.GetDailyLikeTrend(ctx, since)
	if err != nil {
		return nil, err
	}
	points := make([]TrendPoint, 0, len(rows))
	for _, r := range rows {
		points = append(points, TrendPoint{
			Day:   r.Day.Time.Format("2006-01-02"),
			Value: r.Total,
		})
	}
	return points, nil
}

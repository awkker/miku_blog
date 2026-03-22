package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"nanamiku-blog/backend/query"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MomentsService struct {
	q  *query.Queries
	db *pgxpool.Pool
}

func NewMomentsService(db *pgxpool.Pool) *MomentsService {
	return &MomentsService{q: query.New(db), db: db}
}

type MomentItem struct {
	ID            uuid.UUID `json:"id"`
	AuthorName    string    `json:"author_name"`
	Content       string    `json:"content"`
	ImageURLs     []string  `json:"image_urls"`
	LikeCount     int64     `json:"like_count"`
	RepostCount   int64     `json:"repost_count"`
	CommentCount  int64     `json:"comment_count"`
	PublishStatus string    `json:"publish_status,omitempty"`
	PublishedAt   string    `json:"published_at,omitempty"`
	ScheduledAt   string    `json:"scheduled_at,omitempty"`
	CreatedAt     string    `json:"created_at"`
	Liked         bool      `json:"liked"`
	Reposted      bool      `json:"reposted"`
}

type MomentCommentItem struct {
	ID         uuid.UUID `json:"id"`
	AuthorName string    `json:"author_name"`
	Content    string    `json:"content"`
	LikeCount  int64     `json:"like_count"`
	CreatedAt  string    `json:"created_at"`
	Liked      bool      `json:"liked"`
}

func (s *MomentsService) List(ctx context.Context, page, size int, visitorID uuid.UUID) ([]MomentItem, int64, error) {
	total, err := s.q.CountMoments(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("count moments: %w", err)
	}

	rows, err := s.q.ListMoments(ctx, query.ListMomentsParams{
		Limit:  int32(size),
		Offset: int32((page - 1) * size),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("list moments: %w", err)
	}

	ids := make([]uuid.UUID, len(rows))
	for i, r := range rows {
		ids[i] = r.ID
	}

	likeSet, repostSet := s.getInteractionSets(ctx, visitorID, ids)

	items := make([]MomentItem, 0, len(rows))
	for _, r := range rows {
		item := MomentItem{
			ID:            r.ID,
			AuthorName:    r.AuthorName,
			Content:       r.Content,
			ImageURLs:     parseImageURLs(r.ImageUrls),
			LikeCount:     r.LikeCount,
			RepostCount:   r.RepostCount,
			CommentCount:  r.CommentCount,
			PublishStatus: string(r.PublishStatus),
			CreatedAt:     r.CreatedAt.Format("2006-01-02T15:04:05Z"),
			Liked:         likeSet[r.ID],
			Reposted:      repostSet[r.ID],
		}
		if r.PublishedAt.Valid {
			item.PublishedAt = r.PublishedAt.Time.Format(time.RFC3339)
		}
		if r.ScheduledAt.Valid {
			item.ScheduledAt = r.ScheduledAt.Time.Format(time.RFC3339)
		}
		items = append(items, item)
	}

	return items, total, nil
}

func (s *MomentsService) Latest(ctx context.Context, limit int) ([]MomentItem, error) {
	rows, err := s.q.ListLatestMoments(ctx, int32(limit))
	if err != nil {
		return nil, fmt.Errorf("list latest: %w", err)
	}
	items := make([]MomentItem, 0, len(rows))
	for _, r := range rows {
		items = append(items, MomentItem{
			ID:         r.ID,
			AuthorName: r.AuthorName,
			Content:    r.Content,
			ImageURLs:  parseImageURLs(r.ImageUrls),
			CreatedAt:  r.CreatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}
	return items, nil
}

type CreateMomentInput struct {
	AuthorName    string
	Content       string
	ImageURLs     []string
	IPHash        string
	UAHash        string
	PublishStatus string
	ScheduledAt   *time.Time
}

func (s *MomentsService) Create(ctx context.Context, in CreateMomentInput) (*MomentItem, error) {
	imgs, _ := json.Marshal(in.ImageURLs)
	if in.PublishStatus == "" {
		in.PublishStatus = string(query.MomentPublishStatusPublished)
	}

	var scheduledAt interface{}
	if in.ScheduledAt != nil {
		scheduledAt = pgtype.Timestamptz{Time: *in.ScheduledAt, Valid: true}
	}

	row, err := s.q.CreateMoment(ctx, query.CreateMomentParams{
		AuthorName:    in.AuthorName,
		Content:       in.Content,
		ImageUrls:     imgs,
		IpHash:        in.IPHash,
		UaHash:        in.UAHash,
		PublishStatus: query.MomentPublishStatus(in.PublishStatus),
		Column7:       scheduledAt,
	})
	if err != nil {
		return nil, fmt.Errorf("create moment: %w", err)
	}

	item := &MomentItem{
		ID:            row.ID,
		AuthorName:    in.AuthorName,
		Content:       in.Content,
		ImageURLs:     in.ImageURLs,
		PublishStatus: in.PublishStatus,
		CreatedAt:     row.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
	if in.PublishStatus == string(query.MomentPublishStatusPublished) {
		item.PublishedAt = row.CreatedAt.Format(time.RFC3339)
	}
	if in.ScheduledAt != nil {
		item.ScheduledAt = in.ScheduledAt.Format(time.RFC3339)
	}
	return item, nil
}

func (s *MomentsService) Update(ctx context.Context, momentID uuid.UUID, authorName, content string, imageURLs []string, publishStatus string, scheduledAt *time.Time) error {
	imgs, _ := json.Marshal(imageURLs)
	if publishStatus == "" {
		publishStatus = string(query.MomentPublishStatusPublished)
	}

	scheduledAtValue := pgtype.Timestamptz{}
	if scheduledAt != nil {
		scheduledAtValue = pgtype.Timestamptz{Time: *scheduledAt, Valid: true}
	}

	return s.q.UpdateMoment(ctx, query.UpdateMomentParams{
		ID:            momentID,
		AuthorName:    authorName,
		Content:       content,
		ImageUrls:     imgs,
		PublishStatus: query.MomentPublishStatus(publishStatus),
		ScheduledAt:   scheduledAtValue,
	})
}

func (s *MomentsService) ToggleLike(ctx context.Context, momentID, visitorID uuid.UUID) (bool, error) {
	cnt, err := s.q.CheckMomentLike(ctx, query.CheckMomentLikeParams{MomentID: momentID, VisitorID: visitorID})
	if err != nil {
		return false, err
	}
	if cnt > 0 {
		_ = s.q.DeleteMomentLike(ctx, query.DeleteMomentLikeParams{MomentID: momentID, VisitorID: visitorID})
		_ = s.q.DecrementMomentLikeCount(ctx, momentID)
		return false, nil
	}
	_ = s.q.CreateMomentLike(ctx, query.CreateMomentLikeParams{MomentID: momentID, VisitorID: visitorID})
	_ = s.q.IncrementMomentLikeCount(ctx, momentID)
	return true, nil
}

func (s *MomentsService) ToggleRepost(ctx context.Context, momentID, visitorID uuid.UUID) (bool, error) {
	cnt, err := s.q.CheckMomentRepost(ctx, query.CheckMomentRepostParams{MomentID: momentID, VisitorID: visitorID})
	if err != nil {
		return false, err
	}
	if cnt > 0 {
		_ = s.q.DeleteMomentRepost(ctx, query.DeleteMomentRepostParams{MomentID: momentID, VisitorID: visitorID})
		_ = s.q.DecrementMomentRepostCount(ctx, momentID)
		return false, nil
	}
	_ = s.q.CreateMomentRepost(ctx, query.CreateMomentRepostParams{MomentID: momentID, VisitorID: visitorID})
	_ = s.q.IncrementMomentRepostCount(ctx, momentID)
	return true, nil
}

func (s *MomentsService) CreateComment(ctx context.Context, momentID uuid.UUID, authorName, content, ipHash, uaHash string) (*MomentCommentItem, error) {
	row, err := s.q.CreateMomentComment(ctx, query.CreateMomentCommentParams{
		MomentID:   momentID,
		AuthorName: authorName,
		Content:    content,
		IpHash:     ipHash,
		UaHash:     uaHash,
	})
	if err != nil {
		return nil, fmt.Errorf("create comment: %w", err)
	}
	_ = s.q.IncrementMomentCommentCount(ctx, momentID)

	return &MomentCommentItem{
		ID:         row.ID,
		AuthorName: authorName,
		Content:    content,
		LikeCount:  0,
		CreatedAt:  row.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (s *MomentsService) ListComments(ctx context.Context, momentID uuid.UUID, page, size int, visitorID uuid.UUID) ([]MomentCommentItem, int64, error) {
	total, err := s.q.CountMomentComments(ctx, momentID)
	if err != nil {
		return nil, 0, err
	}

	rows, err := s.q.ListMomentComments(ctx, query.ListMomentCommentsParams{
		MomentID: momentID,
		Limit:    int32(size),
		Offset:   int32((page - 1) * size),
	})
	if err != nil {
		return nil, 0, err
	}

	cids := make([]uuid.UUID, len(rows))
	for i, r := range rows {
		cids[i] = r.ID
	}

	likedSet := make(map[uuid.UUID]bool)
	if visitorID != uuid.Nil && len(cids) > 0 {
		liked, _ := s.q.GetVisitorMomentCommentLikes(ctx, query.GetVisitorMomentCommentLikesParams{
			VisitorID: visitorID, Column2: cids,
		})
		for _, id := range liked {
			likedSet[id] = true
		}
	}

	items := make([]MomentCommentItem, 0, len(rows))
	for _, r := range rows {
		items = append(items, MomentCommentItem{
			ID:         r.ID,
			AuthorName: r.AuthorName,
			Content:    r.Content,
			LikeCount:  r.LikeCount,
			CreatedAt:  r.CreatedAt.Format("2006-01-02T15:04:05Z"),
			Liked:      likedSet[r.ID],
		})
	}

	return items, total, nil
}

func (s *MomentsService) ToggleCommentLike(ctx context.Context, commentID, visitorID uuid.UUID) (bool, error) {
	cnt, err := s.q.CheckMomentCommentLike(ctx, query.CheckMomentCommentLikeParams{CommentID: commentID, VisitorID: visitorID})
	if err != nil {
		return false, err
	}
	if cnt > 0 {
		_ = s.q.DeleteMomentCommentLike(ctx, query.DeleteMomentCommentLikeParams{CommentID: commentID, VisitorID: visitorID})
		_ = s.q.DecrementMomentCommentLikeCount(ctx, commentID)
		return false, nil
	}
	_ = s.q.CreateMomentCommentLike(ctx, query.CreateMomentCommentLikeParams{CommentID: commentID, VisitorID: visitorID})
	_ = s.q.IncrementMomentCommentLikeCount(ctx, commentID)
	return true, nil
}

func (s *MomentsService) ListAdmin(ctx context.Context, page, size int) ([]MomentItem, int64, error) {
	total, err := s.q.CountAdminMoments(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("count admin moments: %w", err)
	}

	rows, err := s.q.ListAdminMoments(ctx, query.ListAdminMomentsParams{
		Limit:  int32(size),
		Offset: int32((page - 1) * size),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("list admin moments: %w", err)
	}

	items := make([]MomentItem, 0, len(rows))
	for _, r := range rows {
		item := MomentItem{
			ID:            r.ID,
			AuthorName:    r.AuthorName,
			Content:       r.Content,
			ImageURLs:     parseImageURLs(r.ImageUrls),
			LikeCount:     r.LikeCount,
			RepostCount:   r.RepostCount,
			CommentCount:  r.CommentCount,
			PublishStatus: string(r.PublishStatus),
			CreatedAt:     r.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
		if r.PublishedAt.Valid {
			item.PublishedAt = r.PublishedAt.Time.Format(time.RFC3339)
		}
		if r.ScheduledAt.Valid {
			item.ScheduledAt = r.ScheduledAt.Time.Format(time.RFC3339)
		}
		items = append(items, item)
	}
	return items, total, nil
}

func (s *MomentsService) Publish(ctx context.Context, momentID uuid.UUID) error {
	return s.q.PublishMoment(ctx, momentID)
}

func (s *MomentsService) Schedule(ctx context.Context, momentID uuid.UUID, at time.Time) error {
	return s.q.ScheduleMoment(ctx, query.ScheduleMomentParams{
		ID:          momentID,
		ScheduledAt: pgtype.Timestamptz{Time: at, Valid: true},
	})
}

func (s *MomentsService) Unpublish(ctx context.Context, momentID uuid.UUID) error {
	return s.q.UnpublishMoment(ctx, momentID)
}

func (s *MomentsService) PublishDueScheduled(ctx context.Context) (int, error) {
	rows, err := s.q.ListScheduledMomentsDue(ctx)
	if err != nil {
		return 0, fmt.Errorf("list scheduled moments due: %w", err)
	}
	if len(rows) == 0 {
		return 0, nil
	}

	published := 0
	for _, row := range rows {
		if err := s.q.PublishMoment(ctx, row); err != nil {
			return published, fmt.Errorf("publish scheduled moment %s: %w", row, err)
		}
		published++
	}
	return published, nil
}

func (s *MomentsService) getInteractionSets(ctx context.Context, visitorID uuid.UUID, ids []uuid.UUID) (likes, reposts map[uuid.UUID]bool) {
	likes = make(map[uuid.UUID]bool)
	reposts = make(map[uuid.UUID]bool)
	if visitorID == uuid.Nil || len(ids) == 0 {
		return
	}

	likedRows, _ := s.q.GetVisitorMomentLikes(ctx, query.GetVisitorMomentLikesParams{VisitorID: visitorID, Column2: ids})
	for _, id := range likedRows {
		likes[id] = true
	}

	repostRows, _ := s.q.GetVisitorMomentReposts(ctx, query.GetVisitorMomentRepostsParams{VisitorID: visitorID, Column2: ids})
	for _, id := range repostRows {
		reposts[id] = true
	}
	return
}

func parseImageURLs(raw json.RawMessage) []string {
	var urls []string
	if len(raw) > 0 {
		_ = json.Unmarshal(raw, &urls)
	}
	if urls == nil {
		urls = []string{}
	}
	return urls
}

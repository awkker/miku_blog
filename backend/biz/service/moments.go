package service

import (
	"context"
	"encoding/json"
	"fmt"

	"nanamiku-blog/backend/query"

	"github.com/google/uuid"
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
	ID           uuid.UUID `json:"id"`
	AuthorName   string    `json:"author_name"`
	Content      string    `json:"content"`
	ImageURLs    []string  `json:"image_urls"`
	LikeCount    int64     `json:"like_count"`
	RepostCount  int64     `json:"repost_count"`
	CommentCount int64     `json:"comment_count"`
	CreatedAt    string    `json:"created_at"`
	Liked        bool      `json:"liked"`
	Reposted     bool      `json:"reposted"`
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
		items = append(items, MomentItem{
			ID:           r.ID,
			AuthorName:   r.AuthorName,
			Content:      r.Content,
			ImageURLs:    parseImageURLs(r.ImageUrls),
			LikeCount:    r.LikeCount,
			RepostCount:  r.RepostCount,
			CommentCount: r.CommentCount,
			CreatedAt:    r.CreatedAt.Format("2006-01-02T15:04:05Z"),
			Liked:        likeSet[r.ID],
			Reposted:     repostSet[r.ID],
		})
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

func (s *MomentsService) Create(ctx context.Context, authorName, content string, imageURLs []string, ipHash, uaHash string) (*MomentItem, error) {
	imgs, _ := json.Marshal(imageURLs)

	row, err := s.q.CreateMoment(ctx, query.CreateMomentParams{
		AuthorName: authorName,
		Content:    content,
		ImageUrls:  imgs,
		IpHash:     ipHash,
		UaHash:     uaHash,
	})
	if err != nil {
		return nil, fmt.Errorf("create moment: %w", err)
	}

	return &MomentItem{
		ID:         row.ID,
		AuthorName: authorName,
		Content:    content,
		ImageURLs:  imageURLs,
		CreatedAt:  row.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
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

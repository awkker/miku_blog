package service

import (
	"context"
	"fmt"

	"nanamiku-blog/backend/query"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostCommentsService struct {
	q  *query.Queries
	db *pgxpool.Pool
}

func NewPostCommentsService(db *pgxpool.Pool) *PostCommentsService {
	return &PostCommentsService{q: query.New(db), db: db}
}

type PostCommentItem struct {
	ID         uuid.UUID  `json:"id"`
	PostID     uuid.UUID  `json:"post_id"`
	ParentID   *uuid.UUID `json:"parent_id,omitempty"`
	AuthorName string     `json:"author_name"`
	Content    string     `json:"content"`
	CreatedAt  string     `json:"created_at"`
}

func (s *PostCommentsService) ListComments(ctx context.Context, postID uuid.UUID, page, size int) ([]PostCommentItem, int64, error) {
	total, err := s.q.CountApprovedComments(ctx, postID)
	if err != nil {
		return nil, 0, err
	}

	rows, err := s.q.ListApprovedComments(ctx, query.ListApprovedCommentsParams{
		PostID: postID,
		Limit:  int32(size),
		Offset: int32((page - 1) * size),
	})
	if err != nil {
		return nil, 0, err
	}

	items := make([]PostCommentItem, 0, len(rows))
	for _, r := range rows {
		item := PostCommentItem{
			ID:         r.ID,
			PostID:     r.PostID,
			AuthorName: r.AuthorName,
			Content:    r.Content,
			CreatedAt:  r.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
		if r.ParentID.Valid {
			pid := uuid.UUID(r.ParentID.Bytes)
			item.ParentID = &pid
		}
		items = append(items, item)
	}
	return items, total, nil
}

func (s *PostCommentsService) CreateComment(ctx context.Context, postID uuid.UUID, parentID *uuid.UUID, authorName, authorEmail, content, ipHash, uaHash string, visitorID uuid.UUID) (*PostCommentItem, error) {
	var pid pgtype.UUID
	if parentID != nil {
		pid = pgtype.UUID{Bytes: *parentID, Valid: true}
	}

	row, err := s.q.CreatePostComment(ctx, query.CreatePostCommentParams{
		PostID:      postID,
		ParentID:    pid,
		AuthorName:  authorName,
		AuthorEmail: authorEmail,
		Content:     content,
		IpHash:      ipHash,
		UaHash:      uaHash,
	})
	if err != nil {
		return nil, fmt.Errorf("create comment: %w", err)
	}

	_ = s.q.IncrementPostCommentCount(ctx, postID)

	return &PostCommentItem{
		ID:         row.ID,
		PostID:     postID,
		AuthorName: authorName,
		Content:    content,
		CreatedAt:  row.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

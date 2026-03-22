package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"nanamiku-blog/backend/query"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostsService struct {
	q  *query.Queries
	db *pgxpool.Pool
}

func NewPostsService(db *pgxpool.Pool) *PostsService {
	return &PostsService{q: query.New(db), db: db}
}

type PostListItem struct {
	ID           uuid.UUID `json:"id"`
	Slug         string    `json:"slug"`
	Title        string    `json:"title"`
	Excerpt      string    `json:"excerpt"`
	HeroImageURL string    `json:"hero_image_url"`
	Category     string    `json:"category"`
	PublishedAt  string    `json:"published_at,omitempty"`
	ViewCount    int64     `json:"view_count"`
	LikeCount    int64     `json:"like_count"`
	CommentCount int64     `json:"comment_count"`
	CreatedAt    string    `json:"created_at"`
	Tags         []TagItem `json:"tags,omitempty"`
}

type TagItem struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type PostDetail struct {
	ID              uuid.UUID `json:"id"`
	Slug            string    `json:"slug"`
	Title           string    `json:"title"`
	Excerpt         string    `json:"excerpt"`
	ContentMarkdown string    `json:"content_markdown"`
	HeroImageURL    string    `json:"hero_image_url"`
	Category        string    `json:"category"`
	Status          string    `json:"status"`
	PublishedAt     string    `json:"published_at,omitempty"`
	ScheduledAt     string    `json:"scheduled_at,omitempty"`
	ViewCount       int64     `json:"view_count"`
	LikeCount       int64     `json:"like_count"`
	CommentCount    int64     `json:"comment_count"`
	CreatedAt       string    `json:"created_at"`
	UpdatedAt       string    `json:"updated_at"`
	Tags            []TagItem `json:"tags"`
	Liked           bool      `json:"liked"`
}

type AdminPostListItem struct {
	ID           uuid.UUID `json:"id"`
	Slug         string    `json:"slug"`
	Title        string    `json:"title"`
	Category     string    `json:"category"`
	Status       string    `json:"status"`
	PublishedAt  string    `json:"published_at,omitempty"`
	ScheduledAt  string    `json:"scheduled_at,omitempty"`
	ViewCount    int64     `json:"view_count"`
	LikeCount    int64     `json:"like_count"`
	CommentCount int64     `json:"comment_count"`
	CreatedAt    string    `json:"created_at"`
	UpdatedAt    string    `json:"updated_at"`
}

func (s *PostsService) ListPublished(ctx context.Context, page, size int) ([]PostListItem, int64, error) {
	total, err := s.q.CountPublishedPosts(ctx)
	if err != nil {
		return nil, 0, err
	}

	rows, err := s.q.ListPublishedPosts(ctx, query.ListPublishedPostsParams{
		Limit:  int32(size),
		Offset: int32((page - 1) * size),
	})
	if err != nil {
		return nil, 0, err
	}

	items := make([]PostListItem, 0, len(rows))
	for _, r := range rows {
		tags, _ := s.q.GetPostTagNames(ctx, r.ID)
		item := PostListItem{
			ID:           r.ID,
			Slug:         r.Slug,
			Title:        r.Title,
			Excerpt:      r.Excerpt,
			HeroImageURL: r.HeroImageUrl,
			Category:     r.Category,
			ViewCount:    r.ViewCount,
			LikeCount:    r.LikeCount,
			CommentCount: r.CommentCount,
			CreatedAt:    r.CreatedAt.Format(time.RFC3339),
			Tags:         toTagItems(tags),
		}
		if r.PublishedAt.Valid {
			item.PublishedAt = r.PublishedAt.Time.Format(time.RFC3339)
		}
		items = append(items, item)
	}

	return items, total, nil
}

func (s *PostsService) ListByCategory(ctx context.Context, category string, page, size int) ([]PostListItem, error) {
	rows, err := s.q.ListPostsByCategory(ctx, query.ListPostsByCategoryParams{
		Category: category,
		Limit:    int32(size),
		Offset:   int32((page - 1) * size),
	})
	if err != nil {
		return nil, err
	}

	items := make([]PostListItem, 0, len(rows))
	for _, r := range rows {
		tags, _ := s.q.GetPostTagNames(ctx, r.ID)
		item := PostListItem{
			ID:           r.ID,
			Slug:         r.Slug,
			Title:        r.Title,
			Excerpt:      r.Excerpt,
			HeroImageURL: r.HeroImageUrl,
			Category:     r.Category,
			ViewCount:    r.ViewCount,
			LikeCount:    r.LikeCount,
			CommentCount: r.CommentCount,
			CreatedAt:    r.CreatedAt.Format(time.RFC3339),
			Tags:         toTagItems(tags),
		}
		if r.PublishedAt.Valid {
			item.PublishedAt = r.PublishedAt.Time.Format(time.RFC3339)
		}
		items = append(items, item)
	}
	return items, nil
}

func (s *PostsService) ListHot(ctx context.Context, limit int) ([]PostListItem, error) {
	rows, err := s.q.ListHotPosts(ctx, int32(limit))
	if err != nil {
		return nil, err
	}

	items := make([]PostListItem, 0, len(rows))
	for _, r := range rows {
		item := PostListItem{
			ID:           r.ID,
			Slug:         r.Slug,
			Title:        r.Title,
			Excerpt:      r.Excerpt,
			HeroImageURL: r.HeroImageUrl,
			Category:     r.Category,
			ViewCount:    r.ViewCount,
			LikeCount:    r.LikeCount,
			CommentCount: r.CommentCount,
			CreatedAt:    r.CreatedAt.Format(time.RFC3339),
		}
		if r.PublishedAt.Valid {
			item.PublishedAt = r.PublishedAt.Time.Format(time.RFC3339)
		}
		items = append(items, item)
	}
	return items, nil
}

func (s *PostsService) Search(ctx context.Context, q string, page, size int) ([]PostListItem, error) {
	rows, err := s.q.SearchPosts(ctx, query.SearchPostsParams{
		WebsearchToTsquery: q,
		Limit:              int32(size),
		Offset:             int32((page - 1) * size),
	})
	if err != nil {
		return nil, err
	}

	items := make([]PostListItem, 0, len(rows))
	for _, r := range rows {
		item := PostListItem{
			ID:           r.ID,
			Slug:         r.Slug,
			Title:        r.Title,
			Excerpt:      r.Excerpt,
			HeroImageURL: r.HeroImageUrl,
			Category:     r.Category,
			ViewCount:    r.ViewCount,
			LikeCount:    r.LikeCount,
			CommentCount: r.CommentCount,
		}
		if r.PublishedAt.Valid {
			item.PublishedAt = r.PublishedAt.Time.Format(time.RFC3339)
		}
		items = append(items, item)
	}
	return items, nil
}

func (s *PostsService) GetBySlug(ctx context.Context, slug string, visitorID uuid.UUID) (*PostDetail, error) {
	r, err := s.q.GetPostBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	tags, _ := s.q.GetPostTagNames(ctx, r.ID)

	liked := false
	if visitorID != uuid.Nil {
		cnt, _ := s.q.CheckPostLike(ctx, query.CheckPostLikeParams{PostID: r.ID, VisitorID: visitorID})
		liked = cnt > 0
	}

	_ = s.q.IncrementPostViewCount(ctx, r.ID)
	today := time.Now().Truncate(24 * time.Hour)
	uvIncr := int64(1)
	_ = s.q.UpsertPostViewDaily(ctx, query.UpsertPostViewDailyParams{
		PostID: r.ID,
		Day:    pgtype.Date{Time: today, Valid: true},
		Uv:     uvIncr,
	})

	detail := &PostDetail{
		ID:              r.ID,
		Slug:            r.Slug,
		Title:           r.Title,
		Excerpt:         r.Excerpt,
		ContentMarkdown: r.ContentMarkdown,
		HeroImageURL:    r.HeroImageUrl,
		Category:        r.Category,
		Status:          string(r.Status),
		ViewCount:       r.ViewCount + 1,
		LikeCount:       r.LikeCount,
		CommentCount:    r.CommentCount,
		CreatedAt:       r.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       r.UpdatedAt.Format(time.RFC3339),
		Tags:            toTagItems(tags),
		Liked:           liked,
	}
	if r.PublishedAt.Valid {
		detail.PublishedAt = r.PublishedAt.Time.Format(time.RFC3339)
	}
	if r.ScheduledAt.Valid {
		detail.ScheduledAt = r.ScheduledAt.Time.Format(time.RFC3339)
	}
	return detail, nil
}

func (s *PostsService) GetByID(ctx context.Context, id uuid.UUID) (*PostDetail, error) {
	r, err := s.q.GetPostByID(ctx, id)
	if err != nil {
		return nil, err
	}
	tags, _ := s.q.GetPostTagNames(ctx, r.ID)
	detail := &PostDetail{
		ID:              r.ID,
		Slug:            r.Slug,
		Title:           r.Title,
		Excerpt:         r.Excerpt,
		ContentMarkdown: r.ContentMarkdown,
		HeroImageURL:    r.HeroImageUrl,
		Category:        r.Category,
		Status:          string(r.Status),
		ViewCount:       r.ViewCount,
		LikeCount:       r.LikeCount,
		CommentCount:    r.CommentCount,
		CreatedAt:       r.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       r.UpdatedAt.Format(time.RFC3339),
		Tags:            toTagItems(tags),
	}
	if r.PublishedAt.Valid {
		detail.PublishedAt = r.PublishedAt.Time.Format(time.RFC3339)
	}
	if r.ScheduledAt.Valid {
		detail.ScheduledAt = r.ScheduledAt.Time.Format(time.RFC3339)
	}
	return detail, nil
}

func (s *PostsService) ToggleLike(ctx context.Context, postID, visitorID uuid.UUID) (bool, error) {
	cnt, err := s.q.CheckPostLike(ctx, query.CheckPostLikeParams{PostID: postID, VisitorID: visitorID})
	if err != nil {
		return false, err
	}
	if cnt > 0 {
		_ = s.q.DeletePostLike(ctx, query.DeletePostLikeParams{PostID: postID, VisitorID: visitorID})
		_ = s.q.DecrementPostLikeCount(ctx, postID)
		return false, nil
	}
	_ = s.q.CreatePostLike(ctx, query.CreatePostLikeParams{PostID: postID, VisitorID: visitorID})
	_ = s.q.IncrementPostLikeCount(ctx, postID)
	return true, nil
}

type CreatePostInput struct {
	Slug            string
	Title           string
	Excerpt         string
	ContentMarkdown string
	HeroImageURL    string
	Category        string
	Status          string
	ScheduledAt     *time.Time
	Tags            []string
	AdminID         uuid.UUID
}

func (s *PostsService) Create(ctx context.Context, in CreatePostInput) (uuid.UUID, error) {
	var scheduledAt interface{}
	if in.ScheduledAt != nil {
		scheduledAt = pgtype.Timestamptz{Time: *in.ScheduledAt, Valid: true}
	}

	row, err := s.q.CreatePost(ctx, query.CreatePostParams{
		Slug:            in.Slug,
		Title:           in.Title,
		Excerpt:         in.Excerpt,
		ContentMarkdown: in.ContentMarkdown,
		HeroImageUrl:    in.HeroImageURL,
		Category:        in.Category,
		Status:          query.PostStatus(in.Status),
		CreatedBy:       pgtype.UUID{Bytes: in.AdminID, Valid: true},
		Column9:         scheduledAt,
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("create post: %w", err)
	}

	if err := s.syncTags(ctx, row.ID, in.Tags); err != nil {
		return uuid.Nil, err
	}

	return row.ID, nil
}

type UpdatePostInput struct {
	ID              uuid.UUID
	Slug            string
	Title           string
	Excerpt         string
	ContentMarkdown string
	HeroImageURL    string
	Category        string
	Tags            []string
	AdminID         uuid.UUID
}

func (s *PostsService) Update(ctx context.Context, in UpdatePostInput) error {
	_ = s.q.CreatePostRevision(ctx, query.CreatePostRevisionParams{
		PostID:          in.ID,
		Title:           in.Title,
		Excerpt:         in.Excerpt,
		ContentMarkdown: in.ContentMarkdown,
		HeroImageUrl:    in.HeroImageURL,
		Category:        in.Category,
		EditorID:        pgtype.UUID{Bytes: in.AdminID, Valid: true},
	})

	if err := s.q.UpdatePost(ctx, query.UpdatePostParams{
		ID:              in.ID,
		Slug:            in.Slug,
		Title:           in.Title,
		Excerpt:         in.Excerpt,
		ContentMarkdown: in.ContentMarkdown,
		HeroImageUrl:    in.HeroImageURL,
		Category:        in.Category,
		UpdatedBy:       pgtype.UUID{Bytes: in.AdminID, Valid: true},
	}); err != nil {
		return fmt.Errorf("update post: %w", err)
	}

	return s.syncTags(ctx, in.ID, in.Tags)
}

func (s *PostsService) Publish(ctx context.Context, postID, adminID uuid.UUID) error {
	return s.q.PublishPost(ctx, query.PublishPostParams{
		ID:        postID,
		UpdatedBy: pgtype.UUID{Bytes: adminID, Valid: true},
	})
}

func (s *PostsService) Unpublish(ctx context.Context, postID, adminID uuid.UUID) error {
	return s.q.UnpublishPost(ctx, query.UnpublishPostParams{
		ID:        postID,
		UpdatedBy: pgtype.UUID{Bytes: adminID, Valid: true},
	})
}

func (s *PostsService) Schedule(ctx context.Context, postID, adminID uuid.UUID, at time.Time) error {
	return s.q.SchedulePost(ctx, query.SchedulePostParams{
		ID:          postID,
		ScheduledAt: pgtype.Timestamptz{Time: at, Valid: true},
		UpdatedBy:   pgtype.UUID{Bytes: adminID, Valid: true},
	})
}

func (s *PostsService) Delete(ctx context.Context, postID uuid.UUID) error {
	return s.q.DeletePost(ctx, postID)
}

func (s *PostsService) ListAdmin(ctx context.Context, page, size int) ([]AdminPostListItem, int64, error) {
	total, err := s.q.CountAdminPosts(ctx)
	if err != nil {
		return nil, 0, err
	}

	rows, err := s.q.ListAdminPosts(ctx, query.ListAdminPostsParams{
		Limit:  int32(size),
		Offset: int32((page - 1) * size),
	})
	if err != nil {
		return nil, 0, err
	}

	items := make([]AdminPostListItem, 0, len(rows))
	for _, r := range rows {
		item := AdminPostListItem{
			ID:           r.ID,
			Slug:         r.Slug,
			Title:        r.Title,
			Category:     r.Category,
			Status:       string(r.Status),
			ViewCount:    r.ViewCount,
			LikeCount:    r.LikeCount,
			CommentCount: r.CommentCount,
			CreatedAt:    r.CreatedAt.Format(time.RFC3339),
			UpdatedAt:    r.UpdatedAt.Format(time.RFC3339),
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

func (s *PostsService) PublishDueScheduled(ctx context.Context) (int, error) {
	rows, err := s.q.ListScheduledPostsDue(ctx)
	if err != nil {
		return 0, fmt.Errorf("list scheduled posts due: %w", err)
	}
	if len(rows) == 0 {
		return 0, nil
	}

	published := 0
	for _, row := range rows {
		if err := s.q.PublishPost(ctx, query.PublishPostParams{
			ID:        row.ID,
			UpdatedBy: pgtype.UUID{},
		}); err != nil {
			return published, fmt.Errorf("publish scheduled post %s: %w", row.ID, err)
		}
		published++
	}
	return published, nil
}

func (s *PostsService) syncTags(ctx context.Context, postID uuid.UUID, tagNames []string) error {
	_ = s.q.SetPostTags(ctx, postID)
	for _, name := range tagNames {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		slug := strings.ToLower(strings.ReplaceAll(name, " ", "-"))
		tagID, err := s.q.UpsertTag(ctx, query.UpsertTagParams{Name: name, Slug: slug})
		if err != nil {
			return fmt.Errorf("upsert tag %q: %w", name, err)
		}
		_ = s.q.AddPostTag(ctx, query.AddPostTagParams{PostID: postID, TagID: tagID})
	}
	return nil
}

func toTagItems(rows []query.GetPostTagNamesRow) []TagItem {
	items := make([]TagItem, 0, len(rows))
	for _, r := range rows {
		items = append(items, TagItem{Name: r.Name, Slug: r.Slug})
	}
	return items
}

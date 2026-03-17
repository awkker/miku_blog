package service

import (
	"context"
	"encoding/json"
	"fmt"

	"nanamiku-blog/backend/query"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ModerationService struct {
	q  *query.Queries
	db *pgxpool.Pool
}

func NewModerationService(db *pgxpool.Pool) *ModerationService {
	return &ModerationService{q: query.New(db), db: db}
}

func (s *ModerationService) ListSensitiveWords(ctx context.Context) ([]string, error) {
	return s.q.ListSensitiveWords(ctx)
}

func (s *ModerationService) CheckBlocked(ctx context.Context, ipHash string) (bool, error) {
	cnt, err := s.q.CheckBlockedIP(ctx, ipHash)
	if err != nil {
		return false, err
	}
	return cnt > 0, nil
}

func (s *ModerationService) LogAudit(ctx context.Context, adminID uuid.UUID, action, targetType, targetID string, detail interface{}, ip string) error {
	detailJSON, _ := json.Marshal(detail)
	return s.q.CreateAuditLog(ctx, query.CreateAuditLogParams{
		AdminID:    pgtype.UUID{Bytes: adminID, Valid: adminID != uuid.Nil},
		Action:     action,
		TargetType: targetType,
		TargetID:   targetID,
		Detail:     detailJSON,
		Ip:         ip,
	})
}

type AuditLogItem struct {
	ID            uuid.UUID       `json:"id"`
	Action        string          `json:"action"`
	TargetType    string          `json:"target_type"`
	TargetID      string          `json:"target_id"`
	Detail        json.RawMessage `json:"detail"`
	IP            string          `json:"ip"`
	AdminUsername string          `json:"admin_username,omitempty"`
	CreatedAt     string          `json:"created_at"`
}

func (s *ModerationService) ListAuditLogs(ctx context.Context, page, size int) ([]AuditLogItem, error) {
	rows, err := s.q.ListAuditLogs(ctx, query.ListAuditLogsParams{
		Limit:  int32(size),
		Offset: int32((page - 1) * size),
	})
	if err != nil {
		return nil, fmt.Errorf("list audit logs: %w", err)
	}

	items := make([]AuditLogItem, 0, len(rows))
	for _, r := range rows {
		var username string
		if r.AdminUsername.Valid {
			username = r.AdminUsername.String
		}
		items = append(items, AuditLogItem{
			ID:            r.ID,
			Action:        r.Action,
			TargetType:    r.TargetType,
			TargetID:      r.TargetID,
			Detail:        r.Detail,
			IP:            r.Ip,
			AdminUsername: username,
			CreatedAt:     r.CreatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}
	return items, nil
}

type AdminCommentItem struct {
	ID          uuid.UUID `json:"id"`
	PostID      uuid.UUID `json:"post_id"`
	PostTitle   string    `json:"post_title,omitempty"`
	AuthorName  string    `json:"author_name"`
	AuthorEmail string    `json:"author_email"`
	Content     string    `json:"content"`
	Status      string    `json:"status"`
	IPHash      string    `json:"ip_hash"`
	CreatedAt   string    `json:"created_at"`
}

func (s *ModerationService) ListComments(ctx context.Context, status string, page, size int) ([]AdminCommentItem, int64, error) {
	var statusParam query.NullModerationStatus
	if status != "" {
		statusParam = query.NullModerationStatus{
			ModerationStatus: query.ModerationStatus(status),
			Valid:            true,
		}
	}

	total, err := s.q.CountAdminComments(ctx, statusParam)
	if err != nil {
		return nil, 0, err
	}

	rows, err := s.q.ListAdminComments(ctx, query.ListAdminCommentsParams{
		Status: statusParam,
		Limit:  int32(size),
		Offset: int32((page - 1) * size),
	})
	if err != nil {
		return nil, 0, err
	}

	items := make([]AdminCommentItem, 0, len(rows))
	for _, r := range rows {
		var postTitle string
		if r.PostTitle.Valid {
			postTitle = r.PostTitle.String
		}
		items = append(items, AdminCommentItem{
			ID:          r.ID,
			PostID:      r.PostID,
			PostTitle:   postTitle,
			AuthorName:  r.AuthorName,
			AuthorEmail: r.AuthorEmail,
			Content:     r.Content,
			Status:      string(r.Status),
			IPHash:      r.IpHash,
			CreatedAt:   r.CreatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}
	return items, total, nil
}

func (s *ModerationService) ApproveComment(ctx context.Context, commentID, adminID uuid.UUID) error {
	return s.q.ApproveComment(ctx, query.ApproveCommentParams{ID: commentID, ReviewedBy: pgtype.UUID{Bytes: adminID, Valid: true}})
}

func (s *ModerationService) RejectComment(ctx context.Context, commentID, adminID uuid.UUID) error {
	return s.q.RejectComment(ctx, query.RejectCommentParams{ID: commentID, ReviewedBy: pgtype.UUID{Bytes: adminID, Valid: true}})
}

func (s *ModerationService) DeleteComment(ctx context.Context, commentID uuid.UUID) error {
	return s.q.DeleteComment(ctx, commentID)
}

type AdminGuestbookMessageItem struct {
	ID               uuid.UUID  `json:"id"`
	ParentID         *uuid.UUID `json:"parent_id,omitempty"`
	ParentAuthorName string     `json:"parent_author_name,omitempty"`
	AuthorName       string     `json:"author_name"`
	AuthorWebsite    string     `json:"author_website,omitempty"`
	Content          string     `json:"content"`
	Status           string     `json:"status"`
	IPHash           string     `json:"ip_hash"`
	VoteScore        int32      `json:"vote_score"`
	CreatedAt        string     `json:"created_at"`
}

func (s *ModerationService) ListGuestbookMessages(ctx context.Context, status string, page, size int) ([]AdminGuestbookMessageItem, int64, error) {
	var statusParam query.NullModerationStatus
	if status != "" {
		statusParam = query.NullModerationStatus{
			ModerationStatus: query.ModerationStatus(status),
			Valid:            true,
		}
	}

	total, err := s.q.CountAdminGuestbookMessages(ctx, statusParam)
	if err != nil {
		return nil, 0, err
	}

	rows, err := s.q.ListAdminGuestbookMessages(ctx, query.ListAdminGuestbookMessagesParams{
		Status: statusParam,
		Limit:  int32(size),
		Offset: int32((page - 1) * size),
	})
	if err != nil {
		return nil, 0, err
	}

	items := make([]AdminGuestbookMessageItem, 0, len(rows))
	for _, r := range rows {
		item := AdminGuestbookMessageItem{
			ID:            r.ID,
			AuthorName:    r.AuthorName,
			AuthorWebsite: r.AuthorWebsite,
			Content:       r.Content,
			Status:        string(r.Status),
			IPHash:        r.IpHash,
			VoteScore:     r.VoteScore,
			CreatedAt:     r.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
		if r.ParentID.Valid {
			pid := uuid.UUID(r.ParentID.Bytes)
			item.ParentID = &pid
		}
		if r.ParentAuthorName.Valid {
			item.ParentAuthorName = r.ParentAuthorName.String
		}
		items = append(items, item)
	}
	return items, total, nil
}

func (s *ModerationService) ApproveGuestbookMessage(ctx context.Context, messageID, adminID uuid.UUID) error {
	return s.q.ApproveGuestbookMessage(ctx, query.ApproveGuestbookMessageParams{
		ID:         messageID,
		ReviewedBy: pgtype.UUID{Bytes: adminID, Valid: true},
	})
}

func (s *ModerationService) RejectGuestbookMessage(ctx context.Context, messageID, adminID uuid.UUID) error {
	return s.q.RejectGuestbookMessage(ctx, query.RejectGuestbookMessageParams{
		ID:         messageID,
		ReviewedBy: pgtype.UUID{Bytes: adminID, Valid: true},
	})
}

func (s *ModerationService) DeleteGuestbookMessage(ctx context.Context, messageID uuid.UUID) error {
	return s.q.DeleteGuestbookMessage(ctx, messageID)
}

type AdminFriendItem struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	URL          string    `json:"url"`
	Domain       string    `json:"domain"`
	AvatarURL    string    `json:"avatar_url"`
	Status       string    `json:"status"`
	HealthStatus string    `json:"health_status"`
	SortOrder    int32     `json:"sort_order"`
	CreatedAt    string    `json:"created_at"`
}

func (s *ModerationService) ListAdminFriends(ctx context.Context, page, size int) ([]AdminFriendItem, int64, error) {
	total, err := s.q.CountAdminFriendLinks(ctx)
	if err != nil {
		return nil, 0, err
	}

	rows, err := s.q.ListAdminFriendLinks(ctx, query.ListAdminFriendLinksParams{
		Limit:  int32(size),
		Offset: int32((page - 1) * size),
	})
	if err != nil {
		return nil, 0, err
	}

	items := make([]AdminFriendItem, 0, len(rows))
	for _, r := range rows {
		items = append(items, AdminFriendItem{
			ID:           r.ID,
			Name:         r.Name,
			Description:  r.Description,
			URL:          r.Url,
			Domain:       r.Domain,
			AvatarURL:    r.AvatarUrl,
			Status:       string(r.Status),
			HealthStatus: string(r.HealthStatus),
			SortOrder:    r.SortOrder,
			CreatedAt:    r.CreatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}
	return items, total, nil
}

func (s *ModerationService) CreateFriend(ctx context.Context, name, description, url, domain, avatarURL string, sortOrder int32, adminID uuid.UUID) (uuid.UUID, error) {
	row, err := s.q.CreateFriendLink(ctx, query.CreateFriendLinkParams{
		Name:        name,
		Description: description,
		Url:         url,
		Domain:      domain,
		AvatarUrl:   avatarURL,
		SortOrder:   sortOrder,
		ReviewedBy:  pgtype.UUID{Bytes: adminID, Valid: true},
	})
	if err != nil {
		return uuid.Nil, err
	}
	return row.ID, nil
}

func (s *ModerationService) UpdateFriend(ctx context.Context, id uuid.UUID, name, description, url, domain, avatarURL string, sortOrder int32) error {
	return s.q.UpdateFriendLink(ctx, query.UpdateFriendLinkParams{
		ID:          id,
		Name:        name,
		Description: description,
		Url:         url,
		Domain:      domain,
		AvatarUrl:   avatarURL,
		SortOrder:   sortOrder,
	})
}

func (s *ModerationService) DeleteFriend(ctx context.Context, id uuid.UUID) error {
	return s.q.DeleteFriendLink(ctx, id)
}

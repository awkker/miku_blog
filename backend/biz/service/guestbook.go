package service

import (
	"context"
	"fmt"

	"nanamiku-blog/backend/query"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GuestbookService struct {
	q  *query.Queries
	db *pgxpool.Pool
}

func NewGuestbookService(db *pgxpool.Pool) *GuestbookService {
	return &GuestbookService{q: query.New(db), db: db}
}

type GuestbookMessageItem struct {
	ID            uuid.UUID              `json:"id"`
	ParentID      *uuid.UUID             `json:"parent_id,omitempty"`
	AuthorName    string                 `json:"author_name"`
	AuthorWebsite string                 `json:"author_website,omitempty"`
	Content       string                 `json:"content"`
	VoteScore     int32                  `json:"vote_score"`
	CreatedAt     string                 `json:"created_at"`
	Replies       []GuestbookMessageItem `json:"replies,omitempty"`
	MyVote        *string                `json:"my_vote,omitempty"`
}

func (s *GuestbookService) ListMessages(ctx context.Context, sortBy string, page, size int, visitorID uuid.UUID) ([]GuestbookMessageItem, int64, error) {
	total, err := s.q.CountGuestbookMessages(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("count messages: %w", err)
	}

	rows, err := s.q.ListGuestbookMessages(ctx, query.ListGuestbookMessagesParams{
		SortBy: sortBy,
		Limit:  int32(size),
		Offset: int32((page - 1) * size),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("list messages: %w", err)
	}

	msgIDs := make([]uuid.UUID, len(rows))
	for i, r := range rows {
		msgIDs[i] = r.ID
	}

	var voteMap map[uuid.UUID]string
	if visitorID != uuid.Nil && len(msgIDs) > 0 {
		voteMap = s.getVoteMap(ctx, visitorID, msgIDs)
	}

	items := make([]GuestbookMessageItem, 0, len(rows))
	for _, r := range rows {
		item := GuestbookMessageItem{
			ID:            r.ID,
			AuthorName:    r.AuthorName,
			AuthorWebsite: r.AuthorWebsite,
			Content:       r.Content,
			VoteScore:     r.VoteScore,
			CreatedAt:     r.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
		if v, ok := voteMap[r.ID]; ok {
			item.MyVote = &v
		}

		replies, _ := s.q.ListGuestbookReplies(ctx, pgtype.UUID{Bytes: r.ID, Valid: true})
		replyItems := make([]GuestbookMessageItem, 0, len(replies))
		for _, rp := range replies {
			ri := GuestbookMessageItem{
				ID:            rp.ID,
				AuthorName:    rp.AuthorName,
				AuthorWebsite: rp.AuthorWebsite,
				Content:       rp.Content,
				VoteScore:     rp.VoteScore,
				CreatedAt:     rp.CreatedAt.Format("2006-01-02T15:04:05Z"),
			}
			replyItems = append(replyItems, ri)
		}
		item.Replies = replyItems
		items = append(items, item)
	}

	return items, total, nil
}

func (s *GuestbookService) CreateMessage(ctx context.Context, parentID *uuid.UUID, authorName, authorWebsite, content, ipHash, uaHash string) (*GuestbookMessageItem, error) {
	var pid pgtype.UUID
	if parentID != nil {
		pid = pgtype.UUID{Bytes: *parentID, Valid: true}
	}

	row, err := s.q.CreateGuestbookMessage(ctx, query.CreateGuestbookMessageParams{
		ParentID:      pid,
		AuthorName:    authorName,
		AuthorWebsite: authorWebsite,
		Content:       content,
		IpHash:        ipHash,
		UaHash:        uaHash,
	})
	if err != nil {
		return nil, fmt.Errorf("create message: %w", err)
	}

	return &GuestbookMessageItem{
		ID:            row.ID,
		AuthorName:    authorName,
		AuthorWebsite: authorWebsite,
		Content:       content,
		VoteScore:     0,
		CreatedAt:     row.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (s *GuestbookService) Vote(ctx context.Context, messageID, visitorID uuid.UUID, vote string) error {
	if vote == "" {
		_ = s.q.DeleteGuestbookVote(ctx, query.DeleteGuestbookVoteParams{
			MessageID: messageID,
			VisitorID: visitorID,
		})
	} else {
		vt := query.VoteType(vote)
		_ = s.q.UpsertGuestbookVote(ctx, query.UpsertGuestbookVoteParams{
			MessageID: messageID,
			VisitorID: visitorID,
			Vote:      vt,
		})
	}

	return s.q.RecalcGuestbookVoteScore(ctx, messageID)
}

func (s *GuestbookService) getVoteMap(ctx context.Context, visitorID uuid.UUID, msgIDs []uuid.UUID) map[uuid.UUID]string {
	m := make(map[uuid.UUID]string)
	votes, err := s.q.GetVisitorVotesForMessages(ctx, query.GetVisitorVotesForMessagesParams{
		VisitorID: visitorID,
		Column2:   msgIDs,
	})
	if err != nil {
		return m
	}
	for _, v := range votes {
		m[v.MessageID] = string(v.Vote)
	}
	return m
}

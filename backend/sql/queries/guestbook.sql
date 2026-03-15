-- name: ListGuestbookMessages :many
SELECT m.id, m.parent_id, m.author_name, m.author_website, m.content,
       m.vote_score, m.created_at
FROM guestbook_messages m
WHERE m.status = 'approved' AND m.parent_id IS NULL
ORDER BY
    CASE WHEN @sort_by::text = 'hot' THEN m.vote_score END DESC,
    CASE WHEN @sort_by::text = 'oldest' THEN m.created_at END ASC,
    CASE WHEN @sort_by::text = 'newest' OR @sort_by::text NOT IN ('hot','oldest') THEN m.created_at END DESC
LIMIT $1 OFFSET $2;

-- name: CountGuestbookMessages :one
SELECT count(*) FROM guestbook_messages WHERE status = 'approved' AND parent_id IS NULL;

-- name: ListGuestbookReplies :many
SELECT id, parent_id, author_name, author_website, content, vote_score, created_at
FROM guestbook_messages
WHERE parent_id = $1 AND status = 'approved'
ORDER BY created_at ASC;

-- name: CreateGuestbookMessage :one
INSERT INTO guestbook_messages (parent_id, author_name, author_website, content, ip_hash, ua_hash)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, created_at;

-- name: GetGuestbookVote :one
SELECT vote FROM guestbook_votes
WHERE message_id = $1 AND visitor_id = $2;

-- name: UpsertGuestbookVote :exec
INSERT INTO guestbook_votes (message_id, visitor_id, vote)
VALUES ($1, $2, $3)
ON CONFLICT (message_id, visitor_id) DO UPDATE
SET vote = EXCLUDED.vote, updated_at = now();

-- name: DeleteGuestbookVote :exec
DELETE FROM guestbook_votes WHERE message_id = $1 AND visitor_id = $2;

-- name: RecalcGuestbookVoteScore :exec
UPDATE guestbook_messages
SET vote_score = (
    SELECT coalesce(sum(CASE WHEN vote = 'up' THEN 1 ELSE -1 END), 0)
    FROM guestbook_votes WHERE message_id = $1
)
WHERE id = $1;

-- name: GetVisitorVotesForMessages :many
SELECT message_id, vote
FROM guestbook_votes
WHERE visitor_id = $1 AND message_id = ANY($2::uuid[]);

-- name: ListApprovedComments :many
SELECT id, post_id, parent_id, author_name, author_website, content,
       created_at
FROM post_comments
WHERE post_id = $1 AND status = 'approved'
ORDER BY created_at ASC
LIMIT $2 OFFSET $3;

-- name: CountApprovedComments :one
SELECT count(*) FROM post_comments WHERE post_id = $1 AND status = 'approved';

-- name: CreatePostComment :one
INSERT INTO post_comments (post_id, parent_id, author_name, author_email, author_website, content, ip_hash, ua_hash)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, created_at;

-- name: ListAdminComments :many
SELECT c.id, c.post_id, c.parent_id, c.author_name, c.author_email,
       c.content, c.status, c.ip_hash, c.created_at,
       p.title AS post_title
FROM post_comments c
LEFT JOIN posts p ON p.id = c.post_id
WHERE (sqlc.narg('status')::moderation_status IS NULL OR c.status = sqlc.narg('status')::moderation_status)
ORDER BY c.created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountAdminComments :one
SELECT count(*)
FROM post_comments
WHERE (sqlc.narg('status')::moderation_status IS NULL OR status = sqlc.narg('status')::moderation_status);

-- name: ApproveComment :exec
UPDATE post_comments
SET status = 'approved', approved_at = now(), reviewed_by = $2
WHERE id = $1;

-- name: RejectComment :exec
UPDATE post_comments
SET status = 'rejected', reviewed_by = $2
WHERE id = $1;

-- name: DeleteComment :exec
DELETE FROM post_comments WHERE id = $1;

-- name: CountPendingComments :one
SELECT count(*) FROM post_comments WHERE status = 'pending';

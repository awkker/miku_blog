-- name: ListApprovedFriendLinks :many
SELECT id, name, description, url, avatar_url, health_status
FROM friend_links
WHERE status = 'approved'
ORDER BY sort_order ASC, approved_at DESC;

-- name: ListAdminFriendLinks :many
SELECT id, name, description, url, domain, avatar_url,
       status, health_status, sort_order,
       created_at, approved_at, last_checked_at
FROM friend_links
ORDER BY sort_order ASC, created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountAdminFriendLinks :one
SELECT count(*) FROM friend_links;

-- name: GetFriendLinkByID :one
SELECT id, name, description, url, domain, avatar_url,
       status, health_status, sort_order,
       created_at, approved_at, reviewed_by, last_checked_at
FROM friend_links
WHERE id = $1;

-- name: CreateFriendLink :one
INSERT INTO friend_links (name, description, url, domain, avatar_url, status, sort_order, reviewed_by)
VALUES ($1, $2, $3, $4, $5, 'approved', $6, $7)
RETURNING id, created_at;

-- name: UpdateFriendLink :exec
UPDATE friend_links
SET name = $2, description = $3, url = $4, domain = $5,
    avatar_url = $6, sort_order = $7
WHERE id = $1;

-- name: DeleteFriendLink :exec
DELETE FROM friend_links WHERE id = $1;

-- name: ApproveFriendLink :exec
UPDATE friend_links
SET status = 'approved', approved_at = now(), reviewed_by = $2
WHERE id = $1;

-- name: RejectFriendLink :exec
UPDATE friend_links
SET status = 'rejected', reviewed_by = $2
WHERE id = $1;

-- name: UpdateFriendLinkHealth :exec
UPDATE friend_links
SET health_status = $2, last_checked_at = now()
WHERE id = $1;

-- name: ListFriendLinksForHealthCheck :many
SELECT id, url FROM friend_links
WHERE status = 'approved'
  AND (last_checked_at IS NULL OR last_checked_at < now() - interval '1 hour');

-- name: CreateFriendLinkHealthLog :exec
INSERT INTO friend_link_health_logs (friend_link_id, http_status, latency_ms, result)
VALUES ($1, $2, $3, $4);

-- name: CountApprovedFriendLinks :one
SELECT count(*) FROM friend_links WHERE status = 'approved';

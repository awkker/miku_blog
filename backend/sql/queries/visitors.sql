-- name: GetVisitorByIPUA :one
SELECT id, first_seen_at, last_seen_at, ip_hash, ua_hash
FROM visitors
WHERE ip_hash = $1 AND ua_hash = $2
LIMIT 1;

-- name: CreateVisitor :one
INSERT INTO visitors (ip_hash, ua_hash)
VALUES ($1, $2)
RETURNING id, first_seen_at;

-- name: TouchVisitor :exec
UPDATE visitors SET last_seen_at = now() WHERE id = $1;

-- name: GetVisitorByID :one
SELECT id, first_seen_at, last_seen_at, ip_hash, ua_hash
FROM visitors
WHERE id = $1;

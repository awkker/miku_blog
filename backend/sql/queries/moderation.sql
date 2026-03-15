-- name: ListSensitiveWords :many
SELECT word FROM sensitive_words ORDER BY word;

-- name: CreateAuditLog :exec
INSERT INTO audit_logs (admin_id, action, target_type, target_id, detail, ip)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: ListAuditLogs :many
SELECT a.id, a.action, a.target_type, a.target_id, a.detail, a.ip, a.created_at,
       u.username AS admin_username
FROM audit_logs a
LEFT JOIN admin_users u ON u.id = a.admin_id
ORDER BY a.created_at DESC
LIMIT $1 OFFSET $2;

-- name: CheckBlockedIP :one
SELECT count(*) FROM blocked_ips
WHERE ip_hash = $1 AND (expires_at IS NULL OR expires_at > now());

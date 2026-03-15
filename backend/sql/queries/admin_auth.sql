-- name: GetAdminByUsername :one
SELECT id, username, email, password_hash, role, status, last_login_at, created_at, updated_at
FROM admin_users
WHERE username = $1 AND status = 'active';

-- name: GetAdminByID :one
SELECT id, username, email, role, status, last_login_at, created_at, updated_at
FROM admin_users
WHERE id = $1;

-- name: UpdateAdminLastLogin :exec
UPDATE admin_users SET last_login_at = now() WHERE id = $1;

-- name: CreateRefreshToken :one
INSERT INTO admin_refresh_tokens (admin_user_id, token_hash, expires_at)
VALUES ($1, $2, $3)
RETURNING id, created_at;

-- name: GetRefreshTokenByHash :one
SELECT id, admin_user_id, token_hash, expires_at, revoked_at, created_at, last_used_at
FROM admin_refresh_tokens
WHERE token_hash = $1 AND revoked_at IS NULL AND expires_at > now();

-- name: TouchRefreshToken :exec
UPDATE admin_refresh_tokens SET last_used_at = now() WHERE id = $1;

-- name: RevokeRefreshToken :exec
UPDATE admin_refresh_tokens SET revoked_at = now() WHERE id = $1;

-- name: RevokeAllUserTokens :exec
UPDATE admin_refresh_tokens SET revoked_at = now()
WHERE admin_user_id = $1 AND revoked_at IS NULL;

-- name: CleanExpiredTokens :exec
DELETE FROM admin_refresh_tokens WHERE expires_at < now();

-- name: CreateAdminUser :one
INSERT INTO admin_users (username, email, password_hash, role)
VALUES ($1, $2, $3, $4)
RETURNING id, created_at;

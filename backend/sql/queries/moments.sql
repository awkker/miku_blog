-- name: ListMoments :many
SELECT id, author_name, content, image_urls, like_count, repost_count,
       comment_count, publish_status, published_at, scheduled_at, created_at
FROM moments
WHERE status = 'approved'
  AND publish_status = 'published'
  AND published_at IS NOT NULL
  AND published_at <= now()
ORDER BY published_at DESC, created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountMoments :one
SELECT count(*)
FROM moments
WHERE status = 'approved'
  AND publish_status = 'published'
  AND published_at IS NOT NULL
  AND published_at <= now();

-- name: ListLatestMoments :many
SELECT id, author_name, content, image_urls, created_at
FROM moments
WHERE status = 'approved'
  AND publish_status = 'published'
  AND published_at IS NOT NULL
  AND published_at <= now()
ORDER BY published_at DESC, created_at DESC
LIMIT $1;

-- name: GetMomentByID :one
SELECT id, author_name, content, image_urls, status, like_count,
       repost_count, comment_count, publish_status, published_at, scheduled_at, created_at
FROM moments
WHERE id = $1;

-- name: CreateMoment :one
INSERT INTO moments (author_name, content, image_urls, ip_hash, ua_hash, publish_status, published_at, scheduled_at)
VALUES (
    $1, $2, $3, $4, $5, $6,
    CASE WHEN $6 = 'published'::moment_publish_status THEN now() ELSE NULL END,
    CASE WHEN $6 = 'scheduled'::moment_publish_status THEN $7 ELSE NULL END
)
RETURNING id, created_at;

-- name: UpdateMoment :exec
UPDATE moments
SET author_name = $2, content = $3, image_urls = $4, publish_status = $5,
    published_at = CASE WHEN $5 = 'published'::moment_publish_status THEN COALESCE(published_at, now()) ELSE NULL END,
    scheduled_at = CASE WHEN $5 = 'scheduled'::moment_publish_status THEN $6 ELSE NULL END
WHERE id = $1;

-- name: CheckMomentLike :one
SELECT count(*) FROM moment_likes WHERE moment_id = $1 AND visitor_id = $2;

-- name: CreateMomentLike :exec
INSERT INTO moment_likes (moment_id, visitor_id) VALUES ($1, $2) ON CONFLICT DO NOTHING;

-- name: DeleteMomentLike :exec
DELETE FROM moment_likes WHERE moment_id = $1 AND visitor_id = $2;

-- name: IncrementMomentLikeCount :exec
UPDATE moments SET like_count = like_count + 1 WHERE id = $1;

-- name: DecrementMomentLikeCount :exec
UPDATE moments SET like_count = GREATEST(like_count - 1, 0) WHERE id = $1;

-- name: CheckMomentRepost :one
SELECT count(*) FROM moment_reposts WHERE moment_id = $1 AND visitor_id = $2;

-- name: CreateMomentRepost :exec
INSERT INTO moment_reposts (moment_id, visitor_id) VALUES ($1, $2) ON CONFLICT DO NOTHING;

-- name: DeleteMomentRepost :exec
DELETE FROM moment_reposts WHERE moment_id = $1 AND visitor_id = $2;

-- name: IncrementMomentRepostCount :exec
UPDATE moments SET repost_count = repost_count + 1 WHERE id = $1;

-- name: DecrementMomentRepostCount :exec
UPDATE moments SET repost_count = GREATEST(repost_count - 1, 0) WHERE id = $1;

-- name: ListMomentComments :many
SELECT id, moment_id, author_name, content, like_count, created_at
FROM moment_comments
WHERE moment_id = $1 AND status = 'approved'
ORDER BY created_at ASC
LIMIT $2 OFFSET $3;

-- name: CountMomentComments :one
SELECT count(*) FROM moment_comments WHERE moment_id = $1 AND status = 'approved';

-- name: CreateMomentComment :one
INSERT INTO moment_comments (moment_id, author_name, content, ip_hash, ua_hash)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, created_at;

-- name: IncrementMomentCommentCount :exec
UPDATE moments SET comment_count = comment_count + 1 WHERE id = $1;

-- name: CheckMomentCommentLike :one
SELECT count(*) FROM moment_comment_likes WHERE comment_id = $1 AND visitor_id = $2;

-- name: CreateMomentCommentLike :exec
INSERT INTO moment_comment_likes (comment_id, visitor_id) VALUES ($1, $2) ON CONFLICT DO NOTHING;

-- name: DeleteMomentCommentLike :exec
DELETE FROM moment_comment_likes WHERE comment_id = $1 AND visitor_id = $2;

-- name: IncrementMomentCommentLikeCount :exec
UPDATE moment_comments SET like_count = like_count + 1 WHERE id = $1;

-- name: DecrementMomentCommentLikeCount :exec
UPDATE moment_comments SET like_count = GREATEST(like_count - 1, 0) WHERE id = $1;

-- name: ListAdminMoments :many
SELECT id, author_name, content, image_urls, status, like_count,
       repost_count, comment_count, publish_status, published_at, scheduled_at, created_at
FROM moments
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountAdminMoments :one
SELECT count(*) FROM moments;

-- name: HideMoment :exec
UPDATE moments SET status = 'hidden', reviewed_by = $2 WHERE id = $1;

-- name: HideMomentComment :exec
UPDATE moment_comments SET status = 'hidden', reviewed_by = $2 WHERE id = $1;

-- name: GetVisitorMomentLikes :many
SELECT moment_id FROM moment_likes
WHERE visitor_id = $1 AND moment_id = ANY($2::uuid[]);

-- name: GetVisitorMomentReposts :many
SELECT moment_id FROM moment_reposts
WHERE visitor_id = $1 AND moment_id = ANY($2::uuid[]);

-- name: GetVisitorMomentCommentLikes :many
SELECT comment_id FROM moment_comment_likes
WHERE visitor_id = $1 AND comment_id = ANY($2::uuid[]);

-- name: PublishMoment :exec
UPDATE moments
SET publish_status = 'published', published_at = now(), scheduled_at = NULL
WHERE id = $1;

-- name: ScheduleMoment :exec
UPDATE moments
SET publish_status = 'scheduled', scheduled_at = $2, published_at = NULL
WHERE id = $1;

-- name: UnpublishMoment :exec
UPDATE moments
SET publish_status = 'draft', scheduled_at = NULL, published_at = NULL
WHERE id = $1;

-- name: ListScheduledMomentsDue :many
SELECT id
FROM moments
WHERE status = 'approved' AND publish_status = 'scheduled' AND scheduled_at <= now();

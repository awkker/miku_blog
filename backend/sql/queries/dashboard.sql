-- name: GetTotalPostCount :one
SELECT count(*) FROM posts WHERE status != 'draft';

-- name: GetTotalLikeCount :one
SELECT coalesce(sum(like_count), 0)::bigint FROM posts;

-- name: CountDraftPosts :one
SELECT count(*) FROM posts WHERE status = 'draft';

-- name: GetDailyViewTrend :many
SELECT day, coalesce(sum(pv), 0)::bigint AS pv, coalesce(sum(uv), 0)::bigint AS uv
FROM post_view_daily
WHERE day >= $1
GROUP BY day
ORDER BY day ASC;

-- name: GetDailyCommentTrend :many
SELECT created_at::date AS day, count(*) AS total
FROM post_comments
WHERE created_at >= $1
GROUP BY created_at::date
ORDER BY day ASC;

-- name: GetDailyLikeTrend :many
SELECT created_at::date AS day, count(*) AS total
FROM post_likes
WHERE created_at >= $1
GROUP BY created_at::date
ORDER BY day ASC;

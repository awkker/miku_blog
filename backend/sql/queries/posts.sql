-- name: ListPublishedPosts :many
SELECT id, slug, title, excerpt, hero_image_url, category, status,
       published_at, view_count, like_count, comment_count, created_at
FROM posts
WHERE status = 'published'
  AND published_at IS NOT NULL
  AND published_at <= now()
ORDER BY published_at DESC
LIMIT $1 OFFSET $2;

-- name: CountPublishedPosts :one
SELECT count(*)
FROM posts
WHERE status = 'published'
  AND published_at IS NOT NULL
  AND published_at <= now();

-- name: GetPostBySlug :one
SELECT id, slug, title, excerpt, content_markdown, hero_image_url, category,
       status, published_at, scheduled_at, view_count, like_count, comment_count,
       created_by, updated_by, created_at, updated_at
FROM posts
WHERE slug = $1
  AND status = 'published'
  AND published_at IS NOT NULL
  AND published_at <= now();

-- name: SearchPosts :many
SELECT id, slug, title, excerpt, hero_image_url, category,
       published_at, view_count, like_count, comment_count,
       ts_rank(search_vector, websearch_to_tsquery('simple', $1)) AS rank
FROM posts
WHERE status = 'published'
  AND published_at IS NOT NULL
  AND published_at <= now()
  AND search_vector @@ websearch_to_tsquery('simple', $1)
ORDER BY rank DESC
LIMIT $2 OFFSET $3;

-- name: ListPostsByCategory :many
SELECT id, slug, title, excerpt, hero_image_url, category,
       published_at, view_count, like_count, comment_count, created_at
FROM posts
WHERE status = 'published'
  AND published_at IS NOT NULL
  AND published_at <= now()
  AND category = $1
ORDER BY published_at DESC
LIMIT $2 OFFSET $3;

-- name: ListHotPosts :many
SELECT id, slug, title, excerpt, hero_image_url, category,
       published_at, view_count, like_count, comment_count, created_at
FROM posts
WHERE status = 'published'
  AND published_at IS NOT NULL
  AND published_at <= now()
ORDER BY (view_count + like_count * 3 + comment_count * 5) DESC
LIMIT $1;

-- name: GetPostByID :one
SELECT id, slug, title, excerpt, content_markdown, hero_image_url, category,
       status, published_at, scheduled_at, view_count, like_count, comment_count,
       created_by, updated_by, created_at, updated_at
FROM posts
WHERE id = $1;

-- name: ListAdminPosts :many
SELECT id, slug, title, category, status,
       published_at, scheduled_at, view_count, like_count, comment_count,
       created_at, updated_at
FROM posts
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountAdminPosts :one
SELECT count(*) FROM posts;

-- name: CreatePost :one
INSERT INTO posts (
    slug, title, excerpt, content_markdown, hero_image_url, category, status,
    published_at, scheduled_at, created_by, updated_by
)
VALUES (
    $1, $2, $3, $4, $5, $6, $7,
    CASE WHEN $7 = 'published'::post_status THEN now() ELSE NULL END,
    CASE WHEN $7 = 'scheduled'::post_status THEN $9 ELSE NULL END,
    $8, $8
)
RETURNING id, created_at;

-- name: UpdatePost :exec
UPDATE posts
SET slug = $2, title = $3, excerpt = $4, content_markdown = $5,
    hero_image_url = $6, category = $7, updated_by = $8, updated_at = now()
WHERE id = $1;

-- name: PublishPost :exec
UPDATE posts
SET status = 'published', published_at = now(), scheduled_at = NULL, updated_by = $2
WHERE id = $1;

-- name: SchedulePost :exec
UPDATE posts
SET status = 'scheduled', scheduled_at = $2, published_at = NULL, updated_by = $3
WHERE id = $1;

-- name: UnpublishPost :exec
UPDATE posts
SET status = 'draft', scheduled_at = NULL, updated_by = $2
WHERE id = $1;

-- name: DeletePost :exec
DELETE FROM posts WHERE id = $1;

-- name: IncrementPostViewCount :exec
UPDATE posts SET view_count = view_count + 1 WHERE id = $1;

-- name: IncrementPostLikeCount :exec
UPDATE posts SET like_count = like_count + 1 WHERE id = $1;

-- name: DecrementPostLikeCount :exec
UPDATE posts SET like_count = GREATEST(like_count - 1, 0) WHERE id = $1;

-- name: IncrementPostCommentCount :exec
UPDATE posts SET comment_count = comment_count + 1 WHERE id = $1;

-- name: DecrementPostCommentCount :exec
UPDATE posts SET comment_count = GREATEST(comment_count - 1, 0) WHERE id = $1;

-- name: ListScheduledPostsDue :many
SELECT id, slug, title FROM posts
WHERE status = 'scheduled' AND scheduled_at <= now();

-- name: GetPostTagNames :many
SELECT t.name, t.slug
FROM tags t
JOIN post_tags pt ON pt.tag_id = t.id
WHERE pt.post_id = $1;

-- name: UpsertTag :one
INSERT INTO tags (name, slug) VALUES ($1, $2)
ON CONFLICT (slug) DO UPDATE SET name = EXCLUDED.name
RETURNING id;

-- name: SetPostTags :exec
DELETE FROM post_tags WHERE post_id = $1;

-- name: AddPostTag :exec
INSERT INTO post_tags (post_id, tag_id) VALUES ($1, $2) ON CONFLICT DO NOTHING;

-- name: CreatePostRevision :exec
INSERT INTO post_revisions (post_id, title, excerpt, content_markdown, hero_image_url, category, editor_id)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: ListPostRevisions :many
SELECT id, post_id, title, excerpt, hero_image_url, category, editor_id, created_at
FROM post_revisions
WHERE post_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetPostRevision :one
SELECT id, post_id, title, excerpt, content_markdown, hero_image_url, category, editor_id, created_at
FROM post_revisions
WHERE id = $1;

-- name: CheckPostLike :one
SELECT count(*) FROM post_likes WHERE post_id = $1 AND visitor_id = $2;

-- name: CreatePostLike :exec
INSERT INTO post_likes (post_id, visitor_id) VALUES ($1, $2) ON CONFLICT DO NOTHING;

-- name: DeletePostLike :exec
DELETE FROM post_likes WHERE post_id = $1 AND visitor_id = $2;

-- name: UpsertPostViewDaily :exec
INSERT INTO post_view_daily (post_id, day, pv, uv)
VALUES ($1, $2, 1, 1)
ON CONFLICT (post_id, day) DO UPDATE
SET pv = post_view_daily.pv + 1,
    uv = post_view_daily.uv + $3;

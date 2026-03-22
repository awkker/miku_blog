-- name: ExportBackupPayload :one
SELECT jsonb_build_object(
    'admin_users', COALESCE((SELECT jsonb_agg(to_jsonb(t) ORDER BY t.created_at) FROM admin_users t), '[]'::jsonb),
    'admin_refresh_tokens', COALESCE((SELECT jsonb_agg(to_jsonb(t) ORDER BY t.created_at) FROM admin_refresh_tokens t), '[]'::jsonb),
    'visitors', COALESCE((SELECT jsonb_agg(to_jsonb(t) ORDER BY t.first_seen_at) FROM visitors t), '[]'::jsonb),

    'posts', COALESCE((SELECT jsonb_agg(to_jsonb(t) - 'search_vector' ORDER BY t.created_at) FROM posts t), '[]'::jsonb),
    'tags', COALESCE((SELECT jsonb_agg(to_jsonb(t) ORDER BY t.slug) FROM tags t), '[]'::jsonb),
    'post_tags', COALESCE((SELECT jsonb_agg(to_jsonb(t) ORDER BY t.post_id, t.tag_id) FROM post_tags t), '[]'::jsonb),
    'post_revisions', COALESCE((SELECT jsonb_agg(to_jsonb(t) ORDER BY t.created_at) FROM post_revisions t), '[]'::jsonb),
    'post_likes', COALESCE((SELECT jsonb_agg(to_jsonb(t) ORDER BY t.created_at) FROM post_likes t), '[]'::jsonb),
    'post_view_daily', COALESCE((SELECT jsonb_agg(to_jsonb(t) ORDER BY t.day) FROM post_view_daily t), '[]'::jsonb),
    'post_comments', COALESCE((SELECT jsonb_agg(to_jsonb(t) ORDER BY t.created_at) FROM post_comments t), '[]'::jsonb),

    'guestbook_messages', COALESCE((SELECT jsonb_agg(to_jsonb(t) ORDER BY t.created_at) FROM guestbook_messages t), '[]'::jsonb),
    'guestbook_votes', COALESCE((SELECT jsonb_agg(to_jsonb(t) ORDER BY t.created_at) FROM guestbook_votes t), '[]'::jsonb),

    'moments', COALESCE((SELECT jsonb_agg(to_jsonb(t) ORDER BY t.created_at) FROM moments t), '[]'::jsonb),
    'moment_likes', COALESCE((SELECT jsonb_agg(to_jsonb(t) ORDER BY t.created_at) FROM moment_likes t), '[]'::jsonb),
    'moment_reposts', COALESCE((SELECT jsonb_agg(to_jsonb(t) ORDER BY t.created_at) FROM moment_reposts t), '[]'::jsonb),
    'moment_comments', COALESCE((SELECT jsonb_agg(to_jsonb(t) ORDER BY t.created_at) FROM moment_comments t), '[]'::jsonb),
    'moment_comment_likes', COALESCE((SELECT jsonb_agg(to_jsonb(t) ORDER BY t.created_at) FROM moment_comment_likes t), '[]'::jsonb),

    'friend_links', COALESCE((SELECT jsonb_agg(to_jsonb(t) ORDER BY t.created_at) FROM friend_links t), '[]'::jsonb),
    'friend_link_health_logs', COALESCE((SELECT jsonb_agg(to_jsonb(t) ORDER BY t.checked_at) FROM friend_link_health_logs t), '[]'::jsonb),

    'sensitive_words', COALESCE((SELECT jsonb_agg(to_jsonb(t) ORDER BY t.created_at) FROM sensitive_words t), '[]'::jsonb),
    'blocked_ips', COALESCE((SELECT jsonb_agg(to_jsonb(t) ORDER BY t.blocked_at) FROM blocked_ips t), '[]'::jsonb),
    'audit_logs', COALESCE((
        SELECT jsonb_agg(to_jsonb(t) ORDER BY t.created_at)
        FROM (
            SELECT
                id,
                admin_id,
                action,
                target_type,
                target_id,
                COALESCE(NULLIF(detail, 'null'::jsonb), '{}'::jsonb) AS detail,
                ip,
                created_at
            FROM audit_logs
        ) t
    ), '[]'::jsonb),

    'analytics_sessions', COALESCE((SELECT jsonb_agg(to_jsonb(t) ORDER BY t.started_at) FROM analytics_sessions t), '[]'::jsonb),
    'analytics_pageviews', COALESCE((SELECT jsonb_agg(to_jsonb(t) ORDER BY t.occurred_at) FROM analytics_pageviews t), '[]'::jsonb)
) AS payload;

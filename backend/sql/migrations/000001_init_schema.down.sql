-- Reverse of 000001_init_schema.up.sql

DROP TRIGGER IF EXISTS update_posts_search_vector ON posts;
DROP TRIGGER IF EXISTS set_updated_at_posts ON posts;
DROP TRIGGER IF EXISTS set_updated_at_admin_users ON admin_users;
DROP FUNCTION IF EXISTS trigger_posts_search_vector;
DROP FUNCTION IF EXISTS trigger_set_updated_at;

DROP TABLE IF EXISTS audit_logs CASCADE;
DROP TABLE IF EXISTS blocked_ips CASCADE;
DROP TABLE IF EXISTS sensitive_words CASCADE;
DROP TABLE IF EXISTS friend_link_health_logs CASCADE;
DROP TABLE IF EXISTS friend_links CASCADE;
DROP TABLE IF EXISTS moment_comment_likes CASCADE;
DROP TABLE IF EXISTS moment_comments CASCADE;
DROP TABLE IF EXISTS moment_reposts CASCADE;
DROP TABLE IF EXISTS moment_likes CASCADE;
DROP TABLE IF EXISTS moments CASCADE;
DROP TABLE IF EXISTS guestbook_votes CASCADE;
DROP TABLE IF EXISTS guestbook_messages CASCADE;
DROP TABLE IF EXISTS post_comments CASCADE;
DROP TABLE IF EXISTS post_view_daily CASCADE;
DROP TABLE IF EXISTS post_likes CASCADE;
DROP TABLE IF EXISTS post_revisions CASCADE;
DROP TABLE IF EXISTS post_tags CASCADE;
DROP TABLE IF EXISTS tags CASCADE;
DROP TABLE IF EXISTS posts CASCADE;
DROP TABLE IF EXISTS visitors CASCADE;
DROP TABLE IF EXISTS admin_refresh_tokens CASCADE;
DROP TABLE IF EXISTS admin_users CASCADE;

DROP TYPE IF EXISTS vote_type;
DROP TYPE IF EXISTS health_status;
DROP TYPE IF EXISTS friend_link_status;
DROP TYPE IF EXISTS moderation_status;
DROP TYPE IF EXISTS post_status;

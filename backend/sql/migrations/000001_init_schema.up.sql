-- NanaMiku Blog: Initial Schema
-- All P0 tables as defined in PLAN.md Section 9

-- Extensions
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ========================
-- Enums
-- ========================

CREATE TYPE post_status AS ENUM ('draft', 'published', 'scheduled', 'archived');
CREATE TYPE moderation_status AS ENUM ('pending', 'approved', 'rejected', 'hidden');
CREATE TYPE friend_link_status AS ENUM ('pending', 'approved', 'rejected');
CREATE TYPE health_status AS ENUM ('ok', 'down', 'unknown');
CREATE TYPE vote_type AS ENUM ('up', 'down');

-- ========================
-- Auth & Visitors
-- ========================

CREATE TABLE admin_users (
    id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    username        text NOT NULL UNIQUE,
    email           text NOT NULL UNIQUE,
    password_hash   text NOT NULL,
    role            text NOT NULL DEFAULT 'admin',
    status          text NOT NULL DEFAULT 'active',
    last_login_at   timestamptz,
    created_at      timestamptz NOT NULL DEFAULT now(),
    updated_at      timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE admin_refresh_tokens (
    id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    admin_user_id   uuid NOT NULL REFERENCES admin_users(id) ON DELETE CASCADE,
    token_hash      text NOT NULL UNIQUE,
    expires_at      timestamptz NOT NULL,
    revoked_at      timestamptz,
    created_at      timestamptz NOT NULL DEFAULT now(),
    last_used_at    timestamptz
);

CREATE INDEX idx_admin_refresh_tokens_user ON admin_refresh_tokens(admin_user_id);

CREATE TABLE visitors (
    id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    first_seen_at   timestamptz NOT NULL DEFAULT now(),
    last_seen_at    timestamptz NOT NULL DEFAULT now(),
    ip_hash         text NOT NULL,
    ua_hash         text NOT NULL
);

CREATE INDEX idx_visitors_ip_ua ON visitors(ip_hash, ua_hash);

-- ========================
-- Posts
-- ========================

CREATE TABLE posts (
    id                  uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    slug                text NOT NULL UNIQUE,
    title               text NOT NULL,
    excerpt             text NOT NULL DEFAULT '',
    content_markdown    text NOT NULL DEFAULT '',
    hero_image_url      text NOT NULL DEFAULT '',
    category            text NOT NULL DEFAULT '',
    status              post_status NOT NULL DEFAULT 'draft',
    published_at        timestamptz,
    scheduled_at        timestamptz,
    view_count          bigint NOT NULL DEFAULT 0,
    like_count          bigint NOT NULL DEFAULT 0,
    comment_count       bigint NOT NULL DEFAULT 0,
    created_by          uuid REFERENCES admin_users(id),
    updated_by          uuid REFERENCES admin_users(id),
    created_at          timestamptz NOT NULL DEFAULT now(),
    updated_at          timestamptz NOT NULL DEFAULT now(),
    search_vector       tsvector
);

CREATE INDEX idx_posts_status ON posts(status);
CREATE INDEX idx_posts_published_at ON posts(published_at DESC);
CREATE INDEX idx_posts_category ON posts(category);
CREATE INDEX idx_posts_search ON posts USING GIN(search_vector);

CREATE TABLE tags (
    id      uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name    text NOT NULL UNIQUE,
    slug    text NOT NULL UNIQUE
);

CREATE TABLE post_tags (
    post_id uuid NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    tag_id  uuid NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (post_id, tag_id)
);

CREATE INDEX idx_post_tags_tag ON post_tags(tag_id);

CREATE TABLE post_revisions (
    id                  uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    post_id             uuid NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    title               text NOT NULL,
    excerpt             text NOT NULL DEFAULT '',
    content_markdown    text NOT NULL DEFAULT '',
    hero_image_url      text NOT NULL DEFAULT '',
    category            text NOT NULL DEFAULT '',
    editor_id           uuid REFERENCES admin_users(id),
    created_at          timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_post_revisions_post ON post_revisions(post_id, created_at DESC);

-- ========================
-- Post Interactions
-- ========================

CREATE TABLE post_likes (
    post_id     uuid NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    visitor_id  uuid NOT NULL REFERENCES visitors(id) ON DELETE CASCADE,
    created_at  timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (post_id, visitor_id)
);

CREATE TABLE post_view_daily (
    post_id uuid NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    day     date NOT NULL,
    pv      bigint NOT NULL DEFAULT 0,
    uv      bigint NOT NULL DEFAULT 0,
    PRIMARY KEY (post_id, day)
);

CREATE TABLE post_comments (
    id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    post_id         uuid NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    parent_id       uuid REFERENCES post_comments(id) ON DELETE CASCADE,
    author_name     text NOT NULL,
    author_email    text NOT NULL DEFAULT '',
    author_website  text NOT NULL DEFAULT '',
    content         text NOT NULL,
    status          moderation_status NOT NULL DEFAULT 'pending',
    ip_hash         text NOT NULL DEFAULT '',
    ua_hash         text NOT NULL DEFAULT '',
    created_at      timestamptz NOT NULL DEFAULT now(),
    approved_at     timestamptz,
    reviewed_by     uuid REFERENCES admin_users(id)
);

CREATE INDEX idx_post_comments_post ON post_comments(post_id, created_at);
CREATE INDEX idx_post_comments_status ON post_comments(status);

-- ========================
-- Guestbook
-- ========================

CREATE TABLE guestbook_messages (
    id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    parent_id       uuid REFERENCES guestbook_messages(id) ON DELETE CASCADE,
    author_name     text NOT NULL,
    author_website  text NOT NULL DEFAULT '',
    content         text NOT NULL,
    status          moderation_status NOT NULL DEFAULT 'approved',
    vote_score      int NOT NULL DEFAULT 0,
    ip_hash         text NOT NULL DEFAULT '',
    ua_hash         text NOT NULL DEFAULT '',
    created_at      timestamptz NOT NULL DEFAULT now(),
    reviewed_by     uuid REFERENCES admin_users(id)
);

CREATE INDEX idx_guestbook_status ON guestbook_messages(status);
CREATE INDEX idx_guestbook_parent ON guestbook_messages(parent_id);

CREATE TABLE guestbook_votes (
    message_id  uuid NOT NULL REFERENCES guestbook_messages(id) ON DELETE CASCADE,
    visitor_id  uuid NOT NULL REFERENCES visitors(id) ON DELETE CASCADE,
    vote        vote_type NOT NULL,
    created_at  timestamptz NOT NULL DEFAULT now(),
    updated_at  timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (message_id, visitor_id)
);

-- ========================
-- Moments
-- ========================

CREATE TABLE moments (
    id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    author_name     text NOT NULL,
    content         text NOT NULL,
    image_urls      jsonb NOT NULL DEFAULT '[]'::jsonb,
    status          moderation_status NOT NULL DEFAULT 'approved',
    like_count      bigint NOT NULL DEFAULT 0,
    repost_count    bigint NOT NULL DEFAULT 0,
    comment_count   bigint NOT NULL DEFAULT 0,
    ip_hash         text NOT NULL DEFAULT '',
    ua_hash         text NOT NULL DEFAULT '',
    created_at      timestamptz NOT NULL DEFAULT now(),
    reviewed_by     uuid REFERENCES admin_users(id)
);

CREATE INDEX idx_moments_status ON moments(status);
CREATE INDEX idx_moments_created ON moments(created_at DESC);

CREATE TABLE moment_likes (
    moment_id   uuid NOT NULL REFERENCES moments(id) ON DELETE CASCADE,
    visitor_id  uuid NOT NULL REFERENCES visitors(id) ON DELETE CASCADE,
    created_at  timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (moment_id, visitor_id)
);

CREATE TABLE moment_reposts (
    moment_id   uuid NOT NULL REFERENCES moments(id) ON DELETE CASCADE,
    visitor_id  uuid NOT NULL REFERENCES visitors(id) ON DELETE CASCADE,
    created_at  timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (moment_id, visitor_id)
);

CREATE TABLE moment_comments (
    id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    moment_id       uuid NOT NULL REFERENCES moments(id) ON DELETE CASCADE,
    author_name     text NOT NULL,
    content         text NOT NULL,
    status          moderation_status NOT NULL DEFAULT 'approved',
    like_count      bigint NOT NULL DEFAULT 0,
    ip_hash         text NOT NULL DEFAULT '',
    ua_hash         text NOT NULL DEFAULT '',
    created_at      timestamptz NOT NULL DEFAULT now(),
    reviewed_by     uuid REFERENCES admin_users(id)
);

CREATE INDEX idx_moment_comments_moment ON moment_comments(moment_id, created_at);
CREATE INDEX idx_moment_comments_status ON moment_comments(status);

CREATE TABLE moment_comment_likes (
    comment_id  uuid NOT NULL REFERENCES moment_comments(id) ON DELETE CASCADE,
    visitor_id  uuid NOT NULL REFERENCES visitors(id) ON DELETE CASCADE,
    created_at  timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (comment_id, visitor_id)
);

-- ========================
-- Friends
-- ========================

CREATE TABLE friend_links (
    id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name            text NOT NULL,
    description     text NOT NULL DEFAULT '',
    url             text NOT NULL,
    domain          text NOT NULL DEFAULT '',
    avatar_url      text NOT NULL DEFAULT '',
    status          friend_link_status NOT NULL DEFAULT 'pending',
    health_status   health_status NOT NULL DEFAULT 'unknown',
    sort_order      int NOT NULL DEFAULT 0,
    created_at      timestamptz NOT NULL DEFAULT now(),
    approved_at     timestamptz,
    reviewed_by     uuid REFERENCES admin_users(id),
    last_checked_at timestamptz
);

CREATE INDEX idx_friend_links_status ON friend_links(status);

CREATE TABLE friend_link_health_logs (
    id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    friend_link_id  uuid NOT NULL REFERENCES friend_links(id) ON DELETE CASCADE,
    http_status     int NOT NULL DEFAULT 0,
    latency_ms      int NOT NULL DEFAULT 0,
    result          health_status NOT NULL DEFAULT 'unknown',
    checked_at      timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_fl_health_logs_link ON friend_link_health_logs(friend_link_id, checked_at DESC);

-- ========================
-- Moderation & Audit
-- ========================

CREATE TABLE sensitive_words (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    word        text NOT NULL UNIQUE,
    category    text NOT NULL DEFAULT 'general',
    created_at  timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE blocked_ips (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    ip_hash     text NOT NULL UNIQUE,
    reason      text NOT NULL DEFAULT '',
    blocked_at  timestamptz NOT NULL DEFAULT now(),
    expires_at  timestamptz
);

CREATE TABLE audit_logs (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    admin_id    uuid REFERENCES admin_users(id),
    action      text NOT NULL,
    target_type text NOT NULL DEFAULT '',
    target_id   text NOT NULL DEFAULT '',
    detail      jsonb NOT NULL DEFAULT '{}'::jsonb,
    ip          text NOT NULL DEFAULT '',
    created_at  timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_audit_logs_admin ON audit_logs(admin_id, created_at DESC);
CREATE INDEX idx_audit_logs_action ON audit_logs(action, created_at DESC);

-- ========================
-- Triggers
-- ========================

CREATE OR REPLACE FUNCTION trigger_set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_updated_at_admin_users
    BEFORE UPDATE ON admin_users
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

CREATE TRIGGER set_updated_at_posts
    BEFORE UPDATE ON posts
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

CREATE OR REPLACE FUNCTION trigger_posts_search_vector()
RETURNS TRIGGER AS $$
BEGIN
    NEW.search_vector =
        setweight(to_tsvector('simple', coalesce(NEW.title, '')), 'A') ||
        setweight(to_tsvector('simple', coalesce(NEW.excerpt, '')), 'B') ||
        setweight(to_tsvector('simple', coalesce(NEW.content_markdown, '')), 'C');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_posts_search_vector
    BEFORE INSERT OR UPDATE OF title, excerpt, content_markdown ON posts
    FOR EACH ROW EXECUTE FUNCTION trigger_posts_search_vector();

-- Analytics sessions and pageviews for admin dashboard

CREATE TABLE analytics_sessions (
    id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    session_key     text NOT NULL UNIQUE,
    visitor_id      uuid REFERENCES visitors(id) ON DELETE SET NULL,
    started_at      timestamptz NOT NULL,
    last_seen_at    timestamptz NOT NULL,
    entry_path      text NOT NULL DEFAULT '/',
    exit_path       text NOT NULL DEFAULT '/',
    referrer        text NOT NULL DEFAULT '',
    referrer_host   text NOT NULL DEFAULT '',
    channel         text NOT NULL DEFAULT 'direct',
    browser         text NOT NULL DEFAULT 'Unknown',
    os              text NOT NULL DEFAULT 'Unknown',
    device          text NOT NULL DEFAULT 'Desktop',
    country_code    text NOT NULL DEFAULT 'ZZ',
    region          text NOT NULL DEFAULT 'Unknown',
    city            text NOT NULL DEFAULT 'Unknown',
    timezone        text NOT NULL DEFAULT '',
    language        text NOT NULL DEFAULT '',
    pageviews       int NOT NULL DEFAULT 1,
    created_at      timestamptz NOT NULL DEFAULT now(),
    updated_at      timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_analytics_sessions_started_at ON analytics_sessions(started_at DESC);
CREATE INDEX idx_analytics_sessions_visitor ON analytics_sessions(visitor_id, started_at DESC);
CREATE INDEX idx_analytics_sessions_channel ON analytics_sessions(channel);
CREATE INDEX idx_analytics_sessions_referrer_host ON analytics_sessions(referrer_host);
CREATE INDEX idx_analytics_sessions_country_code ON analytics_sessions(country_code);

CREATE TABLE analytics_pageviews (
    id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id      uuid NOT NULL REFERENCES analytics_sessions(id) ON DELETE CASCADE,
    visitor_id      uuid REFERENCES visitors(id) ON DELETE SET NULL,
    path            text NOT NULL,
    title           text NOT NULL DEFAULT '',
    referrer        text NOT NULL DEFAULT '',
    referrer_host   text NOT NULL DEFAULT '',
    channel         text NOT NULL DEFAULT 'direct',
    browser         text NOT NULL DEFAULT 'Unknown',
    os              text NOT NULL DEFAULT 'Unknown',
    device          text NOT NULL DEFAULT 'Desktop',
    country_code    text NOT NULL DEFAULT 'ZZ',
    region          text NOT NULL DEFAULT 'Unknown',
    city            text NOT NULL DEFAULT 'Unknown',
    timezone        text NOT NULL DEFAULT '',
    language        text NOT NULL DEFAULT '',
    occurred_at     timestamptz NOT NULL,
    created_at      timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_analytics_pageviews_occurred_at ON analytics_pageviews(occurred_at DESC);
CREATE INDEX idx_analytics_pageviews_path ON analytics_pageviews(path);
CREATE INDEX idx_analytics_pageviews_session ON analytics_pageviews(session_id, occurred_at DESC);
CREATE INDEX idx_analytics_pageviews_channel ON analytics_pageviews(channel);
CREATE INDEX idx_analytics_pageviews_referrer_host ON analytics_pageviews(referrer_host);
CREATE INDEX idx_analytics_pageviews_country_code ON analytics_pageviews(country_code);

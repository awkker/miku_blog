-- name: UpsertAnalyticsSession :one
INSERT INTO analytics_sessions (
    session_key,
    visitor_id,
    started_at,
    last_seen_at,
    entry_path,
    exit_path,
    referrer,
    referrer_host,
    channel,
    browser,
    os,
    device,
    country_code,
    region,
    city,
    timezone,
    language,
    pageviews
)
VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9,
    $10, $11, $12, $13, $14, $15, $16, $17, 1
)
ON CONFLICT (session_key) DO UPDATE
SET
    visitor_id = COALESCE(analytics_sessions.visitor_id, EXCLUDED.visitor_id),
    last_seen_at = GREATEST(analytics_sessions.last_seen_at, EXCLUDED.last_seen_at),
    exit_path = EXCLUDED.exit_path,
    referrer = CASE
        WHEN analytics_sessions.referrer = '' THEN EXCLUDED.referrer
        ELSE analytics_sessions.referrer
    END,
    referrer_host = CASE
        WHEN analytics_sessions.referrer_host = '' THEN EXCLUDED.referrer_host
        ELSE analytics_sessions.referrer_host
    END,
    channel = CASE
        WHEN analytics_sessions.channel IN ('', 'direct', 'unknown')
             AND EXCLUDED.channel NOT IN ('', 'direct', 'unknown') THEN EXCLUDED.channel
        ELSE analytics_sessions.channel
    END,
    browser = CASE
        WHEN analytics_sessions.browser = 'Unknown' AND EXCLUDED.browser != 'Unknown' THEN EXCLUDED.browser
        ELSE analytics_sessions.browser
    END,
    os = CASE
        WHEN analytics_sessions.os = 'Unknown' AND EXCLUDED.os != 'Unknown' THEN EXCLUDED.os
        ELSE analytics_sessions.os
    END,
    device = CASE
        WHEN analytics_sessions.device = 'Unknown' AND EXCLUDED.device != 'Unknown' THEN EXCLUDED.device
        ELSE analytics_sessions.device
    END,
    country_code = CASE
        WHEN analytics_sessions.country_code IN ('', 'ZZ') AND EXCLUDED.country_code NOT IN ('', 'ZZ') THEN EXCLUDED.country_code
        ELSE analytics_sessions.country_code
    END,
    region = CASE
        WHEN analytics_sessions.region = 'Unknown' AND EXCLUDED.region != 'Unknown' THEN EXCLUDED.region
        ELSE analytics_sessions.region
    END,
    city = CASE
        WHEN analytics_sessions.city = 'Unknown' AND EXCLUDED.city != 'Unknown' THEN EXCLUDED.city
        ELSE analytics_sessions.city
    END,
    timezone = CASE
        WHEN analytics_sessions.timezone = '' THEN EXCLUDED.timezone
        ELSE analytics_sessions.timezone
    END,
    language = CASE
        WHEN analytics_sessions.language = '' THEN EXCLUDED.language
        ELSE analytics_sessions.language
    END,
    pageviews = analytics_sessions.pageviews + 1,
    updated_at = now()
RETURNING id, started_at, last_seen_at, pageviews;

-- name: CreateAnalyticsPageview :exec
INSERT INTO analytics_pageviews (
    session_id,
    visitor_id,
    path,
    title,
    referrer,
    referrer_host,
    channel,
    browser,
    os,
    device,
    country_code,
    region,
    city,
    timezone,
    language,
    occurred_at
)
VALUES (
    $1, $2, $3, $4, $5, $6, $7,
    $8, $9, $10, $11, $12, $13, $14, $15, $16
);

-- name: CountAnalyticsVisitors :one
SELECT count(DISTINCT coalesce(visitor_id::text, session_key))::bigint AS total
FROM analytics_sessions
WHERE started_at >= $1 AND started_at < $2;

-- name: CountAnalyticsVisits :one
SELECT count(*)::bigint AS total
FROM analytics_sessions
WHERE started_at >= $1 AND started_at < $2;

-- name: CountAnalyticsViews :one
SELECT count(*)::bigint AS total
FROM analytics_pageviews
WHERE occurred_at >= $1 AND occurred_at < $2;

-- name: CountAnalyticsBouncedVisits :one
SELECT count(*)::bigint AS total
FROM analytics_sessions
WHERE started_at >= $1 AND started_at < $2
  AND pageviews <= 1;

-- name: GetAnalyticsAverageVisitDurationSeconds :one
SELECT coalesce(avg(EXTRACT(EPOCH FROM (last_seen_at - started_at))), 0)::double precision AS avg_seconds
FROM analytics_sessions
WHERE started_at >= $1 AND started_at < $2;

-- name: GetAnalyticsTrend :many
SELECT
  CASE
    WHEN $3::text = 'hour' THEN to_char(date_trunc('hour', occurred_at), 'YYYY-MM-DD HH24:00')
    ELSE to_char(date_trunc('day', occurred_at), 'YYYY-MM-DD')
  END AS bucket,
  count(DISTINCT coalesce(visitor_id::text, session_id::text))::bigint AS visitors,
  count(*)::bigint AS views
FROM analytics_pageviews
WHERE occurred_at >= $1 AND occurred_at < $2
GROUP BY 1
ORDER BY 1;

-- name: ListAnalyticsPageStats :many
SELECT
  path,
  count(DISTINCT coalesce(visitor_id::text, session_id::text))::bigint AS visitors,
  count(*)::bigint AS views
FROM analytics_pageviews
WHERE occurred_at >= $1 AND occurred_at < $2
GROUP BY path
ORDER BY visitors DESC, views DESC, path ASC
LIMIT $3;

-- name: ListAnalyticsEntryPaths :many
SELECT
  entry_path AS path,
  count(*)::bigint AS visits
FROM analytics_sessions
WHERE started_at >= $1 AND started_at < $2
GROUP BY entry_path
ORDER BY visits DESC, path ASC
LIMIT $3;

-- name: ListAnalyticsExitPaths :many
SELECT
  exit_path AS path,
  count(*)::bigint AS visits
FROM analytics_sessions
WHERE started_at >= $1 AND started_at < $2
GROUP BY exit_path
ORDER BY visits DESC, path ASC
LIMIT $3;

-- name: ListAnalyticsReferrers :many
SELECT
  CASE WHEN referrer_host = '' THEN 'direct' ELSE referrer_host END AS name,
  count(*)::bigint AS visitors
FROM analytics_sessions
WHERE started_at >= $1 AND started_at < $2
GROUP BY 1
ORDER BY visitors DESC, name ASC
LIMIT $3;

-- name: ListAnalyticsChannels :many
SELECT
  channel AS name,
  count(*)::bigint AS visitors
FROM analytics_sessions
WHERE started_at >= $1 AND started_at < $2
GROUP BY channel
ORDER BY visitors DESC, name ASC
LIMIT $3;

-- name: ListAnalyticsBrowsers :many
SELECT
  browser AS name,
  count(*)::bigint AS visitors
FROM analytics_sessions
WHERE started_at >= $1 AND started_at < $2
GROUP BY browser
ORDER BY visitors DESC, name ASC
LIMIT $3;

-- name: ListAnalyticsOperatingSystems :many
SELECT
  os AS name,
  count(*)::bigint AS visitors
FROM analytics_sessions
WHERE started_at >= $1 AND started_at < $2
GROUP BY os
ORDER BY visitors DESC, name ASC
LIMIT $3;

-- name: ListAnalyticsDevices :many
SELECT
  device AS name,
  count(*)::bigint AS visitors
FROM analytics_sessions
WHERE started_at >= $1 AND started_at < $2
GROUP BY device
ORDER BY visitors DESC, name ASC
LIMIT $3;

-- name: ListAnalyticsCountries :many
SELECT
  country_code,
  count(*)::bigint AS visitors
FROM analytics_sessions
WHERE started_at >= $1 AND started_at < $2
GROUP BY country_code
ORDER BY visitors DESC, country_code ASC
LIMIT $3;

-- name: ListAnalyticsRegions :many
SELECT
  region AS name,
  count(*)::bigint AS visitors
FROM analytics_sessions
WHERE started_at >= $1 AND started_at < $2
GROUP BY region
ORDER BY visitors DESC, name ASC
LIMIT $3;

-- name: ListAnalyticsCities :many
SELECT
  city AS name,
  count(*)::bigint AS visitors
FROM analytics_sessions
WHERE started_at >= $1 AND started_at < $2
GROUP BY city
ORDER BY visitors DESC, name ASC
LIMIT $3;

-- name: GetAnalyticsTrafficHeatmap :many
SELECT
  extract(dow FROM occurred_at)::int AS dow,
  extract(hour FROM occurred_at)::int AS hour,
  count(*)::bigint AS value
FROM analytics_pageviews
WHERE occurred_at >= $1 AND occurred_at < $2
GROUP BY dow, hour
ORDER BY dow, hour;

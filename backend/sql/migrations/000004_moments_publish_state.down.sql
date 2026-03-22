DROP INDEX IF EXISTS idx_moments_published_at;
DROP INDEX IF EXISTS idx_moments_scheduled_at;
DROP INDEX IF EXISTS idx_moments_publish_status;

ALTER TABLE moments
    DROP COLUMN IF EXISTS scheduled_at,
    DROP COLUMN IF EXISTS published_at,
    DROP COLUMN IF EXISTS publish_status;

DROP TYPE IF EXISTS moment_publish_status;

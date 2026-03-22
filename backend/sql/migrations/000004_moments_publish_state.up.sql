-- Add draft/scheduled publish lifecycle for moments

CREATE TYPE moment_publish_status AS ENUM ('draft', 'published', 'scheduled', 'archived');

ALTER TABLE moments
    ADD COLUMN publish_status moment_publish_status NOT NULL DEFAULT 'published',
    ADD COLUMN published_at timestamptz,
    ADD COLUMN scheduled_at timestamptz;

UPDATE moments
SET published_at = created_at
WHERE publish_status = 'published' AND published_at IS NULL;

CREATE INDEX idx_moments_publish_status ON moments(publish_status);
CREATE INDEX idx_moments_scheduled_at ON moments(scheduled_at);
CREATE INDEX idx_moments_published_at ON moments(published_at DESC);

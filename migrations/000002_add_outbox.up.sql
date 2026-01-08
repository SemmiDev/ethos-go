-- Outbox table for reliable event publishing
-- Implements the Outbox Pattern for atomic event publishing alongside domain changes
CREATE TABLE IF NOT EXISTS outbox (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_type VARCHAR(255) NOT NULL,
    aggregate_type VARCHAR(100) NOT NULL,
    aggregate_id VARCHAR(255) NOT NULL,
    payload JSONB NOT NULL,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    published_at TIMESTAMPTZ,
    published BOOLEAN NOT NULL DEFAULT FALSE,
    retry_count INT NOT NULL DEFAULT 0,
    last_error TEXT
);

-- Index for efficient polling of unpublished events
CREATE INDEX IF NOT EXISTS idx_outbox_unpublished ON outbox(published, created_at)
    WHERE published = FALSE;

-- Index for querying by aggregate
CREATE INDEX IF NOT EXISTS idx_outbox_aggregate ON outbox(aggregate_type, aggregate_id);

-- Comment
COMMENT ON TABLE outbox IS 'Outbox table for reliable event publishing using the Outbox Pattern';

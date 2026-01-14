-- Add indexes for pagination performance
-- This optimizes: WHERE status = 'active' AND deleted_at IS NULL ORDER BY created_at DESC

-- Composite index for status-based queries with sorting
CREATE INDEX IF NOT EXISTS idx_listings_status_created_at
    ON listings(status, created_at DESC)
    WHERE deleted_at IS NULL;

-- Index for count queries (can use index-only scan)
CREATE INDEX IF NOT EXISTS idx_listings_status_deleted_at
    ON listings(status)
    WHERE deleted_at IS NULL;
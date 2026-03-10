BEGIN;

-- Pagination & filtering
CREATE INDEX idx_listings_status_created_at
    ON listings (status, created_at)
    WHERE deleted_at IS NULL;

-- Host's listings
CREATE INDEX idx_listings_host
    ON listings (host_id, created_at)
    WHERE deleted_at IS NULL;

COMMIT;

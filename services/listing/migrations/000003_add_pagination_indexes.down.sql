-- Drop indexes for pagination performance
DROP INDEX IF EXISTS idx_listings_status_created_at;
DROP INDEX IF EXISTS idx_listings_status_deleted_at;


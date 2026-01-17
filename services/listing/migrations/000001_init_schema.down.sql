-- Drop in reverse order

DROP INDEX IF EXISTS idx_listings_host;
DROP INDEX IF EXISTS idx_listings_district;
DROP INDEX IF EXISTS idx_listings_province;
DROP INDEX IF EXISTS idx_listings_status_deleted_at;
DROP INDEX IF EXISTS idx_listings_status_created_at;

DROP TABLE IF EXISTS listings;
DROP TABLE IF EXISTS wards;
DROP TABLE IF EXISTS districts;
DROP TABLE IF EXISTS provinces;
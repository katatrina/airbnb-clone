BEGIN;

DROP INDEX idx_bookings_guest;
DROP INDEX idx_bookings_host;
DROP INDEX idx_bookings_listing;

COMMIT;

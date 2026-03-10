BEGIN;

-- Guests view their bookings: WHERE guest_id = ? AND deleted_at IS NULL
CREATE INDEX idx_bookings_guest
    ON bookings (guest_id, status)
    WHERE deleted_at IS NULL;

-- Hosts view their listing's bookings: WHERE host_id = ? AND deleted_at IS NULL
CREATE INDEX idx_bookings_host
    ON bookings (host_id, status)
    WHERE deleted_at IS NULL;

-- Query booking by listing: WHERE listing_id = ? (check availability, list bookings)
CREATE INDEX idx_bookings_listing
    ON bookings (listing_id, status)
    WHERE deleted_at IS NULL;

COMMIT;

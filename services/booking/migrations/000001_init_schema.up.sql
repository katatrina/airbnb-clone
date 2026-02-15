-- Required for exclusion constraint with UUID + daterange
CREATE EXTENSION IF NOT EXISTS btree_gist;

CREATE TABLE IF NOT EXISTS bookings
(
    -- Identity
    id              UUID PRIMARY KEY,
    listing_id      UUID        NOT NULL,
    guest_id        UUID        NOT NULL,
    host_id         UUID        NOT NULL,

    -- Stay period
    check_in_date   DATE        NOT NULL,
    check_out_date  DATE        NOT NULL,
    total_nights    INT         NOT NULL,

    -- Pricing snapshot
    price_per_night BIGINT      NOT NULL,
    total_price     BIGINT      NOT NULL,
    currency        TEXT        NOT NULL DEFAULT 'VND',

    -- Status
    status          TEXT        NOT NULL DEFAULT 'pending',

    -- Timestamps
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ,

    -- Validation constraints
    CONSTRAINT check_dates CHECK (check_out_date > check_in_date),
    CONSTRAINT check_nights CHECK (total_nights > 0),
    CONSTRAINT check_price CHECK (price_per_night > 0 AND total_price > 0)
    );

-- Prevent double booking at database level
ALTER TABLE bookings
    ADD CONSTRAINT no_overlapping_bookings EXCLUDE USING gist (
            listing_id WITH =,
            daterange(check_in_date, check_out_date) WITH &&
        )
        WHERE (status IN ('pending', 'confirmed') AND deleted_at IS NULL);

-- Indexes
CREATE INDEX idx_bookings_guest
    ON bookings (guest_id, status) WHERE deleted_at IS NULL;

CREATE INDEX idx_bookings_host
    ON bookings (host_id, status) WHERE deleted_at IS NULL;

CREATE INDEX idx_bookings_listing
    ON bookings (listing_id, status) WHERE deleted_at IS NULL;
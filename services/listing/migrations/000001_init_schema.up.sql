-- ============================================================================
-- Master Data Tables
-- ============================================================================

-- Provinces (Tỉnh/Thành phố)
CREATE TABLE IF NOT EXISTS provinces
(
    code       TEXT PRIMARY KEY,
    full_name  TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Districts (Quận/Huyện)
CREATE TABLE IF NOT EXISTS districts
(
    code          TEXT PRIMARY KEY,
    full_name     TEXT        NOT NULL,
    province_code TEXT        NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Wards (Phường/Xã/Thị trấn)
CREATE TABLE IF NOT EXISTS wards
(
    code          TEXT PRIMARY KEY,
    full_name     TEXT        NOT NULL,
    district_code TEXT        NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================================
-- Listings Table
-- ============================================================================

CREATE TABLE IF NOT EXISTS listings
(
    -- Identity
    id      UUID PRIMARY KEY,
    host_id UUID NOT NULL,

    -- Basic Info
    title           TEXT   NOT NULL,
    description     TEXT,
    price_per_night BIGINT NOT NULL,
    currency        TEXT   NOT NULL DEFAULT 'VND',

    -- Location (Denormalized - No FK)
    province_code TEXT NOT NULL,
    province_name TEXT NOT NULL,
    district_code TEXT NOT NULL,
    district_name TEXT NOT NULL,
    ward_code     TEXT NOT NULL,
    ward_name     TEXT NOT NULL,
    address_detail TEXT NOT NULL,

    -- Status
    status TEXT NOT NULL DEFAULT 'draft',

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- ============================================================================
-- Indexes
-- ============================================================================

-- Pagination & filtering
CREATE INDEX IF NOT EXISTS idx_listings_status_created_at
    ON listings (status, created_at DESC)
    WHERE deleted_at IS NULL;

-- Count queries
CREATE INDEX IF NOT EXISTS idx_listings_status_deleted_at
    ON listings (status)
    WHERE deleted_at IS NULL;

-- Location-based queries
CREATE INDEX IF NOT EXISTS idx_listings_province
    ON listings (province_code, status)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_listings_district
    ON listings (district_code, status)
    WHERE deleted_at IS NULL;

-- Host's listings
CREATE INDEX IF NOT EXISTS idx_listings_host
    ON listings (host_id, status)
    WHERE deleted_at IS NULL;
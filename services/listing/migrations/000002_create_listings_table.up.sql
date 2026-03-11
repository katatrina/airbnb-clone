CREATE TABLE listings
(
    -- Identity
    id      UUID PRIMARY KEY,
    host_id UUID NOT NULL,

    -- Basic Info
    title           TEXT   NOT NULL,
    description     TEXT   NOT NULL DEFAULT '',
    price_per_night BIGINT NOT NULL,
    currency        TEXT   NOT NULL DEFAULT 'VND',

    -- Location (Denormalized - No FK)
    province_code INTEGER NOT NULL,
    province_name TEXT    NOT NULL,
    district_code INTEGER NOT NULL,
    district_name TEXT    NOT NULL,
    ward_code     INTEGER NOT NULL,
    ward_name     TEXT    NOT NULL,
    address_detail TEXT NOT NULL,

    -- Status
    status TEXT NOT NULL DEFAULT 'draft',

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS listings
(
    id              UUID PRIMARY KEY,
    host_id         UUID        NOT NULL,
    title           TEXT        NOT NULL,
    description     TEXT,
    price_per_night BIGINT      NOT NULL,
    currency        TEXT        NOT NULL DEFAULT 'VND',
    province_code   TEXT        NOT NULL,
    province_name   TEXT        NOT NULL,
    ward_code       TEXT        NOT NULL,
    ward_name       TEXT        NOT NULL,
    address_detail  TEXT        NOT NULL,
    status          TEXT        NOT NULL DEFAULT 'draft', -- draft, active, inactive
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);
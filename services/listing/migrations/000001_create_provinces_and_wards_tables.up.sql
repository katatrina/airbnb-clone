CREATE TABLE IF NOT EXISTS provinces
(
    code       TEXT PRIMARY KEY,
    full_name  TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS wards
(
    code          TEXT PRIMARY KEY,
    full_name     TEXT        NOT NULL,
    province_code TEXT        NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
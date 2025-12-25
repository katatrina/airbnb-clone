CREATE TABLE restaurants
(
    id         UUID PRIMARY KEY,
    name       TEXT        NOT NULL,
    address    TEXT        NOT NULL,
    phone      TEXT        NOT NULL,
    email      TEXT        NOT NULL,
    is_active  BOOLEAN     NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
BEGIN;

-- Provinces (Tỉnh/Thành phố)
CREATE TABLE provinces
(
    code       INTEGER PRIMARY KEY,
    full_name  TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Districts (Quận/Huyện)
CREATE TABLE districts
(
    code          INTEGER PRIMARY KEY,
    full_name     TEXT        NOT NULL,
    province_code INTEGER     NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Wards (Phường/Xã/Thị trấn)
CREATE TABLE wards
(
    code          INTEGER PRIMARY KEY,
    full_name     TEXT        NOT NULL,
    district_code INTEGER     NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


COMMIT;

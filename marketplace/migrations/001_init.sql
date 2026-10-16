CREATE TABLE IF NOT EXISTS listings (
    id              TEXT PRIMARY KEY,
    car_id          TEXT NOT NULL,
    car_brand       TEXT NOT NULL,
    car_model       TEXT NOT NULL,
    car_year        INT  NOT NULL,
    car_mileage     INT  NOT NULL,
    seller_id       TEXT NOT NULL,
    seller_name     TEXT NOT NULL,
    seller_type     TEXT NOT NULL,
    title           TEXT NOT NULL,
    description     TEXT NOT NULL,
    photos          TEXT[] NOT NULL DEFAULT '{}',
    price           INT  NOT NULL,
    status          TEXT NOT NULL,
    has_warranty    BOOLEAN NOT NULL DEFAULT FALSE,
    has_test_drive  BOOLEAN NOT NULL DEFAULT FALSE,
    region          TEXT NOT NULL,
    commission      INT  NOT NULL,
    delivery_id     TEXT NOT NULL DEFAULT '',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_listings_status ON listings(status);

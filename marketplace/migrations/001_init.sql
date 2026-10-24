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
    price_amount    BIGINT NOT NULL,
    status          TEXT NOT NULL,
    has_warranty    BOOLEAN NOT NULL DEFAULT FALSE,
    has_test_drive  BOOLEAN NOT NULL DEFAULT FALSE,
    region          TEXT NOT NULL,
    commission      INT  NOT NULL,
    delivery_id     TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_listings_status ON listings(status);

CREATE TABLE IF NOT EXISTS outbox (
    id           TEXT PRIMARY KEY,
    event_type   TEXT NOT NULL,
    payload      JSONB NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    published_at TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_outbox_unpublished ON outbox(created_at) WHERE published_at IS NULL;

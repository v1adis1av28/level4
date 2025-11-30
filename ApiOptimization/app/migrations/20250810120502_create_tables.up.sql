CREATE TABLE IF NOT EXISTS  orders (
    order_uid         TEXT PRIMARY KEY UNIQUE,
    track_number      TEXT NOT NULL,
    entry             TEXT NOT NULL,
    locale            TEXT,
    internal_signature TEXT,
    customer_id       TEXT,
    delivery_service  TEXT,
    shardkey          TEXT,
    sm_id             INT,
    date_created      TIMESTAMP,
    oof_shard         TEXT
);

CREATE TABLE IF NOT EXISTS delivery (
    order_uid   TEXT PRIMARY KEY REFERENCES orders(order_uid) ON DELETE CASCADE,
    name        TEXT NOT NULL,
    phone       TEXT NOT NULL,
    zip         TEXT,
    city        TEXT,
    address     TEXT,
    region      TEXT,
    email       TEXT
);

CREATE TABLE IF NOT EXISTS payment (
    order_uid      TEXT PRIMARY KEY REFERENCES orders(order_uid) ON DELETE CASCADE,
    transaction    TEXT NOT NULL,
    request_id     TEXT,
    currency       TEXT,
    provider       TEXT,
    amount         NUMERIC(12,2),
    payment_dt     BIGINT,
    bank           TEXT,
    delivery_cost  NUMERIC(12,2),
    goods_total    NUMERIC(12,2),
    custom_fee     NUMERIC(12,2)
);

CREATE TABLE IF NOT EXISTS items (
    id SERIAL PRIMARY KEY,
    order_uid     TEXT REFERENCES orders(order_uid) ON DELETE CASCADE,
    chrt_id       BIGINT,
    track_number  TEXT,
    price         NUMERIC(12,2),
    rid           TEXT,
    name          TEXT,
    sale          NUMERIC(5,2),
    size          TEXT,
    total_price   NUMERIC(12,2),
    nm_id         BIGINT,
    brand         TEXT,
    status        INTEGER
);

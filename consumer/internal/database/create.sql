-- CREATE TABLE orders (
--     order_uid VARCHAR(255) PRIMARY KEY,
--     track_number VARCHAR(255),
--     entry VARCHAR(255),
--     delivery UUID,
--     payment UUID,
--     items UUID[],
--     locale VARCHAR(3),
--     internal_signature VARCHAR(255),
--     customer_id VARCHAR(255),
--     delivery_service VARCHAR(255),
--     shardkey INT,
--     sm_id INT,
--     date_created VARCHAR(255),
--     oof_shard INT
-- );

-- CREATE TABLE payment (
--     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--     transaction VARCHAR(255),
--     request_id VARCHAR(255),
--     currency VARCHAR(3),
--     provider VARCHAR(255),
--     amount INT,
--     payment_dt INT,
--     bank VARCHAR(255),
--     delivery_cost INT,
--     goods_total INT,
--     custom_fee INT
-- );

-- CREATE TABLE items (
--     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--     chrt_id INT,
--     track_number VARCHAR(255),
--     price INT,
--     rid VARCHAR(255),
--     name VARCHAR(255),
--     sale INT,
--     size INT,
--     total_price INT,
--     nm_id INT,
--     brand VARCHAR(255),
--     status VARCHAR(255)
-- );

-- CREATE TABLE delivery (
--     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--     name VARCHAR(255),
--     phone VARCHAR(255),
--     zip INT,
--     city VARCHAR(255),
--     address VARCHAR(255),
--     region VARCHAR(255),
--     email VARCHAR(255)
-- );

CREATE TABLE IF NOT EXISTS orders
(
    order_uid TEXT PRIMARY KEY,
    track_number TEXT NOT NULL,
    entry TEXT NOT NULL,
    locate TEXT NOT NULL,
    internal_signature TEXT,
    customer_id TEXT NOT NULL,
    delivery_service TEXT NOT NULL,
    shard_key TEXT NOT NULL,
    sm_id TEXT NOT NULL,
    date_created TIMESTAMP,
    oof_shard TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS delivery
(
    order_id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    phone TEXT NOT NULL,
    zip TEXT NOT NULL,
    city TEXT NOT NULL,
    address TEXT NOT NULL,
    region TEXT NOT NULL,
    email TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_delivery_order_id ON delivery(order_id);

CREATE TABLE IF NOT EXISTS payments
(
    order_id TEXT PRIMARY KEY ,
    transaction TEXT NOT NULL,
    request_id TEXT NOT NULL,
    currency TEXT NOT NULL,
    provider TEXT NOT NULL,
    amount BIGINT NOT NULL,
    payment_dt BIGINT NOT NULL,
    bank TEXT NOT NULL,
    delivery_cost BIGINT NOT NULL,
    goods_total BIGINT NOT NULL,
    custom_fee BIGINT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_payments_order_id ON payments(order_id);

CREATE TABLE IF NOT EXISTS items
(
    order_id TEXT NOT NULL,
    chrt_id BIGINT NOT NULL,
    track_number TEXT NOT NULL,
    price BIGINT NOT NULL,
    rid TEXT NOT NULL,
    name TEXT NOT NULL,
    sale BIGINT NOT NULL,
    size TEXT NOT NULL,
    total_price BIGINT NOT NULL,
    nm_id BIGINT NOT NULL,
    brand TEXT NOT NULL,
    status BIGINT NOT NULL,

    PRIMARY KEY (order_id, chrt_id)
    );

CREATE INDEX IF NOT EXISTS idx_items_order_id ON items(order_id);
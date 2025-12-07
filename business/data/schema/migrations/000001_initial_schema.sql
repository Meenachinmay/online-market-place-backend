-- +goose Up
CREATE TABLE products (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    price BIGINT NOT NULL,
    buyer_reward_points INTEGER NOT NULL,
    author_reward_points INTEGER NOT NULL
);

CREATE TABLE blogs (
    id TEXT PRIMARY KEY,
    author_id TEXT NOT NULL,
    content TEXT NOT NULL,
    product_id TEXT NOT NULL REFERENCES products(id)
);

CREATE TABLE orders (
    id TEXT PRIMARY KEY,
    buyer_id TEXT NOT NULL,
    product_id TEXT NOT NULL REFERENCES products(id),
    blog_id TEXT NOT NULL,
    amount BIGINT NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE wallets (
    user_id TEXT PRIMARY KEY,
    soda_points BIGINT NOT NULL DEFAULT 0,
    soda_balance BIGINT NOT NULL DEFAULT 0
);

CREATE TABLE transactions (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    type TEXT NOT NULL,
    amount BIGINT NOT NULL,
    related_order_id TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE transactions;
DROP TABLE wallets;
DROP TABLE orders;
DROP TABLE blogs;
DROP TABLE products;

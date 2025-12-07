-- name: GetProduct :one
SELECT * FROM products WHERE id = $1;

-- name: ListProducts :many
SELECT * FROM products;

-- name: CreateProduct :one
INSERT INTO products (id, name, description, price, buyer_reward_points, author_reward_points) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: CreateBlog :one
INSERT INTO blogs (id, author_id, content, product_id) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetBlog :one
SELECT * FROM blogs WHERE id = $1;

-- name: ListBlogs :many
SELECT * FROM blogs;

-- name: CreateOrder :one
INSERT INTO orders (id, buyer_id, product_id, blog_id, amount, status, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: CountOrdersByBuyer :one
SELECT COUNT(*) FROM orders WHERE buyer_id = $1;

-- name: CountOrdersByBuyerAndProduct :one
SELECT COUNT(*) FROM orders WHERE buyer_id = $1 AND product_id = $2;

-- name: GetWallet :one
SELECT * FROM wallets WHERE user_id = $1;

-- name: CreateWallet :one
INSERT INTO wallets (user_id, soda_points, soda_balance) VALUES ($1, 0, 0) ON CONFLICT (user_id) DO NOTHING RETURNING *;

-- name: AddPoints :one
UPDATE wallets SET soda_points = soda_points + sqlc.arg(amount) WHERE user_id = sqlc.arg(user_id) RETURNING *;

-- name: AddBalance :one
UPDATE wallets SET soda_balance = soda_balance + sqlc.arg(amount) WHERE user_id = sqlc.arg(user_id) RETURNING *;

-- name: ConvertPointsToBalance :one
UPDATE wallets 
SET soda_points = soda_points - sqlc.arg(points_deducted),
    soda_balance = soda_balance + sqlc.arg(balance_added)
WHERE user_id = sqlc.arg(user_id) AND soda_points >= sqlc.arg(points_deducted)
RETURNING *;

-- name: CreateTransaction :one
INSERT INTO transactions (id, user_id, type, amount, related_order_id) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: Order :one
SELECT * FROM orders WHERE id = $1 LIMIT 1;

-- name: AddOrder :one
INSERT INTO orders (user_id, status) VALUES ($1, 0) RETURNING *;

-- name: SetAwaitingPayment :exec
UPDATE orders SET status = 1 WHERE id = $1;

-- name: FailOrder :exec
UPDATE orders SET status = 2 WHERE id = $1;

-- name: PayOrder :exec
UPDATE orders SET status = 3 WHERE id = $1;

-- name: CancelOrder :exec
UPDATE orders SET status = 4 WHERE id = $1;

-- name: AddOrderItem :one
INSERT INTO orders_items (order_id, sku, count) VALUES ($1, $2, $3) RETURNING *;

-- name: OrderItems :many
SELECT * FROM orders_items WHERE order_id = $1;

-- name: UpdateStock :exec
UPDATE stocks SET reserved = $2, total_count = $3 WHERE sku = $1;

-- name: Stock :one
SELECT * FROM stocks WHERE sku = $1 LIMIT 1;
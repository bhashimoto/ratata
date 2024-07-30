-- name: CreateTransaction :one
INSERT INTO transactions (created_at, modified_at, description, amount, paid_by)
VALUES (?, ?, ?, ?, ?)
RETURNING *;

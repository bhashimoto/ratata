-- name: CreateTransaction :one
INSERT INTO transactions (account_id, created_at, modified_at, description, amount, paid_by)
VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetTransactionByID :one
SELECT *
FROM transactions
WHERE id = ?;

-- name: GetTransactionsByAccountID :many
SELECT *
FROM transactions
WHERE account_id = ?;

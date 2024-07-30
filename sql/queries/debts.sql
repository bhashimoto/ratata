-- name: CreateDebt :one
INSERT INTO debts (created_at, modified_at, user_id, transaction_id, amount)
VALUES (?,?,?,?,?)
RETURNING *;

-- name: GetDebtsFromTransaction :many
SELECT *
FROM debts
WHERE transaction_id = ?;

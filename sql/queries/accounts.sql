-- name: GetAccounts :many
SELECT *
FROM accounts;

-- name: CreateAccount :one
INSERT INTO accounts (name, created_at, modified_at)
VALUES (?, ?, ?)
RETURNING *;

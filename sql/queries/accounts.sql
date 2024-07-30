-- name: GetAccounts :many
SELECT *
FROM accounts;

-- name: CreateAccount :one
INSERT INTO accounts (name, created_at, modified_at)
VALUES (?, ?, ?)
RETURNING *;

-- name: GetAccountsByUserID :many
SELECT a.*
FROM accounts a
JOIN user_accounts ua ON a.id = ua.account_id
JOIN users u ON ua.user_id = u.id
WHERE u.id = ?;

-- name: GetAccountByID :one
SELECT *
FROM accounts
WHERE id = ?;

-- name: CreateUserAccount :one
INSERT INTO user_accounts(user_id, account_id, created_at, modified_at)
VALUES(?,?,?,?)
RETURNING *;

-- name: GetUsersByAccount :many
SELECT u.*
FROM users u 
JOIN user_accounts ua ON u.id = ua.user_id
WHERE ua.account_id = ?;


-- name: GetAccountsByUser :many
SELECT a.*
FROM accounts a 
JOIN user_accounts ua ON a.id = ua.account_id
WHERE ua.user_id = ?;

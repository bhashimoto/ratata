-- name: GetUsers :many
SELECT *
FROM users;

-- name: CreateUser :one
INSERT INTO users (name, created_at, modified_at)
VALUES (?, ?, ?)
RETURNING *;

-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = ?;

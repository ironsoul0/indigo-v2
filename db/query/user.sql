-- name: CreateUser :one
INSERT INTO users (
    chat_id,
    username,
    password,
    active
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetUsers :many
SELECT * FROM users;

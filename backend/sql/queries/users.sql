-- name: GetUser :one
SELECT id, email, password_hash, created_at, updated_at
FROM users
WHERE id = ? LIMIT 1;

-- name: GetUserByEmail :one
SELECT id, email, password_hash, created_at, updated_at
FROM users
WHERE email = ? LIMIT 1;

-- name: ListUsers :many
SELECT id, email, created_at, updated_at
FROM users
ORDER BY created_at DESC;

-- name: CreateUser :execresult
INSERT INTO users (
    email, password_hash
) VALUES (
    ?, ?
);

-- name: UpdateUser :exec
UPDATE users
SET email = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = ?;
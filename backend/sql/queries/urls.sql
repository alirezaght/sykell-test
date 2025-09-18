-- name: GetUrl :one
SELECT id, user_id, normalized_url, url_hash, domain, created_at, updated_at
FROM urls
WHERE id = ? LIMIT 1;

-- name: GetUrlByUserAndUrl :one
SELECT id, user_id, normalized_url, url_hash, domain, created_at, updated_at
FROM urls
WHERE user_id = ? AND url_hash = UNHEX(MD5(?)) LIMIT 1;

-- name: ListUrlsByUser :many
SELECT id, user_id, normalized_url, url_hash, domain, created_at, updated_at
FROM urls
WHERE user_id = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: CountUrlsByUser :one
SELECT COUNT(*)
FROM urls
WHERE user_id = ?;

-- name: ListUrlsByDomain :many
SELECT id, user_id, normalized_url, url_hash, domain, created_at, updated_at
FROM urls
WHERE domain = ?
ORDER BY created_at DESC;

-- name: CreateUrl :execresult
INSERT INTO urls (
    user_id, normalized_url, domain
) VALUES (
    ?, ?, ?
);

-- name: UpdateUrl :exec
UPDATE urls
SET normalized_url = ?, domain = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: DeleteUrl :exec
DELETE FROM urls
WHERE id = ?;

-- name: DeleteUrlsByUser :exec
DELETE FROM urls
WHERE user_id = ?;
-- name: GetCrawl :one
SELECT id, url_id, status, queued_at, started_at, finished_at, error_message,
       html_version, page_title, h1_count, h2_count, h3_count, h4_count, h5_count, h6_count,
       internal_links_count, external_links_count, inaccessible_links_count, has_login_form,
       created_at, updated_at
FROM crawls
WHERE id = ? LIMIT 1;

-- name: GetCrawlByUrl :one
SELECT id, url_id, status, queued_at, started_at, finished_at, error_message,
       html_version, page_title, h1_count, h2_count, h3_count, h4_count, h5_count, h6_count,
       internal_links_count, external_links_count, inaccessible_links_count, has_login_form,
       created_at, updated_at
FROM crawls
WHERE url_id = ?
ORDER BY created_at DESC LIMIT 1;

-- name: ListCrawlsByUser :many
SELECT c.id, c.url_id, c.status, c.queued_at, c.started_at, c.finished_at, c.error_message,
       c.html_version, c.page_title, c.h1_count, c.h2_count, c.h3_count, c.h4_count, c.h5_count, c.h6_count,
       c.internal_links_count, c.external_links_count, c.inaccessible_links_count, c.has_login_form,
       c.created_at, c.updated_at,
       u.normalized_url, u.domain
FROM crawls c
JOIN urls u ON c.url_id = u.id
WHERE u.user_id = ?
ORDER BY c.created_at DESC
LIMIT ? OFFSET ?;

-- name: CountCrawlsByUser :one
SELECT COUNT(*)
FROM crawls c
JOIN urls u ON c.url_id = u.id
WHERE u.user_id = ?;

-- name: ListCrawlsByStatus :many
SELECT id, url_id, status, queued_at, started_at, finished_at, error_message,
       html_version, page_title, h1_count, h2_count, h3_count, h4_count, h5_count, h6_count,
       internal_links_count, external_links_count, inaccessible_links_count, has_login_form,
       created_at, updated_at
FROM crawls
WHERE status = ?
ORDER BY queued_at ASC;

-- name: ListQueuedCrawls :many
SELECT id, url_id, status, queued_at, started_at, finished_at, error_message,
       html_version, page_title, h1_count, h2_count, h3_count, h4_count, h5_count, h6_count,
       internal_links_count, external_links_count, inaccessible_links_count, has_login_form,
       created_at, updated_at
FROM crawls
WHERE status = 'queued'
ORDER BY queued_at ASC
LIMIT ?;

-- name: ListRunningCrawls :many
SELECT id, url_id, status, queued_at, started_at, finished_at, error_message,
       html_version, page_title, h1_count, h2_count, h3_count, h4_count, h5_count, h6_count,
       internal_links_count, external_links_count, inaccessible_links_count, has_login_form,
       created_at, updated_at
FROM crawls
WHERE status = 'running';

-- name: CreateCrawl :execresult
INSERT INTO crawls (
    url_id, status
) VALUES (
    ?, ?
);

-- name: UpdateCrawlStatus :exec
UPDATE crawls
SET status = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: StartCrawl :exec
UPDATE crawls
SET status = 'running', started_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: FinishCrawl :exec
UPDATE crawls
SET status = 'done', finished_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: ErrorCrawl :exec
UPDATE crawls
SET status = 'error', finished_at = CURRENT_TIMESTAMP, error_message = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: StopCrawl :exec
UPDATE crawls
SET status = 'stopped', finished_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: UpdateCrawlResults :exec
UPDATE crawls
SET html_version = ?, page_title = ?, 
    h1_count = ?, h2_count = ?, h3_count = ?, h4_count = ?, h5_count = ?, h6_count = ?,
    internal_links_count = ?, external_links_count = ?, inaccessible_links_count = ?,
    has_login_form = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: DeleteCrawl :exec
DELETE FROM crawls
WHERE id = ?;

-- name: DeleteCrawlsByUrl :exec
DELETE FROM crawls
WHERE url_id = ?;
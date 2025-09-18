-- name: CountOfActiveCrawlForUrlId :one
SELECT COUNT(*)
FROM crawls
WHERE url_id = ? AND status IN ('queued', 'running');

-- name: GetCrawlByWorkflowID :one
SELECT id, url_id, status, queued_at, started_at, finished_at, error_message, workflow_id,
       html_version, page_title, h1_count, h2_count, h3_count, h4_count, h5_count, h6_count,
       internal_links_count, external_links_count, inaccessible_links_count, has_login_form,
       created_at, updated_at
FROM crawls
WHERE workflow_id = ? 
ORDER BY created_at DESC
LIMIT 1;

-- name: GetActiveCrawlsForUrlId :many
SELECT id, url_id, status, queued_at, started_at, finished_at, error_message, workflow_id,
       html_version, page_title, h1_count, h2_count, h3_count, h4_count, h5_count, h6_count,
       internal_links_count, external_links_count, inaccessible_links_count, has_login_form,
       created_at, updated_at
FROM crawls
WHERE url_id = ? AND status IN ('queued', 'running')
ORDER BY created_at DESC;

-- name: QueueCrawl :execresult
INSERT INTO crawls (
    url_id, status, workflow_id, queued_at
) VALUES (
    ?, 'queued', ?, CURRENT_TIMESTAMP
);


-- name: SetCrawlRunning :exec
UPDATE crawls
SET status='running',
    started_at = IFNULL(started_at, CURRENT_TIMESTAMP),
    updated_at = CURRENT_TIMESTAMP,
    error_message = NULL
WHERE id = ?;

-- name: SetCrawlDone :exec
UPDATE crawls
SET status='done',
    finished_at = IFNULL(finished_at, CURRENT_TIMESTAMP),
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: SetCrawlStopped :exec
UPDATE crawls
SET status='stopped',
    finished_at = IFNULL(finished_at, CURRENT_TIMESTAMP),
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: SetCrawlError :exec
UPDATE crawls
SET status='error',
    finished_at = IFNULL(finished_at, CURRENT_TIMESTAMP),
    updated_at = CURRENT_TIMESTAMP,
    error_message = ?
WHERE id = ?;


-- name: UpdateCrawlResult :exec
UPDATE crawls
SET
    status = ?,
    html_version = ?,
    page_title = ?,
    h1_count = ?,
    h2_count = ?,
    h3_count = ?,
    h4_count = ?,
    h5_count = ?,
    h6_count = ?,
    internal_links_count = ?,
    external_links_count = ?,
    inaccessible_links_count = ?,
    has_login_form = ?,
    error_message = ?,
    finished_at = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?;
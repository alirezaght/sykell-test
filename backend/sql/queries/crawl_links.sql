-- name: GetCrawlLink :one
SELECT id, crawl_id, href, absolute_url, absolute_url_hash, is_internal, status_code, is_accessible, anchor_text, created_at
FROM crawl_links
WHERE id = ? LIMIT 1;

-- name: ListCrawlLinksByCrawl :many
SELECT id, crawl_id, href, absolute_url, absolute_url_hash, is_internal, status_code, is_accessible, anchor_text, created_at
FROM crawl_links
WHERE crawl_id = ?
ORDER BY created_at ASC
LIMIT ? OFFSET ?;

-- name: CountCrawlLinksByCrawl :one
SELECT COUNT(*)
FROM crawl_links
WHERE crawl_id = ?;

-- name: ListInternalLinksByCrawl :many
SELECT id, crawl_id, href, absolute_url, absolute_url_hash, is_internal, status_code, is_accessible, anchor_text, created_at
FROM crawl_links
WHERE crawl_id = ? AND is_internal = TRUE
ORDER BY created_at ASC;

-- name: ListExternalLinksByCrawl :many
SELECT id, crawl_id, href, absolute_url, absolute_url_hash, is_internal, status_code, is_accessible, anchor_text, created_at
FROM crawl_links
WHERE crawl_id = ? AND is_internal = FALSE
ORDER BY created_at ASC;

-- name: ListInaccessibleLinksByCrawl :many
SELECT id, crawl_id, href, absolute_url, absolute_url_hash, is_internal, status_code, is_accessible, anchor_text, created_at
FROM crawl_links
WHERE crawl_id = ? AND is_accessible = FALSE
ORDER BY created_at ASC;

-- name: CountLinksByCrawlAndType :one
SELECT 
    COUNT(CASE WHEN is_internal = TRUE THEN 1 END) as internal_count,
    COUNT(CASE WHEN is_internal = FALSE THEN 1 END) as external_count,
    COUNT(CASE WHEN is_accessible = FALSE THEN 1 END) as inaccessible_count
FROM crawl_links
WHERE crawl_id = ?;

-- name: GetCrawlLinkByUrl :one
SELECT id, crawl_id, href, absolute_url, absolute_url_hash, is_internal, status_code, is_accessible, anchor_text, created_at
FROM crawl_links
WHERE crawl_id = ? AND absolute_url_hash = UNHEX(MD5(?)) LIMIT 1;

-- name: CreateCrawlLink :execresult
INSERT INTO crawl_links (
    crawl_id, href, absolute_url, is_internal, anchor_text
) VALUES (
    ?, ?, ?, ?, ?
);

-- name: CreateCrawlLinkWithStatus :execresult
INSERT INTO crawl_links (
    crawl_id, href, absolute_url, is_internal, status_code, anchor_text
) VALUES (
    ?, ?, ?, ?, ?, ?
);

-- name: UpdateCrawlLinkStatus :exec
UPDATE crawl_links
SET status_code = ?
WHERE id = ?;

-- name: UpdateCrawlLinkAnchorText :exec
UPDATE crawl_links
SET anchor_text = ?
WHERE id = ?;

-- name: DeleteCrawlLink :exec
DELETE FROM crawl_links
WHERE id = ?;

-- name: DeleteCrawlLinksByCrawl :exec
DELETE FROM crawl_links
WHERE crawl_id = ?;

-- name: ListUniqueDomainsFromLinks :many
SELECT DISTINCT 
    SUBSTRING_INDEX(SUBSTRING_INDEX(absolute_url, '/', 3), '://', -1) as domain,
    COUNT(*) as link_count
FROM crawl_links
WHERE crawl_id = ? AND is_internal = FALSE
GROUP BY domain
ORDER BY link_count DESC;
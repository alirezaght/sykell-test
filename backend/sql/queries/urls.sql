-- name: CountUrlsByUser :one
SELECT COUNT(*)
FROM urls
WHERE user_id = ?;

-- name: GetUrlByIdAndUserId :one
SELECT id, user_id, normalized_url, domain, created_at, updated_at
FROM urls
WHERE id = ? AND user_id = ?;

-- name: CountUrlsWithFilter :one
SELECT COUNT(*)
FROM urls u
LEFT JOIN crawls c 
  ON u.id = c.url_id 
 AND c.id = (
    SELECT c2.id 
    FROM crawls c2 
    WHERE c2.url_id = u.id 
    ORDER BY c2.created_at DESC 
    LIMIT 1
)
WHERE u.user_id = sqlc.arg(user_id)
  AND (sqlc.arg(query_filter) = '' OR u.normalized_url LIKE CONCAT('%', sqlc.arg(query_filter), '%') OR c.page_title LIKE CONCAT('%', sqlc.arg(query_filter), '%'));


-- name: GetUrlsWithLatestCrawlsFiltered :many
SELECT 
    u.id as url_id,
    u.normalized_url,
    u.domain,
    u.created_at as url_created_at,
    c.id as crawl_id,
    c.status,
    c.workflow_id,
    c.queued_at,
    c.started_at,
    c.finished_at,
    c.html_version,
    c.page_title,
    c.h1_count,
    c.h2_count,
    c.h3_count,
    c.h4_count,
    c.h5_count,
    c.h6_count,
    c.internal_links_count,
    c.external_links_count,
    c.inaccessible_links_count,
    c.has_login_form,
    c.error_message,
    c.created_at as crawl_created_at,
    c.updated_at as crawl_updated_at
FROM urls u
LEFT JOIN crawls c 
  ON u.id = c.url_id 
 AND c.id = (
    SELECT c2.id 
    FROM crawls c2 
    WHERE c2.url_id = u.id 
    ORDER BY c2.created_at DESC 
    LIMIT 1
)
WHERE u.user_id = sqlc.arg(user_id)
  AND (sqlc.arg(query_filter) = '' OR u.normalized_url LIKE CONCAT('%', sqlc.arg(query_filter), '%') OR c.page_title LIKE CONCAT('%', sqlc.arg(query_filter), '%'))  
ORDER BY
  -- url fields  
  CASE WHEN sqlc.arg(sort_by)='normalized_url'  AND sqlc.arg(sort_dir)='asc'  THEN u.normalized_url END ASC,
  CASE WHEN sqlc.arg(sort_by)='normalized_url'  AND sqlc.arg(sort_dir)='desc' THEN u.normalized_url END DESC,
  CASE WHEN sqlc.arg(sort_by)='domain'          AND sqlc.arg(sort_dir)='asc'  THEN u.domain END ASC,
  CASE WHEN sqlc.arg(sort_by)='domain'          AND sqlc.arg(sort_dir)='desc' THEN u.domain END DESC,
  CASE WHEN sqlc.arg(sort_by)='url_created_at'  AND sqlc.arg(sort_dir)='asc'  THEN u.created_at END ASC,
  CASE WHEN sqlc.arg(sort_by)='url_created_at'  AND sqlc.arg(sort_dir)='desc' THEN u.created_at END DESC,

  -- crawl fields  
  CASE WHEN sqlc.arg(sort_by)='status'          AND sqlc.arg(sort_dir)='asc'  THEN c.status END ASC,
  CASE WHEN sqlc.arg(sort_by)='status'          AND sqlc.arg(sort_dir)='desc' THEN c.status END DESC,  
  CASE WHEN sqlc.arg(sort_by)='html_version'    AND sqlc.arg(sort_dir)='asc'  THEN c.html_version END ASC,
  CASE WHEN sqlc.arg(sort_by)='html_version'    AND sqlc.arg(sort_dir)='desc' THEN c.html_version END DESC,
  CASE WHEN sqlc.arg(sort_by)='page_title'      AND sqlc.arg(sort_dir)='asc'  THEN c.page_title END ASC,
  CASE WHEN sqlc.arg(sort_by)='page_title'      AND sqlc.arg(sort_dir)='desc' THEN c.page_title END DESC,
  CASE WHEN sqlc.arg(sort_by)='h1_count'        AND sqlc.arg(sort_dir)='asc'  THEN c.h1_count END ASC,
  CASE WHEN sqlc.arg(sort_by)='h1_count'        AND sqlc.arg(sort_dir)='desc' THEN c.h1_count END DESC,
  CASE WHEN sqlc.arg(sort_by)='h2_count'        AND sqlc.arg(sort_dir)='asc'  THEN c.h2_count END ASC,
  CASE WHEN sqlc.arg(sort_by)='h2_count'        AND sqlc.arg(sort_dir)='desc' THEN c.h2_count END DESC,
  CASE WHEN sqlc.arg(sort_by)='h3_count'        AND sqlc.arg(sort_dir)='asc'  THEN c.h3_count END ASC,
  CASE WHEN sqlc.arg(sort_by)='h3_count'        AND sqlc.arg(sort_dir)='desc' THEN c.h3_count END DESC,
  CASE WHEN sqlc.arg(sort_by)='h4_count'        AND sqlc.arg(sort_dir)='asc'  THEN c.h4_count END ASC,
  CASE WHEN sqlc.arg(sort_by)='h4_count'        AND sqlc.arg(sort_dir)='desc' THEN c.h4_count END DESC,
  CASE WHEN sqlc.arg(sort_by)='h5_count'        AND sqlc.arg(sort_dir)='asc'  THEN c.h5_count END ASC,
  CASE WHEN sqlc.arg(sort_by)='h5_count'        AND sqlc.arg(sort_dir)='desc' THEN c.h5_count END DESC,
  CASE WHEN sqlc.arg(sort_by)='h6_count'        AND sqlc.arg(sort_dir)='asc'  THEN c.h6_count END ASC,
  CASE WHEN sqlc.arg(sort_by)='h6_count'        AND sqlc.arg(sort_dir)='desc' THEN c.h6_count END DESC,
  CASE WHEN sqlc.arg(sort_by)='internal_links_count'     AND sqlc.arg(sort_dir)='asc'  THEN c.internal_links_count END ASC,
  CASE WHEN sqlc.arg(sort_by)='internal_links_count'     AND sqlc.arg(sort_dir)='desc' THEN c.internal_links_count END DESC,
  CASE WHEN sqlc.arg(sort_by)='external_links_count'     AND sqlc.arg(sort_dir)='asc'  THEN c.external_links_count END ASC,
  CASE WHEN sqlc.arg(sort_by)='external_links_count'     AND sqlc.arg(sort_dir)='desc' THEN c.external_links_count END DESC,
  CASE WHEN sqlc.arg(sort_by)='inaccessible_links_count' AND sqlc.arg(sort_dir)='asc'  THEN c.inaccessible_links_count END ASC,
  CASE WHEN sqlc.arg(sort_by)='inaccessible_links_count' AND sqlc.arg(sort_dir)='desc' THEN c.inaccessible_links_count END DESC,
  CASE WHEN sqlc.arg(sort_by)='has_login_form'           AND sqlc.arg(sort_dir)='asc'  THEN c.has_login_form END ASC,
  CASE WHEN sqlc.arg(sort_by)='has_login_form'           AND sqlc.arg(sort_dir)='desc' THEN c.has_login_form END DESC

LIMIT ? OFFSET ?;


-- name: CreateUrl :execresult
INSERT INTO urls (
    user_id, normalized_url, domain
) VALUES (
    ?, ?, ?
);


-- name: DeleteURLByIdAndUserId :exec
DELETE FROM urls
WHERE id = ? AND user_id = ?;
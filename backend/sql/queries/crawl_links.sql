-- name: CreateInaccessibleLink :execresult
INSERT INTO crawl_links (
    crawl_id, href, absolute_url, is_internal, status_code, anchor_text
) VALUES (
    ?, ?, ?, ?, ?, ?
);

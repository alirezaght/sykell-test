-- Migration to fix character encoding for crawl_links table
-- This ensures the anchor_text column can properly store UTF-8 characters including emojis and special characters

ALTER TABLE crawl_links 
MODIFY COLUMN anchor_text VARCHAR(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL;

-- Also fix other text columns that might have encoding issues
ALTER TABLE crawl_links 
MODIFY COLUMN href VARCHAR(2083) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL;

ALTER TABLE crawl_links 
MODIFY COLUMN absolute_url VARCHAR(2083) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL;

-- Fix crawls table text columns as well
ALTER TABLE crawls 
MODIFY COLUMN error_message TEXT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL;

ALTER TABLE crawls 
MODIFY COLUMN workflow_id VARCHAR(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL;

ALTER TABLE crawls 
MODIFY COLUMN html_version VARCHAR(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL;

ALTER TABLE crawls 
MODIFY COLUMN page_title VARCHAR(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL;
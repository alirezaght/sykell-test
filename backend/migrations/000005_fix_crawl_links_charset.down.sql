-- Rollback migration for crawl_links charset fix
ALTER TABLE crawl_links 
MODIFY COLUMN anchor_text VARCHAR(1024) NULL;

ALTER TABLE crawl_links 
MODIFY COLUMN href VARCHAR(2083) NOT NULL;

ALTER TABLE crawl_links 
MODIFY COLUMN absolute_url VARCHAR(2083) NOT NULL;

-- Rollback crawls table changes
ALTER TABLE crawls 
MODIFY COLUMN error_message TEXT NULL;

ALTER TABLE crawls 
MODIFY COLUMN workflow_id VARCHAR(512) NOT NULL;

ALTER TABLE crawls 
MODIFY COLUMN html_version VARCHAR(32) NULL;

ALTER TABLE crawls 
MODIFY COLUMN page_title VARCHAR(512) NULL;
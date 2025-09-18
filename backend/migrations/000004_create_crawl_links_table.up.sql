CREATE TABLE crawl_links (
  id           CHAR(36) PRIMARY KEY DEFAULT (UUID()),
  crawl_id     CHAR(36) NOT NULL,
  href         VARCHAR(2083) NOT NULL,          
  absolute_url VARCHAR(2083) NOT NULL,          
  absolute_url_hash BINARY(16) AS (UNHEX(MD5(absolute_url))) STORED,
  is_internal  BOOLEAN NOT NULL,
  status_code  INT NULL,                        
  is_accessible BOOLEAN AS (CASE WHEN status_code IS NULL THEN NULL WHEN status_code BETWEEN 400 AND 599 THEN 0 ELSE 1 END) STORED,

  anchor_text  VARCHAR(1024) NULL,

  created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

  CONSTRAINT fk_links_crawl FOREIGN KEY (crawl_id) REFERENCES crawls(id) ON DELETE CASCADE,
  UNIQUE KEY uq_crawl_url (crawl_id, absolute_url_hash),
  KEY idx_links_internal (crawl_id, is_internal),
  KEY idx_links_access   (crawl_id, is_accessible),
  KEY idx_links_status   (crawl_id, status_code)
);
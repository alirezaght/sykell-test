CREATE TABLE crawls (
  id              CHAR(36) PRIMARY KEY DEFAULT (UUID()),
  url_id          CHAR(36) NOT NULL,
  status          ENUM('queued','running','stopped','done','error') NOT NULL DEFAULT 'queued',
  queued_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  started_at      TIMESTAMP NULL,
  finished_at     TIMESTAMP NULL,
  error_message   TEXT NULL,

  html_version    VARCHAR(32) NULL,
  page_title      VARCHAR(512) NULL,
  h1_count        INT UNSIGNED DEFAULT 0,
  h2_count        INT UNSIGNED DEFAULT 0,
  h3_count        INT UNSIGNED DEFAULT 0,
  h4_count        INT UNSIGNED DEFAULT 0,
  h5_count        INT UNSIGNED DEFAULT 0,
  h6_count        INT UNSIGNED DEFAULT 0,
  internal_links_count     INT UNSIGNED DEFAULT 0,
  external_links_count     INT UNSIGNED DEFAULT 0,
  inaccessible_links_count INT UNSIGNED DEFAULT 0,
  has_login_form  BOOLEAN NOT NULL DEFAULT FALSE,

  created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  CONSTRAINT fk_crawls_url FOREIGN KEY (url_id) REFERENCES urls(id) ON DELETE CASCADE,

  KEY idx_crawls_url (url_id),
  KEY idx_crawls_status (status),
  KEY idx_crawls_finished_at (finished_at)
);

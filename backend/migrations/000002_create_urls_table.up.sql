CREATE TABLE urls (
  id            CHAR(36) PRIMARY KEY DEFAULT (UUID()),
  user_id      CHAR(36) NOT NULL,
  normalized_url VARCHAR(2083) NOT NULL,
  url_hash BINARY(16) AS (UNHEX(MD5(normalized_url))) STORED,
  domain        VARCHAR(255)  NOT NULL,
  created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uq_urls_normalized (user_id, url_hash),
  CONSTRAINT fk_urls_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  KEY idx_urls_user_id (user_id),
  KEY idx_urls_domain (domain)
);
-- Initialize the database with proper charset and collation
ALTER DATABASE sykell_db CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- Grant additional privileges to the user
GRANT ALL PRIVILEGES ON sykell_db.* TO 'sykell_user'@'%';
FLUSH PRIVILEGES;
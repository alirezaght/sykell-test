CREATE DATABASE IF NOT EXISTS temporal CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE DATABASE IF NOT EXISTS temporal_visibility CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE USER IF NOT EXISTS 'temporal_user'@'%' IDENTIFIED BY 'temporal_password';
GRANT ALL PRIVILEGES ON temporal.* TO 'temporal_user'@'%';
GRANT ALL PRIVILEGES ON temporal_visibility.* TO 'temporal_user'@'%';
FLUSH PRIVILEGES;
CREATE DATABASE IF NOT EXISTS skogul_test;

USE skogul_test;

CREATE TABLE IF NOT EXISTS metrics (
    id INT AUTO_INCREMENT PRIMARY KEY,
    time DATETIME NOT NULL,
    host VARCHAR(255),
    metric_name VARCHAR(255),
    value DOUBLE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS source_metrics (
    id INT AUTO_INCREMENT PRIMARY KEY,
    time DATETIME NOT NULL,
    host VARCHAR(255),
    metric_name VARCHAR(255),
    value DOUBLE
);

INSERT INTO source_metrics (time, host, metric_name, value) VALUES
    (NOW(), 'host1', 'cpu_usage', 42.5),
    (NOW(), 'host2', 'memory_usage', 78.3),
    (NOW(), 'host1', 'disk_io', 1024.0);

-- Create user that requires X509 certificate authentication
-- The CN in the client cert must match 'testuser'
CREATE USER IF NOT EXISTS 'testuser'@'%' REQUIRE X509;
GRANT ALL PRIVILEGES ON skogul_test.* TO 'testuser'@'%';

FLUSH PRIVILEGES;

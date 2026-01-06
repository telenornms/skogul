-- Create testuser role (must match CN in client certificate)
CREATE USER testuser;
GRANT ALL PRIVILEGES ON DATABASE skogul_test TO testuser;

CREATE TABLE IF NOT EXISTS metrics (
    id SERIAL PRIMARY KEY,
    time TIMESTAMP NOT NULL,
    host VARCHAR(255),
    metric_name VARCHAR(255),
    value DOUBLE PRECISION,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create test table for receiver tests (pre-populated with data)
CREATE TABLE IF NOT EXISTS source_metrics (
    id SERIAL PRIMARY KEY,
    time TIMESTAMP NOT NULL,
    host VARCHAR(255),
    metric_name VARCHAR(255),
    value DOUBLE PRECISION
);

INSERT INTO source_metrics (time, host, metric_name, value) VALUES
    (NOW(), 'host1', 'cpu_usage', 42.5),
    (NOW(), 'host2', 'memory_usage', 78.3),
    (NOW(), 'host1', 'disk_io', 1024.0);

-- Grant privileges on tables to testuser
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO testuser;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO testuser;

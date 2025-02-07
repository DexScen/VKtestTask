CREATE TABLE IF NOT EXISTS pings (
    id SERIAL PRIMARY KEY,
    ip_address VARCHAR(255) NOT NULL,
    ping_time TIMESTAMP NOT NULL,
    last_successful_ping_date TIMESTAMP NOT NULL
);

DROP TABLE IF EXISTS event_db.events_buf;
DROP TABLE IF EXISTS event_db.events;
DROP DATABASE IF EXISTS event_db;
REVOKE ALL PRIVILEGES ON event_db.events_buf FROM event_writer;
REVOKE ALL PRIVILEGES ON event_db.events FROM event_writer;

CREATE DATABASE IF NOT EXISTS event_db;
CREATE TABLE IF NOT EXISTS event_db.events
(
    client_time DateTime64,
    device_id String,
    device_os String,
    session String,
    sequence UInt64,
    event String,
    param_int Int64,    
    param_str String,
    ip IPv4,
    server_time DateTime64
)
ENGINE = MergeTree()
ORDER BY (session, sequence);
CREATE TABLE IF NOT EXISTS event_db.events_buf AS event_db.events ENGINE = Buffer(event_db, events, 16, 10, 100, 10000, 1000000, 10000000, 100000000);

CREATE USER IF NOT EXISTS event_writer IDENTIFIED WITH sh   a256_password BY 'passwd';
GRANT ALL PRIVILEGES ON event_db.events_buf TO event_writer;    
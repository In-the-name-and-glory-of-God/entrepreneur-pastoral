-- Triggers must be dropped before the table.
DROP TRIGGER IF EXISTS set_timestamp_users ON users;
DROP TABLE IF EXISTS users;
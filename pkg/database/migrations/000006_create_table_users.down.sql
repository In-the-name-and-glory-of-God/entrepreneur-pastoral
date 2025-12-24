-- Indexes must be dropped before the table.
DROP INDEX IF EXISTS idx_users_church_id;
DROP INDEX IF EXISTS idx_users_address_id;

-- Triggers must be dropped before the table.
DROP TRIGGER IF EXISTS set_timestamp_users ON users;
DROP TABLE IF EXISTS users;
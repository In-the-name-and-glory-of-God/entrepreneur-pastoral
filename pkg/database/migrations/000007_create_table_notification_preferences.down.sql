-- Triggers must be dropped before the table.
DROP TRIGGER IF EXISTS set_timestamp_notification_preferences ON notification_preferences;
DROP TABLE IF EXISTS notification_preferences;
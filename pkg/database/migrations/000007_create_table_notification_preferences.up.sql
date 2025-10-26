-- Table: notification_preferences
-- Stores user-specific notification settings (one-to-one with users).
CREATE TABLE IF NOT EXISTS notification_preferences (
    user_id UUID PRIMARY KEY, -- This is both PK and FK
    notify_by_email BOOLEAN NOT NULL DEFAULT TRUE,
    notify_by_sms BOOLEAN NOT NULL DEFAULT FALSE,

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- Constraints
    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
        REFERENCES users(id)
        ON DELETE CASCADE -- If the user is deleted, delete their preferences too
        ON UPDATE CASCADE
);

-- Apply the trigger to 'updated_at' column
CREATE TRIGGER set_timestamp_notification_preferences
BEFORE UPDATE ON notification_preferences
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();
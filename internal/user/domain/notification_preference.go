package domain

import (
	"time"

	"github.com/google/uuid"
)

// NotificationPreferences corresponds to the "notification_preferences" table.
type NotificationPreferences struct {
	UserID        uuid.UUID `json:"-" db:"user_id"`
	NotifyByEmail bool      `json:"notify_by_email" db:"notify_by_email"`
	NotifyBySms   bool      `json:"notify_by_sms" db:"notify_by_sms"`
	CreatedAt     time.Time `json:"-" db:"created_at"`
	UpdatedAt     time.Time `json:"-" db:"updated_at"`
}

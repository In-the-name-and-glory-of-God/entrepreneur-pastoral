package domain

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type NotificationPreferencesRepository interface {
	Create(tx *sqlx.Tx, notificationPreferences *NotificationPreferences) error
	Update(tx *sqlx.Tx, notificationPreferences *NotificationPreferences) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*NotificationPreferences, error)
}

package domain

type NotificationPreferencesRepository interface {
	Create(notificationPreferences *NotificationPreferences) error
	Update(notificationPreferences *NotificationPreferences) error
	Delete(userID string) error
	GetByUserID(userID string) (*NotificationPreferences, error)
}

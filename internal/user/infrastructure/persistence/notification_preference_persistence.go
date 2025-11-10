package persistence

import (
	"context"
	"fmt"

	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/user/domain"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// NotificationPreferencesPersistence manages data access for the notification_preferences table.
type NotificationPreferencesPersistence struct {
	db   *sqlx.DB
	psql sq.StatementBuilderType
}

// NewNotificationPreferencesPersistence creates a new NotificationPreferencesPersistence.
func NewNotificationPreferencesPersistence(db *sqlx.DB) *NotificationPreferencesPersistence {
	return &NotificationPreferencesPersistence{
		db:   db,
		psql: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

// Create inserts a new notification preferences record for a user.
func (r *NotificationPreferencesPersistence) Create(tx *sqlx.Tx, prefs *domain.NotificationPreferences) error {
	query, args, err := r.psql.Insert("notification_preferences").
		Columns("user_id").
		Values(prefs.UserID).
		Suffix("ON CONFLICT (user_id) DO NOTHING").
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build create notificationPreferences query: %w", err)
	}

	if _, err := tx.Exec(query, args...); err != nil {
		return fmt.Errorf("failed to execute create notificationPreferences query: %w", err)
	}

	return nil
}

// Update modifies an existing notification preferences record.
func (r *NotificationPreferencesPersistence) Update(tx *sqlx.Tx, prefs *domain.NotificationPreferences) error {
	query, args, err := r.psql.Update("notification_preferences").
		Set("notify_by_email", prefs.NotifyByEmail).
		Set("notify_by_sms", prefs.NotifyBySms).
		Where(sq.Eq{"user_id": prefs.UserID}).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build update notificationPreferences query: %w", err)
	}

	if _, err := tx.Exec(query, args...); err != nil {
		return fmt.Errorf("failed to execute update notificationPreferences query: %w", err)
	}

	return nil
}

// GetByUserID retrieves a single notification preferences record by user ID.
func (r *NotificationPreferencesPersistence) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.NotificationPreferences, error) {
	var prefs domain.NotificationPreferences
	query, args, err := r.psql.Select("*").From("notification_preferences").
		Where(sq.Eq{"user_id": userID}).
		Limit(1).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build get notificationPreferences by userID query: %w", err)
	}

	if err := r.db.GetContext(ctx, &prefs, query, args...); err != nil {
		return nil, fmt.Errorf("failed to execute get notificationPreferences by userID query: %w", err)
	}

	return &prefs, nil
}

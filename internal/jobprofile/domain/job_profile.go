package domain

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// JobProfile corresponds to the "job_profiles" table.
type JobProfile struct {
	UserID     uuid.UUID      `json:"user_id" db:"user_id"`
	OpenToWork bool           `json:"open_to_work" db:"open_to_work"`
	CVPath     sql.NullString `json:"cv_path" db:"cv_path"`
	CreatedAt  time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at" db:"updated_at"`

	// This field would be populated by a custom JOIN query,
	// so it's ignored by default database operations.
	FieldsOfWork []FieldOfWork `json:"fields_of_work,omitempty" db:"-"`
}

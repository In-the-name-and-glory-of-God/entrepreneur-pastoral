package domain

import (
	"database/sql"
	"time"

	adminDomain "github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/admin/domain"
	"github.com/google/uuid"
)

// JobProfile corresponds to the "job_profiles" table.
type JobProfile struct {
	UserID     uuid.UUID      `json:"user_id" db:"user_id"`
	OpenToWork bool           `json:"open_to_work" db:"open_to_work"`
	CVPath     sql.NullString `json:"cv_path" db:"cv_path"`
	CreatedAt  time.Time      `json:"-" db:"created_at"`
	UpdatedAt  time.Time      `json:"-" db:"updated_at"`

	// This field would be populated by a custom JOIN query,
	// so it's ignored by default database operations.
	FieldsOfWork []adminDomain.FieldOfWork `json:"fields_of_work,omitempty" db:"-"`
}

// JobProfileFieldOfWork corresponds to the "job_profile_fields_of_work" junction table.
// This struct is mainly used for insert/delete operations on the many-to-many relationship.
type JobProfileFieldOfWork struct {
	UserID        uuid.UUID `db:"user_id"`
	FieldOfWorkID int16     `db:"field_of_work_id"`
}

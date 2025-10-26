package domain

import "github.com/google/uuid"

// JobProfileFieldOfWork corresponds to the "job_profile_fields_of_work" junction table.
// This struct is mainly used for insert/delete operations on the many-to-many relationship.
type JobProfileFieldOfWork struct {
	UserID        uuid.UUID `json:"user_id" db:"user_id"`
	FieldOfWorkID int16     `json:"field_of_work_id" db:"field_of_work_id"`
}

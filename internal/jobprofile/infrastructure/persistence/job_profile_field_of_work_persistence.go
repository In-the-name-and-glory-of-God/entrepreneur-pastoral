package persistence

import (
	"context"
	"fmt"

	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/job_profile/domain"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// JobProfileFieldOfWorkPersistence manages the junction table.
type JobProfileFieldOfWorkPersistence struct {
	db   *sqlx.DB
	psql sq.StatementBuilderType
}

// NewJobProfileFieldOfWorkPersistence creates a new persistence helper for the junction table.
func NewJobProfileFieldOfWorkPersistence(db *sqlx.DB) *JobProfileFieldOfWorkPersistence {
	return &JobProfileFieldOfWorkPersistence{
		db:   db,
		psql: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

// Create inserts a new link between a user and a field of work.
func (r *JobProfileFieldOfWorkPersistence) Create(ctx context.Context, jobProfileFieldOfWork *domain.JobProfileFieldOfWork) error {
	query, args, err := r.psql.Insert("job_profile_fields_of_work").
		Columns("user_id", "field_of_work_id").
		Values(jobProfileFieldOfWork.UserID, jobProfileFieldOfWork.FieldOfWorkID).
		// Add ON CONFLICT DO NOTHING to prevent errors on duplicate entries
		Suffix("ON CONFLICT (user_id, field_of_work_id) DO NOTHING").
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build create jobProfileFieldOfWork query: %w", err)
	}

	if _, err := r.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to execute create jobProfileFieldOfWork query: %w", err)
	}

	return nil
}

// Delete removes a link between a user and a field of work.
func (r *JobProfileFieldOfWorkPersistence) Delete(ctx context.Context, userID uuid.UUID, fieldOfWorkID int16) error {
	query, args, err := r.psql.Delete("job_profile_fields_of_work").
		Where(sq.Eq{
			"user_id":          userID,
			"field_of_work_id": fieldOfWorkID,
		}).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build delete jobProfileFieldOfWork query: %w", err)
	}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute delete jobProfileFieldOfWork query: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected on delete: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows were deleted")
	}

	return nil
}

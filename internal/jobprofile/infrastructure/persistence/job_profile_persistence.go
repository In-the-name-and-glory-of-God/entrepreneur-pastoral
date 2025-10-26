package persistence

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/job_profile/domain"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// JobProfilePersistence manages data access for the job_profiles table.
type JobProfilePersistence struct {
	db   *sqlx.DB
	psql sq.StatementBuilderType
}

// NewJobProfilePersistence creates a new JobProfilePersistence.
func NewJobProfilePersistence(db *sqlx.DB) *JobProfilePersistence {
	return &JobProfilePersistence{
		db:   db,
		psql: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

// Create inserts a new job profile.
func (r *JobProfilePersistence) Create(ctx context.Context, jobProfile *domain.JobProfile) error {
	query, args, err := r.psql.Insert("job_profiles").
		Columns("user_id", "open_to_work", "cv_path").
		Values(jobProfile.UserID, jobProfile.OpenToWork, jobProfile.CVPath).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build create jobProfile query: %w", err)
	}

	if _, err := r.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to execute create jobProfile query: %w", err)
	}

	return nil
}

// Update modifies an existing job profile.
func (r *JobProfilePersistence) Update(ctx context.Context, jobProfile *domain.JobProfile) error {
	query, args, err := r.psql.Update("job_profiles").
		Set("open_to_work", jobProfile.OpenToWork).
		Set("cv_path", jobProfile.CVPath).
		Where(sq.Eq{"user_id": jobProfile.UserID}).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build update jobProfile query: %w", err)
	}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute update jobProfile query: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected on update: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows were updated")
	}

	return nil
}

// Delete removes a job profile by its user ID.
func (r *JobProfilePersistence) Delete(ctx context.Context, userID uuid.UUID) error {
	query, args, err := r.psql.Delete("job_profiles").
		Where(sq.Eq{"user_id": userID}).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build delete jobProfile query: %w", err)
	}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute delete jobProfile query: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected on update: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows were deleted")
	}
	return nil
}

// GetByUserID retrieves a single job profile by its user ID.
func (r *JobProfilePersistence) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.JobProfile, error) {
	var profile domain.JobProfile
	query, args, err := r.psql.Select("*").From("job_profiles").
		Where(sq.Eq{"user_id": userID}).
		Limit(1).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build get jobProfile by userID query: %w", err)
	}

	if err := r.db.GetContext(ctx, &profile, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}

		return nil, fmt.Errorf("failed to execute get jobProfile by userID query: %w", err)
	}

	return &profile, nil
}

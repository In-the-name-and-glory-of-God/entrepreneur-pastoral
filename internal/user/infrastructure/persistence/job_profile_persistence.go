package persistence

import (
	"context"
	"database/sql"
	"fmt"

	adminDomain "github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/admin/domain"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/user/domain"
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
func (r *JobProfilePersistence) Create(tx *sqlx.Tx, jobProfile *domain.JobProfile) error {
	query, args, err := r.psql.Insert("job_profiles").
		Columns("user_id", "open_to_work", "cv_path").
		Values(jobProfile.UserID, jobProfile.OpenToWork, jobProfile.CVPath).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build create job profile query: %w", err)
	}

	if _, err := tx.Exec(query, args...); err != nil {
		return fmt.Errorf("failed to execute create job profile query: %w", err)
	}

	if err := r.addFieldsOfWork(tx, jobProfile.UserID, jobProfile.FieldsOfWork); err != nil {
		return err
	}

	return nil
}

// Update modifies an existing job profile.
func (r *JobProfilePersistence) Update(tx *sqlx.Tx, jobProfile *domain.JobProfile) error {
	query, args, err := r.psql.Update("job_profiles").
		Set("open_to_work", jobProfile.OpenToWork).
		Set("cv_path", jobProfile.CVPath).
		Where(sq.Eq{"user_id": jobProfile.UserID}).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build update job profile query: %w", err)
	}

	if _, err := tx.Exec(query, args...); err != nil {
		return fmt.Errorf("failed to execute update job profile query: %w", err)
	}

	if err := r.addFieldsOfWork(tx, jobProfile.UserID, jobProfile.FieldsOfWork); err != nil {
		return err
	}

	return nil
}

// GetByUserID retrieves a single job profile by its user ID along with its fields of work.
func (r *JobProfilePersistence) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.JobProfile, error) {
	var profile domain.JobProfile
	query, args, err := r.psql.Select("*").
		From("job_profiles").
		Where(sq.Eq{"user_id": userID}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build get job profile by userID query: %w", err)
	}

	if err := r.db.GetContext(ctx, &profile, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}

		return nil, fmt.Errorf("failed to execute get job profile by userID query: %w", err)
	}

	if profile.OpenToWork {
		fieldsOfWork, err := r.getAllFieldsOfWorkByUserID(ctx, userID)
		if err != nil {
			return nil, err
		}
		profile.FieldsOfWork = fieldsOfWork
	}

	return &profile, nil
}

// GetAllOpenToWork retrieves all job profiles where 'open_to_work' is true.
func (r *JobProfilePersistence) GetAllOpenToWork(ctx context.Context) ([]*domain.JobProfile, error) {
	var profiles []*domain.JobProfile
	query, args, err := r.psql.Select("*").From("job_profiles").
		Where(sq.Eq{"open_to_work": true}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build get all open to work job profiles query: %w", err)
	}

	if err := r.db.SelectContext(ctx, &profiles, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}

		return nil, fmt.Errorf("failed to execute get all open to work job profiles query: %w", err)
	}

	return profiles, nil
}

// getAllFieldsOfWorkByUserID retrieves all fields of work associated with a given user ID.
func (r *JobProfilePersistence) getAllFieldsOfWorkByUserID(ctx context.Context, userID uuid.UUID) ([]adminDomain.FieldOfWork, error) {
	var fieldsOfWork []adminDomain.FieldOfWork
	query, args, err := r.psql.
		Select("fow.id", "fow.name").
		From("job_profile_fields_of_work AS jpfow").
		InnerJoin("fields_of_work AS fow ON jpfow.field_of_work_id = fow.id").
		Where(sq.Eq{"jpfow.user_id": userID}).
		OrderBy("fow.id").
		Limit(3).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build get fields of work query: %w", err)
	}

	if err := r.db.SelectContext(ctx, &fieldsOfWork, query, args...); err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to execute get fields of work query: %w", err)
	}

	return fieldsOfWork, nil
}

// AddFieldOfWork inserts a new link between a user and a field of work.
func (r *JobProfilePersistence) addFieldsOfWork(tx *sqlx.Tx, userID uuid.UUID, fieldsOfWork []adminDomain.FieldOfWork) error {
	if len(fieldsOfWork) == 0 {
		return nil
	}

	if err := r.removeAllFieldsOfWork(tx, userID); err != nil {
		return err
	}

	builder := r.psql.Insert("job_profile_fields_of_work").
		Columns("user_id", "field_of_work_id")

	for _, f := range fieldsOfWork {
		builder = builder.Values(userID, f.ID)
	}

	query, args, err := builder.Suffix("ON CONFLICT (user_id, field_of_work_id) DO NOTHING").ToSql()
	if err != nil {
		return fmt.Errorf("failed to build create job profile field of work query: %w", err)
	}

	if _, err := tx.Exec(query, args...); err != nil {
		return fmt.Errorf("failed to execute create job profile field of work query: %w", err)
	}

	return nil
}

func (r *JobProfilePersistence) removeAllFieldsOfWork(tx *sqlx.Tx, userID uuid.UUID) error {
	query, args, err := r.psql.Delete("job_profile_fields_of_work").
		Where(sq.Eq{"user_id": userID}).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build delete all job profile field of work query: %w", err)
	}

	if _, err = tx.Exec(query, args...); err != nil {
		return fmt.Errorf("failed to execute delete all job profile field of work query: %w", err)
	}

	return nil
}

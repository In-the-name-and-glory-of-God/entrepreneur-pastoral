package persistence

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/entrepreneur/domain"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type JobPersistence struct {
	db   *sqlx.DB
	psql sq.StatementBuilderType
}

func NewJobPersistence(db *sqlx.DB) *JobPersistence {
	return &JobPersistence{
		db:   db,
		psql: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *JobPersistence) Create(tx *sqlx.Tx, job *domain.Job) error {
	query, args, err := r.psql.Insert("jobs").
		Columns(
			"business_id", "title", "description", "type", "location",
			"application_link", "is_open",
		).
		Values(
			job.BusinessID, job.Title, job.Description, job.Type, job.Location,
			job.ApplicationLink, job.IsOpen,
		).
		Suffix("RETURNING id, created_at").
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build create job query: %w", err)
	}

	if err := tx.QueryRowx(query, args...).Scan(&job.ID, &job.CreatedAt); err != nil {
		return fmt.Errorf("failed to execute create job query: %w", err)
	}

	return nil
}

func (r *JobPersistence) Update(tx *sqlx.Tx, job *domain.Job) error {
	query, args, err := r.psql.Update("jobs").
		Set("title", job.Title).
		Set("description", job.Description).
		Set("type", job.Type).
		Set("location", job.Location).
		Set("application_link", job.ApplicationLink).
		Set("is_open", job.IsOpen).
		Where(sq.Eq{"id": job.ID}).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build update job query: %w", err)
	}

	if _, err := tx.Exec(query, args...); err != nil {
		return fmt.Errorf("failed to execute update job query: %w", err)
	}

	return nil
}

func (r *JobPersistence) Delete(tx *sqlx.Tx, id uuid.UUID) error {
	query, args, err := r.psql.Delete("jobs").
		Where(sq.Eq{"id": id}).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build delete job query: %w", err)
	}

	if _, err := tx.Exec(query, args...); err != nil {
		return fmt.Errorf("failed to execute delete job query: %w", err)
	}

	return nil
}

func (r *JobPersistence) GetByID(ctx context.Context, id uuid.UUID) (*domain.Job, error) {
	var job domain.Job
	query, args, err := r.psql.Select("*").From("jobs").
		Where(sq.Eq{"id": id}).
		Limit(1).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build get job by id query: %w", err)
	}

	if err := r.db.GetContext(ctx, &job, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrJobNotFound
		}
		return nil, fmt.Errorf("failed to execute get job by id query: %w", err)
	}

	return &job, nil
}

func (r *JobPersistence) List(ctx context.Context, filter *domain.JobFilters) ([]*domain.Job, error) {
	queryBuilder := r.psql.Select("*").From("jobs")
	queryBuilder = r.buildFilterQuery(queryBuilder, filter)
	queryBuilder = queryBuilder.OrderBy("created_at DESC")

	if filter.Limit != nil {
		queryBuilder = queryBuilder.Limit(uint64(*filter.Limit))
	}
	if filter.Offset != nil {
		queryBuilder = queryBuilder.Offset(uint64(*filter.Offset))
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build list job query: %w", err)
	}

	var jobs []*domain.Job
	if err := r.db.SelectContext(ctx, &jobs, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrJobNotFound
		}
		return nil, fmt.Errorf("failed to execute list job query: %w", err)
	}

	return jobs, nil
}

func (r *JobPersistence) Count(ctx context.Context, filter *domain.JobFilters) (int, error) {
	queryBuilder := r.psql.Select("COUNT(*)").From("jobs")
	queryBuilder = r.buildFilterQuery(queryBuilder, filter)
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build count job query: %w", err)
	}

	var count int
	if err := r.db.GetContext(ctx, &count, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return 0, domain.ErrJobNotFound
		}
		return 0, fmt.Errorf("failed to execute count job query: %w", err)
	}

	return count, nil
}

func (r *JobPersistence) buildFilterQuery(baseQuery sq.SelectBuilder, filter *domain.JobFilters) sq.SelectBuilder {
	if filter.BusinessID != nil {
		baseQuery = baseQuery.Where(sq.Eq{"business_id": *filter.BusinessID})
	}
	if filter.Type != nil {
		baseQuery = baseQuery.Where(sq.Eq{"type": *filter.Type})
	}
	if filter.Location != nil {
		baseQuery = baseQuery.Where(sq.Eq{"location": *filter.Location})
	}
	if filter.IsOpen != nil {
		baseQuery = baseQuery.Where(sq.Eq{"is_open": *filter.IsOpen})
	}
	if filter.TitleContains != nil {
		baseQuery = baseQuery.Where(sq.Like{"title": fmt.Sprintf("%%%s%%", *filter.TitleContains)})
	}
	return baseQuery
}

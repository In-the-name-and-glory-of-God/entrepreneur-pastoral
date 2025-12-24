package persistence

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/admin/domain"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

// FieldOfWorkPersistence manages data access for the fields_of_work table.
type FieldOfWorkPersistence struct {
	db   *sqlx.DB
	psql sq.StatementBuilderType
}

// NewFieldOfWorkPersistence creates a new FieldOfWorkPersistence.
func NewFieldOfWorkPersistence(db *sqlx.DB) *FieldOfWorkPersistence {
	return &FieldOfWorkPersistence{
		db:   db,
		psql: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

// Create inserts a new field of work. The ID is SERIAL, so it's returned and set on the struct.
func (r *FieldOfWorkPersistence) Create(ctx context.Context, fieldOfWork *domain.FieldOfWork) error {
	query, args, err := r.psql.Insert("fields_of_work").
		Columns("key").
		Values(fieldOfWork.Key).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build create fieldOfWork query: %w", err)
	}

	if err := r.db.GetContext(ctx, &fieldOfWork.ID, query, args...); err != nil {
		return fmt.Errorf("failed to execute create fieldOfWork query: %w", err)
	}

	return nil
}

// Update modifies an existing field of work.
func (r *FieldOfWorkPersistence) Update(ctx context.Context, fieldOfWork *domain.FieldOfWork) error {
	query, args, err := r.psql.Update("fields_of_work").
		Set("key", fieldOfWork.Key).
		Where(sq.Eq{"id": fieldOfWork.ID}).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build update fieldOfWork query: %w", err)
	}

	if _, err := r.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to execute update fieldOfWork query: %w", err)
	}

	return nil
}

// Delete removes a field of work by its ID.
func (r *FieldOfWorkPersistence) Delete(ctx context.Context, id int16) error {
	query, args, err := r.psql.Delete("fields_of_work").
		Where(sq.Eq{"id": id}).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build delete field of work query: %w", err)
	}

	if _, err := r.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to execute delete field of work query: %w", err)
	}

	return nil
}

// GetAll retrieves all fields of work, ordered by key.
func (r *FieldOfWorkPersistence) GetAll(ctx context.Context) ([]*domain.FieldOfWork, error) {
	var fields []*domain.FieldOfWork
	query, args, err := r.psql.Select("*").From("fields_of_work").
		OrderBy("key ASC").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build get all fields of work query: %w", err)
	}

	if err := r.db.SelectContext(ctx, &fields, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrFieldOfWorkNotFound
		}

		return nil, fmt.Errorf("failed to execute get all fields of work query: %w", err)
	}

	return fields, nil
}

// GetByID retrieves a single field of work by its ID.
func (r *FieldOfWorkPersistence) GetByID(ctx context.Context, id int16) (*domain.FieldOfWork, error) {
	var field domain.FieldOfWork
	query, args, err := r.psql.Select("*").From("fields_of_work").
		Where(sq.Eq{"id": id}).
		Limit(1).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build get field of work by id query: %w", err)
	}

	if err := r.db.GetContext(ctx, &field, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrFieldOfWorkNotFound
		}

		return nil, fmt.Errorf("failed to execute get field of work by id query: %w", err)
	}

	return &field, nil
}

// GetByKey retrieves a single field of work by its key.
func (r *FieldOfWorkPersistence) GetByKey(ctx context.Context, key string) (*domain.FieldOfWork, error) {
	var field domain.FieldOfWork
	query, args, err := r.psql.Select("*").From("fields_of_work").
		Where(sq.Eq{"key": key}).
		Limit(1).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build get field of work by key query: %w", err)
	}

	if err := r.db.GetContext(ctx, &field, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrFieldOfWorkNotFound
		}

		return nil, fmt.Errorf("failed to execute get field of work by key query: %w", err)
	}

	return &field, nil
}

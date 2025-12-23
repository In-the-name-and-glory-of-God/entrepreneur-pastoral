package persistence

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/admin/domain"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// ChurchPersistence manages data access for the church table.
type ChurchPersistence struct {
	db   *sqlx.DB
	psql sq.StatementBuilderType
}

// NewChurchPersistence creates a new ChurchPersistence.
func NewChurchPersistence(db *sqlx.DB) *ChurchPersistence {
	return &ChurchPersistence{
		db:   db,
		psql: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

// UnitOfWork is a helper function that executes a given function within a database transaction.
func (r *ChurchPersistence) UnitOfWork(ctx context.Context, fn func(*sqlx.Tx) error) error {
	var err error

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		}
	}()

	if err = fn(tx); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Create inserts a new church. The ID is generated and returned.
func (r *ChurchPersistence) Create(tx *sqlx.Tx, church *domain.Church) error {
	query, args, err := r.psql.Insert("church").
		Columns("name", "diocese", "parish_number", "website_url", "phone_number", "address_id", "is_archdiocese", "is_active").
		Values(church.Name, church.Diocese, church.ParishNumber, church.WebsiteURL, church.PhoneNumber, church.AddressID, church.IsArchdiocese, church.IsActive).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build create church query: %w", err)
	}

	if err := tx.GetContext(context.Background(), &church.ID, query, args...); err != nil {
		return fmt.Errorf("failed to execute create church query: %w", err)
	}

	return nil
}

// Update modifies an existing church.
func (r *ChurchPersistence) Update(tx *sqlx.Tx, church *domain.Church) error {
	query, args, err := r.psql.Update("church").
		Set("name", church.Name).
		Set("diocese", church.Diocese).
		Set("parish_number", church.ParishNumber).
		Set("website_url", church.WebsiteURL).
		Set("phone_number", church.PhoneNumber).
		Set("address_id", church.AddressID).
		Set("is_archdiocese", church.IsArchdiocese).
		Set("is_active", church.IsActive).
		Where(sq.Eq{"id": church.ID}).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build update church query: %w", err)
	}

	if _, err := tx.ExecContext(context.Background(), query, args...); err != nil {
		return fmt.Errorf("failed to execute update church query: %w", err)
	}

	return nil
}

// Delete removes a church by its ID.
func (r *ChurchPersistence) Delete(tx *sqlx.Tx, id uuid.UUID) error {
	query, args, err := r.psql.Delete("church").
		Where(sq.Eq{"id": id}).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build delete church query: %w", err)
	}

	if _, err := tx.ExecContext(context.Background(), query, args...); err != nil {
		return fmt.Errorf("failed to execute delete church query: %w", err)
	}

	return nil
}

// GetByID retrieves a single church by its ID.
func (r *ChurchPersistence) GetByID(ctx context.Context, id uuid.UUID) (*domain.Church, error) {
	var church domain.Church
	query, args, err := r.psql.Select("*").From("church").
		Where(sq.Eq{"id": id}).
		Limit(1).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build get church by id query: %w", err)
	}

	if err := r.db.GetContext(ctx, &church, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}

		return nil, fmt.Errorf("failed to execute get church by id query: %w", err)
	}

	return &church, nil
}

// GetByName retrieves a single church by its name.
func (r *ChurchPersistence) GetByName(ctx context.Context, name string) (*domain.Church, error) {
	var church domain.Church
	query, args, err := r.psql.Select("*").From("church").
		Where(sq.Like{"name": fmt.Sprintf("%%%s%%", name)}).
		Limit(1).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build get church by name query: %w", err)
	}

	if err := r.db.GetContext(ctx, &church, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}

		return nil, fmt.Errorf("failed to execute get church by name query: %w", err)
	}

	return &church, nil
}

// List retrieves churches based on filters.
func (r *ChurchPersistence) List(ctx context.Context, filter *domain.ChurchFilters) ([]*domain.Church, error) {
	var churches []*domain.Church
	builder := r.psql.Select("*").From("church")

	if filter != nil {
		if filter.Diocese != nil {
			builder = builder.Where(sq.Eq{"diocese": *filter.Diocese})
		}
		if filter.AddressID != nil {
			builder = builder.Where(sq.Eq{"address_id": *filter.AddressID})
		}
		if filter.IsArchdiocese != nil {
			builder = builder.Where(sq.Eq{"is_archdiocese": *filter.IsArchdiocese})
		}
		if filter.IsActive != nil {
			builder = builder.Where(sq.Eq{"is_active": *filter.IsActive})
		}
		if filter.NameContains != nil {
			builder = builder.Where(sq.Like{"name": fmt.Sprintf("%%%s%%", *filter.NameContains)})
		}
		if filter.Limit != nil {
			builder = builder.Limit(uint64(*filter.Limit))
		}
		if filter.Offset != nil {
			builder = builder.Offset(uint64(*filter.Offset))
		}
	}

	builder = builder.OrderBy("name ASC")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build list churches query: %w", err)
	}

	if err := r.db.SelectContext(ctx, &churches, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}

		return nil, fmt.Errorf("failed to execute list churches query: %w", err)
	}

	return churches, nil
}

// Count returns the total number of churches matching the filter.
func (r *ChurchPersistence) Count(ctx context.Context, filter *domain.ChurchFilters) (int, error) {
	var count int
	builder := r.psql.Select("COUNT(*)").From("church")

	if filter != nil {
		if filter.Diocese != nil {
			builder = builder.Where(sq.Eq{"diocese": *filter.Diocese})
		}
		if filter.AddressID != nil {
			builder = builder.Where(sq.Eq{"address_id": *filter.AddressID})
		}
		if filter.IsArchdiocese != nil {
			builder = builder.Where(sq.Eq{"is_archdiocese": *filter.IsArchdiocese})
		}
		if filter.IsActive != nil {
			builder = builder.Where(sq.Eq{"is_active": *filter.IsActive})
		}
		if filter.NameContains != nil {
			builder = builder.Where(sq.Like{"name": fmt.Sprintf("%%%s%%", *filter.NameContains)})
		}
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build count churches query: %w", err)
	}

	if err := r.db.GetContext(ctx, &count, query, args...); err != nil {
		return 0, fmt.Errorf("failed to execute count churches query: %w", err)
	}

	return count, nil
}

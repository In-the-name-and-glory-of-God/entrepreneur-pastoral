package persistence

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/admin/domain"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type IndustryPersistence struct {
	db   *sqlx.DB
	psql sq.StatementBuilderType
}

func NewIndustryPersistence(db *sqlx.DB) *IndustryPersistence {
	return &IndustryPersistence{
		db:   db,
		psql: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *IndustryPersistence) Create(ctx context.Context, industry *domain.Industry) error {
	query, args, err := r.psql.Insert("industries").
		Columns("key").
		Values(industry.Key).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build create industry query: %w", err)
	}

	if err := r.db.QueryRowxContext(ctx, query, args...).Scan(&industry.ID); err != nil {
		return fmt.Errorf("failed to execute create industry query: %w", err)
	}

	return nil
}

func (r *IndustryPersistence) Update(ctx context.Context, industry *domain.Industry) error {
	query, args, err := r.psql.Update("industries").
		Set("key", industry.Key).
		Where(sq.Eq{"id": industry.ID}).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build update industry query: %w", err)
	}

	if _, err := r.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to execute update industry query: %w", err)
	}

	return nil
}

func (r *IndustryPersistence) Delete(ctx context.Context, id int16) error {
	query, args, err := r.psql.Delete("industries").
		Where(sq.Eq{"id": id}).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build delete industry query: %w", err)
	}

	if _, err := r.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to execute delete industry query: %w", err)
	}

	return nil
}

func (r *IndustryPersistence) GetAll(ctx context.Context) ([]*domain.Industry, error) {
	var industries []*domain.Industry
	query, args, err := r.psql.Select("*").From("industries").OrderBy("key ASC").ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build get all industries query: %w", err)
	}

	if err := r.db.SelectContext(ctx, &industries, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrIndustryNotFound
		}
		return nil, fmt.Errorf("failed to execute get all industries query: %w", err)
	}

	return industries, nil
}

func (r *IndustryPersistence) GetByID(ctx context.Context, id int16) (*domain.Industry, error) {
	var industry domain.Industry
	query, args, err := r.psql.Select("*").From("industries").Where(sq.Eq{"id": id}).Limit(1).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build get industry by id query: %w", err)
	}

	if err := r.db.GetContext(ctx, &industry, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrIndustryNotFound
		}
		return nil, fmt.Errorf("failed to execute get industry by id query: %w", err)
	}

	return &industry, nil
}

func (r *IndustryPersistence) GetByKey(ctx context.Context, key string) (*domain.Industry, error) {
	var industry domain.Industry
	query, args, err := r.psql.Select("*").From("industries").Where(sq.Eq{"key": key}).Limit(1).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build get industry by key query: %w", err)
	}

	if err := r.db.GetContext(ctx, &industry, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrIndustryNotFound
		}
		return nil, fmt.Errorf("failed to execute get industry by key query: %w", err)
	}

	return &industry, nil
}

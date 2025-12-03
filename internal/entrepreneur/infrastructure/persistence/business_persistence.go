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

type BusinessPersistence struct {
	db   *sqlx.DB
	psql sq.StatementBuilderType
}

func NewBusinessPersistence(db *sqlx.DB) *BusinessPersistence {
	return &BusinessPersistence{
		db:   db,
		psql: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *BusinessPersistence) Create(tx *sqlx.Tx, business *domain.Business) error {
	query, args, err := r.psql.Insert("business").
		Columns(
			"user_id", "industry_id", "name", "description", "email",
			"phone_country_code", "phone_number", "website_url", "logo_url", "is_active",
		).
		Values(
			business.UserID, business.IndustryID, business.Name, business.Description, business.Email,
			business.PhoneCountryCode, business.PhoneNumber, business.WebsiteURL, business.LogoURL, business.IsActive,
		).
		Suffix("RETURNING id, created_at").
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build create business query: %w", err)
	}

	if err := tx.QueryRowx(query, args...).Scan(&business.ID, &business.CreatedAt); err != nil {
		return fmt.Errorf("failed to execute create business query: %w", err)
	}

	return nil
}

func (r *BusinessPersistence) Update(tx *sqlx.Tx, business *domain.Business) error {
	query, args, err := r.psql.Update("business").
		Set("industry_id", business.IndustryID).
		Set("name", business.Name).
		Set("description", business.Description).
		Set("email", business.Email).
		Set("phone_country_code", business.PhoneCountryCode).
		Set("phone_number", business.PhoneNumber).
		Set("website_url", business.WebsiteURL).
		Set("logo_url", business.LogoURL).
		Set("is_active", business.IsActive).
		Where(sq.Eq{"id": business.ID}).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build update business query: %w", err)
	}

	if _, err := tx.Exec(query, args...); err != nil {
		return fmt.Errorf("failed to execute update business query: %w", err)
	}

	return nil
}

func (r *BusinessPersistence) Delete(tx *sqlx.Tx, id uuid.UUID) error {
	query, args, err := r.psql.Delete("business").
		Where(sq.Eq{"id": id}).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build delete business query: %w", err)
	}

	if _, err := tx.Exec(query, args...); err != nil {
		return fmt.Errorf("failed to execute delete business query: %w", err)
	}

	return nil
}

func (r *BusinessPersistence) GetByID(ctx context.Context, id uuid.UUID) (*domain.Business, error) {
	var business domain.Business
	query, args, err := r.psql.Select("*").From("business").
		Where(sq.Eq{"id": id}).
		Limit(1).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build get business by id query: %w", err)
	}

	if err := r.db.GetContext(ctx, &business, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrBusinessNotFound
		}
		return nil, fmt.Errorf("failed to execute get business by id query: %w", err)
	}

	return &business, nil
}

func (r *BusinessPersistence) List(ctx context.Context, filter *domain.BusinessFilters) ([]*domain.Business, error) {
	queryBuilder := r.psql.Select("*").From("business")
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
		return nil, fmt.Errorf("failed to build list business query: %w", err)
	}

	var businesses []*domain.Business
	if err := r.db.SelectContext(ctx, &businesses, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrBusinessNotFound
		}
		return nil, fmt.Errorf("failed to execute list business query: %w", err)
	}

	return businesses, nil
}

func (r *BusinessPersistence) Count(ctx context.Context, filter *domain.BusinessFilters) (int, error) {
	queryBuilder := r.psql.Select("COUNT(*)").From("business")
	queryBuilder = r.buildFilterQuery(queryBuilder, filter)
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build count business query: %w", err)
	}

	var count int
	if err := r.db.GetContext(ctx, &count, query, args...); err != nil {
		return 0, fmt.Errorf("failed to execute count business query: %w", err)
	}

	return count, nil
}

func (r *BusinessPersistence) buildFilterQuery(baseQuery sq.SelectBuilder, filter *domain.BusinessFilters) sq.SelectBuilder {
	if filter.UserID != nil {
		baseQuery = baseQuery.Where(sq.Eq{"user_id": *filter.UserID})
	}
	if filter.IndustryID != nil {
		baseQuery = baseQuery.Where(sq.Eq{"industry_id": *filter.IndustryID})
	}
	if filter.IsActive != nil {
		baseQuery = baseQuery.Where(sq.Eq{"is_active": *filter.IsActive})
	}
	if filter.NameContains != nil {
		baseQuery = baseQuery.Where(sq.Like{"name": fmt.Sprintf("%%%s%%", *filter.NameContains)})
	}
	return baseQuery
}

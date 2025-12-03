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

type ServicePersistence struct {
	db   *sqlx.DB
	psql sq.StatementBuilderType
}

func NewServicePersistence(db *sqlx.DB) *ServicePersistence {
	return &ServicePersistence{
		db:   db,
		psql: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *ServicePersistence) Create(tx *sqlx.Tx, service *domain.Service) error {
	query, args, err := r.psql.Insert("services").
		Columns(
			"business_id", "name", "description", "price",
		).
		Values(
			service.BusinessID, service.Name, service.Description, service.Price,
		).
		Suffix("RETURNING id, created_at").
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build create service query: %w", err)
	}

	if err := tx.QueryRowx(query, args...).Scan(&service.ID, &service.CreatedAt); err != nil {
		return fmt.Errorf("failed to execute create service query: %w", err)
	}

	return nil
}

func (r *ServicePersistence) Update(tx *sqlx.Tx, service *domain.Service) error {
	query, args, err := r.psql.Update("services").
		Set("name", service.Name).
		Set("description", service.Description).
		Set("price", service.Price).
		Where(sq.Eq{"id": service.ID}).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build update service query: %w", err)
	}

	if _, err := tx.Exec(query, args...); err != nil {
		return fmt.Errorf("failed to execute update service query: %w", err)
	}
	return nil
}

func (r *ServicePersistence) Delete(tx *sqlx.Tx, id uuid.UUID) error {
	query, args, err := r.psql.Delete("services").
		Where(sq.Eq{"id": id}).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build delete service query: %w", err)
	}

	if _, err := tx.Exec(query, args...); err != nil {
		return fmt.Errorf("failed to execute delete service query: %w", err)
	}

	return nil
}

func (r *ServicePersistence) GetByID(ctx context.Context, id uuid.UUID) (*domain.Service, error) {
	var service domain.Service
	query, args, err := r.psql.Select("*").From("services").
		Where(sq.Eq{"id": id}).
		Limit(1).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build get service by id query: %w", err)
	}

	if err := r.db.GetContext(ctx, &service, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrServiceNotFound
		}
		return nil, fmt.Errorf("failed to execute get service by id query: %w", err)
	}

	return &service, nil
}

func (r *ServicePersistence) List(ctx context.Context, filter *domain.ServiceFilters) ([]*domain.Service, error) {
	queryBuilder := r.psql.Select("*").From("services")
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
		return nil, fmt.Errorf("failed to build list service query: %w", err)
	}

	var services []*domain.Service
	if err := r.db.SelectContext(ctx, &services, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrServiceNotFound
		}
		return nil, fmt.Errorf("failed to execute list service query: %w", err)
	}

	return services, nil
}

func (r *ServicePersistence) Count(ctx context.Context, filter *domain.ServiceFilters) (int, error) {
	queryBuilder := r.psql.Select("COUNT(*)").From("services")
	queryBuilder = r.buildFilterQuery(queryBuilder, filter)
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build count service query: %w", err)
	}

	var count int
	if err := r.db.GetContext(ctx, &count, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return 0, domain.ErrServiceNotFound
		}
		return 0, fmt.Errorf("failed to execute count service query: %w", err)
	}

	return count, nil
}

func (r *ServicePersistence) buildFilterQuery(baseQuery sq.SelectBuilder, filter *domain.ServiceFilters) sq.SelectBuilder {
	if filter.BusinessID != nil {
		baseQuery = baseQuery.Where(sq.Eq{"business_id": *filter.BusinessID})
	}
	if filter.NameContains != nil {
		baseQuery = baseQuery.Where(sq.Like{"name": fmt.Sprintf("%%%s%%", *filter.NameContains)})
	}
	if filter.MinPrice != nil {
		baseQuery = baseQuery.Where(sq.GtOrEq{"price": *filter.MinPrice})
	}
	if filter.MaxPrice != nil {
		baseQuery = baseQuery.Where(sq.LtOrEq{"price": *filter.MaxPrice})
	}
	return baseQuery
}

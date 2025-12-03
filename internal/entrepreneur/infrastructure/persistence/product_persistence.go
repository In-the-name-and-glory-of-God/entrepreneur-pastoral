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

type ProductPersistence struct {
	db   *sqlx.DB
	psql sq.StatementBuilderType
}

func NewProductPersistence(db *sqlx.DB) *ProductPersistence {
	return &ProductPersistence{
		db:   db,
		psql: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *ProductPersistence) Create(tx *sqlx.Tx, product *domain.Product) error {
	query, args, err := r.psql.Insert("products").
		Columns(
			"business_id", "name", "description", "price", "image_url", "is_available",
		).
		Values(
			product.BusinessID, product.Name, product.Description, product.Price, product.ImageURL, product.IsAvailable,
		).
		Suffix("RETURNING id, created_at").
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build create product query: %w", err)
	}

	if err := tx.QueryRowx(query, args...).Scan(&product.ID, &product.CreatedAt); err != nil {
		return fmt.Errorf("failed to execute create product query: %w", err)
	}

	return nil
}

func (r *ProductPersistence) Update(tx *sqlx.Tx, product *domain.Product) error {
	query, args, err := r.psql.Update("products").
		Set("name", product.Name).
		Set("description", product.Description).
		Set("price", product.Price).
		Set("image_url", product.ImageURL).
		Set("is_available", product.IsAvailable).
		Where(sq.Eq{"id": product.ID}).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build update product query: %w", err)
	}

	if _, err := tx.Exec(query, args...); err != nil {
		return fmt.Errorf("failed to execute update product query: %w", err)
	}

	return nil
}

func (r *ProductPersistence) Delete(tx *sqlx.Tx, id uuid.UUID) error {
	query, args, err := r.psql.Delete("products").
		Where(sq.Eq{"id": id}).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build delete product query: %w", err)
	}

	if _, err := tx.Exec(query, args...); err != nil {
		return fmt.Errorf("failed to execute delete product query: %w", err)
	}

	return nil
}

func (r *ProductPersistence) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	var product domain.Product
	query, args, err := r.psql.Select("*").From("products").
		Where(sq.Eq{"id": id}).
		Limit(1).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build get product by id query: %w", err)
	}

	if err := r.db.GetContext(ctx, &product, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrProductNotFound
		}
		return nil, fmt.Errorf("failed to execute get product by id query: %w", err)
	}

	return &product, nil
}

func (r *ProductPersistence) List(ctx context.Context, filter *domain.ProductFilters) ([]*domain.Product, error) {
	queryBuilder := r.psql.Select("*").From("products")
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
		return nil, fmt.Errorf("failed to build list product query: %w", err)
	}

	var products []*domain.Product
	if err := r.db.SelectContext(ctx, &products, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrProductNotFound
		}
		return nil, fmt.Errorf("failed to execute list product query: %w", err)
	}

	return products, nil
}

func (r *ProductPersistence) Count(ctx context.Context, filter *domain.ProductFilters) (int, error) {
	queryBuilder := r.psql.Select("COUNT(*)").From("products")
	queryBuilder = r.buildFilterQuery(queryBuilder, filter)
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build count product query: %w", err)
	}

	var count int
	if err := r.db.GetContext(ctx, &count, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return 0, domain.ErrProductNotFound
		}
		return 0, fmt.Errorf("failed to execute count product query: %w", err)
	}

	return count, nil
}

func (r *ProductPersistence) buildFilterQuery(baseQuery sq.SelectBuilder, filter *domain.ProductFilters) sq.SelectBuilder {
	if filter.BusinessID != nil {
		baseQuery = baseQuery.Where(sq.Eq{"business_id": *filter.BusinessID})
	}
	if filter.IsAvailable != nil {
		baseQuery = baseQuery.Where(sq.Eq{"is_available": *filter.IsAvailable})
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

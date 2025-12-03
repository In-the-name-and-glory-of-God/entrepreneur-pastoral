package domain

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ProductRepository interface {
	Create(tx *sqlx.Tx, product *Product) error
	Update(tx *sqlx.Tx, product *Product) error
	Delete(tx *sqlx.Tx, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*Product, error)
	List(ctx context.Context, filter *ProductFilters) ([]*Product, error)
	Count(ctx context.Context, filter *ProductFilters) (int, error)
}

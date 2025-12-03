package domain

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type BusinessRepository interface {
	Create(tx *sqlx.Tx, business *Business) error
	Update(tx *sqlx.Tx, business *Business) error
	Delete(tx *sqlx.Tx, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*Business, error)
	List(ctx context.Context, filter *BusinessFilters) ([]*Business, error)
	Count(ctx context.Context, filter *BusinessFilters) (int, error)
}

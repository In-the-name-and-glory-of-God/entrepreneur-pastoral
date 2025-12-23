package domain

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ChurchRepository interface {
	UnitOfWork(ctx context.Context, fn func(*sqlx.Tx) error) error
	Create(tx *sqlx.Tx, church *Church) error
	Update(ctx context.Context, church *Church) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*Church, error)
	GetByName(ctx context.Context, name string) (*Church, error)
	List(ctx context.Context, filter *ChurchFilters) ([]*Church, error)
	Count(ctx context.Context, filter *ChurchFilters) (int, error)
}

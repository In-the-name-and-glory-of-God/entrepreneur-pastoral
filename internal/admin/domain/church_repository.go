package domain

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ChurchRepository interface {
	Create(tx *sqlx.Tx, church *Church) error
	Update(tx *sqlx.Tx, church *Church) error
	Delete(tx *sqlx.Tx, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*Church, error)
	GetByName(ctx context.Context, name string) (*Church, error)
	List(ctx context.Context, filter *ChurchFilters) ([]*Church, error)
	Count(ctx context.Context, filter *ChurchFilters) (int, error)
}

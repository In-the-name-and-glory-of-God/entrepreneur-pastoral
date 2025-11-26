package domain

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ServiceRepository interface {
	Create(tx *sqlx.Tx, service *Service) error
	Update(tx *sqlx.Tx, service *Service) error
	Delete(tx *sqlx.Tx, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*Service, error)
	Count(ctx context.Context, filter *ServiceFilters) (int, error)
	List(ctx context.Context, filter *ServiceFilters) ([]*Service, error)
}

package domain

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type AddressRepository interface {
	Create(tx *sqlx.Tx, address *Address) error
	Update(tx *sqlx.Tx, address *Address) error
	Delete(tx *sqlx.Tx, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*Address, error)
}

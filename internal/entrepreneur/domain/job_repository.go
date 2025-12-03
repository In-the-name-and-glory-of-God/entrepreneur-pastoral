package domain

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type JobRepository interface {
	Create(tx *sqlx.Tx, job *Job) error
	Update(tx *sqlx.Tx, job *Job) error
	Delete(tx *sqlx.Tx, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*Job, error)
	Count(ctx context.Context, filter *JobFilters) (int, error)
	List(ctx context.Context, filter *JobFilters) ([]*Job, error)
}

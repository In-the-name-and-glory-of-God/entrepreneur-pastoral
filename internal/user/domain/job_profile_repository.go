package domain

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type JobProfileRepository interface {
	Create(tx *sqlx.Tx, jobProfile *JobProfile) error
	Update(tx *sqlx.Tx, jobProfile *JobProfile) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*JobProfile, error)
	GetAllOpenToWork(ctx context.Context) ([]*JobProfile, error)
}

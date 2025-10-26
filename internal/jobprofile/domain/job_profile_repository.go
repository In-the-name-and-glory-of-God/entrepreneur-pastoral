package domain

import (
	"context"

	"github.com/google/uuid"
)

type JobProfileRepository interface {
	Create(ctx context.Context, jobProfile *JobProfile) error
	Update(ctx context.Context, jobProfile *JobProfile) error
	Delete(ctx context.Context, userID uuid.UUID) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*JobProfile, error)
}

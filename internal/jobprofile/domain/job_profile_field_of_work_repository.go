package domain

import (
	"context"

	"github.com/google/uuid"
)

type JobProfileFieldOfWorkRepository interface {
	Create(ctx context.Context, jobProfileFieldOfWork *JobProfileFieldOfWork) error
	Delete(ctx context.Context, userID uuid.UUID, fieldOfWorkID int16) error
}

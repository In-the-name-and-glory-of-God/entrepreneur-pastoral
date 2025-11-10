package domain

import "context"

type FieldOfWorkRepository interface {
	Create(ctx context.Context, fieldOfWork *FieldOfWork) error
	Update(ctx context.Context, fieldOfWork *FieldOfWork) error
	GetAll(ctx context.Context) ([]*FieldOfWork, error)
	GetByID(ctx context.Context, id int16) (*FieldOfWork, error)
	GetByName(ctx context.Context, name string) (*FieldOfWork, error)
}

package domain

import "context"

type FieldOfWorkRepository interface {
	Create(ctx context.Context, fieldOfWork *FieldOfWork) error
	Update(ctx context.Context, fieldOfWork *FieldOfWork) error
	Delete(ctx context.Context, id int16) error
	GetAll(ctx context.Context) ([]*FieldOfWork, error)
	GetByID(ctx context.Context, id int16) (*FieldOfWork, error)
	GetByName(ctx context.Context, name string) (*FieldOfWork, error)
}

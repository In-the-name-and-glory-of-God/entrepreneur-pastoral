package domain

import "context"

type IndustryRepository interface {
	Create(ctx context.Context, industry *Industry) error
	Update(ctx context.Context, industry *Industry) error
	Delete(ctx context.Context, id int16) error
	GetByID(ctx context.Context, id int16) (*Industry, error)
	GetAll(ctx context.Context) ([]*Industry, error)
	GetByName(ctx context.Context, name string) (*Industry, error)
}

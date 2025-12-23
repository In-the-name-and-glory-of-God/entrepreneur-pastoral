package domain

import "context"

type IndustryRepository interface {
	Create(ctx context.Context, industry *Industry) error
	Update(ctx context.Context, industry *Industry) error
	Delete(ctx context.Context, id int16) error
	GetAll(ctx context.Context) ([]*Industry, error)
	GetByID(ctx context.Context, id int16) (*Industry, error)
	GetByName(ctx context.Context, name string) (*Industry, error)
}

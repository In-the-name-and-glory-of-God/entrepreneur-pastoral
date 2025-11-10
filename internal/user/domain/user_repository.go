package domain

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	UnitOfWork(ctx context.Context, fn func(*sqlx.Tx) error) error
	Create(tx *sqlx.Tx, user *User) error
	Update(tx *sqlx.Tx, user *User) error
	UpdateProperty(ctx context.Context, id uuid.UUID, property UserProperty, value any) error
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByDocumentID(ctx context.Context, documentID string) (*User, error)
	GetAllByRoleID(ctx context.Context, roleID int16) ([]*User, error)
	GetAllByIsActive(ctx context.Context, isActive bool) ([]*User, error)
	GetAllByIsVerified(ctx context.Context, isVerified bool) ([]*User, error)
	GetAllByIsCatholic(ctx context.Context, isCatholic bool) ([]*User, error)
	GetAllByIsEntrepreneur(ctx context.Context, isEntrepreneur bool) ([]*User, error)
	Find(ctx context.Context, filter *UserFilters) ([]*User, error)
	Count(ctx context.Context, filter *UserFilters) (int, error)
}

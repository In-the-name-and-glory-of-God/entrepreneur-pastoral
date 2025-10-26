package domain

import (
	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByDocumentID(ctx context.Context, documentID string) (*User, error)
	GetAllByRoleID(ctx context.Context, roleID int16) ([]*User, error)
	GetAllByIsActive(ctx context.Context, isActive bool) ([]*User, error)
	GetAllByIsVerified(ctx context.Context, isVerified bool) ([]*User, error)
	GetAllByIsCatholic(ctx context.Context, isCatholic bool) ([]*User, error)
	GetAllByIsEntrepreneur(ctx context.Context, isEntrepreneur bool) ([]*User, error)
	Find(ctx context.Context, filter UserFilter) ([]*User, error)
	Count(ctx context.Context, filter UserFilter) (int, error)
}

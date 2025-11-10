package persistence

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/user/domain"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/helper/constants"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// UserPersistence struct holds the database connection (using sqlx)
// and the query builder configured for PostgreSQL.
type UserPersistence struct {
	db   *sqlx.DB
	psql sq.StatementBuilderType
}

// NewUserPersistence creates and returns a new UserPersistence struct.
func NewUserPersistence(db *sqlx.DB) *UserPersistence {
	return &UserPersistence{
		db:   db,
		psql: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

// UnitOfWork is a helper function that executes a given function within a database transaction.
// It handles transaction beginning, committing, and rolling back in case of errors or panics.
func (r *UserPersistence) UnitOfWork(ctx context.Context, fn func(*sqlx.Tx) error) error {
	var err error

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		}
	}()

	if err = fn(tx); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// --- Write Methods ---

// Create inserts a new user into the database.
func (r *UserPersistence) Create(tx *sqlx.Tx, user *domain.User) error {
	query, args, err := r.psql.Insert("users").
		Columns(
			"role_id", "first_name", "last_name", "email", "password",
			"document_id", "phone_country_code", "phone_number",
		).
		Values(
			constants.ROLE_USER, user.FirstName, user.LastName, user.Email, user.Password,
			user.DocumentID, user.PhoneCountryCode, user.PhoneNumber,
		).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build create user query: %w", err)
	}

	if err := tx.Get(&user.ID, query, args...); err != nil {
		return fmt.Errorf("failed to execute create user query: %w", err)
	}

	return nil
}

// Update modifies an existing user in the database.
func (r *UserPersistence) Update(tx *sqlx.Tx, user *domain.User) error {
	query, args, err := r.psql.Update("users").
		Set("role_id", user.RoleID).
		Set("first_name", user.FirstName).
		Set("last_name", user.LastName).
		Set("email", user.Email).
		Set("password", user.Password).
		Set("document_id", user.DocumentID).
		Set("phone_country_code", user.PhoneCountryCode).
		Set("phone_number", user.PhoneNumber).
		Where(sq.Eq{"id": user.ID}).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build update user query: %w", err)
	}

	if _, err := tx.Exec(query, args...); err != nil {
		return fmt.Errorf("failed to execute update user query: %w", err)
	}

	return nil
}

// UpdateProperty modifies a specific property from an existing user in the database.
func (r *UserPersistence) UpdateProperty(ctx context.Context, id uuid.UUID, property domain.UserProperty, value any) error {
	query, args, err := r.psql.Update("users").
		Set(string(property), value).
		Where(sq.Eq{"id": id}).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build update user by property query: %w", err)
	}

	if _, err := r.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to execute update user by property query: %w", err)
	}

	return nil
}

// --- Read Methods ---

// GetByID retrieves a single user by their ID (UUID).
func (r *UserPersistence) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var user domain.User
	query, args, err := r.psql.Select("*").From("users").
		Where(sq.Eq{"id": id}).
		Limit(1).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build get by id query: %w", err)
	}

	if err := r.db.GetContext(ctx, &user, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}

		return nil, fmt.Errorf("failed to execute get by id query: %w", err)
	}

	return &user, nil
}

// GetByEmail retrieves a single user by their email address.
func (r *UserPersistence) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	query, args, err := r.psql.Select("*").From("users").
		Where(sq.Eq{"email": email}).
		Limit(1).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build get by email query: %w", err)
	}

	if err := r.db.GetContext(ctx, &user, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}

		return nil, fmt.Errorf("failed to execute get by email query: %w", err)
	}

	return &user, nil
}

// GetByDocumentID retrieves a single user by their document ID.
func (r *UserPersistence) GetByDocumentID(ctx context.Context, documentID string) (*domain.User, error) {
	var user domain.User
	query, args, err := r.psql.Select("*").From("users").
		Where(sq.Eq{"document_id": documentID}).
		Limit(1).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build get by document id query: %w", err)
	}

	if err := r.db.GetContext(ctx, &user, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}

		return nil, fmt.Errorf("failed to execute get by document id query: %w", err)
	}

	return &user, nil
}

func (r *UserPersistence) Find(ctx context.Context, filter *domain.UserFilters) ([]*domain.User, error) {
	queryBuilder := r.psql.Select("*").From("users")
	queryBuilder = r.buildFilterQuery(queryBuilder, filter)
	queryBuilder = queryBuilder.OrderBy("created_at DESC")

	if filter.Limit != nil {
		queryBuilder = queryBuilder.Limit(uint64(*filter.Limit))
	}
	if filter.Offset != nil {
		queryBuilder = queryBuilder.Offset(uint64(*filter.Offset))
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build find query: %w", err)
	}

	var users []*domain.User
	if err := r.db.SelectContext(ctx, &users, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}

		return nil, fmt.Errorf("failed to execute find query: %w", err)
	}

	return users, nil
}

func (r *UserPersistence) Count(ctx context.Context, filter *domain.UserFilters) (int, error) {
	queryBuilder := r.psql.Select("COUNT(*)").From("users")
	queryBuilder = r.buildFilterQuery(queryBuilder, filter)
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build count query: %w", err)
	}

	var count int
	if err := r.db.GetContext(ctx, &count, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return 0, domain.ErrUserNotFound
		}

		return 0, fmt.Errorf("failed to execute count query: %w", err)
	}

	return count, nil
}

// --- Read Multiple Methods ---

// GetAllByRoleID retrieves all users matching a specific role ID.
func (r *UserPersistence) GetAllByRoleID(ctx context.Context, roleID int16) ([]*domain.User, error) {
	return r.getAllBy(ctx, sq.Eq{"role_id": roleID})
}

// GetAllByIsActive retrieves all users based on their active status.
func (r *UserPersistence) GetAllByIsActive(ctx context.Context, isActive bool) ([]*domain.User, error) {
	return r.getAllBy(ctx, sq.Eq{"is_active": isActive})
}

// GetAllByIsVerified retrieves all users based on their verified status.
func (r *UserPersistence) GetAllByIsVerified(ctx context.Context, isVerified bool) ([]*domain.User, error) {
	return r.getAllBy(ctx, sq.Eq{"is_verified": isVerified})
}

// GetAllByIsCatholic retrieves all users based on their catholic status.
func (r *UserPersistence) GetAllByIsCatholic(ctx context.Context, isCatholic bool) ([]*domain.User, error) {
	return r.getAllBy(ctx, sq.Eq{"is_catholic": isCatholic})
}

// GetAllByIsEntrepreneur retrieves all users based on their entrepreneur status.
func (r *UserPersistence) GetAllByIsEntrepreneur(ctx context.Context, isEntrepreneur bool) ([]*domain.User, error) {
	return r.getAllBy(ctx, sq.Eq{"is_entrepreneur": isEntrepreneur})
}

// helper function to get multiple users based on a condition
func (r *UserPersistence) getAllBy(ctx context.Context, condition any) ([]*domain.User, error) {
	var users []*domain.User
	query, args, err := r.psql.Select("*").From("users").
		Where(condition).
		OrderBy("created_at DESC").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build get all by query: %w", err)
	}

	if err := r.db.SelectContext(ctx, &users, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}

		return nil, fmt.Errorf("failed to execute get all by query: %w", err)
	}

	return users, nil
}

// buildFilterQuery is a helper to apply filter logic for both Find and Count
func (r *UserPersistence) buildFilterQuery(baseQuery sq.SelectBuilder, filter *domain.UserFilters) sq.SelectBuilder {
	// Apply filters
	if filter.RoleID != nil {
		baseQuery = baseQuery.Where(sq.Eq{"role_id": *filter.RoleID})
	}
	if filter.IsActive != nil {
		baseQuery = baseQuery.Where(sq.Eq{"is_active": *filter.IsActive})
	}
	if filter.IsVerified != nil {
		baseQuery = baseQuery.Where(sq.Eq{"is_verified": *filter.IsVerified})
	}
	if filter.IsCatholic != nil {
		baseQuery = baseQuery.Where(sq.Eq{"is_catholic": *filter.IsCatholic})
	}
	if filter.IsEntrepreneur != nil {
		baseQuery = baseQuery.Where(sq.Eq{"is_entrepreneur": *filter.IsEntrepreneur})
	}
	if filter.EmailContains != nil {
		baseQuery = baseQuery.Where(sq.Like{"email": fmt.Sprintf("%%%s%%", *filter.EmailContains)})
	}
	if filter.NameContains != nil {
		nameClause := sq.Or{
			sq.Like{"first_name": fmt.Sprintf("%%%s%%", *filter.NameContains)},
			sq.Like{"last_name": fmt.Sprintf("%%%s%%", *filter.NameContains)},
		}
		baseQuery = baseQuery.Where(nameClause)
	}

	return baseQuery
}

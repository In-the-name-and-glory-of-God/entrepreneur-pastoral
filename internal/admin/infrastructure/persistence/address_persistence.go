package persistence

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/admin/domain"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// AddressPersistence manages data access for the address table.
type AddressPersistence struct {
	db   *sqlx.DB
	psql sq.StatementBuilderType
}

// NewAddressPersistence creates a new AddressPersistence.
func NewAddressPersistence(db *sqlx.DB) *AddressPersistence {
	return &AddressPersistence{
		db:   db,
		psql: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

// Create inserts a new address. The ID is generated and returned.
func (r *AddressPersistence) Create(tx *sqlx.Tx, address *domain.Address) error {
	query, args, err := r.psql.Insert("address").
		Columns("street_line_1", "street_line_2", "city", "state_province", "postal_code", "country").
		Values(address.StreetLine1, address.StreetLine2, address.City, address.StateProvince, address.PostalCode, address.Country).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build create address query: %w", err)
	}

	if err := tx.GetContext(context.Background(), &address.ID, query, args...); err != nil {
		return fmt.Errorf("failed to execute create address query: %w", err)
	}

	return nil
}

// CreateWithContext inserts a new address without a transaction. The ID is generated and returned.
func (r *AddressPersistence) CreateWithContext(ctx context.Context, address *domain.Address) error {
	query, args, err := r.psql.Insert("address").
		Columns("street_line_1", "street_line_2", "city", "state_province", "postal_code", "country").
		Values(address.StreetLine1, address.StreetLine2, address.City, address.StateProvince, address.PostalCode, address.Country).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build create address query: %w", err)
	}

	if err := r.db.GetContext(ctx, &address.ID, query, args...); err != nil {
		return fmt.Errorf("failed to execute create address query: %w", err)
	}

	return nil
}

// Update modifies an existing address.
func (r *AddressPersistence) Update(ctx context.Context, address *domain.Address) error {
	query, args, err := r.psql.Update("address").
		Set("street_line_1", address.StreetLine1).
		Set("street_line_2", address.StreetLine2).
		Set("city", address.City).
		Set("state_province", address.StateProvince).
		Set("postal_code", address.PostalCode).
		Set("country", address.Country).
		Where(sq.Eq{"id": address.ID}).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build update address query: %w", err)
	}

	if _, err := r.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to execute update address query: %w", err)
	}

	return nil
}

// Delete removes an address by its ID.
func (r *AddressPersistence) Delete(ctx context.Context, id uuid.UUID) error {
	query, args, err := r.psql.Delete("address").
		Where(sq.Eq{"id": id}).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build delete address query: %w", err)
	}

	if _, err := r.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to execute delete address query: %w", err)
	}

	return nil
}

// GetByID retrieves a single address by its ID.
func (r *AddressPersistence) GetByID(ctx context.Context, id uuid.UUID) (*domain.Address, error) {
	var address domain.Address
	query, args, err := r.psql.Select("*").From("address").
		Where(sq.Eq{"id": id}).
		Limit(1).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build get address by id query: %w", err)
	}

	if err := r.db.GetContext(ctx, &address, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrAddressNotFound
		}

		return nil, fmt.Errorf("failed to execute get address by id query: %w", err)
	}

	return &address, nil
}

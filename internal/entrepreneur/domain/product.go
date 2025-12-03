package domain

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Product corresponds to the "products" table.
type Product struct {
	ID          uuid.UUID      `json:"id" db:"id"`
	BusinessID  uuid.UUID      `json:"business_id" db:"business_id"`
	Name        string         `json:"name" db:"name"`
	Description string         `json:"description" db:"description"`
	Price       float64        `json:"price" db:"price"`
	ImageURL    sql.NullString `json:"image_url" db:"image_url"`
	IsAvailable bool           `json:"is_available" db:"is_available"`
	CreatedAt   time.Time      `json:"created_at" db:"created_at"`
}

// ProductFilters defines criteria for filtering products.
type ProductFilters struct {
	BusinessID   *uuid.UUID `json:"business_id"`
	IsAvailable  *bool      `json:"is_available"`
	NameContains *string    `json:"name_contains"`
	MinPrice     *float64   `json:"min_price"`
	MaxPrice     *float64   `json:"max_price"`

	// Pagination
	Limit  *int `json:"limit"`
	Offset *int `json:"offset"`
}

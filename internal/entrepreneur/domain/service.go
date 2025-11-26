package domain

import (
	"time"

	"github.com/google/uuid"
)

// Service corresponds to the "services" table.
type Service struct {
	ID          uuid.UUID `json:"id" db:"id"`
	BusinessID  uuid.UUID `json:"business_id" db:"business_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Price       float64   `json:"price" db:"price"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// ServiceFilters defines criteria for filtering services.
type ServiceFilters struct {
	BusinessID   *uuid.UUID `json:"business_id"`
	NameContains *string    `json:"name_contains"`
	MinPrice     *float64   `json:"min_price"`
	MaxPrice     *float64   `json:"max_price"`

	// Pagination
	Limit  *int `json:"limit"`
	Offset *int `json:"offset"`
}

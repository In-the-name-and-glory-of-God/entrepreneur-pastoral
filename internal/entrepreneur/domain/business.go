package domain

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// BusinessProperty used for dynamic update
type BusinessProperty string

const (
	BusinessIsActive BusinessProperty = "is_active"
)

// Business corresponds to the "business" table.
type Business struct {
	ID               uuid.UUID      `json:"id" db:"id"`
	UserID           uuid.UUID      `json:"user_id" db:"user_id"`
	IndustryID       int16          `json:"industry_id" db:"industry_id"`
	Name             string         `json:"name" db:"name"`
	Description      string         `json:"description" db:"description"`
	Email            string         `json:"email" db:"email"`
	PhoneCountryCode sql.NullString `json:"phone_country_code" db:"phone_country_code"`
	PhoneNumber      sql.NullString `json:"phone_number" db:"phone_number"`
	WebsiteURL       sql.NullString `json:"website_url" db:"website_url"`
	LogoURL          sql.NullString `json:"logo_url" db:"logo_url"`
	IsActive         bool           `json:"is_active" db:"is_active"`
	CreatedAt        time.Time      `json:"created_at" db:"created_at"`
}

// BusinessFilters defines criteria for filtering businesses.
type BusinessFilters struct {
	UserID       *uuid.UUID `json:"user_id,omitempty"`
	IndustryID   *int16     `json:"industry_id"`
	IsActive     *bool      `json:"is_active"`
	NameContains *string    `json:"name_contains"`

	// Pagination
	Limit  *int `json:"limit"`
	Offset *int `json:"offset"`
}

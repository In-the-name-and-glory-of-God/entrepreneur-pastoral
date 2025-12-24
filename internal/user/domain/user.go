package domain

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Properties used for dynamic update
type UserProperty string

const (
	RoleID         UserProperty = "role_id"
	Password       UserProperty = "password"
	IsActive       UserProperty = "is_active"
	IsVerified     UserProperty = "is_verified"
	IsCatholic     UserProperty = "is_catholic"
	IsEntrepreneur UserProperty = "is_entrepreneur"
)

// Role corresponds to the "roles" table.
type Role struct {
	ID          int16          `json:"id" db:"id"`
	Name        string         `json:"name" db:"name"`
	Description sql.NullString `json:"description" db:"description"`
}

// User corresponds to the "users" table.
type User struct {
	ID               uuid.UUID      `json:"id" db:"id"`
	RoleID           int16          `json:"role_id" db:"role_id"`
	AddressID        uuid.UUID      `json:"address_id" db:"address_id"`
	ChurchID         uuid.UUID      `json:"church_id" db:"church_id"`
	FirstName        string         `json:"first_name" db:"first_name"`
	LastName         string         `json:"last_name" db:"last_name"`
	Email            string         `json:"email" db:"email"`
	Password         []byte         `json:"-" db:"password"`
	DocumentID       string         `json:"document_id" db:"document_id"`
	PhoneCountryCode sql.NullString `json:"phone_country_code" db:"phone_country_code"`
	PhoneNumber      sql.NullString `json:"phone_number" db:"phone_number"`
	Language         sql.NullString `json:"language" db:"language"`
	IsActive         bool           `json:"is_active" db:"is_active"`
	IsVerified       bool           `json:"is_verified" db:"is_verified"`
	IsCatholic       bool           `json:"is_catholic" db:"is_catholic"`
	IsEntrepreneur   bool           `json:"is_entrepreneur" db:"is_entrepreneur"`
	CreatedAt        time.Time      `json:"-" db:"created_at"`
	UpdatedAt        time.Time      `json:"-" db:"updated_at"`
}

// UserFilters defines the structured criteria for filtering users.
// Using pointers allows us to check if a filter was provided (non-nil)
// or not (nil), which is crucial for boolean and numeric zero-values.
type UserFilters struct {
	RoleID         *int16  `json:"role_id"`
	IsActive       *bool   `json:"is_active"`
	IsVerified     *bool   `json:"is_verified"`
	IsCatholic     *bool   `json:"is_catholic"`
	IsEntrepreneur *bool   `json:"is_entrepreneur"`
	EmailContains  *string `json:"email_contains"` // For LIKE queries
	NameContains   *string `json:"name_contains"`  // For LIKE queries on first/last name

	// Pagination
	Limit  *int `json:"limit"`
	Offset *int `json:"offset"`
}

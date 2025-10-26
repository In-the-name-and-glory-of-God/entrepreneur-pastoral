package domain

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const (
	ErrUserNotFound = "user not found"
	ErrUserExists   = "user already exists"
	ErrInvalidEmail = "invalid email"
	ErrInvalidPhone = "invalid phone number"
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
	FirstName        string         `json:"first_name" db:"first_name"`
	LastName         string         `json:"last_name" db:"last_name"`
	Email            string         `json:"email" db:"email"`
	Password         []byte         `json:"-" db:"password"`
	DocumentID       string         `json:"document_id" db:"document_id"`
	PhoneCountryCode sql.NullString `json:"phone_country_code" db:"phone_country_code"`
	PhoneNumber      sql.NullString `json:"phone_number" db:"phone_number"`
	IsActive         bool           `json:"is_active" db:"is_active"`
	IsVerified       bool           `json:"is_verified" db:"is_verified"`
	IsCatholic       bool           `json:"is_catholic" db:"is_catholic"`
	IsEntrepreneur   bool           `json:"is_entrepreneur" db:"is_entrepreneur"`
	CreatedAt        time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at" db:"updated_at"`
}

// UserFilter defines the structured criteria for filtering users.
// Using pointers allows us to check if a filter was provided (non-nil)
// or not (nil), which is crucial for boolean and numeric zero-values.
type UserFilter struct {
	RoleID         *int16  `json:"role_id,omitempty"`
	IsActive       *bool   `json:"is_active,omitempty"`
	IsVerified     *bool   `json:"is_verified,omitempty"`
	IsCatholic     *bool   `json:"is_catholic,omitempty"`
	IsEntrepreneur *bool   `json:"is_entrepreneur,omitempty"`
	EmailContains  *string `json:"email_contains,omitempty"` // For LIKE queries
	NameContains   *string `json:"name_contains,omitempty"`  // For LIKE queries on first/last name

	// Pagination
	Limit  *int `json:"limit,omitempty"`
	Offset *int `json:"offset,omitempty"`
}

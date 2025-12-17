package domain

import (
	"database/sql"

	"github.com/google/uuid"
)

// Church corresponds to the "church" table.
type Church struct {
	ID            uuid.UUID      `json:"id" db:"id"`
	Name          string         `json:"name" db:"name"`
	Diocese       string         `json:"diocese" db:"diocese"`
	ParishNumber  sql.NullString `json:"parish_number" db:"parish_number"`
	WebsiteURL    sql.NullString `json:"website_url" db:"website_url"`
	PhoneNumber   sql.NullString `json:"phone_number" db:"phone_number"`
	AddressID     uuid.UUID      `json:"address_id" db:"address_id"`
	IsArchdiocese bool           `json:"is_archdiocese" db:"is_archdiocese"`
	IsActive      bool           `json:"is_active" db:"is_active"`
}

// ChurchFilters defines criteria for filtering churches.
type ChurchFilters struct {
	Diocese       *string    `json:"diocese,omitempty"`
	AddressID     *uuid.UUID `json:"address_id,omitempty"`
	IsArchdiocese *bool      `json:"is_archdiocese,omitempty"`
	IsActive      *bool      `json:"is_active,omitempty"`
	NameContains  *string    `json:"name_contains,omitempty"`

	// Pagination
	Limit  *int `json:"limit,omitempty"`
	Offset *int `json:"offset,omitempty"`
}

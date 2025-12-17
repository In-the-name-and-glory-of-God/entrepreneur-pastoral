package domain

import (
	"database/sql"

	"github.com/google/uuid"
)

// Address corresponds to the "address" table.
type Address struct {
	ID            uuid.UUID      `json:"id" db:"id"`
	StreetLine1   string         `json:"street_line_1" db:"street_line_1"`
	StreetLine2   sql.NullString `json:"street_line_2" db:"street_line_2"`
	City          string         `json:"city" db:"city"`
	StateProvince string         `json:"state_province" db:"state_province"`
	PostalCode    string         `json:"postal_code" db:"postal_code"`
	Country       string         `json:"country" db:"country"`
}

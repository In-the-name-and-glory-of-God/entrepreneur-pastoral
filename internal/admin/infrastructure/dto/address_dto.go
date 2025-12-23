package dto

import (
	"github.com/google/uuid"
)

type AddressCreateRequest struct {
	StreetLine1   string `json:"street_line_1"`
	StreetLine2   string `json:"street_line_2"`
	City          string `json:"city"`
	StateProvince string `json:"state_province"`
	PostalCode    string `json:"postal_code"`
	Country       string `json:"country"`
}

type AddressUpdateRequest struct {
	ID            uuid.UUID `json:"id"`
	StreetLine1   string    `json:"street_line_1"`
	StreetLine2   string    `json:"street_line_2"`
	City          string    `json:"city"`
	StateProvince string    `json:"state_province"`
	PostalCode    string    `json:"postal_code"`
	Country       string    `json:"country"`
}

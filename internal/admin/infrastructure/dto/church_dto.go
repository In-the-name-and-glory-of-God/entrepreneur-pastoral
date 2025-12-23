package dto

import (
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/admin/domain"
	"github.com/google/uuid"
)

type ChurchCreateRequest struct {
	Name          string               `json:"name"`
	Diocese       string               `json:"diocese"`
	ParishNumber  string               `json:"parish_number"`
	WebsiteURL    string               `json:"website_url"`
	PhoneNumber   string               `json:"phone_number"`
	Address       AddressCreateRequest `json:"address"`
	IsArchdiocese bool                 `json:"is_archdiocese"`
}

type ChurchUpdateRequest struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Diocese       string    `json:"diocese"`
	ParishNumber  string    `json:"parish_number"`
	WebsiteURL    string    `json:"website_url"`
	PhoneNumber   string    `json:"phone_number"`
	AddressID     uuid.UUID `json:"address_id"`
	IsArchdiocese bool      `json:"is_archdiocese"`
	IsActive      bool      `json:"is_active"`
}

type ChurchListRequest = domain.ChurchFilters

type ChurchListResponse struct {
	Churches []*domain.Church `json:"churches"`
	Count    int              `json:"count"`
	Limit    *int             `json:"limit"`
	Offset   *int             `json:"offset"`
}

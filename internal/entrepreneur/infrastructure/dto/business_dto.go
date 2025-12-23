package dto

import (
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/entrepreneur/domain"
	"github.com/google/uuid"
)

type BusinessCreateRequest struct {
	UserID           uuid.UUID
	IndustryID       int16  `json:"industry_id"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	Email            string `json:"email"`
	PhoneCountryCode string `json:"phone_country_code"`
	PhoneNumber      string `json:"phone_number"`
	WebsiteURL       string `json:"website_url"`
	LogoURL          string `json:"logo_url"`
}

type BusinessUpdateRequest struct {
	ID               uuid.UUID `json:"id"`
	IndustryID       int16     `json:"industry_id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	Email            string    `json:"email"`
	PhoneCountryCode string    `json:"phone_country_code"`
	PhoneNumber      string    `json:"phone_number"`
	WebsiteURL       string    `json:"website_url"`
	LogoURL          string    `json:"logo_url"`
	IsActive         bool      `json:"is_active"`
}

type BusinessListRequest = domain.BusinessFilters

type BusinessListResponse struct {
	Businesses []*domain.Business `json:"businesses"`
	Count      int                `json:"count"`
	Limit      *int               `json:"limit"`
	Offset     *int               `json:"offset"`
}

type BusinessUpdatePropertyRequest struct {
	ID    uuid.UUID `json:"-"`
	Value bool      `json:"value"`
}

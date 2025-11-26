package dto

import (
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/entrepreneur/domain"
	"github.com/google/uuid"
)

type ProductCreateRequest struct {
	BusinessID  uuid.UUID `json:"business_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	ImageURL    string    `json:"image_url"`
	IsAvailable bool      `json:"is_available"`
}

type ProductUpdateRequest struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	ImageURL    string    `json:"image_url"`
	IsAvailable bool      `json:"is_available"`
}

type ProductListRequest = domain.ProductFilters

type ProductListResponse struct {
	Products []*domain.Product `json:"products"`
	Count    int               `json:"count"`
	Limit    *int              `json:"limit"`
	Offset   *int              `json:"offset"`
}

package dto

import (
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/entrepreneur/domain"
	"github.com/google/uuid"
)

type ServiceCreateRequest struct {
	BusinessID  uuid.UUID `json:"business_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
}

type ServiceUpdateRequest struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
}

type ServiceListRequest = domain.ServiceFilters

type ServiceListResponse struct {
	Services []*domain.Service `json:"services"`
	Count    int               `json:"count"`
	Limit    *int              `json:"limit"`
	Offset   *int              `json:"offset"`
}

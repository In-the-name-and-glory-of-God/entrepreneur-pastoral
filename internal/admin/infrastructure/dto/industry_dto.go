package dto

import "github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/admin/domain"

type IndustryCreateRequest struct {
	Name string `json:"name"`
}

type IndustryUpdateRequest struct {
	ID   int16  `json:"id"`
	Name string `json:"name"`
}

type IndustryListResponse struct {
	Industries []*domain.Industry `json:"industries"`
}

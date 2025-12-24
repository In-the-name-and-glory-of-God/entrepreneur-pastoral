package dto

import "github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/admin/domain"

type IndustryCreateRequest struct {
	Key string `json:"key"`
}

type IndustryUpdateRequest struct {
	ID  int16  `json:"id"`
	Key string `json:"key"`
}

type IndustryListResponse struct {
	Industries []*domain.Industry `json:"industries"`
}

package dto

import "github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/entrepreneur/domain"

type IndustryListResponse struct {
	Industries []*domain.Industry `json:"industries"`
}

package dto

import (
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/entrepreneur/domain"
	"github.com/google/uuid"
)

type JobCreateRequest struct {
	BusinessID      uuid.UUID          `json:"business_id"`
	Title           string             `json:"title"`
	Description     string             `json:"description"`
	Type            domain.JobType     `json:"type"`
	Location        domain.JobLocation `json:"location"`
	ApplicationLink string             `json:"application_link"`
	IsOpen          bool               `json:"is_open"`
}

type JobUpdateRequest struct {
	ID              uuid.UUID          `json:"id"`
	Title           string             `json:"title"`
	Description     string             `json:"description"`
	Type            domain.JobType     `json:"type"`
	Location        domain.JobLocation `json:"location"`
	ApplicationLink string             `json:"application_link"`
	IsOpen          bool               `json:"is_open"`
}

type JobListRequest = domain.JobFilters

type JobListResponse struct {
	Jobs   []*domain.Job `json:"jobs"`
	Count  int           `json:"count"`
	Limit  *int          `json:"limit"`
	Offset *int          `json:"offset"`
}

package dto

import "github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/admin/domain"

type FieldOfWorkCreateRequest struct {
	Name string `json:"name"`
}

type FieldOfWorkUpdateRequest struct {
	ID   int16  `json:"id"`
	Name string `json:"name"`
}

type FieldOfWorkListResponse struct {
	FieldsOfWork []*domain.FieldOfWork `json:"fields_of_work"`
}

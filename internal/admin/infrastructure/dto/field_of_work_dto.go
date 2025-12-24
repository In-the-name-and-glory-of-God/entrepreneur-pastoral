package dto

import "github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/admin/domain"

type FieldOfWorkCreateRequest struct {
	Key string `json:"key"`
}

type FieldOfWorkUpdateRequest struct {
	ID  int16  `json:"id"`
	Key string `json:"key"`
}

type FieldOfWorkListResponse struct {
	FieldsOfWork []*domain.FieldOfWork `json:"fields_of_work"`
}

package dto

import (
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/user/domain"
	"github.com/google/uuid"
)

type UserRegisterRequest struct {
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	Email            string `json:"email"`
	Password         string `json:"password"`
	DocumentID       string `json:"document_id"`
	PhoneCountryCode string `json:"phone_country_code"`
	PhoneNumber      string `json:"phone_number"`
	// JobProfile
	OpenToWork   bool                 `json:"open_to_work"`
	CVPath       string               `json:"cv_path"`
	FieldsOfWork []domain.FieldOfWork `json:"fields_of_work"`
}

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginResponse struct {
	Token string `json:"token"`
}

type UserResetPasswordRequest struct {
	ID          uuid.UUID `json:"id"`
	Email       string    `json:"email,omitempty"`
	NewPassword string    `json:"new_password,omitempty"`
}

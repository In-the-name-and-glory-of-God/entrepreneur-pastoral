package dto

import (
	adminDomain "github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/admin/domain"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/user/domain"
	"github.com/google/uuid"
)

type UserAsContext struct {
	ID             uuid.UUID `json:"id"`
	Email          string    `json:"email"`
	RoleID         int16     `json:"role_id"`
	Language       string    `json:"language"`
	IsCatholic     bool      `json:"is_catholic"`
	IsEntrepreneur bool      `json:"is_entrepreneur"`
}

type UserGetResponse struct {
	User                    *domain.User                    `json:"user"`
	NotificationPreferences *domain.NotificationPreferences `json:"notification_preferences"`
	JobProfile              *domain.JobProfile              `json:"job_profile"`
}

type UserUpdateRequest struct {
	ID               uuid.UUID `json:"id"`
	FirstName        string    `json:"first_name"`
	LastName         string    `json:"last_name"`
	Email            string    `json:"email"`
	DocumentID       string    `json:"document_id"`
	PhoneCountryCode string    `json:"phone_country_code"`
	PhoneNumber      string    `json:"phone_number"`
	AddressID        uuid.UUID `json:"address_id"`
	ChurchID         uuid.UUID `json:"church_id"`
	// NotificationPreferences
	NotifyByEmail bool `json:"notify_by_email"`
	NotifyBySms   bool `json:"notify_by_sms"`
	// JobProfile
	OpenToWork   bool                      `json:"open_to_work"`
	CVPath       string                    `json:"cv_path"`
	FieldsOfWork []adminDomain.FieldOfWork `json:"fields_of_work"`
}

type UserUpdateResponse struct {
	Message string `json:"message"`
}

type UserListRequest = domain.UserFilters

type UserListResponse struct {
	Users  []*domain.User `json:"users"`
	Count  int            `json:"count"`
	Limit  *int           `json:"limit"`
	Offset *int           `json:"offset"`
}

type UserUpdatePropertyRequest struct {
	ID    uuid.UUID `json:"-"`
	Value bool      `json:"value"`
}

type UserSetRoleRequest struct {
	ID     uuid.UUID `json:"-"`
	RoleID int16     `json:"role_id"`
}

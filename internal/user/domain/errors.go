package domain

import "errors"

// Authentication errors
var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidPassword    = errors.New("invalid password")
	ErrUserInactive       = errors.New("user account is not active")
	ErrEmailNotVerified   = errors.New("email address is not verified")
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrTokenExpired       = errors.New("token has expired")
	ErrUnauthorized       = errors.New("unauthorized access")
)

// User management errors
var (
	ErrUserAlreadyExists        = errors.New("user already exists")
	ErrEmailAlreadyExists       = errors.New("email address already exists")
	ErrInvalidEmail             = errors.New("invalid email address")
	ErrInvalidPhone             = errors.New("invalid phone number")
	ErrDocumentIDAlreadyExists  = errors.New("document id already registered")
	ErrPhoneNumberAlreadyExists = errors.New("phone number already registered")
	ErrInvalidUserID            = errors.New("invalid user ID")
	ErrUserCreationFailed       = errors.New("failed to create user")
	ErrUserUpdateFailed         = errors.New("failed to update user")
	ErrUserDeletionFailed       = errors.New("failed to delete user")
	ErrUserNotUpdated           = errors.New("user was not updated")
)

// Password errors
var (
	ErrPasswordTooShort   = errors.New("password is too short")
	ErrPasswordTooWeak    = errors.New("password is too weak")
	ErrPasswordMismatch   = errors.New("passwords do not match")
	ErrSamePassword       = errors.New("new password cannot be the same as old password")
	ErrPasswordHashFailed = errors.New("failed to hash password")
	ErrInvalidOldPassword = errors.New("old password is incorrect")
)

// Validation errors
var (
	ErrRequiredField     = errors.New("required field is missing")
	ErrInvalidInput      = errors.New("invalid input provided")
	ErrInvalidFieldValue = errors.New("invalid field value")
	ErrEmptyRequest      = errors.New("request body is empty")
)

// Role and permission errors
var (
	ErrInvalidRole             = errors.New("invalid role")
	ErrInsufficientPermissions = errors.New("insufficient permissions")
	ErrRoleNotFound            = errors.New("role not found")
)

// Profile errors
var (
	ErrFieldOfWorkNotFound = errors.New("field of work not found")
	ErrJobProfileNotFound  = errors.New("job profile not found")
	ErrInvalidProfileData  = errors.New("invalid profile data")
)

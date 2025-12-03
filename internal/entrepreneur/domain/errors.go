package domain

import "errors"

// General errors
var (
	ErrNotFound       = errors.New("resource not found")
	ErrInvalidInput   = errors.New("invalid input provided")
	ErrInternalServer = errors.New("internal server error")
	ErrUnauthorized   = errors.New("unauthorized action")
	ErrForbidden      = errors.New("forbidden action")
)

// Industry errors
var (
	ErrIndustryNotFound = errors.New("industry not found")
)

// Business errors
var (
	ErrBusinessNotFound      = errors.New("business not found")
	ErrBusinessAlreadyExists = errors.New("business already exists")
)

// Product errors
var (
	ErrProductNotFound = errors.New("product not found")
)

// Service errors
var (
	ErrServiceNotFound = errors.New("service not found")
)

// Job errors
var (
	ErrJobNotFound = errors.New("job not found")
)

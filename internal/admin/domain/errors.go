package domain

import "errors"

// Address errors
var (
	ErrAddressNotFound = errors.New("address not found")
)

// Church errors
var (
	ErrChurchNotFound      = errors.New("church not found")
	ErrChurchAlreadyExists = errors.New("church already exists")
)

// FieldOfWork errors
var (
	ErrFieldOfWorkNotFound      = errors.New("field of work not found")
	ErrFieldOfWorkAlreadyExists = errors.New("field of work already exists")
)

// Industry errors
var (
	ErrIndustryNotFound      = errors.New("industry not found")
	ErrIndustryAlreadyExists = errors.New("industry already exists")
)

package domain

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type JobType string

const (
	JobTypeFullTime JobType = "Full-Time"
	JobTypePartTime JobType = "Part-Time"
	JobTypeContract JobType = "Contract"
)

type JobLocation string

const (
	JobLocationRemote JobLocation = "Remote"
	JobLocationOnSite JobLocation = "On Site"
	JobLocationHybrid JobLocation = "Hybrid"
)

// Job corresponds to the "jobs" table.
type Job struct {
	ID              uuid.UUID      `json:"id" db:"id"`
	BusinessID      uuid.UUID      `json:"business_id" db:"business_id"`
	Title           string         `json:"title" db:"title"`
	Description     string         `json:"description" db:"description"`
	Type            JobType        `json:"type" db:"type"`
	Location        JobLocation    `json:"location" db:"location"`
	ApplicationLink sql.NullString `json:"application_link" db:"application_link"`
	IsOpen          bool           `json:"is_open" db:"is_open"`
	CreatedAt       time.Time      `json:"created_at" db:"created_at"`
}

// JobFilters defines criteria for filtering jobs.
type JobFilters struct {
	BusinessID    *uuid.UUID   `json:"business_id"`
	Type          *JobType     `json:"type"`
	Location      *JobLocation `json:"location"`
	IsOpen        *bool        `json:"is_open"`
	TitleContains *string      `json:"title_contains"`

	// Pagination
	Limit  *int `json:"limit"`
	Offset *int `json:"offset"`
}

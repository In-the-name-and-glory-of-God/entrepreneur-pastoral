package application

import (
	"context"
	"database/sql"

	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/entrepreneur/domain"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/entrepreneur/infrastructure/dto"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/helper/response"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type JobService struct {
	logger  *zap.SugaredLogger
	jobRepo domain.JobRepository
}

func NewJobService(logger *zap.SugaredLogger, jobRepo domain.JobRepository) *JobService {
	return &JobService{
		logger:  logger,
		jobRepo: jobRepo,
	}
}

func (s *JobService) Create(ctx context.Context, req *dto.JobCreateRequest) (*domain.Job, error) {
	job := &domain.Job{
		BusinessID:      req.BusinessID,
		Title:           req.Title,
		Description:     req.Description,
		Type:            req.Type,
		Location:        req.Location,
		ApplicationLink: sql.NullString{String: req.ApplicationLink, Valid: req.ApplicationLink != ""},
		IsOpen:          req.IsOpen,
	}

	if err := s.jobRepo.Create(nil, job); err != nil {
		s.logger.Errorw("failed to create job", "error", err)
		return nil, response.ErrInternalServerError
	}

	return job, nil
}

func (s *JobService) Update(ctx context.Context, req *dto.JobUpdateRequest) error {
	job, err := s.jobRepo.GetByID(ctx, req.ID)
	if err != nil {
		return err
	}

	job.Title = req.Title
	job.Description = req.Description
	job.Type = req.Type
	job.Location = req.Location
	job.ApplicationLink = sql.NullString{String: req.ApplicationLink, Valid: req.ApplicationLink != ""}
	job.IsOpen = req.IsOpen

	if err := s.jobRepo.Update(nil, job); err != nil {
		s.logger.Errorw("failed to update job", "id", req.ID, "error", err)
		return response.ErrInternalServerError
	}

	return nil
}

func (s *JobService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.jobRepo.Delete(nil, id); err != nil {
		s.logger.Errorw("failed to delete job", "id", id, "error", err)
		return response.ErrInternalServerError
	}

	return nil
}

func (s *JobService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Job, error) {
	job, err := s.jobRepo.GetByID(ctx, id)
	if err != nil {
		if err == domain.ErrJobNotFound {
			return nil, err
		}

		s.logger.Errorw("failed to get job by ID", "id", id, "error", err)
		return nil, response.ErrInternalServerError
	}

	return job, nil
}

func (s *JobService) List(ctx context.Context, req *dto.JobListRequest) (*dto.JobListResponse, error) {
	jobs, err := s.jobRepo.List(ctx, req)
	if err != nil && err != domain.ErrJobNotFound {
		s.logger.Errorw("failed to list jobs", "error", err)
		return nil, response.ErrInternalServerError
	}

	count := 0
	if len(jobs) > 0 {
		count, err = s.jobRepo.Count(ctx, req)
		if err != nil {
			s.logger.Errorw("failed to count jobs", "error", err)
			return nil, response.ErrInternalServerError
		}
	}

	return &dto.JobListResponse{
		Jobs:   jobs,
		Count:  count,
		Limit:  req.Limit,
		Offset: req.Offset,
	}, nil
}

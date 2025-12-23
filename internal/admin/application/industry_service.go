package application

import (
	"context"
	"errors"

	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/admin/domain"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/admin/infrastructure/dto"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/helper/response"
	"go.uber.org/zap"
)

type IndustryService struct {
	logger       *zap.SugaredLogger
	industryRepo domain.IndustryRepository
}

func NewIndustryService(logger *zap.SugaredLogger, industryRepo domain.IndustryRepository) *IndustryService {
	return &IndustryService{
		logger:       logger,
		industryRepo: industryRepo,
	}
}

func (s *IndustryService) Create(ctx context.Context, req *dto.IndustryCreateRequest) (*domain.Industry, error) {
	// Check if industry with same name already exists
	if _, err := s.industryRepo.GetByName(ctx, req.Name); err == nil {
		return nil, domain.ErrIndustryAlreadyExists
	} else if !errors.Is(err, domain.ErrIndustryNotFound) {
		s.logger.Errorw("failed to check existing industry", "name", req.Name, "error", err)
		return nil, response.ErrInternalServerError
	}

	industry := &domain.Industry{
		Name: req.Name,
	}

	if err := s.industryRepo.Create(ctx, industry); err != nil {
		s.logger.Errorw("failed to create industry", "error", err)
		return nil, response.ErrInternalServerError
	}

	return industry, nil
}

func (s *IndustryService) Update(ctx context.Context, req *dto.IndustryUpdateRequest) error {
	_, err := s.industryRepo.GetByID(ctx, req.ID)
	if err != nil {
		if errors.Is(err, domain.ErrIndustryNotFound) {
			return domain.ErrIndustryNotFound
		}

		s.logger.Errorw("failed to get industry by ID", "id", req.ID, "error", err)
		return response.ErrInternalServerError
	}

	// Check if updating to a name that already exists (and belongs to a different industry)
	if existingIndustry, err := s.industryRepo.GetByName(ctx, req.Name); err == nil && existingIndustry.ID != req.ID {
		return domain.ErrIndustryAlreadyExists
	}

	industry := &domain.Industry{
		ID:   req.ID,
		Name: req.Name,
	}

	if err := s.industryRepo.Update(ctx, industry); err != nil {
		s.logger.Errorw("failed to update industry", "id", req.ID, "error", err)
		return response.ErrInternalServerError
	}

	return nil
}

func (s *IndustryService) Delete(ctx context.Context, id int16) error {
	_, err := s.industryRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrIndustryNotFound) {
			return domain.ErrIndustryNotFound
		}

		s.logger.Errorw("failed to get industry by ID", "id", id, "error", err)
		return response.ErrInternalServerError
	}

	if err := s.industryRepo.Delete(ctx, id); err != nil {
		s.logger.Errorw("failed to delete industry", "id", id, "error", err)
		return response.ErrInternalServerError
	}

	return nil
}

func (s *IndustryService) GetByID(ctx context.Context, id int16) (*domain.Industry, error) {
	industry, err := s.industryRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrIndustryNotFound) {
			return nil, domain.ErrIndustryNotFound
		}

		s.logger.Errorw("failed to get industry by ID", "id", id, "error", err)
		return nil, response.ErrInternalServerError
	}

	return industry, nil
}

func (s *IndustryService) GetAll(ctx context.Context) (*dto.IndustryListResponse, error) {
	industries, err := s.industryRepo.GetAll(ctx)
	if err != nil && !errors.Is(err, domain.ErrIndustryNotFound) {
		s.logger.Errorw("failed to get all industries", "error", err)
		return nil, response.ErrInternalServerError
	}

	return &dto.IndustryListResponse{
		Industries: industries,
	}, nil
}

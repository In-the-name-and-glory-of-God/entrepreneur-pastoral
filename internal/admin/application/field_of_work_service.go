package application

import (
	"context"
	"errors"

	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/admin/domain"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/admin/infrastructure/dto"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/helper/response"
	"go.uber.org/zap"
)

type FieldOfWorkService struct {
	logger          *zap.SugaredLogger
	fieldOfWorkRepo domain.FieldOfWorkRepository
}

func NewFieldOfWorkService(logger *zap.SugaredLogger, fieldOfWorkRepo domain.FieldOfWorkRepository) *FieldOfWorkService {
	return &FieldOfWorkService{
		logger:          logger,
		fieldOfWorkRepo: fieldOfWorkRepo,
	}
}

func (s *FieldOfWorkService) Create(ctx context.Context, req *dto.FieldOfWorkCreateRequest) (*domain.FieldOfWork, error) {
	// Check if field of work with same key already exists
	if _, err := s.fieldOfWorkRepo.GetByKey(ctx, req.Key); err == nil {
		return nil, domain.ErrFieldOfWorkAlreadyExists
	} else if !errors.Is(err, domain.ErrFieldOfWorkNotFound) {
		s.logger.Errorw("failed to check existing field of work", "key", req.Key, "error", err)
		return nil, response.ErrInternalServerError
	}

	fieldOfWork := &domain.FieldOfWork{
		Key: req.Key,
	}

	if err := s.fieldOfWorkRepo.Create(ctx, fieldOfWork); err != nil {
		s.logger.Errorw("failed to create field of work", "error", err)
		return nil, response.ErrInternalServerError
	}

	return fieldOfWork, nil
}

func (s *FieldOfWorkService) Update(ctx context.Context, req *dto.FieldOfWorkUpdateRequest) error {
	_, err := s.fieldOfWorkRepo.GetByID(ctx, req.ID)
	if err != nil {
		if errors.Is(err, domain.ErrFieldOfWorkNotFound) {
			return domain.ErrFieldOfWorkNotFound
		}

		s.logger.Errorw("failed to get field of work by ID", "id", req.ID, "error", err)
		return response.ErrInternalServerError
	}

	// Check if updating to a key that already exists (and belongs to a different field of work)
	if existingFieldOfWork, err := s.fieldOfWorkRepo.GetByKey(ctx, req.Key); err == nil && existingFieldOfWork.ID != req.ID {
		return domain.ErrFieldOfWorkAlreadyExists
	}

	fieldOfWork := &domain.FieldOfWork{
		ID:  req.ID,
		Key: req.Key,
	}

	if err := s.fieldOfWorkRepo.Update(ctx, fieldOfWork); err != nil {
		s.logger.Errorw("failed to update field of work", "id", req.ID, "error", err)
		return response.ErrInternalServerError
	}

	return nil
}

func (s *FieldOfWorkService) Delete(ctx context.Context, id int16) error {
	_, err := s.fieldOfWorkRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrFieldOfWorkNotFound) {
			return domain.ErrFieldOfWorkNotFound
		}

		s.logger.Errorw("failed to get field of work by ID", "id", id, "error", err)
		return response.ErrInternalServerError
	}

	if err := s.fieldOfWorkRepo.Delete(ctx, id); err != nil {
		s.logger.Errorw("failed to delete field of work", "id", id, "error", err)
		return response.ErrInternalServerError
	}

	return nil
}

func (s *FieldOfWorkService) GetByID(ctx context.Context, id int16) (*domain.FieldOfWork, error) {
	fieldOfWork, err := s.fieldOfWorkRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrFieldOfWorkNotFound) {
			return nil, domain.ErrFieldOfWorkNotFound
		}

		s.logger.Errorw("failed to get field of work by ID", "id", id, "error", err)
		return nil, response.ErrInternalServerError
	}

	return fieldOfWork, nil
}

func (s *FieldOfWorkService) GetAll(ctx context.Context) (*dto.FieldOfWorkListResponse, error) {
	fieldsOfWork, err := s.fieldOfWorkRepo.GetAll(ctx)
	if err != nil && !errors.Is(err, domain.ErrFieldOfWorkNotFound) {
		s.logger.Errorw("failed to get all fields of work", "error", err)
		return nil, response.ErrInternalServerError
	}

	return &dto.FieldOfWorkListResponse{
		FieldsOfWork: fieldsOfWork,
	}, nil
}

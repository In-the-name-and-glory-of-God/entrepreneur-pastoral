package application

import (
	"context"

	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/entrepreneur/domain"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/entrepreneur/infrastructure/dto"
	userDto "github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/user/infrastructure/dto"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/helper/auth"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/helper/response"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ServiceService struct {
	logger       *zap.SugaredLogger
	serviceRepo  domain.ServiceRepository
	businessRepo domain.BusinessRepository
}

func NewServiceService(logger *zap.SugaredLogger, serviceRepo domain.ServiceRepository, businessRepo domain.BusinessRepository) *ServiceService {
	return &ServiceService{
		logger:       logger,
		serviceRepo:  serviceRepo,
		businessRepo: businessRepo,
	}
}

func (s *ServiceService) Create(ctx context.Context, req *dto.ServiceCreateRequest) (*domain.Service, error) {
	userCtx := ctx.Value(auth.UserContextKey).(*userDto.UserAsContext)
	// Check if business belongs to user
	business, err := s.businessRepo.GetByID(ctx, req.BusinessID)
	if err != nil {
		if err == domain.ErrBusinessNotFound {
			return nil, domain.ErrBusinessNotFound
		}
		return nil, err
	}

	if business.UserID != userCtx.ID {
		return nil, domain.ErrUnauthorized
	}

	service := &domain.Service{
		BusinessID:  req.BusinessID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	}

	if err := s.serviceRepo.Create(nil, service); err != nil {
		s.logger.Errorw("failed to create service", "error", err)
		return nil, response.ErrInternalServerError
	}

	return service, nil
}

func (s *ServiceService) Update(ctx context.Context, req *dto.ServiceUpdateRequest) error {
	userCtx := ctx.Value(auth.UserContextKey).(*userDto.UserAsContext)
	service, err := s.serviceRepo.GetByID(ctx, req.ID)
	if err != nil {
		return err
	}

	// Check if business belongs to user
	business, err := s.businessRepo.GetByID(ctx, service.BusinessID)
	if err != nil {
		return err
	}

	if business.UserID != userCtx.ID {
		return domain.ErrUnauthorized
	}

	service.Name = req.Name
	service.Description = req.Description
	service.Price = req.Price

	if err := s.serviceRepo.Update(nil, service); err != nil {
		s.logger.Errorw("failed to update service", "id", req.ID, "error", err)
		return response.ErrInternalServerError
	}

	return nil
}

func (s *ServiceService) Delete(ctx context.Context, id uuid.UUID) error {
	userCtx := ctx.Value(auth.UserContextKey).(*userDto.UserAsContext)
	service, err := s.serviceRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Check if business belongs to user
	business, err := s.businessRepo.GetByID(ctx, service.BusinessID)
	if err != nil {
		return err
	}

	if business.UserID != userCtx.ID {
		return domain.ErrUnauthorized
	}

	if err := s.serviceRepo.Delete(nil, id); err != nil {
		s.logger.Errorw("failed to delete service", "id", id, "error", err)
		return response.ErrInternalServerError
	}

	return nil
}

func (s *ServiceService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Service, error) {
	service, err := s.serviceRepo.GetByID(ctx, id)
	if err != nil {
		if err == domain.ErrServiceNotFound {
			return nil, err
		}

		s.logger.Errorw("failed to get service by ID", "id", id, "error", err)
		return nil, response.ErrInternalServerError
	}

	return service, nil
}

func (s *ServiceService) List(ctx context.Context, req *dto.ServiceListRequest) (*dto.ServiceListResponse, error) {
	services, err := s.serviceRepo.List(ctx, req)
	if err != nil && err != domain.ErrServiceNotFound {
		s.logger.Errorw("failed to list services", "error", err)
		return nil, response.ErrInternalServerError
	}

	count := 0
	if len(services) > 0 {
		count, err = s.serviceRepo.Count(ctx, req)
		if err != nil {
			s.logger.Errorw("failed to count services", "error", err)
			return nil, response.ErrInternalServerError
		}
	}

	return &dto.ServiceListResponse{
		Services: services,
		Count:    count,
		Limit:    req.Limit,
		Offset:   req.Offset,
	}, nil
}

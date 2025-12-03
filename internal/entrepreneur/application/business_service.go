package application

import (
	"context"
	"database/sql"

	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/entrepreneur/domain"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/entrepreneur/infrastructure/dto"
	userDto "github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/user/infrastructure/dto"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/helper/auth"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/helper/response"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type BusinessService struct {
	logger       *zap.SugaredLogger
	businessRepo domain.BusinessRepository
}

func NewBusinessService(logger *zap.SugaredLogger, businessRepo domain.BusinessRepository) *BusinessService {
	return &BusinessService{
		logger:       logger,
		businessRepo: businessRepo,
	}
}

func (s *BusinessService) Create(ctx context.Context, req *dto.BusinessCreateRequest) (*domain.Business, error) {
	userCtx := ctx.Value(auth.UserContextKey).(*userDto.UserAsContext)
	business := &domain.Business{
		UserID:           userCtx.ID,
		IndustryID:       req.IndustryID,
		Name:             req.Name,
		Description:      req.Description,
		Email:            req.Email,
		PhoneCountryCode: sql.NullString{String: req.PhoneCountryCode, Valid: req.PhoneCountryCode != ""},
		PhoneNumber:      sql.NullString{String: req.PhoneNumber, Valid: req.PhoneNumber != ""},
		WebsiteURL:       sql.NullString{String: req.WebsiteURL, Valid: req.WebsiteURL != ""},
		LogoURL:          sql.NullString{String: req.LogoURL, Valid: req.LogoURL != ""},
		IsActive:         true,
	}

	if err := s.businessRepo.Create(nil, business); err != nil {
		s.logger.Errorw("failed to create business", "error", err)
		return nil, response.ErrInternalServerError
	}

	return business, nil
}

func (s *BusinessService) Update(ctx context.Context, req *dto.BusinessUpdateRequest) error {
	userCtx := ctx.Value(auth.UserContextKey).(*userDto.UserAsContext)
	business, err := s.businessRepo.GetByID(ctx, req.ID)
	if err != nil {
		return err
	}

	if business.UserID != userCtx.ID {
		return domain.ErrUnauthorized
	}

	business.IndustryID = req.IndustryID
	business.Name = req.Name
	business.Description = req.Description
	business.Email = req.Email
	business.PhoneCountryCode = sql.NullString{String: req.PhoneCountryCode, Valid: req.PhoneCountryCode != ""}
	business.PhoneNumber = sql.NullString{String: req.PhoneNumber, Valid: req.PhoneNumber != ""}
	business.WebsiteURL = sql.NullString{String: req.WebsiteURL, Valid: req.WebsiteURL != ""}
	business.LogoURL = sql.NullString{String: req.LogoURL, Valid: req.LogoURL != ""}
	business.IsActive = req.IsActive

	if err := s.businessRepo.Update(nil, business); err != nil {
		s.logger.Errorw("failed to update business", "id", req.ID, "error", err)
		return response.ErrInternalServerError
	}

	return nil
}

func (s *BusinessService) Delete(ctx context.Context, id uuid.UUID) error {
	userCtx := ctx.Value(auth.UserContextKey).(*userDto.UserAsContext)
	business, err := s.businessRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if business.UserID != userCtx.ID {
		return domain.ErrUnauthorized
	}

	if err := s.businessRepo.Delete(nil, id); err != nil {
		s.logger.Errorw("failed to delete business", "id", id, "error", err)
		return response.ErrInternalServerError
	}

	return nil
}

func (s *BusinessService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Business, error) {
	business, err := s.businessRepo.GetByID(ctx, id)
	if err != nil {
		if err == domain.ErrBusinessNotFound {
			return nil, err
		}

		s.logger.Errorw("failed to get business by ID", "id", id, "error", err)
		return nil, response.ErrInternalServerError
	}

	return business, nil
}

func (s *BusinessService) List(ctx context.Context, req *dto.BusinessListRequest) (*dto.BusinessListResponse, error) {
	businesses, err := s.businessRepo.List(ctx, req)
	if err != nil && err != domain.ErrBusinessNotFound {
		s.logger.Errorw("failed to list businesses", "error", err)
		return nil, response.ErrInternalServerError
	}

	count := 0
	if len(businesses) > 0 {
		count, err = s.businessRepo.Count(ctx, req)
		if err != nil {
			s.logger.Errorw("failed to count businesses", "error", err)
			return nil, response.ErrInternalServerError
		}
	}

	return &dto.BusinessListResponse{
		Businesses: businesses,
		Count:      count,
		Limit:      req.Limit,
		Offset:     req.Offset,
	}, nil
}

package application

import (
	"context"
	"database/sql"
	"errors"

	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/admin/domain"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/admin/infrastructure/dto"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/helper/response"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type ChurchService struct {
	logger      *zap.SugaredLogger
	churchRepo  domain.ChurchRepository
	addressRepo domain.AddressRepository
}

func NewChurchService(logger *zap.SugaredLogger, churchRepo domain.ChurchRepository, addressRepo domain.AddressRepository) *ChurchService {
	return &ChurchService{
		logger:      logger,
		churchRepo:  churchRepo,
		addressRepo: addressRepo,
	}
}

func (s *ChurchService) Create(ctx context.Context, req *dto.ChurchCreateRequest) (*domain.Church, error) {
	// Check if church with same name already exists
	if _, err := s.churchRepo.GetByName(ctx, req.Name); err == nil {
		return nil, domain.ErrChurchAlreadyExists
	} else if !errors.Is(err, domain.ErrChurchNotFound) {
		s.logger.Errorw("failed to check existing church", "name", req.Name, "error", err)
		return nil, response.ErrInternalServerError
	}

	var church *domain.Church

	err := s.churchRepo.UnitOfWork(ctx, func(tx *sqlx.Tx) error {
		// 1. Create the Address
		address := &domain.Address{
			StreetLine1:   req.Address.StreetLine1,
			StreetLine2:   sql.NullString{String: req.Address.StreetLine2, Valid: req.Address.StreetLine2 != ""},
			City:          req.Address.City,
			StateProvince: req.Address.StateProvince,
			PostalCode:    req.Address.PostalCode,
			Country:       req.Address.Country,
		}
		if err := s.addressRepo.Create(tx, address); err != nil {
			s.logger.Errorw("failed to create address for church", "error", err)
			return response.ErrInternalServerError
		}

		// 2. Create the Church with the Address ID
		church = &domain.Church{
			Name:          req.Name,
			Diocese:       req.Diocese,
			ParishNumber:  sql.NullString{String: req.ParishNumber, Valid: req.ParishNumber != ""},
			WebsiteURL:    sql.NullString{String: req.WebsiteURL, Valid: req.WebsiteURL != ""},
			PhoneNumber:   sql.NullString{String: req.PhoneNumber, Valid: req.PhoneNumber != ""},
			AddressID:     address.ID,
			IsArchdiocese: req.IsArchdiocese,
			IsActive:      true,
		}
		if err := s.churchRepo.Create(tx, church); err != nil {
			s.logger.Errorw("failed to create church", "error", err)
			return response.ErrInternalServerError
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return church, nil
}

func (s *ChurchService) Update(ctx context.Context, req *dto.ChurchUpdateRequest) error {
	church, err := s.churchRepo.GetByID(ctx, req.ID)
	if err != nil {
		if errors.Is(err, domain.ErrChurchNotFound) {
			return domain.ErrChurchNotFound
		}

		s.logger.Errorw("failed to get church by ID", "id", req.ID, "error", err)
		return response.ErrInternalServerError
	}

	// Check if updating to a name that already exists (and belongs to a different church)
	if existingChurch, err := s.churchRepo.GetByName(ctx, req.Name); err == nil && existingChurch.ID != req.ID {
		return domain.ErrChurchAlreadyExists
	}

	church.Name = req.Name
	church.Diocese = req.Diocese
	church.ParishNumber = sql.NullString{String: req.ParishNumber, Valid: req.ParishNumber != ""}
	church.WebsiteURL = sql.NullString{String: req.WebsiteURL, Valid: req.WebsiteURL != ""}
	church.PhoneNumber = sql.NullString{String: req.PhoneNumber, Valid: req.PhoneNumber != ""}
	church.AddressID = req.AddressID
	church.IsArchdiocese = req.IsArchdiocese
	church.IsActive = req.IsActive

	if err := s.churchRepo.Update(ctx, church); err != nil {
		s.logger.Errorw("failed to update church", "id", req.ID, "error", err)
		return response.ErrInternalServerError
	}

	return nil
}

func (s *ChurchService) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.churchRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrChurchNotFound) {
			return domain.ErrChurchNotFound
		}

		s.logger.Errorw("failed to get church by ID", "id", id, "error", err)
		return response.ErrInternalServerError
	}

	if err := s.churchRepo.Delete(ctx, id); err != nil {
		s.logger.Errorw("failed to delete church", "id", id, "error", err)
		return response.ErrInternalServerError
	}

	return nil
}

func (s *ChurchService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Church, error) {
	church, err := s.churchRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrChurchNotFound) {
			return nil, domain.ErrChurchNotFound
		}

		s.logger.Errorw("failed to get church by ID", "id", id, "error", err)
		return nil, response.ErrInternalServerError
	}

	return church, nil
}

func (s *ChurchService) List(ctx context.Context, req *dto.ChurchListRequest) (*dto.ChurchListResponse, error) {
	churches, err := s.churchRepo.List(ctx, req)
	if err != nil && !errors.Is(err, domain.ErrChurchNotFound) {
		s.logger.Errorw("failed to list churches", "error", err)
		return nil, response.ErrInternalServerError
	}

	count := 0
	if len(churches) > 0 {
		count, err = s.churchRepo.Count(ctx, req)
		if err != nil {
			s.logger.Errorw("failed to count churches", "error", err)
			return nil, response.ErrInternalServerError
		}
	}

	return &dto.ChurchListResponse{
		Churches: churches,
		Count:    count,
		Limit:    req.Limit,
		Offset:   req.Offset,
	}, nil
}

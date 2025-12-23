package application

import (
	"context"
	"database/sql"
	"errors"

	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/admin/domain"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/admin/infrastructure/dto"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/helper/response"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AddressService struct {
	logger      *zap.SugaredLogger
	addressRepo domain.AddressRepository
}

func NewAddressService(logger *zap.SugaredLogger, addressRepo domain.AddressRepository) *AddressService {
	return &AddressService{
		logger:      logger,
		addressRepo: addressRepo,
	}
}

func (s *AddressService) Create(ctx context.Context, req *dto.AddressCreateRequest) (*domain.Address, error) {
	address := &domain.Address{
		StreetLine1:   req.StreetLine1,
		StreetLine2:   sql.NullString{String: req.StreetLine2, Valid: req.StreetLine2 != ""},
		City:          req.City,
		StateProvince: req.StateProvince,
		PostalCode:    req.PostalCode,
		Country:       req.Country,
	}

	if err := s.addressRepo.CreateWithContext(ctx, address); err != nil {
		s.logger.Errorw("failed to create address", "error", err)
		return nil, response.ErrInternalServerError
	}

	return address, nil
}

func (s *AddressService) Update(ctx context.Context, req *dto.AddressUpdateRequest) error {
	address, err := s.addressRepo.GetByID(ctx, req.ID)
	if err != nil {
		if errors.Is(err, domain.ErrAddressNotFound) {
			return domain.ErrAddressNotFound
		}

		s.logger.Errorw("failed to get address by ID", "id", req.ID, "error", err)
		return response.ErrInternalServerError
	}

	address.StreetLine1 = req.StreetLine1
	address.StreetLine2 = sql.NullString{String: req.StreetLine2, Valid: req.StreetLine2 != ""}
	address.City = req.City
	address.StateProvince = req.StateProvince
	address.PostalCode = req.PostalCode
	address.Country = req.Country

	if err := s.addressRepo.Update(ctx, address); err != nil {
		s.logger.Errorw("failed to update address", "id", req.ID, "error", err)
		return response.ErrInternalServerError
	}

	return nil
}

func (s *AddressService) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.addressRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrAddressNotFound) {
			return domain.ErrAddressNotFound
		}

		s.logger.Errorw("failed to get address by ID", "id", id, "error", err)
		return response.ErrInternalServerError
	}

	if err := s.addressRepo.Delete(ctx, id); err != nil {
		s.logger.Errorw("failed to delete address", "id", id, "error", err)
		return response.ErrInternalServerError
	}

	return nil
}

func (s *AddressService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Address, error) {
	address, err := s.addressRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrAddressNotFound) {
			return nil, domain.ErrAddressNotFound
		}

		s.logger.Errorw("failed to get address by ID", "id", id, "error", err)
		return nil, response.ErrInternalServerError
	}

	return address, nil
}

package application

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/admin/domain"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/admin/infrastructure/dto"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/helper/response"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// MockAddressRepository
type MockAddressRepository struct {
	mock.Mock
}

func (m *MockAddressRepository) Create(tx *sqlx.Tx, address *domain.Address) error {
	args := m.Called(tx, address)
	if args.Error(0) == nil {
		address.ID = uuid.New()
	}
	return args.Error(0)
}

func (m *MockAddressRepository) CreateWithContext(ctx context.Context, address *domain.Address) error {
	args := m.Called(ctx, address)
	if args.Error(0) == nil {
		address.ID = uuid.New()
	}
	return args.Error(0)
}

func (m *MockAddressRepository) Update(ctx context.Context, address *domain.Address) error {
	args := m.Called(ctx, address)
	return args.Error(0)
}

func (m *MockAddressRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockAddressRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Address, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Address), args.Error(1)
}

func TestAddressService_Create(t *testing.T) {
	logger := zap.NewNop().Sugar()
	mockRepo := new(MockAddressRepository)
	service := NewAddressService(logger, mockRepo)
	ctx := context.Background()

	req := &dto.AddressCreateRequest{
		StreetLine1:   "123 Main St",
		StreetLine2:   "Apt 4B",
		City:          "New York",
		StateProvince: "NY",
		PostalCode:    "10001",
		Country:       "USA",
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("CreateWithContext", ctx, mock.AnythingOfType("*domain.Address")).Return(nil)

		result, err := service.Create(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, req.StreetLine1, result.StreetLine1)
		assert.Equal(t, req.City, result.City)
		assert.Equal(t, req.StateProvince, result.StateProvince)
		assert.Equal(t, req.PostalCode, result.PostalCode)
		assert.Equal(t, req.Country, result.Country)
		assert.True(t, result.StreetLine2.Valid)
		assert.Equal(t, req.StreetLine2, result.StreetLine2.String)
		mockRepo.AssertExpectations(t)
	})

	t.Run("SuccessWithoutOptionalFields", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		reqNoOptional := &dto.AddressCreateRequest{
			StreetLine1:   "456 Oak Ave",
			City:          "Los Angeles",
			StateProvince: "CA",
			PostalCode:    "90001",
			Country:       "USA",
		}
		mockRepo.On("CreateWithContext", ctx, mock.MatchedBy(func(a *domain.Address) bool {
			return !a.StreetLine2.Valid
		})).Return(nil)

		result, err := service.Create(ctx, reqNoOptional)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.False(t, result.StreetLine2.Valid)
		mockRepo.AssertExpectations(t)
	})

	t.Run("CreateFailure", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("CreateWithContext", ctx, mock.AnythingOfType("*domain.Address")).Return(errors.New("db error"))

		result, err := service.Create(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, response.ErrInternalServerError, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestAddressService_Update(t *testing.T) {
	logger := zap.NewNop().Sugar()
	mockRepo := new(MockAddressRepository)
	service := NewAddressService(logger, mockRepo)
	ctx := context.Background()

	addressID := uuid.New()
	req := &dto.AddressUpdateRequest{
		ID:            addressID,
		StreetLine1:   "789 Updated St",
		StreetLine2:   "Suite 200",
		City:          "Chicago",
		StateProvince: "IL",
		PostalCode:    "60601",
		Country:       "USA",
	}

	existingAddress := &domain.Address{
		ID:            addressID,
		StreetLine1:   "123 Original St",
		StreetLine2:   sql.NullString{String: "Apt 1", Valid: true},
		City:          "Original City",
		StateProvince: "OC",
		PostalCode:    "00000",
		Country:       "USA",
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("GetByID", ctx, req.ID).Return(existingAddress, nil)
		mockRepo.On("Update", ctx, mock.AnythingOfType("*domain.Address")).Return(nil)

		err := service.Update(ctx, req)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("GetByID", ctx, req.ID).Return(nil, domain.ErrAddressNotFound)

		err := service.Update(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, domain.ErrAddressNotFound, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("UpdateFailure", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("GetByID", ctx, req.ID).Return(existingAddress, nil)
		mockRepo.On("Update", ctx, mock.AnythingOfType("*domain.Address")).Return(errors.New("db error"))

		err := service.Update(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, response.ErrInternalServerError, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("GetByIDInternalError", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("GetByID", ctx, req.ID).Return(nil, errors.New("db error"))

		err := service.Update(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, response.ErrInternalServerError, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("UpdateRemovesOptionalField", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		reqNoOptional := &dto.AddressUpdateRequest{
			ID:            addressID,
			StreetLine1:   "789 Updated St",
			StreetLine2:   "", // Empty string should result in Invalid NullString
			City:          "Chicago",
			StateProvince: "IL",
			PostalCode:    "60601",
			Country:       "USA",
		}
		mockRepo.On("GetByID", ctx, reqNoOptional.ID).Return(existingAddress, nil)
		mockRepo.On("Update", ctx, mock.MatchedBy(func(a *domain.Address) bool {
			return !a.StreetLine2.Valid
		})).Return(nil)

		err := service.Update(ctx, reqNoOptional)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestAddressService_Delete(t *testing.T) {
	logger := zap.NewNop().Sugar()
	mockRepo := new(MockAddressRepository)
	service := NewAddressService(logger, mockRepo)
	ctx := context.Background()

	id := uuid.New()
	existingAddress := &domain.Address{
		ID:            id,
		StreetLine1:   "123 Main St",
		City:          "Test City",
		StateProvince: "TS",
		PostalCode:    "12345",
		Country:       "USA",
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("GetByID", ctx, id).Return(existingAddress, nil)
		mockRepo.On("Delete", ctx, id).Return(nil)

		err := service.Delete(ctx, id)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("GetByID", ctx, id).Return(nil, domain.ErrAddressNotFound)

		err := service.Delete(ctx, id)

		assert.Error(t, err)
		assert.Equal(t, domain.ErrAddressNotFound, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("DeleteFailure", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("GetByID", ctx, id).Return(existingAddress, nil)
		mockRepo.On("Delete", ctx, id).Return(errors.New("db error"))

		err := service.Delete(ctx, id)

		assert.Error(t, err)
		assert.Equal(t, response.ErrInternalServerError, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("GetByIDInternalError", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("GetByID", ctx, id).Return(nil, errors.New("db error"))

		err := service.Delete(ctx, id)

		assert.Error(t, err)
		assert.Equal(t, response.ErrInternalServerError, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestAddressService_GetByID(t *testing.T) {
	logger := zap.NewNop().Sugar()
	mockRepo := new(MockAddressRepository)
	service := NewAddressService(logger, mockRepo)
	ctx := context.Background()

	id := uuid.New()
	expectedAddress := &domain.Address{
		ID:            id,
		StreetLine1:   "123 Main St",
		StreetLine2:   sql.NullString{String: "Apt 4B", Valid: true},
		City:          "New York",
		StateProvince: "NY",
		PostalCode:    "10001",
		Country:       "USA",
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("GetByID", ctx, id).Return(expectedAddress, nil)

		result, err := service.GetByID(ctx, id)

		assert.NoError(t, err)
		assert.Equal(t, expectedAddress, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("GetByID", ctx, id).Return(nil, domain.ErrAddressNotFound)

		result, err := service.GetByID(ctx, id)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, domain.ErrAddressNotFound, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("InternalError", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("GetByID", ctx, id).Return(nil, errors.New("db error"))

		result, err := service.GetByID(ctx, id)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, response.ErrInternalServerError, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestAddressService_CreateWithEmptyOptionalFields(t *testing.T) {
	logger := zap.NewNop().Sugar()
	mockRepo := new(MockAddressRepository)
	service := NewAddressService(logger, mockRepo)
	ctx := context.Background()

	req := &dto.AddressCreateRequest{
		StreetLine1:   "Simple Address",
		StreetLine2:   "", // Empty optional field
		City:          "Simple City",
		StateProvince: "SC",
		PostalCode:    "00000",
		Country:       "USA",
	}

	t.Run("EmptyStreetLine2ResultsInInvalidNullString", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("CreateWithContext", ctx, mock.MatchedBy(func(a *domain.Address) bool {
			return a.StreetLine1 == req.StreetLine1 &&
				!a.StreetLine2.Valid &&
				a.City == req.City
		})).Return(nil)

		result, err := service.Create(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.False(t, result.StreetLine2.Valid)
		mockRepo.AssertExpectations(t)
	})
}

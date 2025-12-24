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

// MockChurchRepository
type MockChurchRepository struct {
	mock.Mock
}

func (m *MockChurchRepository) UnitOfWork(ctx context.Context, fn func(*sqlx.Tx) error) error {
	args := m.Called(ctx, fn)
	// Execute the function with nil tx if the mock expects success
	if args.Error(0) == nil {
		return fn(nil)
	}
	return args.Error(0)
}

func (m *MockChurchRepository) Create(tx *sqlx.Tx, church *domain.Church) error {
	args := m.Called(tx, church)
	if args.Error(0) == nil {
		church.ID = uuid.New()
	}
	return args.Error(0)
}

func (m *MockChurchRepository) Update(ctx context.Context, church *domain.Church) error {
	args := m.Called(ctx, church)
	return args.Error(0)
}

func (m *MockChurchRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockChurchRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Church, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Church), args.Error(1)
}

func (m *MockChurchRepository) GetByName(ctx context.Context, name string) (*domain.Church, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Church), args.Error(1)
}

func (m *MockChurchRepository) List(ctx context.Context, filter *domain.ChurchFilters) ([]*domain.Church, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Church), args.Error(1)
}

func (m *MockChurchRepository) Count(ctx context.Context, filter *domain.ChurchFilters) (int, error) {
	args := m.Called(ctx, filter)
	return args.Int(0), args.Error(1)
}

// MockAddressRepository for church service tests
type MockAddressRepositoryForChurch struct {
	mock.Mock
}

func (m *MockAddressRepositoryForChurch) Create(tx *sqlx.Tx, address *domain.Address) error {
	args := m.Called(tx, address)
	if args.Error(0) == nil {
		address.ID = uuid.New()
	}
	return args.Error(0)
}

func (m *MockAddressRepositoryForChurch) CreateWithContext(ctx context.Context, address *domain.Address) error {
	args := m.Called(ctx, address)
	if args.Error(0) == nil {
		address.ID = uuid.New()
	}
	return args.Error(0)
}

func (m *MockAddressRepositoryForChurch) Update(ctx context.Context, address *domain.Address) error {
	args := m.Called(ctx, address)
	return args.Error(0)
}

func (m *MockAddressRepositoryForChurch) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockAddressRepositoryForChurch) GetByID(ctx context.Context, id uuid.UUID) (*domain.Address, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Address), args.Error(1)
}

func TestChurchService_Create(t *testing.T) {
	logger := zap.NewNop().Sugar()
	mockChurchRepo := new(MockChurchRepository)
	mockAddressRepo := new(MockAddressRepositoryForChurch)
	service := NewChurchService(logger, mockChurchRepo, mockAddressRepo)
	ctx := context.Background()

	req := &dto.ChurchCreateRequest{
		Name:    "St. Mary's Cathedral",
		Diocese: "Los Angeles",
		Address: dto.AddressCreateRequest{
			StreetLine1:   "123 Main St",
			City:          "Los Angeles",
			StateProvince: "CA",
			PostalCode:    "90001",
			Country:       "USA",
		},
		IsArchdiocese: true,
	}

	t.Run("Success", func(t *testing.T) {
		mockChurchRepo.ExpectedCalls = nil
		mockAddressRepo.ExpectedCalls = nil
		mockChurchRepo.On("GetByName", ctx, req.Name).Return(nil, domain.ErrChurchNotFound)
		mockChurchRepo.On("UnitOfWork", ctx, mock.Anything).Return(nil)
		mockAddressRepo.On("Create", (*sqlx.Tx)(nil), mock.AnythingOfType("*domain.Address")).Return(nil)
		mockChurchRepo.On("Create", (*sqlx.Tx)(nil), mock.AnythingOfType("*domain.Church")).Return(nil)

		result, err := service.Create(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, req.Name, result.Name)
		assert.Equal(t, req.Diocese, result.Diocese)
		mockChurchRepo.AssertExpectations(t)
		mockAddressRepo.AssertExpectations(t)
	})

	t.Run("AlreadyExists", func(t *testing.T) {
		mockChurchRepo.ExpectedCalls = nil
		mockAddressRepo.ExpectedCalls = nil
		existingChurch := &domain.Church{ID: uuid.New(), Name: req.Name}
		mockChurchRepo.On("GetByName", ctx, req.Name).Return(existingChurch, nil)

		result, err := service.Create(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, domain.ErrChurchAlreadyExists, err)
		mockChurchRepo.AssertExpectations(t)
	})

	t.Run("GetByNameInternalError", func(t *testing.T) {
		mockChurchRepo.ExpectedCalls = nil
		mockAddressRepo.ExpectedCalls = nil
		mockChurchRepo.On("GetByName", ctx, req.Name).Return(nil, errors.New("db error"))

		result, err := service.Create(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, response.ErrInternalServerError, err)
		mockChurchRepo.AssertExpectations(t)
	})
}

func TestChurchService_Update(t *testing.T) {
	logger := zap.NewNop().Sugar()
	mockChurchRepo := new(MockChurchRepository)
	mockAddressRepo := new(MockAddressRepositoryForChurch)
	service := NewChurchService(logger, mockChurchRepo, mockAddressRepo)
	ctx := context.Background()

	churchID := uuid.New()
	addressID := uuid.New()

	req := &dto.ChurchUpdateRequest{
		ID:            churchID,
		Name:          "Updated Church Name",
		Diocese:       "Updated Diocese",
		AddressID:     addressID,
		IsArchdiocese: false,
		IsActive:      true,
	}

	existingChurch := &domain.Church{
		ID:            churchID,
		Name:          "Original Church",
		Diocese:       "Original Diocese",
		AddressID:     addressID,
		IsArchdiocese: true,
		IsActive:      true,
	}

	t.Run("Success", func(t *testing.T) {
		mockChurchRepo.ExpectedCalls = nil
		mockAddressRepo.ExpectedCalls = nil
		mockChurchRepo.On("GetByID", ctx, req.ID).Return(existingChurch, nil)
		mockChurchRepo.On("GetByName", ctx, req.Name).Return(nil, domain.ErrChurchNotFound)
		mockChurchRepo.On("Update", ctx, mock.AnythingOfType("*domain.Church")).Return(nil)

		err := service.Update(ctx, req)

		assert.NoError(t, err)
		mockChurchRepo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockChurchRepo.ExpectedCalls = nil
		mockAddressRepo.ExpectedCalls = nil
		mockChurchRepo.On("GetByID", ctx, req.ID).Return(nil, domain.ErrChurchNotFound)

		err := service.Update(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, domain.ErrChurchNotFound, err)
		mockChurchRepo.AssertExpectations(t)
	})

	t.Run("NameAlreadyExists", func(t *testing.T) {
		mockChurchRepo.ExpectedCalls = nil
		mockAddressRepo.ExpectedCalls = nil
		differentChurch := &domain.Church{ID: uuid.New(), Name: req.Name}
		mockChurchRepo.On("GetByID", ctx, req.ID).Return(existingChurch, nil)
		mockChurchRepo.On("GetByName", ctx, req.Name).Return(differentChurch, nil)

		err := service.Update(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, domain.ErrChurchAlreadyExists, err)
		mockChurchRepo.AssertExpectations(t)
	})

	t.Run("UpdateSameNameSameID", func(t *testing.T) {
		mockChurchRepo.ExpectedCalls = nil
		mockAddressRepo.ExpectedCalls = nil
		sameChurch := &domain.Church{ID: req.ID, Name: req.Name}
		mockChurchRepo.On("GetByID", ctx, req.ID).Return(existingChurch, nil)
		mockChurchRepo.On("GetByName", ctx, req.Name).Return(sameChurch, nil)
		mockChurchRepo.On("Update", ctx, mock.AnythingOfType("*domain.Church")).Return(nil)

		err := service.Update(ctx, req)

		assert.NoError(t, err)
		mockChurchRepo.AssertExpectations(t)
	})

	t.Run("UpdateFailure", func(t *testing.T) {
		mockChurchRepo.ExpectedCalls = nil
		mockAddressRepo.ExpectedCalls = nil
		mockChurchRepo.On("GetByID", ctx, req.ID).Return(existingChurch, nil)
		mockChurchRepo.On("GetByName", ctx, req.Name).Return(nil, domain.ErrChurchNotFound)
		mockChurchRepo.On("Update", ctx, mock.AnythingOfType("*domain.Church")).Return(errors.New("db error"))

		err := service.Update(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, response.ErrInternalServerError, err)
		mockChurchRepo.AssertExpectations(t)
	})

	t.Run("GetByIDInternalError", func(t *testing.T) {
		mockChurchRepo.ExpectedCalls = nil
		mockAddressRepo.ExpectedCalls = nil
		mockChurchRepo.On("GetByID", ctx, req.ID).Return(nil, errors.New("db error"))

		err := service.Update(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, response.ErrInternalServerError, err)
		mockChurchRepo.AssertExpectations(t)
	})
}

func TestChurchService_Delete(t *testing.T) {
	logger := zap.NewNop().Sugar()
	mockChurchRepo := new(MockChurchRepository)
	mockAddressRepo := new(MockAddressRepositoryForChurch)
	service := NewChurchService(logger, mockChurchRepo, mockAddressRepo)
	ctx := context.Background()

	id := uuid.New()
	existingChurch := &domain.Church{ID: id, Name: "Test Church"}

	t.Run("Success", func(t *testing.T) {
		mockChurchRepo.ExpectedCalls = nil
		mockAddressRepo.ExpectedCalls = nil
		mockChurchRepo.On("GetByID", ctx, id).Return(existingChurch, nil)
		mockChurchRepo.On("Delete", ctx, id).Return(nil)

		err := service.Delete(ctx, id)

		assert.NoError(t, err)
		mockChurchRepo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockChurchRepo.ExpectedCalls = nil
		mockAddressRepo.ExpectedCalls = nil
		mockChurchRepo.On("GetByID", ctx, id).Return(nil, domain.ErrChurchNotFound)

		err := service.Delete(ctx, id)

		assert.Error(t, err)
		assert.Equal(t, domain.ErrChurchNotFound, err)
		mockChurchRepo.AssertExpectations(t)
	})

	t.Run("DeleteFailure", func(t *testing.T) {
		mockChurchRepo.ExpectedCalls = nil
		mockAddressRepo.ExpectedCalls = nil
		mockChurchRepo.On("GetByID", ctx, id).Return(existingChurch, nil)
		mockChurchRepo.On("Delete", ctx, id).Return(errors.New("db error"))

		err := service.Delete(ctx, id)

		assert.Error(t, err)
		assert.Equal(t, response.ErrInternalServerError, err)
		mockChurchRepo.AssertExpectations(t)
	})

	t.Run("GetByIDInternalError", func(t *testing.T) {
		mockChurchRepo.ExpectedCalls = nil
		mockAddressRepo.ExpectedCalls = nil
		mockChurchRepo.On("GetByID", ctx, id).Return(nil, errors.New("db error"))

		err := service.Delete(ctx, id)

		assert.Error(t, err)
		assert.Equal(t, response.ErrInternalServerError, err)
		mockChurchRepo.AssertExpectations(t)
	})
}

func TestChurchService_GetByID(t *testing.T) {
	logger := zap.NewNop().Sugar()
	mockChurchRepo := new(MockChurchRepository)
	mockAddressRepo := new(MockAddressRepositoryForChurch)
	service := NewChurchService(logger, mockChurchRepo, mockAddressRepo)
	ctx := context.Background()

	id := uuid.New()
	expectedChurch := &domain.Church{
		ID:      id,
		Name:    "Test Church",
		Diocese: "Test Diocese",
	}

	t.Run("Success", func(t *testing.T) {
		mockChurchRepo.ExpectedCalls = nil
		mockAddressRepo.ExpectedCalls = nil
		mockChurchRepo.On("GetByID", ctx, id).Return(expectedChurch, nil)

		result, err := service.GetByID(ctx, id)

		assert.NoError(t, err)
		assert.Equal(t, expectedChurch, result)
		mockChurchRepo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockChurchRepo.ExpectedCalls = nil
		mockAddressRepo.ExpectedCalls = nil
		mockChurchRepo.On("GetByID", ctx, id).Return(nil, domain.ErrChurchNotFound)

		result, err := service.GetByID(ctx, id)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, domain.ErrChurchNotFound, err)
		mockChurchRepo.AssertExpectations(t)
	})

	t.Run("InternalError", func(t *testing.T) {
		mockChurchRepo.ExpectedCalls = nil
		mockAddressRepo.ExpectedCalls = nil
		mockChurchRepo.On("GetByID", ctx, id).Return(nil, errors.New("db error"))

		result, err := service.GetByID(ctx, id)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, response.ErrInternalServerError, err)
		mockChurchRepo.AssertExpectations(t)
	})
}

func TestChurchService_List(t *testing.T) {
	logger := zap.NewNop().Sugar()
	mockChurchRepo := new(MockChurchRepository)
	mockAddressRepo := new(MockAddressRepositoryForChurch)
	service := NewChurchService(logger, mockChurchRepo, mockAddressRepo)
	ctx := context.Background()

	limit := 10
	offset := 0
	req := &dto.ChurchListRequest{
		Limit:  &limit,
		Offset: &offset,
	}

	expectedChurches := []*domain.Church{
		{ID: uuid.New(), Name: "Church 1", Diocese: "Diocese 1"},
		{ID: uuid.New(), Name: "Church 2", Diocese: "Diocese 2"},
	}

	t.Run("Success", func(t *testing.T) {
		mockChurchRepo.ExpectedCalls = nil
		mockAddressRepo.ExpectedCalls = nil
		mockChurchRepo.On("List", ctx, req).Return(expectedChurches, nil)
		mockChurchRepo.On("Count", ctx, req).Return(2, nil)

		result, err := service.List(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 2, result.Count)
		assert.Len(t, result.Churches, 2)
		mockChurchRepo.AssertExpectations(t)
	})

	t.Run("EmptyList", func(t *testing.T) {
		mockChurchRepo.ExpectedCalls = nil
		mockAddressRepo.ExpectedCalls = nil
		mockChurchRepo.On("List", ctx, req).Return([]*domain.Church{}, nil)

		result, err := service.List(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 0, result.Count)
		assert.Len(t, result.Churches, 0)
		mockChurchRepo.AssertExpectations(t)
	})

	t.Run("ListFailure", func(t *testing.T) {
		mockChurchRepo.ExpectedCalls = nil
		mockAddressRepo.ExpectedCalls = nil
		mockChurchRepo.On("List", ctx, req).Return(nil, errors.New("db error"))

		result, err := service.List(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, response.ErrInternalServerError, err)
		mockChurchRepo.AssertExpectations(t)
	})

	t.Run("CountFailure", func(t *testing.T) {
		mockChurchRepo.ExpectedCalls = nil
		mockAddressRepo.ExpectedCalls = nil
		mockChurchRepo.On("List", ctx, req).Return(expectedChurches, nil)
		mockChurchRepo.On("Count", ctx, req).Return(0, errors.New("db error"))

		result, err := service.List(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, response.ErrInternalServerError, err)
		mockChurchRepo.AssertExpectations(t)
	})

	t.Run("WithFilters", func(t *testing.T) {
		mockChurchRepo.ExpectedCalls = nil
		mockAddressRepo.ExpectedCalls = nil
		diocese := "Los Angeles"
		isActive := true
		filterReq := &dto.ChurchListRequest{
			Diocese:  &diocese,
			IsActive: &isActive,
			Limit:    &limit,
			Offset:   &offset,
		}
		mockChurchRepo.On("List", ctx, filterReq).Return(expectedChurches, nil)
		mockChurchRepo.On("Count", ctx, filterReq).Return(2, nil)

		result, err := service.List(ctx, filterReq)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 2, result.Count)
		mockChurchRepo.AssertExpectations(t)
	})
}

func TestChurchService_CreateWithOptionalFields(t *testing.T) {
	logger := zap.NewNop().Sugar()
	mockChurchRepo := new(MockChurchRepository)
	mockAddressRepo := new(MockAddressRepositoryForChurch)
	service := NewChurchService(logger, mockChurchRepo, mockAddressRepo)
	ctx := context.Background()

	req := &dto.ChurchCreateRequest{
		Name:         "Simple Church",
		Diocese:      "Test Diocese",
		ParishNumber: "P-123",
		WebsiteURL:   "https://church.example.com",
		PhoneNumber:  "+1234567890",
		Address: dto.AddressCreateRequest{
			StreetLine1:   "456 Church St",
			StreetLine2:   "Suite 100",
			City:          "Test City",
			StateProvince: "TS",
			PostalCode:    "12345",
			Country:       "USA",
		},
		IsArchdiocese: false,
	}

	t.Run("SuccessWithAllOptionalFields", func(t *testing.T) {
		mockChurchRepo.ExpectedCalls = nil
		mockAddressRepo.ExpectedCalls = nil
		mockChurchRepo.On("GetByName", ctx, req.Name).Return(nil, domain.ErrChurchNotFound)
		mockChurchRepo.On("UnitOfWork", ctx, mock.Anything).Return(nil)
		mockAddressRepo.On("Create", (*sqlx.Tx)(nil), mock.AnythingOfType("*domain.Address")).Return(nil)
		mockChurchRepo.On("Create", (*sqlx.Tx)(nil), mock.MatchedBy(func(c *domain.Church) bool {
			return c.Name == req.Name &&
				c.ParishNumber.Valid &&
				c.ParishNumber.String == req.ParishNumber &&
				c.WebsiteURL.Valid &&
				c.WebsiteURL.String == req.WebsiteURL &&
				c.PhoneNumber.Valid &&
				c.PhoneNumber.String == req.PhoneNumber
		})).Return(nil)

		result, err := service.Create(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, sql.NullString{String: req.ParishNumber, Valid: true}, result.ParishNumber)
		assert.Equal(t, sql.NullString{String: req.WebsiteURL, Valid: true}, result.WebsiteURL)
		assert.Equal(t, sql.NullString{String: req.PhoneNumber, Valid: true}, result.PhoneNumber)
		mockChurchRepo.AssertExpectations(t)
		mockAddressRepo.AssertExpectations(t)
	})
}

package application

import (
	"context"
	"errors"
	"testing"

	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/entrepreneur/domain"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/entrepreneur/infrastructure/dto"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/helper/response"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// MockServiceRepository
type MockServiceRepository struct {
	mock.Mock
}

func (m *MockServiceRepository) Create(tx *sqlx.Tx, service *domain.Service) error {
	args := m.Called(tx, service)
	if args.Get(0) == nil {
		if service.ID == uuid.Nil {
			service.ID = uuid.New()
		}
	}
	return args.Error(0)
}

func (m *MockServiceRepository) Update(tx *sqlx.Tx, service *domain.Service) error {
	args := m.Called(tx, service)
	return args.Error(0)
}

func (m *MockServiceRepository) Delete(tx *sqlx.Tx, id uuid.UUID) error {
	args := m.Called(tx, id)
	return args.Error(0)
}

func (m *MockServiceRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Service, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Service), args.Error(1)
}

func (m *MockServiceRepository) List(ctx context.Context, filter *domain.ServiceFilters) ([]*domain.Service, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Service), args.Error(1)
}

func (m *MockServiceRepository) Count(ctx context.Context, filter *domain.ServiceFilters) (int, error) {
	args := m.Called(ctx, filter)
	return args.Int(0), args.Error(1)
}

func TestServiceService_Create(t *testing.T) {
	logger := zap.NewNop().Sugar()
	mockRepo := new(MockServiceRepository)
	service := NewServiceService(logger, mockRepo)
	ctx := context.Background()

	req := &dto.ServiceCreateRequest{
		BusinessID:  uuid.New(),
		Name:        "Test Service",
		Description: "Description",
		Price:       50.0,
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("Create", (*sqlx.Tx)(nil), mock.AnythingOfType("*domain.Service")).Return(nil)

		result, err := service.Create(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, req.Name, result.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failure", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("Create", (*sqlx.Tx)(nil), mock.AnythingOfType("*domain.Service")).Return(errors.New("db error"))

		result, err := service.Create(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, response.ErrInternalServerError, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestServiceService_Update(t *testing.T) {
	logger := zap.NewNop().Sugar()
	mockRepo := new(MockServiceRepository)
	service := NewServiceService(logger, mockRepo)
	ctx := context.Background()

	id := uuid.New()
	req := &dto.ServiceUpdateRequest{
		ID:          id,
		Name:        "Updated Service",
		Description: "Updated Desc",
		Price:       75.0,
	}

	existingService := &domain.Service{
		ID:          id,
		Name:        "Old Service",
		Description: "Old Desc",
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("GetByID", ctx, id).Return(existingService, nil)
		mockRepo.On("Update", (*sqlx.Tx)(nil), mock.AnythingOfType("*domain.Service")).Return(nil)

		err := service.Update(ctx, req)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("GetByID", ctx, id).Return(nil, domain.ErrServiceNotFound)

		err := service.Update(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, domain.ErrServiceNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestServiceService_Delete(t *testing.T) {
	logger := zap.NewNop().Sugar()
	mockRepo := new(MockServiceRepository)
	service := NewServiceService(logger, mockRepo)
	ctx := context.Background()
	id := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("Delete", (*sqlx.Tx)(nil), id).Return(nil)

		err := service.Delete(ctx, id)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failure", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("Delete", (*sqlx.Tx)(nil), id).Return(errors.New("db error"))

		err := service.Delete(ctx, id)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestServiceService_GetByID(t *testing.T) {
	logger := zap.NewNop().Sugar()
	mockRepo := new(MockServiceRepository)
	service := NewServiceService(logger, mockRepo)
	ctx := context.Background()
	id := uuid.New()

	expectedService := &domain.Service{
		ID:   id,
		Name: "Test Service",
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("GetByID", ctx, id).Return(expectedService, nil)

		result, err := service.GetByID(ctx, id)

		assert.NoError(t, err)
		assert.Equal(t, expectedService, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("GetByID", ctx, id).Return(nil, domain.ErrServiceNotFound)

		result, err := service.GetByID(ctx, id)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, domain.ErrServiceNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestServiceService_List(t *testing.T) {
	logger := zap.NewNop().Sugar()
	mockRepo := new(MockServiceRepository)
	service := NewServiceService(logger, mockRepo)
	ctx := context.Background()

	req := &dto.ServiceListRequest{
		Limit:  func() *int { i := 10; return &i }(),
		Offset: func() *int { i := 0; return &i }(),
	}

	expectedServices := []*domain.Service{
		{ID: uuid.New(), Name: "Service 1"},
		{ID: uuid.New(), Name: "Service 2"},
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("List", ctx, req).Return(expectedServices, nil)
		mockRepo.On("Count", ctx, req).Return(2, nil)

		result, err := service.List(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 2, result.Count)
		assert.Len(t, result.Services, 2)
		mockRepo.AssertExpectations(t)
	})

	t.Run("ListFailure", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("List", ctx, req).Return(nil, errors.New("db error"))

		result, err := service.List(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, response.ErrInternalServerError, err)
		mockRepo.AssertExpectations(t)
	})
}

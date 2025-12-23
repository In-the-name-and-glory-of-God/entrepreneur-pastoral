package application

import (
	"context"
	"errors"
	"testing"

	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/admin/domain"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/admin/infrastructure/dto"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/helper/response"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// MockIndustryRepository
type MockIndustryRepository struct {
	mock.Mock
}

func (m *MockIndustryRepository) Create(ctx context.Context, industry *domain.Industry) error {
	args := m.Called(ctx, industry)
	if args.Error(0) == nil {
		industry.ID = 1
	}
	return args.Error(0)
}

func (m *MockIndustryRepository) Update(ctx context.Context, industry *domain.Industry) error {
	args := m.Called(ctx, industry)
	return args.Error(0)
}

func (m *MockIndustryRepository) Delete(ctx context.Context, id int16) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockIndustryRepository) GetAll(ctx context.Context) ([]*domain.Industry, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Industry), args.Error(1)
}

func (m *MockIndustryRepository) GetByID(ctx context.Context, id int16) (*domain.Industry, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Industry), args.Error(1)
}

func (m *MockIndustryRepository) GetByName(ctx context.Context, name string) (*domain.Industry, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Industry), args.Error(1)
}

func TestIndustryService_Create(t *testing.T) {
	logger := zap.NewNop().Sugar()
	mockRepo := new(MockIndustryRepository)
	service := NewIndustryService(logger, mockRepo)
	ctx := context.Background()

	req := &dto.IndustryCreateRequest{
		Name: "Technology",
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("GetByName", ctx, req.Name).Return(nil, domain.ErrIndustryNotFound)
		mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.Industry")).Return(nil)

		result, err := service.Create(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, req.Name, result.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("AlreadyExists", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		existingIndustry := &domain.Industry{ID: 1, Name: req.Name}
		mockRepo.On("GetByName", ctx, req.Name).Return(existingIndustry, nil)

		result, err := service.Create(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, domain.ErrIndustryAlreadyExists, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("CreateFailure", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("GetByName", ctx, req.Name).Return(nil, domain.ErrIndustryNotFound)
		mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.Industry")).Return(errors.New("db error"))

		result, err := service.Create(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, response.ErrInternalServerError, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestIndustryService_Update(t *testing.T) {
	logger := zap.NewNop().Sugar()
	mockRepo := new(MockIndustryRepository)
	service := NewIndustryService(logger, mockRepo)
	ctx := context.Background()

	req := &dto.IndustryUpdateRequest{
		ID:   1,
		Name: "Updated Technology",
	}

	existingIndustry := &domain.Industry{
		ID:   1,
		Name: "Technology",
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("GetByID", ctx, req.ID).Return(existingIndustry, nil)
		mockRepo.On("GetByName", ctx, req.Name).Return(nil, domain.ErrIndustryNotFound)
		mockRepo.On("Update", ctx, mock.AnythingOfType("*domain.Industry")).Return(nil)

		err := service.Update(ctx, req)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("GetByID", ctx, req.ID).Return(nil, domain.ErrIndustryNotFound)

		err := service.Update(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, domain.ErrIndustryNotFound, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("NameAlreadyExists", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		differentIndustry := &domain.Industry{ID: 2, Name: req.Name}
		mockRepo.On("GetByID", ctx, req.ID).Return(existingIndustry, nil)
		mockRepo.On("GetByName", ctx, req.Name).Return(differentIndustry, nil)

		err := service.Update(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, domain.ErrIndustryAlreadyExists, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("UpdateSameNameSameID", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		sameIndustry := &domain.Industry{ID: req.ID, Name: req.Name}
		mockRepo.On("GetByID", ctx, req.ID).Return(existingIndustry, nil)
		mockRepo.On("GetByName", ctx, req.Name).Return(sameIndustry, nil)
		mockRepo.On("Update", ctx, mock.AnythingOfType("*domain.Industry")).Return(nil)

		err := service.Update(ctx, req)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("UpdateFailure", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("GetByID", ctx, req.ID).Return(existingIndustry, nil)
		mockRepo.On("GetByName", ctx, req.Name).Return(nil, domain.ErrIndustryNotFound)
		mockRepo.On("Update", ctx, mock.AnythingOfType("*domain.Industry")).Return(errors.New("db error"))

		err := service.Update(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, response.ErrInternalServerError, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestIndustryService_Delete(t *testing.T) {
	logger := zap.NewNop().Sugar()
	mockRepo := new(MockIndustryRepository)
	service := NewIndustryService(logger, mockRepo)
	ctx := context.Background()

	id := int16(1)
	existingIndustry := &domain.Industry{ID: id, Name: "Technology"}

	t.Run("Success", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("GetByID", ctx, id).Return(existingIndustry, nil)
		mockRepo.On("Delete", ctx, id).Return(nil)

		err := service.Delete(ctx, id)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("GetByID", ctx, id).Return(nil, domain.ErrIndustryNotFound)

		err := service.Delete(ctx, id)

		assert.Error(t, err)
		assert.Equal(t, domain.ErrIndustryNotFound, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("DeleteFailure", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("GetByID", ctx, id).Return(existingIndustry, nil)
		mockRepo.On("Delete", ctx, id).Return(errors.New("db error"))

		err := service.Delete(ctx, id)

		assert.Error(t, err)
		assert.Equal(t, response.ErrInternalServerError, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestIndustryService_GetByID(t *testing.T) {
	logger := zap.NewNop().Sugar()
	mockRepo := new(MockIndustryRepository)
	service := NewIndustryService(logger, mockRepo)
	ctx := context.Background()

	id := int16(1)
	expectedIndustry := &domain.Industry{ID: id, Name: "Technology"}

	t.Run("Success", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("GetByID", ctx, id).Return(expectedIndustry, nil)

		result, err := service.GetByID(ctx, id)

		assert.NoError(t, err)
		assert.Equal(t, expectedIndustry, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("GetByID", ctx, id).Return(nil, domain.ErrIndustryNotFound)

		result, err := service.GetByID(ctx, id)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, domain.ErrIndustryNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestIndustryService_GetAll(t *testing.T) {
	logger := zap.NewNop().Sugar()
	mockRepo := new(MockIndustryRepository)
	service := NewIndustryService(logger, mockRepo)
	ctx := context.Background()

	expectedIndustries := []*domain.Industry{
		{ID: 1, Name: "Technology"},
		{ID: 2, Name: "Healthcare"},
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("GetAll", ctx).Return(expectedIndustries, nil)

		result, err := service.GetAll(ctx)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Industries, 2)
		mockRepo.AssertExpectations(t)
	})

	t.Run("EmptyList", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("GetAll", ctx).Return([]*domain.Industry{}, nil)

		result, err := service.GetAll(ctx)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Industries, 0)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failure", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("GetAll", ctx).Return(nil, errors.New("db error"))

		result, err := service.GetAll(ctx)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, response.ErrInternalServerError, err)
		mockRepo.AssertExpectations(t)
	})
}

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

// MockFieldOfWorkRepository
type MockFieldOfWorkRepository struct {
	mock.Mock
}

func (m *MockFieldOfWorkRepository) Create(ctx context.Context, fieldOfWork *domain.FieldOfWork) error {
	args := m.Called(ctx, fieldOfWork)
	if args.Error(0) == nil {
		fieldOfWork.ID = 1
	}
	return args.Error(0)
}

func (m *MockFieldOfWorkRepository) Update(ctx context.Context, fieldOfWork *domain.FieldOfWork) error {
	args := m.Called(ctx, fieldOfWork)
	return args.Error(0)
}

func (m *MockFieldOfWorkRepository) Delete(ctx context.Context, id int16) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockFieldOfWorkRepository) GetAll(ctx context.Context) ([]*domain.FieldOfWork, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.FieldOfWork), args.Error(1)
}

func (m *MockFieldOfWorkRepository) GetByID(ctx context.Context, id int16) (*domain.FieldOfWork, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.FieldOfWork), args.Error(1)
}

func (m *MockFieldOfWorkRepository) GetByName(ctx context.Context, name string) (*domain.FieldOfWork, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.FieldOfWork), args.Error(1)
}

func TestFieldOfWorkService_Create(t *testing.T) {
	logger := zap.NewNop().Sugar()
	mockRepo := new(MockFieldOfWorkRepository)
	service := NewFieldOfWorkService(logger, mockRepo)
	ctx := context.Background()

	req := &dto.FieldOfWorkCreateRequest{
		Name: "Engineering",
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("GetByName", ctx, req.Name).Return(nil, domain.ErrFieldOfWorkNotFound)
		mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.FieldOfWork")).Return(nil)

		result, err := service.Create(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, req.Name, result.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("AlreadyExists", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		existingFieldOfWork := &domain.FieldOfWork{ID: 1, Name: req.Name}
		mockRepo.On("GetByName", ctx, req.Name).Return(existingFieldOfWork, nil)

		result, err := service.Create(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, domain.ErrFieldOfWorkAlreadyExists, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("CreateFailure", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("GetByName", ctx, req.Name).Return(nil, domain.ErrFieldOfWorkNotFound)
		mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.FieldOfWork")).Return(errors.New("db error"))

		result, err := service.Create(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, response.ErrInternalServerError, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("GetByNameInternalError", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("GetByName", ctx, req.Name).Return(nil, errors.New("db error"))

		result, err := service.Create(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, response.ErrInternalServerError, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestFieldOfWorkService_Update(t *testing.T) {
	logger := zap.NewNop().Sugar()
	mockRepo := new(MockFieldOfWorkRepository)
	service := NewFieldOfWorkService(logger, mockRepo)
	ctx := context.Background()

	req := &dto.FieldOfWorkUpdateRequest{
		ID:   1,
		Name: "Updated Engineering",
	}

	existingFieldOfWork := &domain.FieldOfWork{
		ID:   1,
		Name: "Engineering",
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("GetByID", ctx, req.ID).Return(existingFieldOfWork, nil)
		mockRepo.On("GetByName", ctx, req.Name).Return(nil, domain.ErrFieldOfWorkNotFound)
		mockRepo.On("Update", ctx, mock.AnythingOfType("*domain.FieldOfWork")).Return(nil)

		err := service.Update(ctx, req)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("GetByID", ctx, req.ID).Return(nil, domain.ErrFieldOfWorkNotFound)

		err := service.Update(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, domain.ErrFieldOfWorkNotFound, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("NameAlreadyExists", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		differentFieldOfWork := &domain.FieldOfWork{ID: 2, Name: req.Name}
		mockRepo.On("GetByID", ctx, req.ID).Return(existingFieldOfWork, nil)
		mockRepo.On("GetByName", ctx, req.Name).Return(differentFieldOfWork, nil)

		err := service.Update(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, domain.ErrFieldOfWorkAlreadyExists, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("UpdateSameNameSameID", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		sameFieldOfWork := &domain.FieldOfWork{ID: req.ID, Name: req.Name}
		mockRepo.On("GetByID", ctx, req.ID).Return(existingFieldOfWork, nil)
		mockRepo.On("GetByName", ctx, req.Name).Return(sameFieldOfWork, nil)
		mockRepo.On("Update", ctx, mock.AnythingOfType("*domain.FieldOfWork")).Return(nil)

		err := service.Update(ctx, req)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("UpdateFailure", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("GetByID", ctx, req.ID).Return(existingFieldOfWork, nil)
		mockRepo.On("GetByName", ctx, req.Name).Return(nil, domain.ErrFieldOfWorkNotFound)
		mockRepo.On("Update", ctx, mock.AnythingOfType("*domain.FieldOfWork")).Return(errors.New("db error"))

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
}

func TestFieldOfWorkService_Delete(t *testing.T) {
	logger := zap.NewNop().Sugar()
	mockRepo := new(MockFieldOfWorkRepository)
	service := NewFieldOfWorkService(logger, mockRepo)
	ctx := context.Background()

	id := int16(1)
	existingFieldOfWork := &domain.FieldOfWork{ID: id, Name: "Engineering"}

	t.Run("Success", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("GetByID", ctx, id).Return(existingFieldOfWork, nil)
		mockRepo.On("Delete", ctx, id).Return(nil)

		err := service.Delete(ctx, id)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("GetByID", ctx, id).Return(nil, domain.ErrFieldOfWorkNotFound)

		err := service.Delete(ctx, id)

		assert.Error(t, err)
		assert.Equal(t, domain.ErrFieldOfWorkNotFound, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("DeleteFailure", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("GetByID", ctx, id).Return(existingFieldOfWork, nil)
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

func TestFieldOfWorkService_GetByID(t *testing.T) {
	logger := zap.NewNop().Sugar()
	mockRepo := new(MockFieldOfWorkRepository)
	service := NewFieldOfWorkService(logger, mockRepo)
	ctx := context.Background()

	id := int16(1)
	expectedFieldOfWork := &domain.FieldOfWork{ID: id, Name: "Engineering"}

	t.Run("Success", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("GetByID", ctx, id).Return(expectedFieldOfWork, nil)

		result, err := service.GetByID(ctx, id)

		assert.NoError(t, err)
		assert.Equal(t, expectedFieldOfWork, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("GetByID", ctx, id).Return(nil, domain.ErrFieldOfWorkNotFound)

		result, err := service.GetByID(ctx, id)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, domain.ErrFieldOfWorkNotFound, err)
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

func TestFieldOfWorkService_GetAll(t *testing.T) {
	logger := zap.NewNop().Sugar()
	mockRepo := new(MockFieldOfWorkRepository)
	service := NewFieldOfWorkService(logger, mockRepo)
	ctx := context.Background()

	expectedFieldsOfWork := []*domain.FieldOfWork{
		{ID: 1, Name: "Engineering"},
		{ID: 2, Name: "Design"},
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("GetAll", ctx).Return(expectedFieldsOfWork, nil)

		result, err := service.GetAll(ctx)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.FieldsOfWork, 2)
		mockRepo.AssertExpectations(t)
	})

	t.Run("EmptyList", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("GetAll", ctx).Return([]*domain.FieldOfWork{}, nil)

		result, err := service.GetAll(ctx)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.FieldsOfWork, 0)
		mockRepo.AssertExpectations(t)
	})

	t.Run("NotFoundReturnsEmptyList", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("GetAll", ctx).Return(nil, domain.ErrFieldOfWorkNotFound)

		result, err := service.GetAll(ctx)

		assert.NoError(t, err)
		assert.NotNil(t, result)
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

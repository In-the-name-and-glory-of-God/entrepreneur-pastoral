package application

import (
	"context"
	"errors"
	"testing"

	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/entrepreneur/domain"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/entrepreneur/infrastructure/dto"
	userDto "github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/user/infrastructure/dto"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/helper/auth"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/helper/response"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// MockBusinessRepository
type MockBusinessRepository struct {
	mock.Mock
}

func (m *MockBusinessRepository) Create(tx *sqlx.Tx, business *domain.Business) error {
	args := m.Called(tx, business)
	if args.Get(0) == nil {
		if business.ID == uuid.Nil {
			business.ID = uuid.New()
		}
	}
	return args.Error(0)
}

func (m *MockBusinessRepository) Update(tx *sqlx.Tx, business *domain.Business) error {
	args := m.Called(tx, business)
	return args.Error(0)
}

func (m *MockBusinessRepository) Delete(tx *sqlx.Tx, id uuid.UUID) error {
	args := m.Called(tx, id)
	return args.Error(0)
}

func (m *MockBusinessRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Business, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Business), args.Error(1)
}

func (m *MockBusinessRepository) List(ctx context.Context, filter *domain.BusinessFilters) ([]*domain.Business, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Business), args.Error(1)
}

func (m *MockBusinessRepository) Count(ctx context.Context, filter *domain.BusinessFilters) (int, error) {
	args := m.Called(ctx, filter)
	return args.Int(0), args.Error(1)
}

func TestBusinessService_Create(t *testing.T) {
	logger := zap.NewNop().Sugar()
	mockRepo := new(MockBusinessRepository)
	service := NewBusinessService(logger, mockRepo)

	userID := uuid.New()
	userCtx := &userDto.UserAsContext{ID: userID}
	ctx := context.WithValue(context.Background(), auth.UserContextKey, userCtx)

	req := &dto.BusinessCreateRequest{
		UserID:      userID,
		IndustryID:  1,
		Name:        "Test Business",
		Description: "Description",
		Email:       "test@business.com",
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("Create", (*sqlx.Tx)(nil), mock.AnythingOfType("*domain.Business")).Return(nil)

		result, err := service.Create(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, req.Name, result.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failure", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("Create", (*sqlx.Tx)(nil), mock.AnythingOfType("*domain.Business")).Return(errors.New("db error"))

		result, err := service.Create(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, response.ErrInternalServerError, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestBusinessService_Update(t *testing.T) {
	logger := zap.NewNop().Sugar()
	mockRepo := new(MockBusinessRepository)
	service := NewBusinessService(logger, mockRepo)

	userID := uuid.New()
	userCtx := &userDto.UserAsContext{ID: userID}
	ctx := context.WithValue(context.Background(), auth.UserContextKey, userCtx)

	id := uuid.New()
	req := &dto.BusinessUpdateRequest{
		ID:          id,
		IndustryID:  2,
		Name:        "Updated Name",
		Description: "Updated Desc",
		Email:       "updated@email.com",
		IsActive:    true,
	}

	existingBusiness := &domain.Business{
		ID:          id,
		UserID:      userID,
		Name:        "Old Name",
		Description: "Old Desc",
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("GetByID", ctx, id).Return(existingBusiness, nil)
		mockRepo.On("Update", (*sqlx.Tx)(nil), mock.AnythingOfType("*domain.Business")).Return(nil)

		err := service.Update(ctx, req)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("GetByID", ctx, id).Return(nil, domain.ErrBusinessNotFound)

		err := service.Update(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, domain.ErrBusinessNotFound, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		existingBusiness.UserID = uuid.New()
		mockRepo.On("GetByID", ctx, id).Return(existingBusiness, nil)

		err := service.Update(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, domain.ErrUnauthorized, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestBusinessService_Delete(t *testing.T) {
	logger := zap.NewNop().Sugar()
	mockRepo := new(MockBusinessRepository)
	service := NewBusinessService(logger, mockRepo)

	userID := uuid.New()
	userCtx := &userDto.UserAsContext{ID: userID}
	ctx := context.WithValue(context.Background(), auth.UserContextKey, userCtx)

	id := uuid.New()
	existingBusiness := &domain.Business{
		ID:     id,
		UserID: userID,
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("GetByID", ctx, id).Return(existingBusiness, nil)
		mockRepo.On("Delete", (*sqlx.Tx)(nil), id).Return(nil)

		err := service.Delete(ctx, id)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failure", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("GetByID", ctx, id).Return(existingBusiness, nil)
		mockRepo.On("Delete", (*sqlx.Tx)(nil), id).Return(errors.New("db error"))

		err := service.Delete(ctx, id)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		existingBusiness.UserID = uuid.New()
		mockRepo.On("GetByID", ctx, id).Return(existingBusiness, nil)

		err := service.Delete(ctx, id)

		assert.Error(t, err)
		assert.Equal(t, domain.ErrUnauthorized, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestBusinessService_GetByID(t *testing.T) {
	logger := zap.NewNop().Sugar()
	mockRepo := new(MockBusinessRepository)
	service := NewBusinessService(logger, mockRepo)
	ctx := context.Background()
	id := uuid.New()

	expectedBusiness := &domain.Business{
		ID:   id,
		Name: "Test Business",
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("GetByID", ctx, id).Return(expectedBusiness, nil)

		result, err := service.GetByID(ctx, id)

		assert.NoError(t, err)
		assert.Equal(t, expectedBusiness, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("GetByID", ctx, id).Return(nil, domain.ErrBusinessNotFound)

		result, err := service.GetByID(ctx, id)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, domain.ErrBusinessNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestBusinessService_List(t *testing.T) {
	logger := zap.NewNop().Sugar()
	mockRepo := new(MockBusinessRepository)
	service := NewBusinessService(logger, mockRepo)
	ctx := context.Background()

	req := &dto.BusinessListRequest{
		Limit:  func() *int { i := 10; return &i }(),
		Offset: func() *int { i := 0; return &i }(),
	}

	expectedBusinesses := []*domain.Business{
		{ID: uuid.New(), Name: "Business 1"},
		{ID: uuid.New(), Name: "Business 2"},
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("List", ctx, req).Return(expectedBusinesses, nil)
		mockRepo.On("Count", ctx, req).Return(2, nil)

		result, err := service.List(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 2, result.Count)
		assert.Len(t, result.Businesses, 2)
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

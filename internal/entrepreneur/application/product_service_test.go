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

// MockProductRepository
type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) Create(tx *sqlx.Tx, product *domain.Product) error {
	args := m.Called(tx, product)
	if args.Get(0) == nil {
		if product.ID == uuid.Nil {
			product.ID = uuid.New()
		}
	}
	return args.Error(0)
}

func (m *MockProductRepository) Update(tx *sqlx.Tx, product *domain.Product) error {
	args := m.Called(tx, product)
	return args.Error(0)
}

func (m *MockProductRepository) Delete(tx *sqlx.Tx, id uuid.UUID) error {
	args := m.Called(tx, id)
	return args.Error(0)
}

func (m *MockProductRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *MockProductRepository) List(ctx context.Context, filter *domain.ProductFilters) ([]*domain.Product, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Product), args.Error(1)
}

func (m *MockProductRepository) Count(ctx context.Context, filter *domain.ProductFilters) (int, error) {
	args := m.Called(ctx, filter)
	return args.Int(0), args.Error(1)
}

func TestProductService_Create(t *testing.T) {
	logger := zap.NewNop().Sugar()
	mockRepo := new(MockProductRepository)
	mockBusinessRepo := new(MockBusinessRepository)
	service := NewProductService(logger, mockRepo, mockBusinessRepo)

	userID := uuid.New()
	userCtx := &userDto.UserAsContext{ID: userID}
	ctx := context.WithValue(context.Background(), auth.UserContextKey, userCtx)

	req := &dto.ProductCreateRequest{
		BusinessID:  uuid.New(),
		Name:        "Test Product",
		Description: "Description",
		Price:       100.0,
		IsAvailable: true,
	}

	t.Run("Success", func(t *testing.T) {
		mockBusinessRepo.On("GetByID", ctx, req.BusinessID).Return(&domain.Business{ID: req.BusinessID, UserID: userID}, nil)
		mockRepo.On("Create", (*sqlx.Tx)(nil), mock.AnythingOfType("*domain.Product")).Return(nil)

		result, err := service.Create(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, req.Name, result.Name)
		mockRepo.AssertExpectations(t)
		mockBusinessRepo.AssertExpectations(t)
	})

	t.Run("Failure_BusinessNotFound", func(t *testing.T) {
		mockBusinessRepo.ExpectedCalls = nil
		mockRepo.ExpectedCalls = nil
		mockBusinessRepo.On("GetByID", ctx, req.BusinessID).Return(nil, domain.ErrBusinessNotFound)

		result, err := service.Create(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, domain.ErrBusinessNotFound, err)
		mockBusinessRepo.AssertExpectations(t)
	})

	t.Run("Failure_Unauthorized", func(t *testing.T) {
		mockBusinessRepo.ExpectedCalls = nil
		mockRepo.ExpectedCalls = nil
		mockBusinessRepo.On("GetByID", ctx, req.BusinessID).Return(&domain.Business{ID: req.BusinessID, UserID: uuid.New()}, nil)

		result, err := service.Create(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, domain.ErrUnauthorized, err)
		mockBusinessRepo.AssertExpectations(t)
	})

	t.Run("Failure_RepoError", func(t *testing.T) {
		mockBusinessRepo.ExpectedCalls = nil
		mockRepo.ExpectedCalls = nil
		mockBusinessRepo.On("GetByID", ctx, req.BusinessID).Return(&domain.Business{ID: req.BusinessID, UserID: userID}, nil)
		mockRepo.On("Create", (*sqlx.Tx)(nil), mock.AnythingOfType("*domain.Product")).Return(errors.New("db error"))

		result, err := service.Create(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, response.ErrInternalServerError, err)
		mockRepo.AssertExpectations(t)
		mockBusinessRepo.AssertExpectations(t)
	})
}

func TestProductService_Update(t *testing.T) {
	logger := zap.NewNop().Sugar()
	mockRepo := new(MockProductRepository)
	mockBusinessRepo := new(MockBusinessRepository)
	service := NewProductService(logger, mockRepo, mockBusinessRepo)

	userID := uuid.New()
	userCtx := &userDto.UserAsContext{ID: userID}
	ctx := context.WithValue(context.Background(), auth.UserContextKey, userCtx)

	id := uuid.New()
	businessID := uuid.New()
	req := &dto.ProductUpdateRequest{
		ID:          id,
		Name:        "Updated Product",
		Description: "Updated Desc",
		Price:       150.0,
		IsAvailable: false,
	}

	existingProduct := &domain.Product{
		ID:          id,
		BusinessID:  businessID,
		Name:        "Old Product",
		Description: "Old Desc",
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("GetByID", ctx, id).Return(existingProduct, nil)
		mockBusinessRepo.On("GetByID", ctx, businessID).Return(&domain.Business{ID: businessID, UserID: userID}, nil)
		mockRepo.On("Update", (*sqlx.Tx)(nil), mock.AnythingOfType("*domain.Product")).Return(nil)

		err := service.Update(ctx, req)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockBusinessRepo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockBusinessRepo.ExpectedCalls = nil
		mockRepo.On("GetByID", ctx, id).Return(nil, domain.ErrProductNotFound)

		err := service.Update(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, domain.ErrProductNotFound, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockBusinessRepo.ExpectedCalls = nil
		mockRepo.On("GetByID", ctx, id).Return(existingProduct, nil)
		mockBusinessRepo.On("GetByID", ctx, businessID).Return(&domain.Business{ID: businessID, UserID: uuid.New()}, nil)

		err := service.Update(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, domain.ErrUnauthorized, err)
		mockRepo.AssertExpectations(t)
		mockBusinessRepo.AssertExpectations(t)
	})
}

func TestProductService_Delete(t *testing.T) {
	logger := zap.NewNop().Sugar()
	mockRepo := new(MockProductRepository)
	mockBusinessRepo := new(MockBusinessRepository)
	service := NewProductService(logger, mockRepo, mockBusinessRepo)

	userID := uuid.New()
	userCtx := &userDto.UserAsContext{ID: userID}
	ctx := context.WithValue(context.Background(), auth.UserContextKey, userCtx)

	id := uuid.New()
	businessID := uuid.New()
	existingProduct := &domain.Product{
		ID:         id,
		BusinessID: businessID,
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("GetByID", ctx, id).Return(existingProduct, nil)
		mockBusinessRepo.On("GetByID", ctx, businessID).Return(&domain.Business{ID: businessID, UserID: userID}, nil)
		mockRepo.On("Delete", (*sqlx.Tx)(nil), id).Return(nil)

		err := service.Delete(ctx, id)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockBusinessRepo.AssertExpectations(t)
	})

	t.Run("Failure", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockBusinessRepo.ExpectedCalls = nil
		mockRepo.On("GetByID", ctx, id).Return(existingProduct, nil)
		mockBusinessRepo.On("GetByID", ctx, businessID).Return(&domain.Business{ID: businessID, UserID: userID}, nil)
		mockRepo.On("Delete", (*sqlx.Tx)(nil), id).Return(errors.New("db error"))

		err := service.Delete(ctx, id)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
		mockBusinessRepo.AssertExpectations(t)
	})
}

func TestProductService_GetByID(t *testing.T) {
	logger := zap.NewNop().Sugar()
	mockRepo := new(MockProductRepository)
	mockBusinessRepo := new(MockBusinessRepository)
	service := NewProductService(logger, mockRepo, mockBusinessRepo)
	ctx := context.Background()
	id := uuid.New()

	expectedProduct := &domain.Product{
		ID:   id,
		Name: "Test Product",
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("GetByID", ctx, id).Return(expectedProduct, nil)

		result, err := service.GetByID(ctx, id)

		assert.NoError(t, err)
		assert.Equal(t, expectedProduct, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("GetByID", ctx, id).Return(nil, domain.ErrProductNotFound)

		result, err := service.GetByID(ctx, id)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, domain.ErrProductNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestProductService_List(t *testing.T) {
	logger := zap.NewNop().Sugar()
	mockRepo := new(MockProductRepository)
	mockBusinessRepo := new(MockBusinessRepository)
	service := NewProductService(logger, mockRepo, mockBusinessRepo)
	ctx := context.Background()

	req := &dto.ProductListRequest{
		Limit:  func() *int { i := 10; return &i }(),
		Offset: func() *int { i := 0; return &i }(),
	}

	expectedProducts := []*domain.Product{
		{ID: uuid.New(), Name: "Product 1"},
		{ID: uuid.New(), Name: "Product 2"},
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("List", ctx, req).Return(expectedProducts, nil)
		mockRepo.On("Count", ctx, req).Return(2, nil)

		result, err := service.List(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 2, result.Count)
		assert.Len(t, result.Products, 2)
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

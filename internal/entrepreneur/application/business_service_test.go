package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/entrepreneur/domain"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/entrepreneur/infrastructure/dto"
	userDto "github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/user/infrastructure/dto"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/helper/auth"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/helper/response"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/storage"
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

func (m *MockBusinessRepository) UpdateProperty(ctx context.Context, id uuid.UUID, property domain.BusinessProperty, value any) error {
	args := m.Called(ctx, id, property, value)
	return args.Error(0)
}

// MockCacheStorage
type MockCacheStorage struct {
	mock.Mock
}

func (m *MockCacheStorage) BuildKey(prefix storage.CachePrefix, data ...string) string {
	args := m.Called(prefix, data)
	return args.String(0)
}

func (m *MockCacheStorage) Get(ctx context.Context, key string, dest any) error {
	args := m.Called(ctx, key, dest)
	return args.Error(0)
}

func (m *MockCacheStorage) GetAndDel(ctx context.Context, key string, dest any) error {
	args := m.Called(ctx, key, dest)
	return args.Error(0)
}

func (m *MockCacheStorage) GetString(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *MockCacheStorage) GetStringAndDel(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *MockCacheStorage) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func (m *MockCacheStorage) SetString(ctx context.Context, key string, value string, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func (m *MockCacheStorage) Del(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockCacheStorage) Scan(ctx context.Context, match string) ([]string, error) {
	args := m.Called(ctx, match)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockCacheStorage) Exists(ctx context.Context, key string) (bool, error) {
	args := m.Called(ctx, key)
	return args.Bool(0), args.Error(1)
}

func TestBusinessService_Create(t *testing.T) {
	logger := zap.NewNop().Sugar()
	mockRepo := new(MockBusinessRepository)
	mockCache := new(MockCacheStorage)
	service := NewBusinessService(logger, mockCache, mockRepo)

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
		mockCache.On("BuildKey", storage.CACHE_PREFIX_BUSINESS_LIST, mock.Anything).Return("business_list:*")
		mockCache.On("Scan", ctx, "business_list:*").Return([]string{}, nil)

		result, err := service.Create(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, req.Name, result.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failure", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockCache.ExpectedCalls = nil
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
	mockCache := new(MockCacheStorage)
	service := NewBusinessService(logger, mockCache, mockRepo)

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
		mockCache.On("BuildKey", storage.CACHE_PREFIX_BUSINESS, mock.Anything).Return("business:" + id.String())
		mockCache.On("Del", ctx, "business:"+id.String()).Return(nil)
		mockCache.On("BuildKey", storage.CACHE_PREFIX_BUSINESS_LIST, mock.Anything).Return("business_list:*")
		mockCache.On("Scan", ctx, "business_list:*").Return([]string{}, nil)

		err := service.Update(ctx, req)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockCache.ExpectedCalls = nil
		mockRepo.On("GetByID", ctx, id).Return(nil, domain.ErrBusinessNotFound)

		err := service.Update(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, domain.ErrBusinessNotFound, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockCache.ExpectedCalls = nil
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
	mockCache := new(MockCacheStorage)
	service := NewBusinessService(logger, mockCache, mockRepo)

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
		mockCache.On("BuildKey", storage.CACHE_PREFIX_BUSINESS, mock.Anything).Return("business:" + id.String())
		mockCache.On("Del", ctx, "business:"+id.String()).Return(nil)
		mockCache.On("BuildKey", storage.CACHE_PREFIX_BUSINESS_LIST, mock.Anything).Return("business_list:*")
		mockCache.On("Scan", ctx, "business_list:*").Return([]string{}, nil)

		err := service.Delete(ctx, id)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failure", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockCache.ExpectedCalls = nil
		mockRepo.On("GetByID", ctx, id).Return(existingBusiness, nil)
		mockRepo.On("Delete", (*sqlx.Tx)(nil), id).Return(errors.New("db error"))

		err := service.Delete(ctx, id)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockCache.ExpectedCalls = nil
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
	mockCache := new(MockCacheStorage)
	service := NewBusinessService(logger, mockCache, mockRepo)
	ctx := context.Background()
	id := uuid.New()

	expectedBusiness := &domain.Business{
		ID:   id,
		Name: "Test Business",
	}

	t.Run("Success", func(t *testing.T) {
		cacheKey := "business:" + id.String()
		mockCache.On("BuildKey", storage.CACHE_PREFIX_BUSINESS, mock.Anything).Return(cacheKey)
		mockCache.On("Get", ctx, cacheKey, mock.AnythingOfType("*domain.Business")).Return(errors.New("cache miss"))
		mockRepo.On("GetByID", ctx, id).Return(expectedBusiness, nil)
		mockCache.On("Set", ctx, cacheKey, expectedBusiness, mock.Anything).Return(nil)

		result, err := service.GetByID(ctx, id)

		assert.NoError(t, err)
		assert.Equal(t, expectedBusiness, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockCache.ExpectedCalls = nil
		cacheKey := "business:" + id.String()
		mockCache.On("BuildKey", storage.CACHE_PREFIX_BUSINESS, mock.Anything).Return(cacheKey)
		mockCache.On("Get", ctx, cacheKey, mock.AnythingOfType("*domain.Business")).Return(errors.New("cache miss"))
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
	mockCache := new(MockCacheStorage)
	service := NewBusinessService(logger, mockCache, mockRepo)
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
		cacheKey := "business_list:10:0"
		mockCache.On("BuildKey", storage.CACHE_PREFIX_BUSINESS_LIST, mock.Anything).Return(cacheKey)
		mockCache.On("Get", ctx, cacheKey, mock.AnythingOfType("*dto.BusinessListResponse")).Return(errors.New("cache miss"))
		mockRepo.On("List", ctx, req).Return(expectedBusinesses, nil)
		mockRepo.On("Count", ctx, req).Return(2, nil)
		mockCache.On("Set", ctx, cacheKey, mock.AnythingOfType("*dto.BusinessListResponse"), mock.Anything).Return(nil)

		result, err := service.List(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 2, result.Count)
		assert.Len(t, result.Businesses, 2)
		mockRepo.AssertExpectations(t)
	})

	t.Run("ListFailure", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockCache.ExpectedCalls = nil
		cacheKey := "business_list:10:0"
		mockCache.On("BuildKey", storage.CACHE_PREFIX_BUSINESS_LIST, mock.Anything).Return(cacheKey)
		mockCache.On("Get", ctx, cacheKey, mock.AnythingOfType("*dto.BusinessListResponse")).Return(errors.New("cache miss"))
		mockRepo.On("List", ctx, req).Return(nil, errors.New("db error"))

		result, err := service.List(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, response.ErrInternalServerError, err)
		mockRepo.AssertExpectations(t)
	})
}

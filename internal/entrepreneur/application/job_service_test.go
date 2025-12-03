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

// MockJobRepository
type MockJobRepository struct {
	mock.Mock
}

func (m *MockJobRepository) Create(tx *sqlx.Tx, job *domain.Job) error {
	args := m.Called(tx, job)
	if args.Get(0) == nil {
		if job.ID == uuid.Nil {
			job.ID = uuid.New()
		}
	}
	return args.Error(0)
}

func (m *MockJobRepository) Update(tx *sqlx.Tx, job *domain.Job) error {
	args := m.Called(tx, job)
	return args.Error(0)
}

func (m *MockJobRepository) Delete(tx *sqlx.Tx, id uuid.UUID) error {
	args := m.Called(tx, id)
	return args.Error(0)
}

func (m *MockJobRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Job, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Job), args.Error(1)
}

func (m *MockJobRepository) List(ctx context.Context, filter *domain.JobFilters) ([]*domain.Job, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Job), args.Error(1)
}

func (m *MockJobRepository) Count(ctx context.Context, filter *domain.JobFilters) (int, error) {
	args := m.Called(ctx, filter)
	return args.Int(0), args.Error(1)
}

func TestJobService_Create(t *testing.T) {
	logger := zap.NewNop().Sugar()
	mockRepo := new(MockJobRepository)
	service := NewJobService(logger, mockRepo)
	ctx := context.Background()

	req := &dto.JobCreateRequest{
		BusinessID:  uuid.New(),
		Title:       "Test Job",
		Description: "Description",
		Location:    domain.JobLocationRemote,
		Type:        domain.JobTypeFullTime,
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("Create", (*sqlx.Tx)(nil), mock.AnythingOfType("*domain.Job")).Return(nil)

		result, err := service.Create(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, req.Title, result.Title)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failure", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("Create", (*sqlx.Tx)(nil), mock.AnythingOfType("*domain.Job")).Return(errors.New("db error"))

		result, err := service.Create(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, response.ErrInternalServerError, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestJobService_Update(t *testing.T) {
	logger := zap.NewNop().Sugar()
	mockRepo := new(MockJobRepository)
	service := NewJobService(logger, mockRepo)
	ctx := context.Background()

	id := uuid.New()
	req := &dto.JobUpdateRequest{
		ID:          id,
		Title:       "Updated Job",
		Description: "Updated Desc",
	}

	existingJob := &domain.Job{
		ID:          id,
		Title:       "Old Job",
		Description: "Old Desc",
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("GetByID", ctx, id).Return(existingJob, nil)
		mockRepo.On("Update", (*sqlx.Tx)(nil), mock.AnythingOfType("*domain.Job")).Return(nil)

		err := service.Update(ctx, req)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("GetByID", ctx, id).Return(nil, domain.ErrJobNotFound)

		err := service.Update(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, domain.ErrJobNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestJobService_Delete(t *testing.T) {
	logger := zap.NewNop().Sugar()
	mockRepo := new(MockJobRepository)
	service := NewJobService(logger, mockRepo)
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

func TestJobService_GetByID(t *testing.T) {
	logger := zap.NewNop().Sugar()
	mockRepo := new(MockJobRepository)
	service := NewJobService(logger, mockRepo)
	ctx := context.Background()
	id := uuid.New()

	expectedJob := &domain.Job{
		ID:    id,
		Title: "Test Job",
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("GetByID", ctx, id).Return(expectedJob, nil)

		result, err := service.GetByID(ctx, id)

		assert.NoError(t, err)
		assert.Equal(t, expectedJob, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("GetByID", ctx, id).Return(nil, domain.ErrJobNotFound)

		result, err := service.GetByID(ctx, id)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, domain.ErrJobNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestJobService_List(t *testing.T) {
	logger := zap.NewNop().Sugar()
	mockRepo := new(MockJobRepository)
	service := NewJobService(logger, mockRepo)
	ctx := context.Background()

	req := &dto.JobListRequest{
		Limit:  func() *int { i := 10; return &i }(),
		Offset: func() *int { i := 0; return &i }(),
	}

	expectedJobs := []*domain.Job{
		{ID: uuid.New(), Title: "Job 1"},
		{ID: uuid.New(), Title: "Job 2"},
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("List", ctx, req).Return(expectedJobs, nil)
		mockRepo.On("Count", ctx, req).Return(2, nil)

		result, err := service.List(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 2, result.Count)
		assert.Len(t, result.Jobs, 2)
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

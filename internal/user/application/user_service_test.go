package application

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	adminDomain "github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/admin/domain"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/user/domain"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/user/infrastructure/dto"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/helper/response"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// Mock repositories
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) UnitOfWork(ctx context.Context, fn func(*sqlx.Tx) error) error {
	args := m.Called(ctx, fn)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	// Execute the function with a nil tx for testing purposes
	return fn(nil)
}

func (m *MockUserRepository) Create(tx *sqlx.Tx, user *domain.User) error {
	args := m.Called(tx, user)
	if args.Get(0) == nil {
		// Simulate ID generation
		if user.ID == uuid.Nil {
			user.ID = uuid.New()
		}
	}
	return args.Error(0)
}

func (m *MockUserRepository) Update(tx *sqlx.Tx, user *domain.User) error {
	args := m.Called(tx, user)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateProperty(ctx context.Context, id uuid.UUID, property domain.UserProperty, value any) error {
	args := m.Called(ctx, id, property, value)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByDocumentID(ctx context.Context, documentID string) (*domain.User, error) {
	args := m.Called(ctx, documentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetAllByRoleID(ctx context.Context, roleID int16) ([]*domain.User, error) {
	args := m.Called(ctx, roleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetAllByIsActive(ctx context.Context, isActive bool) ([]*domain.User, error) {
	args := m.Called(ctx, isActive)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetAllByIsVerified(ctx context.Context, isVerified bool) ([]*domain.User, error) {
	args := m.Called(ctx, isVerified)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetAllByIsCatholic(ctx context.Context, isCatholic bool) ([]*domain.User, error) {
	args := m.Called(ctx, isCatholic)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetAllByIsEntrepreneur(ctx context.Context, isEntrepreneur bool) ([]*domain.User, error) {
	args := m.Called(ctx, isEntrepreneur)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.User), args.Error(1)
}

func (m *MockUserRepository) List(ctx context.Context, filter *domain.UserFilters) ([]*domain.User, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.User), args.Error(1)
}

func (m *MockUserRepository) Count(ctx context.Context, filter *domain.UserFilters) (int, error) {
	args := m.Called(ctx, filter)
	return args.Int(0), args.Error(1)
}

type MockNotificationPreferencesRepository struct {
	mock.Mock
}

func (m *MockNotificationPreferencesRepository) Create(tx *sqlx.Tx, notificationPreferences *domain.NotificationPreferences) error {
	args := m.Called(tx, notificationPreferences)
	return args.Error(0)
}

func (m *MockNotificationPreferencesRepository) Update(tx *sqlx.Tx, notificationPreferences *domain.NotificationPreferences) error {
	args := m.Called(tx, notificationPreferences)
	return args.Error(0)
}

func (m *MockNotificationPreferencesRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.NotificationPreferences, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.NotificationPreferences), args.Error(1)
}

type MockJobProfileRepository struct {
	mock.Mock
}

func (m *MockJobProfileRepository) Create(tx *sqlx.Tx, jobProfile *domain.JobProfile) error {
	args := m.Called(tx, jobProfile)
	return args.Error(0)
}

func (m *MockJobProfileRepository) Update(tx *sqlx.Tx, jobProfile *domain.JobProfile) error {
	args := m.Called(tx, jobProfile)
	return args.Error(0)
}

func (m *MockJobProfileRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.JobProfile, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.JobProfile), args.Error(1)
}

func (m *MockJobProfileRepository) GetAllOpenToWork(ctx context.Context) ([]*domain.JobProfile, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.JobProfile), args.Error(1)
}

// Test helpers
func setupTest() (*UserService, *MockUserRepository, *MockNotificationPreferencesRepository, *MockJobProfileRepository) {
	logger := zap.NewNop().Sugar()
	mockUserRepo := new(MockUserRepository)
	mockNotifPrefRepo := new(MockNotificationPreferencesRepository)
	mockJobProfileRepo := new(MockJobProfileRepository)

	service := NewUserService(logger, mockUserRepo, mockNotifPrefRepo, mockJobProfileRepo)

	return service, mockUserRepo, mockNotifPrefRepo, mockJobProfileRepo
}

// Test Create
func TestUserService_Create_Success(t *testing.T) {
	service, mockUserRepo, mockNotifPrefRepo, mockJobProfileRepo := setupTest()
	ctx := context.Background()

	req := &dto.UserRegisterRequest{
		FirstName:        "John",
		LastName:         "Doe",
		Email:            "john.doe@example.com",
		Password:         "SecurePassword123!",
		DocumentID:       "123456789",
		PhoneCountryCode: "+1",
		PhoneNumber:      "5551234567",
		OpenToWork:       true,
		CVPath:           "/path/to/cv.pdf",
		FieldsOfWork: []adminDomain.FieldOfWork{
			{ID: 1, Name: "Engineering"},
		},
	}

	// Mock expectations
	mockUserRepo.On("GetByEmail", ctx, req.Email).Return(nil, domain.ErrUserNotFound)
	mockUserRepo.On("GetByDocumentID", ctx, req.DocumentID).Return(nil, domain.ErrUserNotFound)
	mockUserRepo.On("UnitOfWork", ctx, mock.AnythingOfType("func(*sqlx.Tx) error")).Return(nil)
	mockUserRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)
	mockNotifPrefRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.NotificationPreferences")).Return(nil)
	mockJobProfileRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.JobProfile")).Return(nil)

	err := service.Create(ctx, req)

	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
	mockNotifPrefRepo.AssertExpectations(t)
	mockJobProfileRepo.AssertExpectations(t)
}

func TestUserService_Create_EmailAlreadyExists(t *testing.T) {
	service, mockUserRepo, _, _ := setupTest()
	ctx := context.Background()

	req := &dto.UserRegisterRequest{
		Email:      "existing@example.com",
		DocumentID: "123456789",
	}

	existingUser := &domain.User{
		ID:    uuid.New(),
		Email: req.Email,
	}

	mockUserRepo.On("GetByEmail", ctx, req.Email).Return(existingUser, nil)

	err := service.Create(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrEmailAlreadyExists, err)
	mockUserRepo.AssertExpectations(t)
}

func TestUserService_Create_DocumentIDAlreadyExists(t *testing.T) {
	service, mockUserRepo, _, _ := setupTest()
	ctx := context.Background()

	req := &dto.UserRegisterRequest{
		Email:      "john@example.com",
		DocumentID: "123456789",
	}

	existingUser := &domain.User{
		ID:         uuid.New(),
		DocumentID: req.DocumentID,
	}

	mockUserRepo.On("GetByEmail", ctx, req.Email).Return(nil, domain.ErrUserNotFound)
	mockUserRepo.On("GetByDocumentID", ctx, req.DocumentID).Return(existingUser, nil)

	err := service.Create(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrDocumentIDAlreadyExists, err)
	mockUserRepo.AssertExpectations(t)
}

func TestUserService_Create_DatabaseError(t *testing.T) {
	service, mockUserRepo, _, _ := setupTest()
	ctx := context.Background()

	req := &dto.UserRegisterRequest{
		Email:      "john@example.com",
		DocumentID: "123456789",
	}

	mockUserRepo.On("GetByEmail", ctx, req.Email).Return(nil, errors.New("database error"))

	err := service.Create(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, response.ErrInternalServerError, err)
	mockUserRepo.AssertExpectations(t)
}

func TestUserService_Create_UserCreationFailed(t *testing.T) {
	service, mockUserRepo, _, _ := setupTest()
	ctx := context.Background()

	req := &dto.UserRegisterRequest{
		FirstName:  "John",
		LastName:   "Doe",
		Email:      "john@example.com",
		Password:   "SecurePassword123!",
		DocumentID: "123456789",
	}

	mockUserRepo.On("GetByEmail", ctx, req.Email).Return(nil, domain.ErrUserNotFound)
	mockUserRepo.On("GetByDocumentID", ctx, req.DocumentID).Return(nil, domain.ErrUserNotFound)
	mockUserRepo.On("UnitOfWork", ctx, mock.AnythingOfType("func(*sqlx.Tx) error")).
		Return(errors.New("database error"))

	err := service.Create(ctx, req)

	assert.Error(t, err)
	mockUserRepo.AssertExpectations(t)
}

// Test Update
func TestUserService_Update_Success(t *testing.T) {
	service, mockUserRepo, mockNotifPrefRepo, mockJobProfileRepo := setupTest()
	ctx := context.Background()

	userID := uuid.New()
	req := &dto.UserUpdateRequest{
		ID:               userID,
		FirstName:        "Jane",
		LastName:         "Doe",
		Email:            "jane.doe@example.com",
		DocumentID:       "987654321",
		PhoneCountryCode: "+1",
		PhoneNumber:      "5559876543",
		NotifyByEmail:    true,
		NotifyBySms:      false,
		OpenToWork:       true,
		CVPath:           "/path/to/new_cv.pdf",
		FieldsOfWork: []adminDomain.FieldOfWork{
			{ID: 2, Name: "Marketing"},
		},
	}

	existingUser := &domain.User{
		ID:         userID,
		FirstName:  "John",
		LastName:   "Doe",
		Email:      "john.doe@example.com",
		DocumentID: "123456789",
	}

	mockUserRepo.On("GetByID", ctx, userID).Return(existingUser, nil)
	mockUserRepo.On("UnitOfWork", ctx, mock.AnythingOfType("func(*sqlx.Tx) error")).Return(nil)
	mockUserRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)
	mockNotifPrefRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.NotificationPreferences")).Return(nil)
	mockJobProfileRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.JobProfile")).Return(nil)

	err := service.Update(ctx, req)

	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
	mockNotifPrefRepo.AssertExpectations(t)
	mockJobProfileRepo.AssertExpectations(t)
}

func TestUserService_Update_UserNotFound(t *testing.T) {
	service, mockUserRepo, _, _ := setupTest()
	ctx := context.Background()

	userID := uuid.New()
	req := &dto.UserUpdateRequest{
		ID: userID,
	}

	mockUserRepo.On("GetByID", ctx, userID).Return(nil, domain.ErrUserNotFound)

	err := service.Update(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrUserNotFound, err)
	mockUserRepo.AssertExpectations(t)
}

func TestUserService_Update_UpdateFailed(t *testing.T) {
	service, mockUserRepo, _, _ := setupTest()
	ctx := context.Background()

	userID := uuid.New()
	req := &dto.UserUpdateRequest{
		ID: userID,
	}

	existingUser := &domain.User{
		ID: userID,
	}

	mockUserRepo.On("GetByID", ctx, userID).Return(existingUser, nil)
	mockUserRepo.On("UnitOfWork", ctx, mock.AnythingOfType("func(*sqlx.Tx) error")).
		Return(errors.New("database error"))

	err := service.Update(ctx, req)

	assert.Error(t, err)
	mockUserRepo.AssertExpectations(t)
}

// Test GetByID
func TestUserService_GetByID_Success(t *testing.T) {
	service, mockUserRepo, mockNotifPrefRepo, mockJobProfileRepo := setupTest()
	ctx := context.Background()

	userID := uuid.New()
	expectedUser := &domain.User{
		ID:        userID,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
	}
	expectedNotifPref := &domain.NotificationPreferences{
		UserID:        userID,
		NotifyByEmail: true,
		NotifyBySms:   false,
	}
	expectedJobProfile := &domain.JobProfile{
		UserID:     userID,
		OpenToWork: true,
		CVPath:     sql.NullString{String: "/path/to/cv.pdf", Valid: true},
	}

	mockUserRepo.On("GetByID", ctx, userID).Return(expectedUser, nil)
	mockNotifPrefRepo.On("GetByUserID", ctx, userID).Return(expectedNotifPref, nil)
	mockJobProfileRepo.On("GetByUserID", ctx, userID).Return(expectedJobProfile, nil)

	result, err := service.GetByID(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedUser, result.User)
	assert.Equal(t, expectedNotifPref, result.NotificationPreferences)
	assert.Equal(t, expectedJobProfile, result.JobProfile)
	mockUserRepo.AssertExpectations(t)
	mockNotifPrefRepo.AssertExpectations(t)
	mockJobProfileRepo.AssertExpectations(t)
}

func TestUserService_GetByID_UserNotFound(t *testing.T) {
	service, mockUserRepo, _, _ := setupTest()
	ctx := context.Background()

	userID := uuid.New()

	mockUserRepo.On("GetByID", ctx, userID).Return(nil, domain.ErrUserNotFound)

	result, err := service.GetByID(ctx, userID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrUserNotFound, err)
	mockUserRepo.AssertExpectations(t)
}

func TestUserService_GetByID_NotificationPreferencesError(t *testing.T) {
	service, mockUserRepo, mockNotifPrefRepo, _ := setupTest()
	ctx := context.Background()

	userID := uuid.New()
	expectedUser := &domain.User{
		ID:        userID,
		FirstName: "John",
		LastName:  "Doe",
	}

	mockUserRepo.On("GetByID", ctx, userID).Return(expectedUser, nil)
	mockNotifPrefRepo.On("GetByUserID", ctx, userID).Return(nil, errors.New("database error"))

	result, err := service.GetByID(ctx, userID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, response.ErrInternalServerError, err)
	mockUserRepo.AssertExpectations(t)
	mockNotifPrefRepo.AssertExpectations(t)
}

func TestUserService_GetByID_JobProfileError(t *testing.T) {
	service, mockUserRepo, mockNotifPrefRepo, mockJobProfileRepo := setupTest()
	ctx := context.Background()

	userID := uuid.New()
	expectedUser := &domain.User{
		ID: userID,
	}
	expectedNotifPref := &domain.NotificationPreferences{
		UserID: userID,
	}

	mockUserRepo.On("GetByID", ctx, userID).Return(expectedUser, nil)
	mockNotifPrefRepo.On("GetByUserID", ctx, userID).Return(expectedNotifPref, nil)
	mockJobProfileRepo.On("GetByUserID", ctx, userID).Return(nil, errors.New("database error"))

	result, err := service.GetByID(ctx, userID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, response.ErrInternalServerError, err)
	mockUserRepo.AssertExpectations(t)
	mockNotifPrefRepo.AssertExpectations(t)
	mockJobProfileRepo.AssertExpectations(t)
}

// Test List
func TestUserService_List_Success(t *testing.T) {
	service, mockUserRepo, _, _ := setupTest()
	ctx := context.Background()

	limit := 10
	offset := 0
	filter := &dto.UserListRequest{
		Limit:  &limit,
		Offset: &offset,
	}

	expectedUsers := []*domain.User{
		{ID: uuid.New(), FirstName: "John", LastName: "Doe"},
		{ID: uuid.New(), FirstName: "Jane", LastName: "Smith"},
	}
	expectedCount := 2

	mockUserRepo.On("List", ctx, filter).Return(expectedUsers, nil)
	mockUserRepo.On("Count", ctx, filter).Return(expectedCount, nil)

	result, err := service.List(ctx, filter)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Users, 2)
	assert.Equal(t, expectedCount, result.Count)
	mockUserRepo.AssertExpectations(t)
}

func TestUserService_List_EmptyResults(t *testing.T) {
	service, mockUserRepo, _, _ := setupTest()
	ctx := context.Background()

	filter := &dto.UserListRequest{}
	expectedUsers := []*domain.User{}

	mockUserRepo.On("List", ctx, filter).Return(expectedUsers, domain.ErrUserNotFound)

	result, err := service.List(ctx, filter)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Users, 0)
	assert.Equal(t, 0, result.Count)
	mockUserRepo.AssertExpectations(t)
}

func TestUserService_List_ListError(t *testing.T) {
	service, mockUserRepo, _, _ := setupTest()
	ctx := context.Background()

	filter := &dto.UserListRequest{}

	mockUserRepo.On("List", ctx, filter).Return(nil, errors.New("database error"))

	result, err := service.List(ctx, filter)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, response.ErrInternalServerError, err)
	mockUserRepo.AssertExpectations(t)
}

func TestUserService_List_CountError(t *testing.T) {
	service, mockUserRepo, _, _ := setupTest()
	ctx := context.Background()

	filter := &dto.UserListRequest{}
	expectedUsers := []*domain.User{
		{ID: uuid.New()},
	}

	mockUserRepo.On("List", ctx, filter).Return(expectedUsers, nil)
	mockUserRepo.On("Count", ctx, filter).Return(0, errors.New("database error"))

	result, err := service.List(ctx, filter)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, response.ErrInternalServerError, err)
	mockUserRepo.AssertExpectations(t)
}

// Test UpdateActiveStatus
func TestUserService_UpdateActiveStatus_Success(t *testing.T) {
	service, mockUserRepo, _, _ := setupTest()
	ctx := context.Background()

	userID := uuid.New()
	req := &dto.UserUpdatePropertyRequest{
		ID:    userID,
		Value: true,
	}

	existingUser := &domain.User{
		ID:       userID,
		IsActive: false,
	}

	mockUserRepo.On("GetByID", ctx, userID).Return(existingUser, nil)
	mockUserRepo.On("UpdateProperty", ctx, userID, domain.IsActive, true).Return(nil)

	err := service.UpdateActiveStatus(ctx, req)

	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
}

func TestUserService_UpdateActiveStatus_UserNotFound(t *testing.T) {
	service, mockUserRepo, _, _ := setupTest()
	ctx := context.Background()

	userID := uuid.New()
	req := &dto.UserUpdatePropertyRequest{
		ID:    userID,
		Value: true,
	}

	mockUserRepo.On("GetByID", ctx, userID).Return(nil, domain.ErrUserNotFound)

	err := service.UpdateActiveStatus(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrUserNotFound, err)
	mockUserRepo.AssertExpectations(t)
}

func TestUserService_UpdateActiveStatus_UpdateError(t *testing.T) {
	service, mockUserRepo, _, _ := setupTest()
	ctx := context.Background()

	userID := uuid.New()
	req := &dto.UserUpdatePropertyRequest{
		ID:    userID,
		Value: true,
	}

	existingUser := &domain.User{
		ID: userID,
	}

	mockUserRepo.On("GetByID", ctx, userID).Return(existingUser, nil)
	mockUserRepo.On("UpdateProperty", ctx, userID, domain.IsActive, true).Return(errors.New("database error"))

	err := service.UpdateActiveStatus(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, response.ErrInternalServerError, err)
	mockUserRepo.AssertExpectations(t)
}

// Test VerifyUser
func TestUserService_VerifyEmail_Success(t *testing.T) {
	service, mockUserRepo, _, _ := setupTest()
	ctx := context.Background()

	userID := uuid.New()
	existingUser := &domain.User{
		ID:         userID,
		IsVerified: false,
	}

	mockUserRepo.On("GetByID", ctx, userID).Return(existingUser, nil)
	mockUserRepo.On("UpdateProperty", ctx, userID, domain.IsVerified, true).Return(nil)

	err := service.VerifyEmail(ctx, userID)

	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
}

// Test UpdateCatholicStatus
func TestUserService_UpdateCatholicStatus_Success(t *testing.T) {
	service, mockUserRepo, _, _ := setupTest()
	ctx := context.Background()

	userID := uuid.New()
	req := &dto.UserUpdatePropertyRequest{
		ID:    userID,
		Value: true,
	}

	existingUser := &domain.User{
		ID:         userID,
		IsCatholic: false,
	}

	mockUserRepo.On("GetByID", ctx, userID).Return(existingUser, nil)
	mockUserRepo.On("UpdateProperty", ctx, userID, domain.IsCatholic, true).Return(nil)

	err := service.UpdateCatholicStatus(ctx, req)

	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
}

// Test UpdateEntrepreneurStatus
func TestUserService_UpdateEntrepreneurStatus_Success(t *testing.T) {
	service, mockUserRepo, _, _ := setupTest()
	ctx := context.Background()

	userID := uuid.New()
	req := &dto.UserUpdatePropertyRequest{
		ID:    userID,
		Value: true,
	}

	existingUser := &domain.User{
		ID:             userID,
		IsEntrepreneur: false,
	}

	mockUserRepo.On("GetByID", ctx, userID).Return(existingUser, nil)
	mockUserRepo.On("UpdateProperty", ctx, userID, domain.IsEntrepreneur, true).Return(nil)

	err := service.UpdateEntrepreneurStatus(ctx, req)

	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
}

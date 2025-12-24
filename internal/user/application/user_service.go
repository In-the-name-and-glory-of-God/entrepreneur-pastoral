package application

import (
	"context"
	"database/sql"
	"errors"

	adminDomain "github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/admin/domain"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/user/domain"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/user/infrastructure/dto"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/helper/auth"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/helper/response"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// UserService implements the application logic for managing users.
// It orchestrates calls to different domain repositories.
type UserService struct {
	logger         *zap.SugaredLogger
	userRepo       domain.UserRepository
	notifPrefRepo  domain.NotificationPreferencesRepository
	jobProfileRepo domain.JobProfileRepository
	addressRepo    adminDomain.AddressRepository
}

// NewUserService creates a new UserService with its dependencies.
func NewUserService(
	logger *zap.SugaredLogger,
	userRepo domain.UserRepository,
	notifPrefRepo domain.NotificationPreferencesRepository,
	jobProfileRepo domain.JobProfileRepository,
	addressRepo adminDomain.AddressRepository,
) *UserService {
	return &UserService{
		logger:         logger,
		userRepo:       userRepo,
		notifPrefRepo:  notifPrefRepo,
		jobProfileRepo: jobProfileRepo,
		addressRepo:    addressRepo,
	}
}

func (s *UserService) Create(ctx context.Context, req *dto.UserRegisterRequest) (*domain.User, error) {
	// Check if email already exists
	if _, err := s.userRepo.GetByEmail(ctx, req.Email); err == nil {
		// User found - email already exists
		return nil, domain.ErrEmailAlreadyExists
	} else if !errors.Is(err, domain.ErrUserNotFound) {
		// Unexpected error
		s.logger.Errorw("failed to check existing email", "email", req.Email, "error", err)
		return nil, response.ErrInternalServerError
	}

	// Check if document ID already exists
	if _, err := s.userRepo.GetByDocumentID(ctx, req.DocumentID); err == nil {
		// User found - document ID already exists
		return nil, domain.ErrDocumentIDAlreadyExists
	} else if !errors.Is(err, domain.ErrUserNotFound) {
		// Unexpected error
		s.logger.Errorw("failed to check existing document ID", "documentID", req.DocumentID, "error", err)
		return nil, response.ErrInternalServerError
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		s.logger.Errorw("failed to hash password", "error", err)
		return nil, domain.ErrPasswordHashFailed
	}

	var newUser *domain.User

	err = s.userRepo.UnitOfWork(ctx, func(tx *sqlx.Tx) error {
		// 1. Create the Address
		address := &adminDomain.Address{
			StreetLine1:   req.Address.StreetLine1,
			StreetLine2:   sql.NullString{String: req.Address.StreetLine2, Valid: req.Address.StreetLine2 != ""},
			City:          req.Address.City,
			StateProvince: req.Address.StateProvince,
			PostalCode:    req.Address.PostalCode,
			Country:       req.Address.Country,
		}
		if err := s.addressRepo.Create(tx, address); err != nil {
			s.logger.Errorw("failed to create address for user", "error", err)
			return response.ErrInternalServerError
		}

		// 2. Create the User with the Address ID
		newUser = &domain.User{
			FirstName:        req.FirstName,
			LastName:         req.LastName,
			Email:            req.Email,
			Password:         hashedPassword,
			DocumentID:       req.DocumentID,
			PhoneCountryCode: sql.NullString{String: req.PhoneCountryCode, Valid: req.PhoneCountryCode != ""},
			PhoneNumber:      sql.NullString{String: req.PhoneNumber, Valid: req.PhoneNumber != ""},
			AddressID:        address.ID,
			ChurchID:         req.ChurchID,
		}
		if err := s.userRepo.Create(tx, newUser); err != nil {
			s.logger.Errorw("failed to create user", "email", req.Email, "error", err)
			return response.ErrInternalServerError
		}

		// 3. Create the default NotificationPreferences
		if err := s.notifPrefRepo.Create(tx, &domain.NotificationPreferences{UserID: newUser.ID}); err != nil {
			s.logger.Errorw("failed to create notification preferences", "error", err)
			return response.ErrInternalServerError
		}

		// 4. Create the JobProfile
		newJobProfile := &domain.JobProfile{
			UserID:       newUser.ID,
			OpenToWork:   req.OpenToWork,
			CVPath:       sql.NullString{String: req.CVPath, Valid: req.CVPath != ""},
			FieldsOfWork: req.FieldsOfWork,
		}
		if err := s.jobProfileRepo.Create(tx, newJobProfile); err != nil {
			s.logger.Errorw("failed to create job profile", "error", err)
			return response.ErrInternalServerError
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return newUser, nil
}

// Update modifies an existing user's basic information.
func (s *UserService) Update(ctx context.Context, req *dto.UserUpdateRequest) error {
	user, err := s.userRepo.GetByID(ctx, req.ID)
	if err != nil {
		return domain.ErrUserNotFound
	}

	// 2. Map DTO fields to the entity
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Email = req.Email
	user.DocumentID = req.DocumentID
	user.PhoneCountryCode = sql.NullString{String: req.PhoneCountryCode, Valid: req.PhoneCountryCode != ""}
	user.PhoneNumber = sql.NullString{String: req.PhoneNumber, Valid: req.PhoneNumber != ""}
	user.AddressID = req.AddressID
	user.ChurchID = req.ChurchID

	return s.userRepo.UnitOfWork(ctx, func(tx *sqlx.Tx) error {
		// 3. Update the User
		if err := s.userRepo.Update(tx, user); err != nil {
			s.logger.Errorw("failed to update user", "userID", req.ID, "error", err)
			return response.ErrInternalServerError
		}

		// 4. Update the NotificationPreferences
		updateNotifPref := &domain.NotificationPreferences{
			UserID:        req.ID,
			NotifyByEmail: req.NotifyByEmail,
			NotifyBySms:   req.NotifyBySms,
		}
		if err := s.notifPrefRepo.Update(tx, updateNotifPref); err != nil {
			s.logger.Errorw("failed to update notification preferences", "userID", req.ID, "error", err)
			return response.ErrInternalServerError
		}

		// 5. Update the JobProfile
		updateJobProfile := &domain.JobProfile{
			UserID:       req.ID,
			OpenToWork:   req.OpenToWork,
			CVPath:       sql.NullString{String: req.CVPath, Valid: req.CVPath != ""},
			FieldsOfWork: req.FieldsOfWork,
		}
		if err := s.jobProfileRepo.Update(tx, updateJobProfile); err != nil {
			s.logger.Errorw("failed to update job profile", "userID", req.ID, "error", err)
			return response.ErrInternalServerError
		}

		return nil
	})
}

// GetByID retrieves a single user by their ID along with related entities.
func (s *UserService) GetByID(ctx context.Context, userID uuid.UUID) (*dto.UserGetResponse, error) {
	var resp dto.UserGetResponse

	// 1. Get the User
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrUserNotFound
		}

		s.logger.Errorw("failed to get user by ID", "userID", userID, "error", err)
		return nil, response.ErrInternalServerError
	}
	resp.User = user

	// 2. Get the NotificationPreferences
	notifPrefs, err := s.notifPrefRepo.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.Errorw("failed to get user by ID", "userID", userID, "error", err)
		return nil, response.ErrInternalServerError
	}
	resp.NotificationPreferences = notifPrefs

	// 3. Get the JobProfile
	jobProfile, err := s.jobProfileRepo.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.Errorw("failed to get user by ID", "userID", userID, "error", err)
		return nil, response.ErrInternalServerError
	}
	resp.JobProfile = jobProfile

	return &resp, nil
}

// List retrieves a paginated and filtered list of users.
func (s *UserService) List(ctx context.Context, filter *dto.UserListRequest) (*dto.UserListResponse, error) {
	// 1. Get the list of users
	users, err := s.userRepo.List(ctx, filter)
	if err != nil && err != domain.ErrUserNotFound {
		s.logger.Errorw("failed to list users", "error", err)
		return nil, response.ErrInternalServerError
	}

	// 2. Get the total count for pagination
	count := 0
	if len(users) > 0 {
		count, err = s.userRepo.Count(ctx, filter)
		if err != nil {
			s.logger.Errorw("failed to count users", "error", err)
			return nil, response.ErrInternalServerError
		}
	}

	return &dto.UserListResponse{
		Users: users,
		Count: count,
	}, nil
}

// --- Flag Management Methods ---

// VerifyEmail sets the user's is_verified flag.
func (s *UserService) VerifyEmail(ctx context.Context, userID uuid.UUID) error {
	return s.updateUserProperty(ctx, userID, func(user *domain.User) (domain.UserProperty, any) {
		user.IsVerified = true

		return domain.IsVerified, true
	})
}

// UpdateActiveStatus sets the user's is_active flag.
func (s *UserService) UpdateActiveStatus(ctx context.Context, req *dto.UserUpdatePropertyRequest) error {
	return s.updateUserProperty(ctx, req.ID, func(user *domain.User) (domain.UserProperty, any) {
		user.IsActive = req.Value

		return domain.IsActive, req.Value
	})
}

// UpdateCatholicStatus sets the user's is_catholic flag.
func (s *UserService) UpdateCatholicStatus(ctx context.Context, req *dto.UserUpdatePropertyRequest) error {
	return s.updateUserProperty(ctx, req.ID, func(user *domain.User) (domain.UserProperty, any) {
		user.IsCatholic = req.Value

		return domain.IsCatholic, req.Value
	})
}

// UpdateEntrepreneurStatus sets the user's is_entrepreneur flag.
func (s *UserService) UpdateEntrepreneurStatus(ctx context.Context, req *dto.UserUpdatePropertyRequest) error {
	return s.updateUserProperty(ctx, req.ID, func(user *domain.User) (domain.UserProperty, any) {
		user.IsEntrepreneur = req.Value

		return domain.IsEntrepreneur, req.Value
	})
}

// helper function to update a single user property
func (s *UserService) updateUserProperty(ctx context.Context, userID uuid.UUID, getPropertyAndValueFunc func(user *domain.User) (domain.UserProperty, any)) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return domain.ErrUserNotFound
		}

		s.logger.Errorw("failed to get user by ID", "userID", userID, "error", err)
		return response.ErrInternalServerError
	}

	// Apply the specific update
	property, value := getPropertyAndValueFunc(user)

	if err := s.userRepo.UpdateProperty(ctx, userID, property, value); err != nil {
		s.logger.Errorw("failed to update user property", "userID", userID, "property", property, "value", value, "error", err)
		return response.ErrInternalServerError
	}

	return nil
}

// SetRole sets the user's role.
func (s *UserService) SetRole(ctx context.Context, req *dto.UserSetRoleRequest) error {
	_, err := s.userRepo.GetByID(ctx, req.ID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return domain.ErrUserNotFound
		}

		s.logger.Errorw("failed to get user by ID", "userID", req.ID, "error", err)
		return response.ErrInternalServerError
	}

	if err := s.userRepo.UpdateProperty(ctx, req.ID, domain.RoleID, req.RoleID); err != nil {
		s.logger.Errorw("failed to set user role", "userID", req.ID, "roleID", req.RoleID, "error", err)
		return response.ErrInternalServerError
	}

	return nil
}

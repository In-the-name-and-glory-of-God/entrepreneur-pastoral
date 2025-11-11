package application

import (
	"context"
	"errors"

	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/user/domain"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/user/infrastructure/dto"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/helper/auth"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/helper/response"
	"go.uber.org/zap"
)

type AuthService struct {
	logger       *zap.SugaredLogger
	tokenManager *auth.TokenManager
	userRepo     domain.UserRepository
}

func NewAuthService(logger *zap.SugaredLogger, tokenManager *auth.TokenManager, userRepo domain.UserRepository) *AuthService {
	return &AuthService{
		logger:       logger,
		tokenManager: tokenManager,
		userRepo:     userRepo,
	}
}

func (s *AuthService) Login(ctx context.Context, req *dto.UserLoginRequest) (string, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return "", domain.ErrUserNotFound
		}

		s.logger.Errorw("failed to get user by email", "email", req.Email, "error", err)
		return "", response.ErrInternalServerError
	}

	if !user.IsActive {
		return "", domain.ErrUserInactive
	}

	if !user.IsVerified {
		return "", domain.ErrEmailNotVerified
	}

	if err := auth.VerifyPassword(user.Password, req.Password); err != nil {
		s.logger.Errorw("failed to verify password", "userID", user.ID, "error", err)
		return "", domain.ErrInvalidPassword
	}

	return s.tokenManager.GenerateToken(user.ID.String())
}

func (s *AuthService) UpdatePassword(ctx context.Context, req *dto.UserResetPasswordRequest) error {
	_, err := s.userRepo.GetByID(ctx, req.ID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return err
		}

		s.logger.Errorw("failed to get user by ID", "userID", req.ID, "error", err)
		return response.ErrInternalServerError
	}

	hashedPassword, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		s.logger.Errorw("failed to hash password", "userID", req.ID, "error", err)
		return domain.ErrPasswordHashFailed
	}

	if err := s.userRepo.UpdateProperty(ctx, req.ID, domain.Password, hashedPassword); err != nil {
		s.logger.Errorw("failed to update user password", "userID", req.ID, "error", err)
		return response.ErrInternalServerError
	}

	return nil
}

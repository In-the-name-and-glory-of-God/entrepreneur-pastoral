package application

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/user/domain"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/user/infrastructure/dto"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/config"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/helper/auth"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/helper/constants"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/helper/response"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/i18n"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/storage"
	"go.uber.org/zap"
)

const (
	emailVerificationExpiry = 24 * time.Hour
	passwordResetExpiry     = 1 * time.Hour
)

// NotificationPayload represents the payload sent to the notification queue
type NotificationPayload struct {
	From         string   `json:"from"`
	To           []string `json:"to"`
	Subject      string   `json:"subject"`
	TemplateName string   `json:"template_name"`
	Data         any      `json:"data"`
}

type AuthService struct {
	logger       *zap.SugaredLogger
	config       config.Config
	cache        storage.CacheStorage
	queue        storage.QueueStorage
	tokenManager *auth.TokenManager
	userRepo     domain.UserRepository
}

func NewAuthService(logger *zap.SugaredLogger, cfg config.Config, cache storage.CacheStorage, queue storage.QueueStorage, tokenManager *auth.TokenManager, userRepo domain.UserRepository) *AuthService {
	return &AuthService{
		logger:       logger,
		config:       cfg,
		cache:        cache,
		queue:        queue,
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

// GetUserByEmail retrieves a user by their email address
func (s *AuthService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrUserNotFound
		}

		s.logger.Errorw("failed to get user by email", "email", email, "error", err)
		return nil, response.ErrInternalServerError
	}

	return user, nil
}

// SendVerificationEmail generates a verification token and sends a verification email to the user
func (s *AuthService) SendVerificationEmail(ctx context.Context, user *domain.User) error {
	// Generate a random token for email verification
	token, err := auth.GenerateRandomToken(32)
	if err != nil {
		return err
	}

	// Store token in cache with user ID as value
	cacheKey := s.cache.BuildKey(storage.CACHE_PREFIX_EMAIL_VERIFICATION, token)
	if err := s.cache.SetString(ctx, cacheKey, user.ID.String(), emailVerificationExpiry); err != nil {
		return err
	}

	// Determine user's language preference
	lang := i18n.GetLanguage(ctx)
	if user.Language.Valid && user.Language.String != "" {
		lang = i18n.Language(user.Language.String)
	}

	// Build verification link
	verificationLink := fmt.Sprintf("%s:%d/api/v1/auth/verify-email/%s", s.config.API.Host, s.config.API.Port, token)

	// Create notification payload with translated strings
	payload := NotificationPayload{
		From:         s.config.SMTP.From,
		To:           []string{user.Email},
		Subject:      i18n.Translate(lang, "email.verify_account.subject"),
		TemplateName: constants.EMAIL_TEMPLATE_VERIFY_ACCOUNT,
		Data: map[string]string{
			"Lang":             string(lang),
			"Brand":            i18n.Translate(lang, "email.common.brand"),
			"Title":            i18n.Translate(lang, "email.verify_account.title"),
			"Greeting":         i18n.TranslateWithParams(lang, "email.verify_account.greeting", map[string]string{"name": user.FirstName}),
			"Message":          i18n.Translate(lang, "email.verify_account.message"),
			"Button":           i18n.Translate(lang, "email.verify_account.button"),
			"LinkFallback":     i18n.Translate(lang, "email.verify_account.link_fallback"),
			"Expiry":           i18n.Translate(lang, "email.verify_account.expiry"),
			"Footer":           i18n.Translate(lang, "email.verify_account.footer"),
			"Copyright":        i18n.Translate(lang, "email.common.copyright"),
			"VerificationLink": verificationLink,
		},
	}

	// Publish to notification queue
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return s.queue.Publish(ctx, "", constants.QUEUE_NOTIFICATIONS, payloadBytes)
}

// SendPasswordResetEmail generates a reset token and sends a password reset email to the user
func (s *AuthService) SendPasswordResetEmail(ctx context.Context, user *domain.User) error {
	// Generate a random token for password reset
	token, err := auth.GenerateRandomToken(32)
	if err != nil {
		return err
	}

	// Store token in cache with user ID as value
	cacheKey := s.cache.BuildKey(storage.CACHE_PREFIX_PASSWORD_RESET, token)
	if err := s.cache.SetString(ctx, cacheKey, user.ID.String(), passwordResetExpiry); err != nil {
		return err
	}

	// Determine user's language preference
	lang := i18n.GetLanguage(ctx)
	if user.Language.Valid && user.Language.String != "" {
		lang = i18n.Language(user.Language.String)
	}

	// Build reset link
	resetLink := fmt.Sprintf("%s:%d/api/v1/auth/reset-password/%s/%s", s.config.API.Host, s.config.API.Port, user.ID.String(), token)

	// Create notification payload with translated strings
	payload := NotificationPayload{
		From:         s.config.SMTP.From,
		To:           []string{user.Email},
		Subject:      i18n.Translate(lang, "email.password_reset.subject"),
		TemplateName: constants.EMAIL_TEMPLATE_PASSWORD_RESET,
		Data: map[string]string{
			"Lang":            string(lang),
			"Brand":           i18n.Translate(lang, "email.common.brand"),
			"Title":           i18n.Translate(lang, "email.password_reset.title"),
			"Greeting":        i18n.TranslateWithParams(lang, "email.password_reset.greeting", map[string]string{"name": user.FirstName}),
			"Message":         i18n.Translate(lang, "email.password_reset.message"),
			"Button":          i18n.Translate(lang, "email.password_reset.button"),
			"LinkFallback":    i18n.Translate(lang, "email.password_reset.link_fallback"),
			"Expiry":          i18n.Translate(lang, "email.password_reset.expiry"),
			"SecurityTip":     i18n.Translate(lang, "email.password_reset.security_tip"),
			"SecurityMessage": i18n.Translate(lang, "email.password_reset.security_message"),
			"Footer":          i18n.Translate(lang, "email.password_reset.footer"),
			"Copyright":       i18n.Translate(lang, "email.common.copyright"),
			"ResetLink":       resetLink,
		},
	}

	// Publish to notification queue
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return s.queue.Publish(ctx, "", constants.QUEUE_NOTIFICATIONS, payloadBytes)
}

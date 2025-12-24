package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/user/application"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/user/domain"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/user/infrastructure/dto"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/helper/auth"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/helper/response"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/storage"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AuthHandler struct {
	logger      *zap.SugaredLogger
	cache       storage.CacheStorage
	authService *application.AuthService
	userService *application.UserService
}

func NewAuthHandler(logger *zap.SugaredLogger, cache storage.CacheStorage, authService *application.AuthService, userService *application.UserService) *AuthHandler {
	return &AuthHandler{
		logger:      logger,
		cache:       cache,
		authService: authService,
		userService: userService,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req dto.UserRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequestT(ctx, w, "error.invalid_request_body", nil)
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	if err := auth.IsValidEmail(req.Email); err != nil {
		response.BadRequestT(ctx, w, "error.valid_email_required", nil)
		return
	}

	if err := auth.IsStrongPassword(req.Password); err != nil {
		response.BadRequestT(ctx, w, "error.password_strength", nil)
		return
	}

	user, err := h.userService.Create(ctx, &req)
	if err != nil {
		if errors.Is(err, domain.ErrEmailAlreadyExists) {
			response.ConflictT(ctx, w, "error.email_already_exists", nil)
			return
		}
		if errors.Is(err, domain.ErrDocumentIDAlreadyExists) {
			response.ConflictT(ctx, w, "error.document_id_already_exists", nil)
			return
		}

		h.logger.Errorf("Failed to register user", "error", err)
		response.InternalServerErrorT(ctx, w, "error.failed_register_user")
		return
	}

	// Send verification email
	if err := h.authService.SendVerificationEmail(ctx, user); err != nil {
		h.logger.Errorw("Failed to send verification email", "userID", user.ID, "error", err)
		// Don't fail registration if email sending fails, just log the error
	}

	response.CreatedT(ctx, w, "success.user_registered", nil)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req dto.UserLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequestT(ctx, w, "error.invalid_request_body", nil)
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	if req.Email == "" {
		response.BadRequestT(ctx, w, "error.email_required", nil)
		return
	}

	if req.Password == "" {
		response.BadRequestT(ctx, w, "error.password_required", nil)
		return
	}

	token, err := h.authService.Login(ctx, &req)
	if err != nil {
		if errors.Is(err, domain.ErrUserInactive) {
			response.ForbiddenT(ctx, w, "error.user_inactive")
			return
		} else if errors.Is(err, domain.ErrEmailNotVerified) {
			response.ForbiddenT(ctx, w, "error.email_not_verified")
			return
		} else if errors.Is(err, domain.ErrInvalidPassword) || errors.Is(err, domain.ErrUserNotFound) {
			response.UnauthorizedT(ctx, w, "error.invalid_credentials")
			return
		}

		h.logger.Errorf("Failed to login user", "error", err)
		response.InternalServerErrorT(ctx, w, "error.failed_login_user")
		return
	}

	// TODO: Set refresh token cookie logic here
	// auth.SetRefreshTokenCookie(w, user.ID.String())

	response.OKT(ctx, w, "success.login", dto.UserLoginResponse{
		Token: token,
	})
}

// RequestPasswordReset handles the initial password reset request (sends email with reset link)
func (h *AuthHandler) RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req dto.UserResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequestT(ctx, w, "error.invalid_request_body", nil)
		return
	}

	if req.Email == "" {
		response.BadRequestT(ctx, w, "error.email_required", nil)
		return
	}

	if err := auth.IsValidEmail(req.Email); err != nil {
		response.BadRequestT(ctx, w, "error.valid_email_required", nil)
		return
	}

	// Get user by email - don't reveal if email exists or not for security
	user, err := h.authService.GetUserByEmail(ctx, req.Email)
	if err != nil {
		// Return success even if user not found to prevent email enumeration
		response.OKT(ctx, w, "success.password_reset_sent", nil)
		return
	}

	// Check if user is active
	if !user.IsActive {
		response.OKT(ctx, w, "success.password_reset_sent", nil)
		return
	}

	// Send password reset email
	if err := h.authService.SendPasswordResetEmail(ctx, user); err != nil {
		h.logger.Errorw("Failed to send password reset email", "email", req.Email, "error", err)
	}

	response.OKT(ctx, w, "success.password_reset_sent", nil)
}

// ConfirmPasswordReset handles the actual password reset (with token validation)
func (h *AuthHandler) ConfirmPasswordReset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req dto.UserResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequestT(ctx, w, "error.invalid_request_body", nil)
		return
	}

	id := chi.URLParam(r, "id")
	token := chi.URLParam(r, "token")
	if token == "" {
		response.BadRequestT(ctx, w, "error.missing_token", nil)
		return
	}

	userID, err := uuid.Parse(id)
	if err != nil {
		response.BadRequestT(ctx, w, "error.invalid_user_id", nil)
		return
	}
	req.ID = userID

	// Validate token from cache
	cacheKey := h.cache.BuildKey(storage.CACHE_PREFIX_PASSWORD_RESET, token)
	cachedUserID, err := h.cache.GetStringAndDel(ctx, cacheKey)
	if err != nil {
		if errors.Is(err, storage.ErrCacheMiss) {
			response.BadRequestT(ctx, w, "error.invalid_token", nil)
			return
		}
		h.logger.Errorf("Failed to retrieve password reset token", "error", err)
		response.InternalServerErrorT(ctx, w, "error.failed_reset_password")
		return
	}

	// Verify the token belongs to the correct user
	if cachedUserID != userID.String() {
		response.BadRequestT(ctx, w, "error.invalid_token_for_user", nil)
		return
	}

	// Validate user status before allowing password reset
	userResp, err := h.userService.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			response.NotFoundT(ctx, w, "error.user_not_found")
			return
		}
		h.logger.Errorf("Failed to get user for password reset", "error", err)
		response.InternalServerErrorT(ctx, w, "error.failed_reset_password")
		return
	}

	if !userResp.User.IsActive {
		response.ForbiddenT(ctx, w, "error.user_inactive")
		return
	}

	if err := auth.IsStrongPassword(req.NewPassword); err != nil {
		response.BadRequestT(ctx, w, "error.new_password_strength", nil)
		return
	}

	if err := h.authService.UpdatePassword(ctx, &req); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			response.NotFoundT(ctx, w, "error.user_not_found")
			return
		} else if errors.Is(err, domain.ErrInvalidPassword) {
			response.UnauthorizedT(ctx, w, "error.old_password_incorrect")
			return
		}

		h.logger.Errorf("Failed to update password", "error", err)
		response.InternalServerErrorT(ctx, w, "error.failed_update_password")
		return
	}

	response.OKT(ctx, w, "success.password_updated", nil)
}

func (h *AuthHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	token := chi.URLParam(r, "token")
	if token == "" {
		response.BadRequestT(ctx, w, "error.missing_token", nil)
		return
	}

	id, err := h.cache.GetStringAndDel(ctx, h.cache.BuildKey(storage.CACHE_PREFIX_EMAIL_VERIFICATION, token))
	if err != nil {
		if errors.Is(err, storage.ErrCacheMiss) {
			response.BadRequestT(ctx, w, "error.invalid_token", nil)
			return
		}

		h.logger.Errorf("Failed to retrieve email verification token", "error", err)
		response.InternalServerErrorT(ctx, w, "error.failed_verify_email")
		return
	}

	userID, err := uuid.Parse(id)
	if err != nil {
		response.BadRequestT(ctx, w, "error.invalid_user_id_in_token", nil)
		return
	}

	if err := h.userService.VerifyEmail(ctx, userID); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			response.NotFoundT(ctx, w, "error.user_not_found")
			return
		}

		h.logger.Errorf("Failed to verify email", "error", err)
		response.InternalServerErrorT(ctx, w, "error.failed_verify_email")
		return
	}

	response.OKT(ctx, w, "success.email_verified", nil)
}

// TODO: Set proper refresh token expiration and security flags
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	response.NotImplementedT(r.Context(), w, "error.not_implemented")
}

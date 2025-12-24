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
	var req dto.UserRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body", nil)
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	if err := auth.IsValidEmail(req.Email); err != nil {
		response.BadRequest(w, "Valid email is required: "+err.Error(), nil)
		return
	}

	if err := auth.IsStrongPassword(req.Password); err != nil {
		response.BadRequest(w, "Password does not meet strength requirements: "+err.Error(), nil)
		return
	}

	user, err := h.userService.Create(r.Context(), &req)
	if err != nil {
		if errors.Is(err, domain.ErrEmailAlreadyExists) {
			response.Conflict(w, "Email already exists", nil)
			return
		}
		if errors.Is(err, domain.ErrDocumentIDAlreadyExists) {
			response.Conflict(w, "Document ID already exists", nil)
			return
		}

		h.logger.Errorf("Failed to register user", "error", err)
		response.InternalServerError(w, "Failed to register user")
		return
	}

	// Send verification email
	if err := h.authService.SendVerificationEmail(r.Context(), user); err != nil {
		h.logger.Errorw("Failed to send verification email", "userID", user.ID, "error", err)
		// Don't fail registration if email sending fails, just log the error
	}

	response.Created(w, "User registered successfully. Please check your email to verify your account.", nil)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.UserLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body", nil)
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	if req.Email == "" {
		response.BadRequest(w, "Email is required", nil)
		return
	}

	if req.Password == "" {
		response.BadRequest(w, "Password is required", nil)
		return
	}

	token, err := h.authService.Login(r.Context(), &req)
	if err != nil {
		if errors.Is(err, domain.ErrUserInactive) {
			response.Forbidden(w, "User account is inactive")
			return
		} else if errors.Is(err, domain.ErrEmailNotVerified) {
			response.Forbidden(w, "Email address is not verified")
			return
		} else if errors.Is(err, domain.ErrInvalidPassword) || errors.Is(err, domain.ErrUserNotFound) {
			response.Unauthorized(w, "Invalid credentials")
			return
		}

		h.logger.Errorf("Failed to login user", "error", err)
		response.InternalServerError(w, "Failed to login user")
		return
	}

	// TODO: Set refresh token cookie logic here
	// auth.SetRefreshTokenCookie(w, user.ID.String())

	response.OK(w, "Login successful", dto.UserLoginResponse{
		Token: token,
	})
}

// RequestPasswordReset handles the initial password reset request (sends email with reset link)
func (h *AuthHandler) RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
	var req dto.UserResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body", nil)
		return
	}

	if req.Email == "" {
		response.BadRequest(w, "Email is required", nil)
		return
	}

	if err := auth.IsValidEmail(req.Email); err != nil {
		response.BadRequest(w, "Valid email is required: "+err.Error(), nil)
		return
	}

	// Get user by email - don't reveal if email exists or not for security
	user, err := h.authService.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		// Return success even if user not found to prevent email enumeration
		response.OK(w, "If the email exists, a password reset link has been sent", nil)
		return
	}

	// Check if user is active
	if !user.IsActive {
		response.OK(w, "If the email exists, a password reset link has been sent", nil)
		return
	}

	// Send password reset email
	if err := h.authService.SendPasswordResetEmail(r.Context(), user); err != nil {
		h.logger.Errorw("Failed to send password reset email", "email", req.Email, "error", err)
	}

	response.OK(w, "If the email exists, a password reset link has been sent", nil)
}

// ConfirmPasswordReset handles the actual password reset (with token validation)
func (h *AuthHandler) ConfirmPasswordReset(w http.ResponseWriter, r *http.Request) {
	var req dto.UserResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body", nil)
		return
	}

	id := chi.URLParam(r, "id")
	token := chi.URLParam(r, "token")
	if token == "" {
		response.BadRequest(w, "Missing token", nil)
		return
	}

	userID, err := uuid.Parse(id)
	if err != nil {
		response.BadRequest(w, "Invalid user ID", nil)
		return
	}
	req.ID = userID

	// Validate token from cache
	cacheKey := h.cache.BuildKey(storage.CACHE_PREFIX_PASSWORD_RESET, token)
	cachedUserID, err := h.cache.GetStringAndDel(r.Context(), cacheKey)
	if err != nil {
		if errors.Is(err, storage.ErrCacheMiss) {
			response.BadRequest(w, "Invalid or expired token", nil)
			return
		}
		h.logger.Errorf("Failed to retrieve password reset token", "error", err)
		response.InternalServerError(w, "Failed to reset password")
		return
	}

	// Verify the token belongs to the correct user
	if cachedUserID != userID.String() {
		response.BadRequest(w, "Invalid token for this user", nil)
		return
	}

	// Validate user status before allowing password reset
	userResp, err := h.userService.GetByID(r.Context(), userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			response.NotFound(w, "User not found")
			return
		}
		h.logger.Errorf("Failed to get user for password reset", "error", err)
		response.InternalServerError(w, "Failed to reset password")
		return
	}

	if !userResp.User.IsActive {
		response.Forbidden(w, "User account is inactive")
		return
	}

	if err := auth.IsStrongPassword(req.NewPassword); err != nil {
		response.BadRequest(w, "New password does not meet strength requirements: "+err.Error(), nil)
		return
	}

	if err := h.authService.UpdatePassword(r.Context(), &req); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			response.NotFound(w, "User not found")
			return
		} else if errors.Is(err, domain.ErrInvalidPassword) {
			response.Unauthorized(w, "Old password is incorrect")
			return
		}

		h.logger.Errorf("Failed to update password", "error", err)
		response.InternalServerError(w, "Failed to update password")
		return
	}

	response.OK(w, "Password updated successfully", nil)
}

func (h *AuthHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	if token == "" {
		response.BadRequest(w, "Missing token", nil)
		return
	}

	id, err := h.cache.GetStringAndDel(r.Context(), h.cache.BuildKey(storage.CACHE_PREFIX_EMAIL_VERIFICATION, token))
	if err != nil {
		if errors.Is(err, storage.ErrCacheMiss) {
			response.BadRequest(w, "Invalid or expired token", nil)
			return
		}

		h.logger.Errorf("Failed to retrieve email verification token", "error", err)
		response.InternalServerError(w, "Failed to verify email")
		return
	}

	userID, err := uuid.Parse(id)
	if err != nil {
		response.BadRequest(w, "Invalid user ID in token", nil)
		return
	}

	if err := h.userService.VerifyEmail(r.Context(), userID); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			response.NotFound(w, "User not found")
			return
		}

		h.logger.Errorf("Failed to verify email", "error", err)
		response.InternalServerError(w, "Failed to verify email")
		return
	}
}

// TODO: Set proper refresh token expiration and security flags
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	response.NotImplemented(w, "Refresh token functionality not implemented yet")
}

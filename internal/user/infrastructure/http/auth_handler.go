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

func (h *AuthHandler) RegisterRoutes(r chi.Router) {
	r.Put("/register", h.Register)
	r.Post("/login", h.Login)
	r.Patch("/update-password", h.UpdatePassword)
	r.Patch("/verify-email", h.VerifyEmail)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.UserRegisterRequest
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

	if err := auth.IsStrongPassword(req.Password); err != nil {
		response.BadRequest(w, "Password does not meet strength requirements: "+err.Error(), nil)
		return
	}

	if err := h.userService.Create(r.Context(), &req); err != nil {
		if errors.Is(err, domain.ErrEmailAlreadyExists) {
			response.Conflict(w, "Email already exists", nil)
			return
		}
		if errors.Is(err, domain.ErrDocumentIDAlreadyExists) {
			response.Conflict(w, "Document ID already exists", nil)
			return
		}

		h.logger.Errorf("Failed to register user: %v", err)
		response.InternalServerError(w, "Failed to register user")
		return
	}

	response.Created(w, "User registered successfully", nil)
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
		} else if errors.Is(err, domain.ErrInvalidPassword) || errors.Is(err, domain.ErrUserNotFound) {
			response.Unauthorized(w, "Invalid credentials")
			return
		}

		h.logger.Errorf("Failed to login user: %v", err)
		response.InternalServerError(w, "Failed to login user")
		return
	}

	response.OK(w, "Login successful", dto.UserLoginResponse{
		Token: token,
	})
}

func (h *AuthHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	// TODO: Update password logic
}

func (h *AuthHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	// TODO: Email verification logic
}

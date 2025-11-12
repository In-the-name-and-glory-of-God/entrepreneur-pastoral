package http

import (
	"encoding/json"
	"net/http"

	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/user/application"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/user/domain"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/user/infrastructure/dto"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/helper/response"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserHandler struct {
	logger      *zap.SugaredLogger
	userService *application.UserService
}

func NewUserHandler(logger *zap.SugaredLogger, userService *application.UserService) *UserHandler {
	return &UserHandler{
		logger:      logger,
		userService: userService,
	}
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	uuid, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "Invalid user ID", nil)
		return
	}

	user, err := h.userService.GetByID(r.Context(), uuid)
	if err != nil {
		if err == domain.ErrUserNotFound {
			response.NotFound(w, "User not found")
			return
		}

		h.logger.Errorw("failed to get user by ID", "userID", uuid, "error", err)
		response.InternalServerError(w, "Failed to get user")
		return
	}

	response.OK(w, "User retrieved successfully", user)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	uuid, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "Invalid user ID", nil)
		return
	}

	var req dto.UserUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body", nil)
		return
	}
	req.ID = uuid

	if err := h.userService.Update(r.Context(), &req); err != nil {
		if err == domain.ErrUserNotFound {
			response.NotFound(w, "User not found")
			return
		}

		h.logger.Errorw("failed to update user", "userID", uuid, "error", err)
		response.InternalServerError(w, "Failed to update user")
		return
	}

	response.OK(w, "User updated successfully", nil)
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	var req dto.UserListRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body", nil)
		return
	}

	list, err := h.userService.List(r.Context(), &req)
	if err != nil {
		if err == domain.ErrUserNotFound {
			response.NotFound(w, "Users not found")
			return
		}

		h.logger.Errorw("failed to list users", "error", err)
		response.InternalServerError(w, "Failed to list users")
		return
	}

	response.OK(w, "Users listed successfully", list)
}

func (h *UserHandler) SetIsActive(w http.ResponseWriter, r *http.Request) {
	uuid, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "Invalid user ID", nil)
		return
	}

	var req dto.UserUpdatePropertyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body", nil)
		return
	}
	req.ID = uuid

	if err := h.userService.UpdateActiveStatus(r.Context(), &req); err != nil {
		if err == domain.ErrUserNotFound {
			response.NotFound(w, "User not found")
			return
		}

		h.logger.Errorw("failed to set is_active flag", "userID", uuid, "error", err)
		response.InternalServerError(w, "Failed to update is_active flag")
		return
	}

	response.OK(w, "User patched successfully", nil)
}

func (h *UserHandler) SetIsCatholic(w http.ResponseWriter, r *http.Request) {
	uuid, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "Invalid user ID", nil)
		return
	}

	var req dto.UserUpdatePropertyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body", nil)
		return
	}
	req.ID = uuid

	if err := h.userService.UpdateCatholicStatus(r.Context(), &req); err != nil {
		if err == domain.ErrUserNotFound {
			response.NotFound(w, "User not found")
			return
		}

		h.logger.Errorw("failed to set is_catholic flag", "userID", uuid, "error", err)
		response.InternalServerError(w, "Failed to update is_catholic flag")
		return
	}

	response.OK(w, "User patched successfully", nil)
}

func (h *UserHandler) SetIsEntrepreneur(w http.ResponseWriter, r *http.Request) {
	uuid, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "Invalid user ID", nil)
		return
	}

	var req dto.UserUpdatePropertyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body", nil)
		return
	}
	req.ID = uuid

	if err := h.userService.UpdateEntrepreneurStatus(r.Context(), &req); err != nil {
		if err == domain.ErrUserNotFound {
			response.NotFound(w, "User not found")
			return
		}

		h.logger.Errorw("failed to set is_entrepreneur flag", "userID", uuid, "error", err)
		response.InternalServerError(w, "Failed to update is_entrepreneur flag")
		return
	}

	response.OK(w, "User patched successfully", nil)
}

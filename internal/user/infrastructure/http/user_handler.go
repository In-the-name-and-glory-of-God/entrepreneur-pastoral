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
	ctx := r.Context()
	uuid, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequestT(ctx, w, "error.invalid_user_id", nil)
		return
	}

	user, err := h.userService.GetByID(ctx, uuid)
	if err != nil {
		if err == domain.ErrUserNotFound {
			response.NotFoundT(ctx, w, "error.user_not_found")
			return
		}

		h.logger.Errorw("failed to get user by ID", "userID", uuid, "error", err)
		response.InternalServerErrorT(ctx, w, "error.failed_get_user")
		return
	}

	response.OKT(ctx, w, "success.user_retrieved", user)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	uuid, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequestT(ctx, w, "error.invalid_user_id", nil)
		return
	}

	var req dto.UserUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequestT(ctx, w, "error.invalid_request_body", nil)
		return
	}
	req.ID = uuid

	if err := h.userService.Update(ctx, &req); err != nil {
		if err == domain.ErrUserNotFound {
			response.NotFoundT(ctx, w, "error.user_not_found")
			return
		}

		h.logger.Errorw("failed to update user", "userID", uuid, "error", err)
		response.InternalServerErrorT(ctx, w, "error.failed_update_user")
		return
	}

	response.OKT(ctx, w, "success.user_updated", nil)
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req dto.UserListRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequestT(ctx, w, "error.invalid_request_body", nil)
		return
	}

	list, err := h.userService.List(ctx, &req)
	if err != nil {
		if err == domain.ErrUserNotFound {
			response.NotFoundT(ctx, w, "error.users_not_found")
			return
		}

		h.logger.Errorw("failed to list users", "error", err)
		response.InternalServerErrorT(ctx, w, "error.failed_list_users")
		return
	}

	response.OKT(ctx, w, "success.users_listed", list)
}

func (h *UserHandler) SetIsActive(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	uuid, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequestT(ctx, w, "error.invalid_user_id", nil)
		return
	}

	var req dto.UserUpdatePropertyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequestT(ctx, w, "error.invalid_request_body", nil)
		return
	}
	req.ID = uuid

	if err := h.userService.UpdateActiveStatus(ctx, &req); err != nil {
		if err == domain.ErrUserNotFound {
			response.NotFoundT(ctx, w, "error.user_not_found")
			return
		}

		h.logger.Errorw("failed to set is_active flag", "userID", uuid, "error", err)
		response.InternalServerErrorT(ctx, w, "error.failed_update_active_flag")
		return
	}

	response.OKT(ctx, w, "success.user_active_updated", nil)
}

func (h *UserHandler) SetIsCatholic(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	uuid, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequestT(ctx, w, "error.invalid_user_id", nil)
		return
	}

	var req dto.UserUpdatePropertyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequestT(ctx, w, "error.invalid_request_body", nil)
		return
	}
	req.ID = uuid

	if err := h.userService.UpdateCatholicStatus(ctx, &req); err != nil {
		if err == domain.ErrUserNotFound {
			response.NotFoundT(ctx, w, "error.user_not_found")
			return
		}

		h.logger.Errorw("failed to set is_catholic flag", "userID", uuid, "error", err)
		response.InternalServerErrorT(ctx, w, "error.failed_update_catholic_flag")
		return
	}

	response.OKT(ctx, w, "success.user_catholic_updated", nil)
}

func (h *UserHandler) SetIsEntrepreneur(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	uuid, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequestT(ctx, w, "error.invalid_user_id", nil)
		return
	}

	var req dto.UserUpdatePropertyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequestT(ctx, w, "error.invalid_request_body", nil)
		return
	}
	req.ID = uuid

	if err := h.userService.UpdateEntrepreneurStatus(ctx, &req); err != nil {
		if err == domain.ErrUserNotFound {
			response.NotFoundT(ctx, w, "error.user_not_found")
			return
		}

		h.logger.Errorw("failed to set is_entrepreneur flag", "userID", uuid, "error", err)
		response.InternalServerErrorT(ctx, w, "error.failed_update_entrepreneur_flag")
		return
	}

	response.OKT(ctx, w, "success.user_entrepreneur_updated", nil)
}

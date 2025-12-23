package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/admin/application"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/admin/domain"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/admin/infrastructure/dto"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/helper/response"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type FieldOfWorkHandler struct {
	logger             *zap.SugaredLogger
	fieldOfWorkService *application.FieldOfWorkService
}

func NewFieldOfWorkHandler(logger *zap.SugaredLogger, fieldOfWorkService *application.FieldOfWorkService) *FieldOfWorkHandler {
	return &FieldOfWorkHandler{
		logger:             logger,
		fieldOfWorkService: fieldOfWorkService,
	}
}

func (h *FieldOfWorkHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.FieldOfWorkCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body", nil)
		return
	}

	fieldOfWork, err := h.fieldOfWorkService.Create(r.Context(), &req)
	if err != nil {
		if err == domain.ErrFieldOfWorkAlreadyExists {
			response.Conflict(w, "Field of work with this name already exists", nil)
			return
		}

		h.logger.Errorw("failed to create field of work", "error", err)
		response.InternalServerError(w, "Failed to create field of work")
		return
	}

	response.Created(w, "Field of work created successfully", fieldOfWork)
}

func (h *FieldOfWorkHandler) Update(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 16)
	if err != nil {
		response.BadRequest(w, "Invalid field of work ID", nil)
		return
	}

	var req dto.FieldOfWorkUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body", nil)
		return
	}
	req.ID = int16(id)

	if err := h.fieldOfWorkService.Update(r.Context(), &req); err != nil {
		if err == domain.ErrFieldOfWorkNotFound {
			response.NotFound(w, "Field of work not found")
			return
		}
		if err == domain.ErrFieldOfWorkAlreadyExists {
			response.Conflict(w, "Field of work with this name already exists", nil)
			return
		}

		h.logger.Errorw("failed to update field of work", "id", id, "error", err)
		response.InternalServerError(w, "Failed to update field of work")
		return
	}

	response.OK(w, "Field of work updated successfully", nil)
}

func (h *FieldOfWorkHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 16)
	if err != nil {
		response.BadRequest(w, "Invalid field of work ID", nil)
		return
	}

	if err := h.fieldOfWorkService.Delete(r.Context(), int16(id)); err != nil {
		if err == domain.ErrFieldOfWorkNotFound {
			response.NotFound(w, "Field of work not found")
			return
		}

		h.logger.Errorw("failed to delete field of work", "id", id, "error", err)
		response.InternalServerError(w, "Failed to delete field of work")
		return
	}

	response.OK(w, "Field of work deleted successfully", nil)
}

func (h *FieldOfWorkHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 16)
	if err != nil {
		response.BadRequest(w, "Invalid field of work ID", nil)
		return
	}

	fieldOfWork, err := h.fieldOfWorkService.GetByID(r.Context(), int16(id))
	if err != nil {
		if err == domain.ErrFieldOfWorkNotFound {
			response.NotFound(w, "Field of work not found")
			return
		}

		h.logger.Errorw("failed to get field of work by ID", "id", id, "error", err)
		response.InternalServerError(w, "Failed to get field of work")
		return
	}

	response.OK(w, "Field of work retrieved successfully", fieldOfWork)
}

func (h *FieldOfWorkHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	list, err := h.fieldOfWorkService.GetAll(r.Context())
	if err != nil {
		h.logger.Errorw("failed to get all fields of work", "error", err)
		response.InternalServerError(w, "Failed to get fields of work")
		return
	}

	response.OK(w, "Fields of work retrieved successfully", list)
}

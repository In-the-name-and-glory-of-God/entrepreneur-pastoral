package http

import (
	"encoding/json"
	"net/http"

	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/admin/application"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/admin/domain"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/admin/infrastructure/dto"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/helper/response"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ChurchHandler struct {
	logger        *zap.SugaredLogger
	churchService *application.ChurchService
}

func NewChurchHandler(logger *zap.SugaredLogger, churchService *application.ChurchService) *ChurchHandler {
	return &ChurchHandler{
		logger:        logger,
		churchService: churchService,
	}
}

func (h *ChurchHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.ChurchCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body", nil)
		return
	}

	church, err := h.churchService.Create(r.Context(), &req)
	if err != nil {
		if err == domain.ErrChurchAlreadyExists {
			response.Conflict(w, "Church with this name already exists", nil)
			return
		}

		h.logger.Errorw("failed to create church", "error", err)
		response.InternalServerError(w, "Failed to create church")
		return
	}

	response.Created(w, "Church created successfully", church)
}

func (h *ChurchHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "Invalid church ID", nil)
		return
	}

	var req dto.ChurchUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body", nil)
		return
	}
	req.ID = id

	if err := h.churchService.Update(r.Context(), &req); err != nil {
		if err == domain.ErrChurchNotFound {
			response.NotFound(w, "Church not found")
			return
		}
		if err == domain.ErrChurchAlreadyExists {
			response.Conflict(w, "Church with this name already exists", nil)
			return
		}

		h.logger.Errorw("failed to update church", "id", id, "error", err)
		response.InternalServerError(w, "Failed to update church")
		return
	}

	response.OK(w, "Church updated successfully", nil)
}

func (h *ChurchHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "Invalid church ID", nil)
		return
	}

	if err := h.churchService.Delete(r.Context(), id); err != nil {
		if err == domain.ErrChurchNotFound {
			response.NotFound(w, "Church not found")
			return
		}

		h.logger.Errorw("failed to delete church", "id", id, "error", err)
		response.InternalServerError(w, "Failed to delete church")
		return
	}

	response.OK(w, "Church deleted successfully", nil)
}

func (h *ChurchHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "Invalid church ID", nil)
		return
	}

	church, err := h.churchService.GetByID(r.Context(), id)
	if err != nil {
		if err == domain.ErrChurchNotFound {
			response.NotFound(w, "Church not found")
			return
		}

		h.logger.Errorw("failed to get church by ID", "id", id, "error", err)
		response.InternalServerError(w, "Failed to get church")
		return
	}

	response.OK(w, "Church retrieved successfully", church)
}

func (h *ChurchHandler) List(w http.ResponseWriter, r *http.Request) {
	var req dto.ChurchListRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body", nil)
		return
	}

	list, err := h.churchService.List(r.Context(), &req)
	if err != nil {
		h.logger.Errorw("failed to list churches", "error", err)
		response.InternalServerError(w, "Failed to list churches")
		return
	}

	response.OK(w, "Churches listed successfully", list)
}

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

type IndustryHandler struct {
	logger          *zap.SugaredLogger
	industryService *application.IndustryService
}

func NewIndustryHandler(logger *zap.SugaredLogger, industryService *application.IndustryService) *IndustryHandler {
	return &IndustryHandler{
		logger:          logger,
		industryService: industryService,
	}
}

func (h *IndustryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.IndustryCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body", nil)
		return
	}

	industry, err := h.industryService.Create(r.Context(), &req)
	if err != nil {
		if err == domain.ErrIndustryAlreadyExists {
			response.Conflict(w, "Industry with this name already exists", nil)
			return
		}

		h.logger.Errorw("failed to create industry", "error", err)
		response.InternalServerError(w, "Failed to create industry")
		return
	}

	response.Created(w, "Industry created successfully", industry)
}

func (h *IndustryHandler) Update(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 16)
	if err != nil {
		response.BadRequest(w, "Invalid industry ID", nil)
		return
	}

	var req dto.IndustryUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body", nil)
		return
	}
	req.ID = int16(id)

	if err := h.industryService.Update(r.Context(), &req); err != nil {
		if err == domain.ErrIndustryNotFound {
			response.NotFound(w, "Industry not found")
			return
		}
		if err == domain.ErrIndustryAlreadyExists {
			response.Conflict(w, "Industry with this name already exists", nil)
			return
		}

		h.logger.Errorw("failed to update industry", "id", id, "error", err)
		response.InternalServerError(w, "Failed to update industry")
		return
	}

	response.OK(w, "Industry updated successfully", nil)
}

func (h *IndustryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 16)
	if err != nil {
		response.BadRequest(w, "Invalid industry ID", nil)
		return
	}

	if err := h.industryService.Delete(r.Context(), int16(id)); err != nil {
		if err == domain.ErrIndustryNotFound {
			response.NotFound(w, "Industry not found")
			return
		}

		h.logger.Errorw("failed to delete industry", "id", id, "error", err)
		response.InternalServerError(w, "Failed to delete industry")
		return
	}

	response.OK(w, "Industry deleted successfully", nil)
}

func (h *IndustryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 16)
	if err != nil {
		response.BadRequest(w, "Invalid industry ID", nil)
		return
	}

	industry, err := h.industryService.GetByID(r.Context(), int16(id))
	if err != nil {
		if err == domain.ErrIndustryNotFound {
			response.NotFound(w, "Industry not found")
			return
		}

		h.logger.Errorw("failed to get industry by ID", "id", id, "error", err)
		response.InternalServerError(w, "Failed to get industry")
		return
	}

	response.OK(w, "Industry retrieved successfully", industry)
}

func (h *IndustryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	list, err := h.industryService.GetAll(r.Context())
	if err != nil {
		h.logger.Errorw("failed to get all industries", "error", err)
		response.InternalServerError(w, "Failed to get industries")
		return
	}

	response.OK(w, "Industries retrieved successfully", list)
}

package http

import (
	"encoding/json"
	"net/http"

	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/entrepreneur/application"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/entrepreneur/domain"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/entrepreneur/infrastructure/dto"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/helper/response"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type BusinessHandler struct {
	logger          *zap.SugaredLogger
	businessService *application.BusinessService
}

func NewBusinessHandler(logger *zap.SugaredLogger, businessService *application.BusinessService) *BusinessHandler {
	return &BusinessHandler{
		logger:          logger,
		businessService: businessService,
	}
}

func (h *BusinessHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.BusinessCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body", nil)
		return
	}

	business, err := h.businessService.Create(r.Context(), &req)
	if err != nil {
		h.logger.Errorw("failed to create business", "error", err)
		response.InternalServerError(w, "Failed to create business")
		return
	}

	response.Created(w, "Business created successfully", business)
}

func (h *BusinessHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "Invalid business ID", nil)
		return
	}

	var req dto.BusinessUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body", nil)
		return
	}
	req.ID = id

	if err := h.businessService.Update(r.Context(), &req); err != nil {
		if err == domain.ErrBusinessNotFound {
			response.NotFound(w, "Business not found")
			return
		}
		if err == domain.ErrUnauthorized {
			response.Unauthorized(w, "Unauthorized to update business")
			return
		}
		h.logger.Errorw("failed to update business", "id", id, "error", err)
		response.InternalServerError(w, "Failed to update business")
		return
	}

	response.OK(w, "Business updated successfully", nil)
}

func (h *BusinessHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "Invalid business ID", nil)
		return
	}

	if err := h.businessService.Delete(r.Context(), id); err != nil {
		if err == domain.ErrUnauthorized {
			response.Unauthorized(w, "Unauthorized to delete business")
			return
		}
		h.logger.Errorw("failed to delete business", "id", id, "error", err)
		response.InternalServerError(w, "Failed to delete business")
		return
	}

	response.OK(w, "Business deleted successfully", nil)
}

func (h *BusinessHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "Invalid business ID", nil)
		return
	}

	business, err := h.businessService.GetByID(r.Context(), id)
	if err != nil {
		if err == domain.ErrBusinessNotFound {
			response.NotFound(w, "Business not found")
			return
		}
		h.logger.Errorw("failed to get business", "id", id, "error", err)
		response.InternalServerError(w, "Failed to get business")
		return
	}

	response.OK(w, "Business retrieved successfully", business)
}

func (h *BusinessHandler) List(w http.ResponseWriter, r *http.Request) {
	var req dto.BusinessListRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body", nil)
		return
	}

	result, err := h.businessService.List(r.Context(), &req)
	if err != nil {
		h.logger.Errorw("failed to list businesses", "error", err)
		response.InternalServerError(w, "Failed to list businesses")
		return
	}

	response.OK(w, "Businesses retrieved successfully", result)
}

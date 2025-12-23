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

		h.logger.Errorw("failed to get business by ID", "businessID", id, "error", err)
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

	list, err := h.businessService.List(r.Context(), &req)
	if err != nil {
		h.logger.Errorw("failed to list businesses", "error", err)
		response.InternalServerError(w, "Failed to list businesses")
		return
	}

	response.OK(w, "Businesses listed successfully", list)
}

func (h *BusinessHandler) SetIsActive(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "Invalid business ID", nil)
		return
	}

	var req dto.BusinessUpdatePropertyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body", nil)
		return
	}
	req.ID = id

	if err := h.businessService.UpdateActiveStatus(r.Context(), &req); err != nil {
		if err == domain.ErrBusinessNotFound {
			response.NotFound(w, "Business not found")
			return
		}

		h.logger.Errorw("failed to set is_active flag", "businessID", id, "error", err)
		response.InternalServerError(w, "Failed to update is_active flag")
		return
	}

	response.OK(w, "Business active status updated successfully", nil)
}

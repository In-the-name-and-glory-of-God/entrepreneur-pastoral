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

type ServiceHandler struct {
	logger         *zap.SugaredLogger
	serviceService *application.ServiceService
}

func NewServiceHandler(logger *zap.SugaredLogger, serviceService *application.ServiceService) *ServiceHandler {
	return &ServiceHandler{
		logger:         logger,
		serviceService: serviceService,
	}
}

func (h *ServiceHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.ServiceCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body", nil)
		return
	}

	service, err := h.serviceService.Create(r.Context(), &req)
	if err != nil {
		if err == domain.ErrUnauthorized {
			response.Unauthorized(w, "Unauthorized to create service for this business")
			return
		}
		if err == domain.ErrBusinessNotFound {
			response.NotFound(w, "Business not found")
			return
		}
		h.logger.Errorw("failed to create service", "error", err)
		response.InternalServerError(w, "Failed to create service")
		return
	}

	response.Created(w, "Service created successfully", service)
}

func (h *ServiceHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "Invalid service ID", nil)
		return
	}

	var req dto.ServiceUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body", nil)
		return
	}
	req.ID = id

	if err := h.serviceService.Update(r.Context(), &req); err != nil {
		if err == domain.ErrServiceNotFound {
			response.NotFound(w, "Service not found")
			return
		}
		if err == domain.ErrUnauthorized {
			response.Unauthorized(w, "Unauthorized to update service")
			return
		}
		h.logger.Errorw("failed to update service", "id", id, "error", err)
		response.InternalServerError(w, "Failed to update service")
		return
	}

	response.OK(w, "Service updated successfully", nil)
}

func (h *ServiceHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "Invalid service ID", nil)
		return
	}

	if err := h.serviceService.Delete(r.Context(), id); err != nil {
		if err == domain.ErrUnauthorized {
			response.Unauthorized(w, "Unauthorized to delete service")
			return
		}
		h.logger.Errorw("failed to delete service", "id", id, "error", err)
		response.InternalServerError(w, "Failed to delete service")
		return
	}

	response.OK(w, "Service deleted successfully", nil)
}

func (h *ServiceHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "Invalid service ID", nil)
		return
	}

	service, err := h.serviceService.GetByID(r.Context(), id)
	if err != nil {
		if err == domain.ErrServiceNotFound {
			response.NotFound(w, "Service not found")
			return
		}
		h.logger.Errorw("failed to get service", "id", id, "error", err)
		response.InternalServerError(w, "Failed to get service")
		return
	}

	response.OK(w, "Service retrieved successfully", service)
}

func (h *ServiceHandler) List(w http.ResponseWriter, r *http.Request) {
	var req dto.ServiceListRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body", nil)
		return
	}

	result, err := h.serviceService.List(r.Context(), &req)
	if err != nil {
		h.logger.Errorw("failed to list services", "error", err)
		response.InternalServerError(w, "Failed to list services")
		return
	}

	response.OK(w, "Services retrieved successfully", result)
}

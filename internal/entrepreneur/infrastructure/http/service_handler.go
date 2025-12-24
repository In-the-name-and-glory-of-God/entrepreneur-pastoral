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
	ctx := r.Context()
	var req dto.ServiceCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequestT(ctx, w, "error.invalid_request_body", nil)
		return
	}

	service, err := h.serviceService.Create(ctx, &req)
	if err != nil {
		if err == domain.ErrUnauthorized {
			response.UnauthorizedT(ctx, w, "error.unauthorized_create_service")
			return
		}
		if err == domain.ErrBusinessNotFound {
			response.NotFoundT(ctx, w, "error.business_not_found")
			return
		}
		h.logger.Errorw("failed to create service", "error", err)
		response.InternalServerErrorT(ctx, w, "error.failed_create_service")
		return
	}

	response.CreatedT(ctx, w, "success.service_created", service)
}

func (h *ServiceHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequestT(ctx, w, "error.invalid_service_id", nil)
		return
	}

	var req dto.ServiceUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequestT(ctx, w, "error.invalid_request_body", nil)
		return
	}
	req.ID = id

	if err := h.serviceService.Update(ctx, &req); err != nil {
		if err == domain.ErrServiceNotFound {
			response.NotFoundT(ctx, w, "error.service_not_found")
			return
		}
		if err == domain.ErrUnauthorized {
			response.UnauthorizedT(ctx, w, "error.unauthorized_update_service")
			return
		}
		h.logger.Errorw("failed to update service", "id", id, "error", err)
		response.InternalServerErrorT(ctx, w, "error.failed_update_service")
		return
	}

	response.OKT(ctx, w, "success.service_updated", nil)
}

func (h *ServiceHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequestT(ctx, w, "error.invalid_service_id", nil)
		return
	}

	if err := h.serviceService.Delete(ctx, id); err != nil {
		if err == domain.ErrUnauthorized {
			response.UnauthorizedT(ctx, w, "error.unauthorized_delete_service")
			return
		}
		h.logger.Errorw("failed to delete service", "id", id, "error", err)
		response.InternalServerErrorT(ctx, w, "error.failed_delete_service")
		return
	}

	response.OKT(ctx, w, "success.service_deleted", nil)
}

func (h *ServiceHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequestT(ctx, w, "error.invalid_service_id", nil)
		return
	}

	service, err := h.serviceService.GetByID(ctx, id)
	if err != nil {
		if err == domain.ErrServiceNotFound {
			response.NotFoundT(ctx, w, "error.service_not_found")
			return
		}
		h.logger.Errorw("failed to get service", "id", id, "error", err)
		response.InternalServerErrorT(ctx, w, "error.failed_get_service")
		return
	}

	response.OKT(ctx, w, "success.service_retrieved", service)
}

func (h *ServiceHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req dto.ServiceListRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequestT(ctx, w, "error.invalid_request_body", nil)
		return
	}

	result, err := h.serviceService.List(ctx, &req)
	if err != nil {
		h.logger.Errorw("failed to list services", "error", err)
		response.InternalServerErrorT(ctx, w, "error.failed_list_services")
		return
	}

	response.OKT(ctx, w, "success.services_listed", result)
}

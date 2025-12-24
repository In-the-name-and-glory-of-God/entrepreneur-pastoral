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
	ctx := r.Context()
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequestT(ctx, w, "error.invalid_business_id", nil)
		return
	}

	business, err := h.businessService.GetByID(ctx, id)
	if err != nil {
		if err == domain.ErrBusinessNotFound {
			response.NotFoundT(ctx, w, "error.business_not_found")
			return
		}

		h.logger.Errorw("failed to get business by ID", "businessID", id, "error", err)
		response.InternalServerErrorT(ctx, w, "error.failed_get_business")
		return
	}

	response.OKT(ctx, w, "success.business_retrieved", business)
}

func (h *BusinessHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req dto.BusinessListRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequestT(ctx, w, "error.invalid_request_body", nil)
		return
	}

	list, err := h.businessService.List(ctx, &req)
	if err != nil {
		h.logger.Errorw("failed to list businesses", "error", err)
		response.InternalServerErrorT(ctx, w, "error.failed_list_businesses")
		return
	}

	response.OKT(ctx, w, "success.businesses_listed", list)
}

func (h *BusinessHandler) SetIsActive(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequestT(ctx, w, "error.invalid_business_id", nil)
		return
	}

	var req dto.BusinessUpdatePropertyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequestT(ctx, w, "error.invalid_request_body", nil)
		return
	}
	req.ID = id

	if err := h.businessService.UpdateActiveStatus(ctx, &req); err != nil {
		if err == domain.ErrBusinessNotFound {
			response.NotFoundT(ctx, w, "error.business_not_found")
			return
		}

		h.logger.Errorw("failed to set is_active flag", "businessID", id, "error", err)
		response.InternalServerErrorT(ctx, w, "error.failed_update_business_status")
		return
	}

	response.OKT(ctx, w, "success.business_status_updated", nil)
}

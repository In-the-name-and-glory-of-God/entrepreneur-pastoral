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
	ctx := r.Context()
	var req dto.IndustryCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequestT(ctx, w, "error.invalid_request_body", nil)
		return
	}

	industry, err := h.industryService.Create(ctx, &req)
	if err != nil {
		if err == domain.ErrIndustryAlreadyExists {
			response.ConflictT(ctx, w, "error.industry_already_exists", nil)
			return
		}

		h.logger.Errorw("failed to create industry", "error", err)
		response.InternalServerErrorT(ctx, w, "error.failed_create_industry")
		return
	}

	response.CreatedT(ctx, w, "success.industry_created", industry)
}

func (h *IndustryHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 16)
	if err != nil {
		response.BadRequestT(ctx, w, "error.invalid_industry_id", nil)
		return
	}

	var req dto.IndustryUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequestT(ctx, w, "error.invalid_request_body", nil)
		return
	}
	req.ID = int16(id)

	if err := h.industryService.Update(ctx, &req); err != nil {
		if err == domain.ErrIndustryNotFound {
			response.NotFoundT(ctx, w, "error.industry_not_found")
			return
		}
		if err == domain.ErrIndustryAlreadyExists {
			response.ConflictT(ctx, w, "error.industry_already_exists", nil)
			return
		}

		h.logger.Errorw("failed to update industry", "id", id, "error", err)
		response.InternalServerErrorT(ctx, w, "error.failed_update_industry")
		return
	}

	response.OKT(ctx, w, "success.industry_updated", nil)
}

func (h *IndustryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 16)
	if err != nil {
		response.BadRequestT(ctx, w, "error.invalid_industry_id", nil)
		return
	}

	if err := h.industryService.Delete(ctx, int16(id)); err != nil {
		if err == domain.ErrIndustryNotFound {
			response.NotFoundT(ctx, w, "error.industry_not_found")
			return
		}

		h.logger.Errorw("failed to delete industry", "id", id, "error", err)
		response.InternalServerErrorT(ctx, w, "error.failed_delete_industry")
		return
	}

	response.OKT(ctx, w, "success.industry_deleted", nil)
}

func (h *IndustryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 16)
	if err != nil {
		response.BadRequestT(ctx, w, "error.invalid_industry_id", nil)
		return
	}

	industry, err := h.industryService.GetByID(ctx, int16(id))
	if err != nil {
		if err == domain.ErrIndustryNotFound {
			response.NotFoundT(ctx, w, "error.industry_not_found")
			return
		}

		h.logger.Errorw("failed to get industry by ID", "id", id, "error", err)
		response.InternalServerErrorT(ctx, w, "error.failed_get_industry")
		return
	}

	response.OKT(ctx, w, "success.industry_retrieved", industry)
}

func (h *IndustryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	list, err := h.industryService.GetAll(ctx)
	if err != nil {
		h.logger.Errorw("failed to get all industries", "error", err)
		response.InternalServerErrorT(ctx, w, "error.failed_get_industries")
		return
	}

	response.OKT(ctx, w, "success.industries_retrieved", list)
}

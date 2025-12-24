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
	ctx := r.Context()
	var req dto.ChurchCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequestT(ctx, w, "error.invalid_request_body", nil)
		return
	}

	church, err := h.churchService.Create(ctx, &req)
	if err != nil {
		if err == domain.ErrChurchAlreadyExists {
			response.ConflictT(ctx, w, "error.church_already_exists", nil)
			return
		}

		h.logger.Errorw("failed to create church", "error", err)
		response.InternalServerErrorT(ctx, w, "error.failed_create_church")
		return
	}

	response.CreatedT(ctx, w, "success.church_created", church)
}

func (h *ChurchHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequestT(ctx, w, "error.invalid_church_id", nil)
		return
	}

	var req dto.ChurchUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequestT(ctx, w, "error.invalid_request_body", nil)
		return
	}
	req.ID = id

	if err := h.churchService.Update(ctx, &req); err != nil {
		if err == domain.ErrChurchNotFound {
			response.NotFoundT(ctx, w, "error.church_not_found")
			return
		}
		if err == domain.ErrChurchAlreadyExists {
			response.ConflictT(ctx, w, "error.church_already_exists", nil)
			return
		}

		h.logger.Errorw("failed to update church", "id", id, "error", err)
		response.InternalServerErrorT(ctx, w, "error.failed_update_church")
		return
	}

	response.OKT(ctx, w, "success.church_updated", nil)
}

func (h *ChurchHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequestT(ctx, w, "error.invalid_church_id", nil)
		return
	}

	if err := h.churchService.Delete(ctx, id); err != nil {
		if err == domain.ErrChurchNotFound {
			response.NotFoundT(ctx, w, "error.church_not_found")
			return
		}

		h.logger.Errorw("failed to delete church", "id", id, "error", err)
		response.InternalServerErrorT(ctx, w, "error.failed_delete_church")
		return
	}

	response.OKT(ctx, w, "success.church_deleted", nil)
}

func (h *ChurchHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequestT(ctx, w, "error.invalid_church_id", nil)
		return
	}

	church, err := h.churchService.GetByID(ctx, id)
	if err != nil {
		if err == domain.ErrChurchNotFound {
			response.NotFoundT(ctx, w, "error.church_not_found")
			return
		}

		h.logger.Errorw("failed to get church by ID", "id", id, "error", err)
		response.InternalServerErrorT(ctx, w, "error.failed_get_church")
		return
	}

	response.OKT(ctx, w, "success.church_retrieved", church)
}

func (h *ChurchHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req dto.ChurchListRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequestT(ctx, w, "error.invalid_request_body", nil)
		return
	}

	list, err := h.churchService.List(ctx, &req)
	if err != nil {
		h.logger.Errorw("failed to list churches", "error", err)
		response.InternalServerErrorT(ctx, w, "error.failed_list_churches")
		return
	}

	response.OKT(ctx, w, "success.churches_listed", list)
}

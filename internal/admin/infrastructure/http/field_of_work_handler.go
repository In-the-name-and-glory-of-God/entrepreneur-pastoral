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
	ctx := r.Context()
	var req dto.FieldOfWorkCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequestT(ctx, w, "error.invalid_request_body", nil)
		return
	}

	fieldOfWork, err := h.fieldOfWorkService.Create(ctx, &req)
	if err != nil {
		if err == domain.ErrFieldOfWorkAlreadyExists {
			response.ConflictT(ctx, w, "error.field_of_work_already_exists", nil)
			return
		}

		h.logger.Errorw("failed to create field of work", "error", err)
		response.InternalServerErrorT(ctx, w, "error.failed_create_field_of_work")
		return
	}

	response.CreatedT(ctx, w, "success.field_of_work_created", fieldOfWork)
}

func (h *FieldOfWorkHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 16)
	if err != nil {
		response.BadRequestT(ctx, w, "error.invalid_field_of_work_id", nil)
		return
	}

	var req dto.FieldOfWorkUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequestT(ctx, w, "error.invalid_request_body", nil)
		return
	}
	req.ID = int16(id)

	if err := h.fieldOfWorkService.Update(ctx, &req); err != nil {
		if err == domain.ErrFieldOfWorkNotFound {
			response.NotFoundT(ctx, w, "error.field_of_work_not_found")
			return
		}
		if err == domain.ErrFieldOfWorkAlreadyExists {
			response.ConflictT(ctx, w, "error.field_of_work_already_exists", nil)
			return
		}

		h.logger.Errorw("failed to update field of work", "id", id, "error", err)
		response.InternalServerErrorT(ctx, w, "error.failed_update_field_of_work")
		return
	}

	response.OKT(ctx, w, "success.field_of_work_updated", nil)
}

func (h *FieldOfWorkHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 16)
	if err != nil {
		response.BadRequestT(ctx, w, "error.invalid_field_of_work_id", nil)
		return
	}

	if err := h.fieldOfWorkService.Delete(ctx, int16(id)); err != nil {
		if err == domain.ErrFieldOfWorkNotFound {
			response.NotFoundT(ctx, w, "error.field_of_work_not_found")
			return
		}

		h.logger.Errorw("failed to delete field of work", "id", id, "error", err)
		response.InternalServerErrorT(ctx, w, "error.failed_delete_field_of_work")
		return
	}

	response.OKT(ctx, w, "success.field_of_work_deleted", nil)
}

func (h *FieldOfWorkHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 16)
	if err != nil {
		response.BadRequestT(ctx, w, "error.invalid_field_of_work_id", nil)
		return
	}

	fieldOfWork, err := h.fieldOfWorkService.GetByID(ctx, int16(id))
	if err != nil {
		if err == domain.ErrFieldOfWorkNotFound {
			response.NotFoundT(ctx, w, "error.field_of_work_not_found")
			return
		}

		h.logger.Errorw("failed to get field of work by ID", "id", id, "error", err)
		response.InternalServerErrorT(ctx, w, "error.failed_get_field_of_work")
		return
	}

	response.OKT(ctx, w, "success.field_of_work_retrieved", fieldOfWork)
}

func (h *FieldOfWorkHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	list, err := h.fieldOfWorkService.GetAll(ctx)
	if err != nil {
		h.logger.Errorw("failed to get all fields of work", "error", err)
		response.InternalServerErrorT(ctx, w, "error.failed_get_fields_of_work")
		return
	}

	response.OKT(ctx, w, "success.fields_of_work_retrieved", list)
}

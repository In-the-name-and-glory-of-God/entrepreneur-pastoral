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

type JobHandler struct {
	logger     *zap.SugaredLogger
	jobService *application.JobService
}

func NewJobHandler(logger *zap.SugaredLogger, jobService *application.JobService) *JobHandler {
	return &JobHandler{
		logger:     logger,
		jobService: jobService,
	}
}

func (h *JobHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req dto.JobCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequestT(ctx, w, "error.invalid_request_body", nil)
		return
	}

	job, err := h.jobService.Create(ctx, &req)
	if err != nil {
		if err == domain.ErrUnauthorized {
			response.UnauthorizedT(ctx, w, "error.unauthorized_create_job")
			return
		}
		if err == domain.ErrBusinessNotFound {
			response.NotFoundT(ctx, w, "error.business_not_found")
			return
		}
		h.logger.Errorw("failed to create job", "error", err)
		response.InternalServerErrorT(ctx, w, "error.failed_create_job")
		return
	}

	response.CreatedT(ctx, w, "success.job_created", job)
}

func (h *JobHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequestT(ctx, w, "error.invalid_job_id", nil)
		return
	}

	var req dto.JobUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequestT(ctx, w, "error.invalid_request_body", nil)
		return
	}
	req.ID = id

	if err := h.jobService.Update(ctx, &req); err != nil {
		if err == domain.ErrJobNotFound {
			response.NotFoundT(ctx, w, "error.job_not_found")
			return
		}
		if err == domain.ErrUnauthorized {
			response.UnauthorizedT(ctx, w, "error.unauthorized_update_job")
			return
		}
		h.logger.Errorw("failed to update job", "id", id, "error", err)
		response.InternalServerErrorT(ctx, w, "error.failed_update_job")
		return
	}

	response.OKT(ctx, w, "success.job_updated", nil)
}

func (h *JobHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequestT(ctx, w, "error.invalid_job_id", nil)
		return
	}

	if err := h.jobService.Delete(ctx, id); err != nil {
		if err == domain.ErrUnauthorized {
			response.UnauthorizedT(ctx, w, "error.unauthorized_delete_job")
			return
		}
		h.logger.Errorw("failed to delete job", "id", id, "error", err)
		response.InternalServerErrorT(ctx, w, "error.failed_delete_job")
		return
	}

	response.OKT(ctx, w, "success.job_deleted", nil)
}

func (h *JobHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequestT(ctx, w, "error.invalid_job_id", nil)
		return
	}

	job, err := h.jobService.GetByID(ctx, id)
	if err != nil {
		if err == domain.ErrJobNotFound {
			response.NotFoundT(ctx, w, "error.job_not_found")
			return
		}
		h.logger.Errorw("failed to get job", "id", id, "error", err)
		response.InternalServerErrorT(ctx, w, "error.failed_get_job")
		return
	}

	response.OKT(ctx, w, "success.job_retrieved", job)
}

func (h *JobHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req dto.JobListRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequestT(ctx, w, "error.invalid_request_body", nil)
		return
	}

	result, err := h.jobService.List(ctx, &req)
	if err != nil {
		h.logger.Errorw("failed to list jobs", "error", err)
		response.InternalServerErrorT(ctx, w, "error.failed_list_jobs")
		return
	}

	response.OKT(ctx, w, "success.jobs_listed", result)
}

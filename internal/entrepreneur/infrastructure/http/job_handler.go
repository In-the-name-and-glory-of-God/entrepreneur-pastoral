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
	var req dto.JobCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body", nil)
		return
	}

	job, err := h.jobService.Create(r.Context(), &req)
	if err != nil {
		if err == domain.ErrUnauthorized {
			response.Unauthorized(w, "Unauthorized to create job for this business")
			return
		}
		if err == domain.ErrBusinessNotFound {
			response.NotFound(w, "Business not found")
			return
		}
		h.logger.Errorw("failed to create job", "error", err)
		response.InternalServerError(w, "Failed to create job")
		return
	}

	response.Created(w, "Job created successfully", job)
}

func (h *JobHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "Invalid job ID", nil)
		return
	}

	var req dto.JobUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body", nil)
		return
	}
	req.ID = id

	if err := h.jobService.Update(r.Context(), &req); err != nil {
		if err == domain.ErrJobNotFound {
			response.NotFound(w, "Job not found")
			return
		}
		if err == domain.ErrUnauthorized {
			response.Unauthorized(w, "Unauthorized to update job")
			return
		}
		h.logger.Errorw("failed to update job", "id", id, "error", err)
		response.InternalServerError(w, "Failed to update job")
		return
	}

	response.OK(w, "Job updated successfully", nil)
}

func (h *JobHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "Invalid job ID", nil)
		return
	}

	if err := h.jobService.Delete(r.Context(), id); err != nil {
		if err == domain.ErrUnauthorized {
			response.Unauthorized(w, "Unauthorized to delete job")
			return
		}
		h.logger.Errorw("failed to delete job", "id", id, "error", err)
		response.InternalServerError(w, "Failed to delete job")
		return
	}

	response.OK(w, "Job deleted successfully", nil)
}

func (h *JobHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "Invalid job ID", nil)
		return
	}

	job, err := h.jobService.GetByID(r.Context(), id)
	if err != nil {
		if err == domain.ErrJobNotFound {
			response.NotFound(w, "Job not found")
			return
		}
		h.logger.Errorw("failed to get job", "id", id, "error", err)
		response.InternalServerError(w, "Failed to get job")
		return
	}

	response.OK(w, "Job retrieved successfully", job)
}

func (h *JobHandler) List(w http.ResponseWriter, r *http.Request) {
	var req dto.JobListRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body", nil)
		return
	}

	result, err := h.jobService.List(r.Context(), &req)
	if err != nil {
		h.logger.Errorw("failed to list jobs", "error", err)
		response.InternalServerError(w, "Failed to list jobs")
		return
	}

	response.OK(w, "Jobs retrieved successfully", result)
}

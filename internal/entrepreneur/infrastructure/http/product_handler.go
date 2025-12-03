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

type ProductHandler struct {
	logger         *zap.SugaredLogger
	productService *application.ProductService
}

func NewProductHandler(logger *zap.SugaredLogger, productService *application.ProductService) *ProductHandler {
	return &ProductHandler{
		logger:         logger,
		productService: productService,
	}
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.ProductCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body", nil)
		return
	}

	product, err := h.productService.Create(r.Context(), &req)
	if err != nil {
		if err == domain.ErrUnauthorized {
			response.Unauthorized(w, "Unauthorized to create product for this business")
			return
		}
		if err == domain.ErrBusinessNotFound {
			response.NotFound(w, "Business not found")
			return
		}
		h.logger.Errorw("failed to create product", "error", err)
		response.InternalServerError(w, "Failed to create product")
		return
	}

	response.Created(w, "Product created successfully", product)
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "Invalid product ID", nil)
		return
	}

	var req dto.ProductUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body", nil)
		return
	}
	req.ID = id

	if err := h.productService.Update(r.Context(), &req); err != nil {
		if err == domain.ErrProductNotFound {
			response.NotFound(w, "Product not found")
			return
		}
		if err == domain.ErrUnauthorized {
			response.Unauthorized(w, "Unauthorized to update product")
			return
		}
		h.logger.Errorw("failed to update product", "id", id, "error", err)
		response.InternalServerError(w, "Failed to update product")
		return
	}

	response.OK(w, "Product updated successfully", nil)
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "Invalid product ID", nil)
		return
	}

	if err := h.productService.Delete(r.Context(), id); err != nil {
		if err == domain.ErrUnauthorized {
			response.Unauthorized(w, "Unauthorized to delete product")
			return
		}
		h.logger.Errorw("failed to delete product", "id", id, "error", err)
		response.InternalServerError(w, "Failed to delete product")
		return
	}

	response.OK(w, "Product deleted successfully", nil)
}

func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "Invalid product ID", nil)
		return
	}

	product, err := h.productService.GetByID(r.Context(), id)
	if err != nil {
		if err == domain.ErrProductNotFound {
			response.NotFound(w, "Product not found")
			return
		}
		h.logger.Errorw("failed to get product", "id", id, "error", err)
		response.InternalServerError(w, "Failed to get product")
		return
	}

	response.OK(w, "Product retrieved successfully", product)
}

func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	var req dto.ProductListRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body", nil)
		return
	}

	result, err := h.productService.List(r.Context(), &req)
	if err != nil {
		h.logger.Errorw("failed to list products", "error", err)
		response.InternalServerError(w, "Failed to list products")
		return
	}

	response.OK(w, "Products retrieved successfully", result)
}

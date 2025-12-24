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
	ctx := r.Context()
	var req dto.ProductCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequestT(ctx, w, "error.invalid_request_body", nil)
		return
	}

	product, err := h.productService.Create(ctx, &req)
	if err != nil {
		if err == domain.ErrUnauthorized {
			response.UnauthorizedT(ctx, w, "error.unauthorized_create_product")
			return
		}
		if err == domain.ErrBusinessNotFound {
			response.NotFoundT(ctx, w, "error.business_not_found")
			return
		}
		h.logger.Errorw("failed to create product", "error", err)
		response.InternalServerErrorT(ctx, w, "error.failed_create_product")
		return
	}

	response.CreatedT(ctx, w, "success.product_created", product)
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequestT(ctx, w, "error.invalid_product_id", nil)
		return
	}

	var req dto.ProductUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequestT(ctx, w, "error.invalid_request_body", nil)
		return
	}
	req.ID = id

	if err := h.productService.Update(ctx, &req); err != nil {
		if err == domain.ErrProductNotFound {
			response.NotFoundT(ctx, w, "error.product_not_found")
			return
		}
		if err == domain.ErrUnauthorized {
			response.UnauthorizedT(ctx, w, "error.unauthorized_update_product")
			return
		}
		h.logger.Errorw("failed to update product", "id", id, "error", err)
		response.InternalServerErrorT(ctx, w, "error.failed_update_product")
		return
	}

	response.OKT(ctx, w, "success.product_updated", nil)
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequestT(ctx, w, "error.invalid_product_id", nil)
		return
	}

	if err := h.productService.Delete(ctx, id); err != nil {
		if err == domain.ErrUnauthorized {
			response.UnauthorizedT(ctx, w, "error.unauthorized_delete_product")
			return
		}
		h.logger.Errorw("failed to delete product", "id", id, "error", err)
		response.InternalServerErrorT(ctx, w, "error.failed_delete_product")
		return
	}

	response.OKT(ctx, w, "success.product_deleted", nil)
}

func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequestT(ctx, w, "error.invalid_product_id", nil)
		return
	}

	product, err := h.productService.GetByID(ctx, id)
	if err != nil {
		if err == domain.ErrProductNotFound {
			response.NotFoundT(ctx, w, "error.product_not_found")
			return
		}
		h.logger.Errorw("failed to get product", "id", id, "error", err)
		response.InternalServerErrorT(ctx, w, "error.failed_get_product")
		return
	}

	response.OKT(ctx, w, "success.product_retrieved", product)
}

func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req dto.ProductListRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequestT(ctx, w, "error.invalid_request_body", nil)
		return
	}

	result, err := h.productService.List(ctx, &req)
	if err != nil {
		h.logger.Errorw("failed to list products", "error", err)
		response.InternalServerErrorT(ctx, w, "error.failed_list_products")
		return
	}

	response.OKT(ctx, w, "success.products_listed", result)
}

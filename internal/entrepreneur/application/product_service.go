package application

import (
	"context"
	"database/sql"

	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/entrepreneur/domain"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/entrepreneur/infrastructure/dto"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/helper/response"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ProductService struct {
	logger      *zap.SugaredLogger
	productRepo domain.ProductRepository
}

func NewProductService(logger *zap.SugaredLogger, productRepo domain.ProductRepository) *ProductService {
	return &ProductService{
		logger:      logger,
		productRepo: productRepo,
	}
}

func (s *ProductService) Create(ctx context.Context, req *dto.ProductCreateRequest) (*domain.Product, error) {
	product := &domain.Product{
		BusinessID:  req.BusinessID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		ImageURL:    sql.NullString{String: req.ImageURL, Valid: req.ImageURL != ""},
		IsAvailable: req.IsAvailable,
	}

	if err := s.productRepo.Create(nil, product); err != nil {
		s.logger.Errorw("failed to create product", "error", err)
		return nil, response.ErrInternalServerError
	}

	return product, nil
}

func (s *ProductService) Update(ctx context.Context, req *dto.ProductUpdateRequest) error {
	product, err := s.productRepo.GetByID(ctx, req.ID)
	if err != nil {
		return err
	}

	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price
	product.ImageURL = sql.NullString{String: req.ImageURL, Valid: req.ImageURL != ""}
	product.IsAvailable = req.IsAvailable

	if err := s.productRepo.Update(nil, product); err != nil {
		s.logger.Errorw("failed to update product", "id", req.ID, "error", err)
		return response.ErrInternalServerError
	}

	return nil
}

func (s *ProductService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.productRepo.Delete(nil, id); err != nil {
		s.logger.Errorw("failed to delete product", "id", id, "error", err)
		return response.ErrInternalServerError
	}

	return nil
}

func (s *ProductService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil && err != domain.ErrProductNotFound {
		s.logger.Errorw("failed to get product by ID", "id", id, "error", err)
		return nil, response.ErrInternalServerError
	}

	return product, nil
}

func (s *ProductService) List(ctx context.Context, req *dto.ProductListRequest) (*dto.ProductListResponse, error) {
	products, err := s.productRepo.List(ctx, req)
	if err != nil && err != domain.ErrProductNotFound {
		s.logger.Errorw("failed to list products", "error", err)
		return nil, response.ErrInternalServerError
	}

	count := 0
	if len(products) > 0 {
		count, err = s.productRepo.Count(ctx, req)
		if err != nil {
			s.logger.Errorw("failed to count products", "error", err)
			return nil, response.ErrInternalServerError
		}
	}

	return &dto.ProductListResponse{
		Products: products,
		Count:    count,
		Limit:    req.Limit,
		Offset:   req.Offset,
	}, nil
}

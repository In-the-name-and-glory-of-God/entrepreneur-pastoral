package application

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"

	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/entrepreneur/domain"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/entrepreneur/infrastructure/dto"
	userDto "github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/user/infrastructure/dto"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/helper/auth"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/helper/response"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/storage"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	businessCacheTTL     = 15 * time.Minute
	businessListCacheTTL = 5 * time.Minute
)

type BusinessService struct {
	logger       *zap.SugaredLogger
	cache        storage.CacheStorage
	businessRepo domain.BusinessRepository
}

func NewBusinessService(logger *zap.SugaredLogger, cache storage.CacheStorage, businessRepo domain.BusinessRepository) *BusinessService {
	return &BusinessService{
		logger:       logger,
		cache:        cache,
		businessRepo: businessRepo,
	}
}

func (s *BusinessService) Create(ctx context.Context, req *dto.BusinessCreateRequest) (*domain.Business, error) {
	userCtx := ctx.Value(auth.UserContextKey).(*userDto.UserAsContext)
	business := &domain.Business{
		UserID:           userCtx.ID,
		IndustryID:       req.IndustryID,
		Name:             req.Name,
		Description:      req.Description,
		Email:            req.Email,
		PhoneCountryCode: sql.NullString{String: req.PhoneCountryCode, Valid: req.PhoneCountryCode != ""},
		PhoneNumber:      sql.NullString{String: req.PhoneNumber, Valid: req.PhoneNumber != ""},
		WebsiteURL:       sql.NullString{String: req.WebsiteURL, Valid: req.WebsiteURL != ""},
		LogoURL:          sql.NullString{String: req.LogoURL, Valid: req.LogoURL != ""},
		IsActive:         true,
	}

	if err := s.businessRepo.Create(nil, business); err != nil {
		s.logger.Errorw("failed to create business", "error", err)
		return nil, response.ErrInternalServerError
	}

	// Invalidate list cache
	s.invalidateListCache(ctx)

	return business, nil
}

func (s *BusinessService) Update(ctx context.Context, req *dto.BusinessUpdateRequest) error {
	userCtx := ctx.Value(auth.UserContextKey).(*userDto.UserAsContext)
	business, err := s.businessRepo.GetByID(ctx, req.ID)
	if err != nil {
		return err
	}

	if business.UserID != userCtx.ID {
		return domain.ErrUnauthorized
	}

	business.IndustryID = req.IndustryID
	business.Name = req.Name
	business.Description = req.Description
	business.Email = req.Email
	business.PhoneCountryCode = sql.NullString{String: req.PhoneCountryCode, Valid: req.PhoneCountryCode != ""}
	business.PhoneNumber = sql.NullString{String: req.PhoneNumber, Valid: req.PhoneNumber != ""}
	business.WebsiteURL = sql.NullString{String: req.WebsiteURL, Valid: req.WebsiteURL != ""}
	business.LogoURL = sql.NullString{String: req.LogoURL, Valid: req.LogoURL != ""}
	business.IsActive = req.IsActive

	if err := s.businessRepo.Update(nil, business); err != nil {
		s.logger.Errorw("failed to update business", "id", req.ID, "error", err)
		return response.ErrInternalServerError
	}

	// Invalidate caches
	s.invalidateBusinessCache(ctx, req.ID)
	s.invalidateListCache(ctx)

	return nil
}

func (s *BusinessService) Delete(ctx context.Context, id uuid.UUID) error {
	userCtx := ctx.Value(auth.UserContextKey).(*userDto.UserAsContext)
	business, err := s.businessRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if business.UserID != userCtx.ID {
		return domain.ErrUnauthorized
	}

	if err := s.businessRepo.Delete(nil, id); err != nil {
		s.logger.Errorw("failed to delete business", "id", id, "error", err)
		return response.ErrInternalServerError
	}

	// Invalidate caches
	s.invalidateBusinessCache(ctx, id)
	s.invalidateListCache(ctx)

	return nil
}

func (s *BusinessService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Business, error) {
	// Try to get from cache first
	cacheKey := s.cache.BuildKey(storage.CACHE_PREFIX_BUSINESS, id.String())
	var business domain.Business
	if err := s.cache.Get(ctx, cacheKey, &business); err == nil {
		return &business, nil
	}

	// Cache miss - get from database
	businessFromDB, err := s.businessRepo.GetByID(ctx, id)
	if err != nil {
		if err == domain.ErrBusinessNotFound {
			return nil, err
		}

		s.logger.Errorw("failed to get business by ID", "id", id, "error", err)
		return nil, response.ErrInternalServerError
	}

	// Store in cache
	if err := s.cache.Set(ctx, cacheKey, businessFromDB, businessCacheTTL); err != nil {
		s.logger.Warnw("failed to cache business", "id", id, "error", err)
	}

	return businessFromDB, nil
}

func (s *BusinessService) List(ctx context.Context, req *dto.BusinessListRequest) (*dto.BusinessListResponse, error) {
	// Generate cache key based on filter parameters
	cacheKey := s.buildListCacheKey(req)

	// Try to get from cache
	var cachedResponse dto.BusinessListResponse
	if err := s.cache.Get(ctx, cacheKey, &cachedResponse); err == nil {
		return &cachedResponse, nil
	}

	// Cache miss - get from database
	businesses, err := s.businessRepo.List(ctx, req)
	if err != nil && err != domain.ErrBusinessNotFound {
		s.logger.Errorw("failed to list businesses", "error", err)
		return nil, response.ErrInternalServerError
	}

	count := 0
	if len(businesses) > 0 {
		count, err = s.businessRepo.Count(ctx, req)
		if err != nil {
			s.logger.Errorw("failed to count businesses", "error", err)
			return nil, response.ErrInternalServerError
		}
	}

	resp := &dto.BusinessListResponse{
		Businesses: businesses,
		Count:      count,
		Limit:      req.Limit,
		Offset:     req.Offset,
	}

	// Store in cache
	if err := s.cache.Set(ctx, cacheKey, resp, businessListCacheTTL); err != nil {
		s.logger.Warnw("failed to cache business list", "error", err)
	}

	return resp, nil
}

func (s *BusinessService) UpdateActiveStatus(ctx context.Context, req *dto.BusinessUpdatePropertyRequest) error {
	business, err := s.businessRepo.GetByID(ctx, req.ID)
	if err != nil {
		if err == domain.ErrBusinessNotFound {
			return domain.ErrBusinessNotFound
		}

		s.logger.Errorw("failed to get business by ID", "id", req.ID, "error", err)
		return response.ErrInternalServerError
	}

	business.IsActive = req.Value

	if err := s.businessRepo.UpdateProperty(ctx, req.ID, domain.BusinessIsActive, req.Value); err != nil {
		s.logger.Errorw("failed to update business active status", "id", req.ID, "error", err)
		return response.ErrInternalServerError
	}

	// Invalidate caches
	s.invalidateBusinessCache(ctx, req.ID)
	s.invalidateListCache(ctx)

	return nil
}

// Cache helper methods

// buildListCacheKey generates a unique cache key based on filter parameters
func (s *BusinessService) buildListCacheKey(req *dto.BusinessListRequest) string {
	// Serialize the filter to JSON and hash it for a consistent key
	filterBytes, _ := json.Marshal(req)
	hash := sha256.Sum256(filterBytes)
	return s.cache.BuildKey(storage.CACHE_PREFIX_BUSINESS_LIST, hex.EncodeToString(hash[:8]))
}

// invalidateBusinessCache removes a specific business from cache
func (s *BusinessService) invalidateBusinessCache(ctx context.Context, id uuid.UUID) {
	cacheKey := s.cache.BuildKey(storage.CACHE_PREFIX_BUSINESS, id.String())
	if err := s.cache.Del(ctx, cacheKey); err != nil && !errors.Is(err, storage.ErrCacheMiss) {
		s.logger.Warnw("failed to invalidate business cache", "id", id, "error", err)
	}
}

// invalidateListCache removes all business list caches using scan and delete
func (s *BusinessService) invalidateListCache(ctx context.Context) {
	pattern := s.cache.BuildKey(storage.CACHE_PREFIX_BUSINESS_LIST, "*")
	keys, err := s.cache.Scan(ctx, pattern)
	if err != nil {
		s.logger.Warnw("failed to scan business list cache keys", "error", err)
		return
	}

	for _, key := range keys {
		if err := s.cache.Del(ctx, key); err != nil && !errors.Is(err, storage.ErrCacheMiss) {
			s.logger.Warnw("failed to invalidate business list cache", "key", key, "error", err)
		}
	}
}

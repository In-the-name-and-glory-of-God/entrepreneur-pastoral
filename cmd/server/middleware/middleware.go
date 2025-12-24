package middleware

import (
	"context"
	"net/http"
	"slices"
	"strings"

	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/user/domain"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/user/infrastructure/dto"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/helper/auth"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/helper/response"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/i18n"
	"github.com/google/uuid"
)

type Middleware struct {
	UserPersistence domain.UserRepository
	TokenManager    *auth.TokenManager
}

func NewMiddleware(userRepo domain.UserRepository, tokenManager *auth.TokenManager) *Middleware {
	return &Middleware{
		UserPersistence: userRepo,
		TokenManager:    tokenManager,
	}
}

// Language middleware parses Accept-Language header and sets language in context
func (m *Middleware) Language(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		acceptLang := r.Header.Get("Accept-Language")
		lang := i18n.ParseAcceptLanguage(acceptLang)

		ctx := i18n.SetLanguage(r.Context(), lang)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.UnauthorizedT(ctx, w, "error.missing_authorization_token")
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.UnauthorizedT(ctx, w, "error.invalid_authorization_header")
			return
		}

		token := parts[1]

		claims, err := m.TokenManager.ParseToken(token)
		if err != nil {
			if err == auth.ErrTokenExpired {
				response.UnauthorizedT(ctx, w, "error.token_expired")
				return
			}

			response.UnauthorizedT(ctx, w, "error.invalid_token")
			return
		}

		userID, err := uuid.Parse(claims.UserID)
		if err != nil {
			response.BadRequestT(ctx, w, "error.invalid_user_id", nil)
			return
		}

		user, err := m.UserPersistence.GetByID(ctx, userID)
		if err != nil {
			response.NotFoundT(ctx, w, "error.user_not_found")
			return
		}

		if !user.IsActive {
			response.UnauthorizedT(ctx, w, "error.user_inactive")
			return
		}

		if !user.IsVerified {
			response.UnauthorizedT(ctx, w, "error.email_not_verified")
			return
		}

		// Update language in context if user has a preferred language
		userLang := ""
		if user.Language.Valid && user.Language.String != "" {
			userLang = user.Language.String
			ctx = i18n.SetLanguage(ctx, i18n.Language(userLang))
		}

		ctx = context.WithValue(
			ctx,
			auth.UserContextKey,
			&dto.UserAsContext{
				ID:             user.ID,
				Email:          user.Email,
				RoleID:         user.RoleID,
				Language:       userLang,
				IsCatholic:     user.IsCatholic,
				IsEntrepreneur: user.IsEntrepreneur,
			},
		)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Middleware) Authorize(allowedRoles ...int16) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			userCtx := ctx.Value(auth.UserContextKey)
			if userCtx == nil {
				response.UnauthorizedT(ctx, w, "error.unauthorized")
				return
			}

			user, ok := userCtx.(*dto.UserAsContext)
			if !ok {
				response.UnauthorizedT(ctx, w, "error.unauthorized")
				return
			}

			// Check if user's role is in allowedRoles
			if slices.Contains(allowedRoles, user.RoleID) {
				next.ServeHTTP(w, r)
				return
			}

			response.ForbiddenT(ctx, w, "error.unauthorized")
		})
	}
}

func (m *Middleware) UserIsEntrepreneur(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userCtx := ctx.Value(auth.UserContextKey)
		if userCtx == nil {
			response.UnauthorizedT(ctx, w, "error.unauthorized")
			return
		}

		user, ok := userCtx.(*dto.UserAsContext)
		if !ok {
			response.UnauthorizedT(ctx, w, "error.unauthorized")
			return
		}

		if !user.IsEntrepreneur {
			response.ForbiddenT(ctx, w, "error.unauthorized")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) UserIsCatholic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userCtx := ctx.Value(auth.UserContextKey)
		if userCtx == nil {
			response.UnauthorizedT(ctx, w, "error.unauthorized")
			return
		}

		user, ok := userCtx.(*dto.UserAsContext)
		if !ok {
			response.UnauthorizedT(ctx, w, "error.unauthorized")
			return
		}

		if !user.IsCatholic {
			response.ForbiddenT(ctx, w, "error.unauthorized")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) RequireRoles(allowedRoles ...int16) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			userCtx := ctx.Value(auth.UserContextKey)
			if userCtx == nil {
				response.UnauthorizedT(ctx, w, "error.unauthorized")
				return
			}

			user, ok := userCtx.(*dto.UserAsContext)
			if !ok {
				response.UnauthorizedT(ctx, w, "error.unauthorized")
				return
			}

			// Check if user's role is in allowedRoles
			if slices.Contains(allowedRoles, user.RoleID) {
				next.ServeHTTP(w, r)
				return
			}

			response.ForbiddenT(ctx, w, "error.unauthorized")
		})
	}
}

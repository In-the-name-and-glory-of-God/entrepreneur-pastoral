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

func (m *Middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.Unauthorized(w, "Missing authorization token")
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(w, "Invalid authorization header format")
			return
		}

		token := parts[1]

		claims, err := m.TokenManager.ParseToken(token)
		if err != nil {
			if err == auth.ErrTokenExpired {
				response.Unauthorized(w, "Authorization token has expired")
				return
			}

			response.Unauthorized(w, "Invalid authorization token")
			return
		}

		userID, err := uuid.Parse(claims.UserID)
		if err != nil {
			response.BadRequest(w, "Invalid user ID", nil)
			return
		}

		user, err := m.UserPersistence.GetByID(r.Context(), userID)
		if err != nil {
			response.NotFound(w, "User not found")
			return
		}

		if !user.IsActive {
			response.Unauthorized(w, "User is inactive")
			return
		}

		if !user.IsVerified {
			response.Unauthorized(w, "Email not verified")
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(
			r.Context(),
			auth.UserContextKey,
			&dto.UserAsContext{
				ID:             user.ID,
				Email:          user.Email,
				RoleID:         user.RoleID,
				IsCatholic:     user.IsCatholic,
				IsEntrepreneur: user.IsEntrepreneur,
			},
		)))
	})
}

func (m *Middleware) Authorize(allowedRoles ...int16) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userCtx := r.Context().Value(auth.UserContextKey)
			if userCtx == nil {
				response.Unauthorized(w, "User not authenticated")
				return
			}

			user, ok := userCtx.(*dto.UserAsContext)
			if !ok {
				response.Unauthorized(w, "User not authenticated")
				return
			}

			// Check if user's role is in allowedRoles
			if slices.Contains(allowedRoles, user.RoleID) {
				next.ServeHTTP(w, r)
				return
			}

			response.Forbidden(w, "User does not have permission to access this resource")
		})
	}
}

func (m *Middleware) UserIsEntrepreneur(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userCtx := r.Context().Value(auth.UserContextKey)
		if userCtx == nil {
			response.Unauthorized(w, "User not authenticated")
			return
		}

		user, ok := userCtx.(*dto.UserAsContext)
		if !ok {
			response.Unauthorized(w, "User not authenticated")
			return
		}

		if !user.IsEntrepreneur {
			response.Forbidden(w, "User is not an entrepreneur")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) UserIsCatholic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userCtx := r.Context().Value(auth.UserContextKey)
		if userCtx == nil {
			response.Unauthorized(w, "User context missing")
			return
		}

		user, ok := userCtx.(*dto.UserAsContext)
		if !ok {
			response.Unauthorized(w, "User not authenticated")
			return
		}

		if !user.IsCatholic {
			response.Forbidden(w, "User is not catholic")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) RequireRoles(allowedRoles ...int16) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userCtx := r.Context().Value(auth.UserContextKey)
			if userCtx == nil {
				response.Unauthorized(w, "User not authenticated")
				return
			}

			user, ok := userCtx.(*dto.UserAsContext)
			if !ok {
				response.Unauthorized(w, "User not authenticated")
				return
			}

			// Check if user's role is in allowedRoles
			if slices.Contains(allowedRoles, user.RoleID) {
				next.ServeHTTP(w, r)
				return
			}

			response.Forbidden(w, "User does not have permission to access this resource")
		})
	}
}

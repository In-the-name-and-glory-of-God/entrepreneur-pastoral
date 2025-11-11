package middleware

import (
	"net/http"

	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/user/domain"
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
		token := r.Header.Get("Authorization")
		if token == "" {
			response.Unauthorized(w, "Missing authorization token")
			return
		}

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

		// Authentication logic here
		next.ServeHTTP(w, r)
	})
}

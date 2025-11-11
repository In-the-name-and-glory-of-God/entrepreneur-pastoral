package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrTokenExpired = jwt.ErrTokenExpired
	ErrInvalidToken = errors.New("invalid or expired token")
	ErrUnauthorized = errors.New("unauthorized access")
)

// Claims defines the structure of the JWT claims.
type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// TokenManager
type TokenManager struct {
	secret string
}

func NewTokenManager(secret string) *TokenManager {
	return &TokenManager{secret: secret}
}

// GenerateToken creates a new JWT token for the given user ID, email, and role ID.
func (t *TokenManager) GenerateToken(userID string) (string, error) {
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(t.secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ParseToken parses and validates a JWT token string, returning the claims if valid.
func (t *TokenManager) ParseToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(t.secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword []byte, password string) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
}

func IsValidEmail(email string) error {
	// Simple email validation logic
	if email == "" {
		return fmt.Errorf("email cannot be empty")
	}

	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return fmt.Errorf("invalid email format")
	}

	return nil
}

func IsStrongPassword(password string) error {
	// Implement password strength validation logic here
	// For example, check length, presence of numbers/special characters, etc.
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	} else if !strings.ContainsAny(password, "0123456789") {
		return fmt.Errorf("password must contain at least one number")
	} else if !strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		return fmt.Errorf("password must contain at least one uppercase letter")
	} else if !strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz") {
		return fmt.Errorf("password must contain at least one lowercase letter")
	} else if !strings.ContainsAny(password, "!@#$%^&*()-_=+[]{}|;:',.<>?/`~\"\\") {
		return fmt.Errorf("password must contain at least one special character")
	}

	return nil
}

func GenerateRandomToken(length int) (string, error) {
	return uuid.NewString(), nil
}

func SetRefreshTokenCookie(w http.ResponseWriter, value string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "rt",
		Value:    value,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/api/v1/auth",
		MaxAge:   604800,
	})
}

func GetRefreshTokenCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("rt")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

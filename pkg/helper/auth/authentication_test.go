package auth

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func TestNewTokenManager(t *testing.T) {
	secret := "test-secret"
	tm := NewTokenManager(secret)

	if tm == nil {
		t.Fatal("NewTokenManager returned nil")
	}

	if tm.secret != secret {
		t.Errorf("Expected secret %s, got %s", secret, tm.secret)
	}
}

func TestTokenManager_GenerateToken(t *testing.T) {
	tm := NewTokenManager("test-secret")
	userID := "user-123"

	token, err := tm.GenerateToken(userID)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	if token == "" {
		t.Error("Generated token is empty")
	}

	// Verify token can be parsed
	claims, err := tm.ParseToken(token)
	if err != nil {
		t.Fatalf("Failed to parse generated token: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("Expected UserID %s, got %s", userID, claims.UserID)
	}
}

func TestTokenManager_ParseToken(t *testing.T) {
	tm := NewTokenManager("test-secret")
	userID := "user-456"

	tests := []struct {
		name        string
		setupToken  func() string
		expectError bool
		errorCheck  func(error) bool
	}{
		{
			name: "valid token",
			setupToken: func() string {
				token, _ := tm.GenerateToken(userID)
				return token
			},
			expectError: false,
		},
		{
			name: "invalid token string",
			setupToken: func() string {
				return "invalid.token.string"
			},
			expectError: true,
		},
		{
			name: "expired token",
			setupToken: func() string {
				claims := &Claims{
					UserID: userID,
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
						IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
					},
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, _ := token.SignedString([]byte("test-secret"))
				return tokenString
			},
			expectError: true,
			errorCheck: func(err error) bool {
				return strings.Contains(err.Error(), "expired")
			},
		},
		{
			name: "token with wrong secret",
			setupToken: func() string {
				wrongTM := NewTokenManager("wrong-secret")
				token, _ := wrongTM.GenerateToken(userID)
				return token
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := tt.setupToken()
			claims, err := tm.ParseToken(token)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				if tt.errorCheck != nil && !tt.errorCheck(err) {
					t.Errorf("Error check failed for error: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if claims == nil {
					t.Error("Claims should not be nil for valid token")
				}
				if claims != nil && claims.UserID != userID {
					t.Errorf("Expected UserID %s, got %s", userID, claims.UserID)
				}
			}
		})
	}
}

func TestHashPassword(t *testing.T) {
	password := "mySecurePassword123!"

	hashed, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if len(hashed) == 0 {
		t.Error("Hashed password is empty")
	}

	// Verify the hash is a valid bcrypt hash
	err = bcrypt.CompareHashAndPassword(hashed, []byte(password))
	if err != nil {
		t.Error("Hashed password does not match original password")
	}
}

func TestVerifyPassword(t *testing.T) {
	password := "mySecurePassword123!"
	wrongPassword := "wrongPassword123!"

	hashed, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	tests := []struct {
		name        string
		password    string
		expectError bool
	}{
		{
			name:        "correct password",
			password:    password,
			expectError: false,
		},
		{
			name:        "incorrect password",
			password:    wrongPassword,
			expectError: true,
		},
		{
			name:        "empty password",
			password:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := VerifyPassword(hashed, tt.password)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		name        string
		email       string
		expectError bool
	}{
		{
			name:        "valid email",
			email:       "test@example.com",
			expectError: false,
		},
		{
			name:        "valid email with subdomain",
			email:       "user@mail.example.com",
			expectError: false,
		},
		{
			name:        "empty email",
			email:       "",
			expectError: true,
		},
		{
			name:        "missing @",
			email:       "testexample.com",
			expectError: true,
		},
		{
			name:        "missing domain",
			email:       "test@",
			expectError: true,
		},
		{
			name:        "missing dot - valid per RFC 5322",
			email:       "test@example",
			expectError: false, // RFC 5322 allows domains without dots (e.g., localhost)
		},
		{
			name:        "missing username",
			email:       "@example.com",
			expectError: true, // RFC 5322 requires a local part before @
		},
		{
			name:        "valid email with plus sign",
			email:       "user+tag@example.com",
			expectError: false,
		},
		{
			name:        "valid email with dots in local part",
			email:       "first.last@example.com",
			expectError: false,
		},
		{
			name:        "invalid email with spaces",
			email:       "test @example.com",
			expectError: true,
		},
		{
			name:        "invalid email with double @",
			email:       "test@@example.com",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := IsValidEmail(tt.email)

			if tt.expectError && err == nil {
				t.Errorf("Expected error for email %q but got none", tt.email)
			}

			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error for email %q: %v", tt.email, err)
			}
		})
	}
}

func TestIsStrongPassword(t *testing.T) {
	tests := []struct {
		name        string
		password    string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid strong password",
			password:    "MyPassword123!",
			expectError: false,
		},
		{
			name:        "too short",
			password:    "Pass1!",
			expectError: true,
			errorMsg:    "at least 8 characters",
		},
		{
			name:        "no numbers",
			password:    "MyPassword!",
			expectError: true,
			errorMsg:    "at least one number",
		},
		{
			name:        "no uppercase",
			password:    "mypassword123!",
			expectError: true,
			errorMsg:    "at least one uppercase letter",
		},
		{
			name:        "no lowercase",
			password:    "MYPASSWORD123!",
			expectError: true,
			errorMsg:    "at least one lowercase letter",
		},
		{
			name:        "no special character",
			password:    "MyPassword123",
			expectError: true,
			errorMsg:    "at least one special character",
		},
		{
			name:        "all requirements met",
			password:    "Abcdefgh1!",
			expectError: false,
		},
		{
			name:        "complex password",
			password:    "C0mpl3x!P@ssw0rd",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := IsStrongPassword(tt.password)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for password %q but got none", tt.password)
				} else if tt.errorMsg != "" && !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error message to contain %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for password %q: %v", tt.password, err)
				}
			}
		})
	}
}

func TestGenerateRandomToken(t *testing.T) {
	token1, err := GenerateRandomToken(32)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	if token1 == "" {
		t.Error("Generated token is empty")
	}

	// Verify the token length is 2 * requested length (hex encoding)
	expectedLen := 32 * 2 // hex encoding doubles the length
	if len(token1) != expectedLen {
		t.Errorf("Expected token length %d, got %d", expectedLen, len(token1))
	}

	// Verify it's a valid hex string
	for _, c := range token1 {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			t.Errorf("Token contains invalid hex character: %c", c)
			break
		}
	}

	// Generate another token to ensure they're different
	token2, err := GenerateRandomToken(32)
	if err != nil {
		t.Fatalf("Failed to generate second token: %v", err)
	}

	if token1 == token2 {
		t.Error("Generated tokens should be unique")
	}

	// Test different lengths
	token16, err := GenerateRandomToken(16)
	if err != nil {
		t.Fatalf("Failed to generate 16-byte token: %v", err)
	}
	if len(token16) != 32 { // 16 * 2 = 32 hex chars
		t.Errorf("Expected 16-byte token to have length 32, got %d", len(token16))
	}
}

func TestSetRefreshTokenCookie(t *testing.T) {
	w := httptest.NewRecorder()
	tokenValue := "refresh-token-value"

	SetRefreshTokenCookie(w, tokenValue)

	cookies := w.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatal("No cookies were set")
	}

	var refreshCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "rt" {
			refreshCookie = cookie
			break
		}
	}

	if refreshCookie == nil {
		t.Fatal("Refresh token cookie 'rt' was not found")
	}

	tests := []struct {
		name     string
		check    func() bool
		errorMsg string
	}{
		{
			name:     "cookie value",
			check:    func() bool { return refreshCookie.Value == tokenValue },
			errorMsg: "Cookie value mismatch",
		},
		{
			name:     "HttpOnly flag",
			check:    func() bool { return refreshCookie.HttpOnly },
			errorMsg: "Cookie should be HttpOnly",
		},
		{
			name:     "Secure flag",
			check:    func() bool { return refreshCookie.Secure },
			errorMsg: "Cookie should be Secure",
		},
		{
			name:     "SameSite attribute",
			check:    func() bool { return refreshCookie.SameSite == http.SameSiteStrictMode },
			errorMsg: "Cookie should have SameSite=Strict",
		},
		{
			name:     "Path attribute",
			check:    func() bool { return refreshCookie.Path == "/api/v1/auth" },
			errorMsg: "Cookie path should be /api/v1/auth",
		},
		{
			name:     "MaxAge attribute",
			check:    func() bool { return refreshCookie.MaxAge == 604800 },
			errorMsg: "Cookie MaxAge should be 604800 (7 days)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.check() {
				t.Error(tt.errorMsg)
			}
		})
	}
}

func TestGetRefreshTokenCookie(t *testing.T) {
	tests := []struct {
		name        string
		setupReq    func() *http.Request
		expectValue string
		expectError bool
	}{
		{
			name: "cookie exists",
			setupReq: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				req.AddCookie(&http.Cookie{
					Name:  "rt",
					Value: "my-refresh-token",
				})
				return req
			},
			expectValue: "my-refresh-token",
			expectError: false,
		},
		{
			name: "cookie does not exist",
			setupReq: func() *http.Request {
				return httptest.NewRequest(http.MethodGet, "/", nil)
			},
			expectValue: "",
			expectError: true,
		},
		{
			name: "wrong cookie name",
			setupReq: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				req.AddCookie(&http.Cookie{
					Name:  "other-cookie",
					Value: "other-value",
				})
				return req
			},
			expectValue: "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := tt.setupReq()
			value, err := GetRefreshTokenCookie(req)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if value != tt.expectValue {
					t.Errorf("Expected value %q, got %q", tt.expectValue, value)
				}
			}
		})
	}
}

func TestTokenManager_Integration(t *testing.T) {
	// Integration test: Generate, parse, and verify token lifecycle
	tm := NewTokenManager("integration-test-secret")
	userID := "integration-user-789"

	// Generate token
	token, err := tm.GenerateToken(userID)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Parse token
	claims, err := tm.ParseToken(token)
	if err != nil {
		t.Fatalf("Failed to parse token: %v", err)
	}

	// Verify claims
	if claims.UserID != userID {
		t.Errorf("Expected UserID %s, got %s", userID, claims.UserID)
	}

	if claims.ExpiresAt.Before(time.Now()) {
		t.Error("Token should not be expired")
	}

	if claims.IssuedAt.After(time.Now()) {
		t.Error("Token IssuedAt should be in the past")
	}

	// Verify token expires in approximately 24 hours
	expectedExpiry := time.Now().Add(24 * time.Hour)
	timeDiff := claims.ExpiresAt.Time.Sub(expectedExpiry)
	if timeDiff > time.Minute || timeDiff < -time.Minute {
		t.Errorf("Token expiry time is off by %v", timeDiff)
	}
}

func TestPasswordHashingIntegration(t *testing.T) {
	// Integration test: Hash and verify password lifecycle
	password := "IntegrationTest123!"

	// Hash password
	hashed, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	// Verify correct password
	err = VerifyPassword(hashed, password)
	if err != nil {
		t.Errorf("Failed to verify correct password: %v", err)
	}

	// Verify incorrect password
	err = VerifyPassword(hashed, "WrongPassword123!")
	if err == nil {
		t.Error("Should have failed to verify incorrect password")
	}

	// Hash the same password again - should produce different hash
	hashed2, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password second time: %v", err)
	}

	if string(hashed) == string(hashed2) {
		t.Error("Same password should produce different hashes due to salt")
	}

	// But both hashes should verify the same password
	err = VerifyPassword(hashed2, password)
	if err != nil {
		t.Error("Second hash should also verify the password")
	}
}

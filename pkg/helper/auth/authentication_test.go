package auth

import (
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestNewTokenManager(t *testing.T) {
	secret := "test-secret-key"
	tm := NewTokenManager(secret)

	if tm == nil {
		t.Fatal("Expected TokenManager to be created, got nil")
	}

	if tm.secret != secret {
		t.Errorf("Expected secret '%s', got '%s'", secret, tm.secret)
	}

	if tm.claims == nil {
		t.Error("Expected claims to be initialized")
	}
}

func TestGenerateToken(t *testing.T) {
	tests := []struct {
		name   string
		secret string
		userID string
		email  string
		roleID int16
	}{
		{
			name:   "valid token generation",
			secret: "test-secret-key",
			userID: "123e4567-e89b-12d3-a456-426614174000",
			email:  "test@example.com",
			roleID: 1,
		},
		{
			name:   "different user",
			secret: "another-secret",
			userID: "987e6543-e89b-12d3-a456-426614174999",
			email:  "user@example.com",
			roleID: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := NewTokenManager(tt.secret)
			token, err := tm.GenerateToken(tt.userID, tt.email, tt.roleID)

			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if token == "" {
				t.Error("Expected token string, got empty string")
			}

			// Verify token has 3 parts (header.payload.signature)
			parts := strings.Split(token, ".")
			if len(parts) != 3 {
				t.Errorf("Expected 3 token parts, got %d", len(parts))
			}

			// Verify claims were set correctly
			if tm.claims.UserID != tt.userID {
				t.Errorf("Expected UserID '%s', got '%s'", tt.userID, tm.claims.UserID)
			}

			if tm.claims.Email != tt.email {
				t.Errorf("Expected Email '%s', got '%s'", tt.email, tm.claims.Email)
			}

			if tm.claims.RoleID != tt.roleID {
				t.Errorf("Expected RoleID %d, got %d", tt.roleID, tm.claims.RoleID)
			}

			// Verify expiration is set (24 hours from now)
			if tm.claims.ExpiresAt == nil {
				t.Error("Expected ExpiresAt to be set")
			}

			// Verify IssuedAt is set
			if tm.claims.IssuedAt == nil {
				t.Error("Expected IssuedAt to be set")
			}
		})
	}
}

func TestParseToken(t *testing.T) {
	secret := "test-secret-key"
	userID := "123e4567-e89b-12d3-a456-426614174000"
	email := "test@example.com"
	roleID := int16(1)

	tests := []struct {
		name          string
		setupToken    func() string
		expectedError bool
		errorContains string
	}{
		{
			name: "valid token",
			setupToken: func() string {
				tm := NewTokenManager(secret)
				token, _ := tm.GenerateToken(userID, email, roleID)
				return token
			},
			expectedError: false,
		},
		{
			name: "invalid token format",
			setupToken: func() string {
				return "invalid.token"
			},
			expectedError: true,
			errorContains: "failed to parse token",
		},
		{
			name: "expired token",
			setupToken: func() string {
				claims := &Claims{
					UserID: userID,
					Email:  email,
					RoleID: roleID,
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
						IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
					},
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, _ := token.SignedString([]byte(secret))
				return tokenString
			},
			expectedError: true,
			errorContains: "failed to parse token",
		},
		{
			name: "wrong secret",
			setupToken: func() string {
				tm := NewTokenManager("different-secret")
				token, _ := tm.GenerateToken(userID, email, roleID)
				return token
			},
			expectedError: true,
			errorContains: "failed to parse token",
		},
		{
			name: "tampered token",
			setupToken: func() string {
				tm := NewTokenManager(secret)
				token, _ := tm.GenerateToken(userID, email, roleID)
				// Tamper with the signature
				parts := strings.Split(token, ".")
				if len(parts) == 3 {
					parts[2] = "tampered"
					return strings.Join(parts, ".")
				}
				return token
			},
			expectedError: true,
			errorContains: "failed to parse token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := NewTokenManager(secret)
			token := tt.setupToken()
			err := tm.ParseToken(token)

			if tt.expectedError {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error to contain '%s', got '%s'", tt.errorContains, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}

				// Verify claims were parsed correctly
				if tm.claims.UserID != userID {
					t.Errorf("Expected UserID '%s', got '%s'", userID, tm.claims.UserID)
				}

				if tm.claims.Email != email {
					t.Errorf("Expected Email '%s', got '%s'", email, tm.claims.Email)
				}

				if tm.claims.RoleID != roleID {
					t.Errorf("Expected RoleID %d, got %d", roleID, tm.claims.RoleID)
				}
			}
		})
	}
}

func TestGenerateAndParseTokenIntegration(t *testing.T) {
	secret := "integration-test-secret"
	userID := "test-user-id"
	email := "integration@example.com"
	roleID := int16(3)

	// Generate token
	tm1 := NewTokenManager(secret)
	token, err := tm1.GenerateToken(userID, email, roleID)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Parse token with new instance
	tm2 := NewTokenManager(secret)
	err = tm2.ParseToken(token)
	if err != nil {
		t.Fatalf("Failed to parse token: %v", err)
	}

	// Verify claims match
	if tm2.claims.UserID != userID {
		t.Errorf("Expected UserID '%s', got '%s'", userID, tm2.claims.UserID)
	}

	if tm2.claims.Email != email {
		t.Errorf("Expected Email '%s', got '%s'", email, tm2.claims.Email)
	}

	if tm2.claims.RoleID != roleID {
		t.Errorf("Expected RoleID %d, got %d", roleID, tm2.claims.RoleID)
	}
}

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
	}{
		{
			name:     "simple password",
			password: "password123",
		},
		{
			name:     "complex password",
			password: "P@ssw0rd!Complex#2024",
		},
		{
			name:     "short password",
			password: "abc",
		},
		{
			name:     "long password",
			password: "this-is-a-very-long-password-with-many-characters-1234567890",
		},
		{
			name:     "special characters",
			password: "!@#$%^&*()_+-=[]{}|;:',.<>?/`~",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hashed, err := HashPassword(tt.password)

			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if len(hashed) == 0 {
				t.Error("Expected hashed password, got empty byte slice")
			}

			// Verify hashed password is different from original
			if string(hashed) == tt.password {
				t.Error("Hashed password should be different from original")
			}

			// Verify hash starts with bcrypt identifier
			hashedStr := string(hashed)
			if !strings.HasPrefix(hashedStr, "$2a$") && !strings.HasPrefix(hashedStr, "$2b$") {
				t.Error("Expected bcrypt hash format")
			}
		})
	}
}

func TestVerifyPassword(t *testing.T) {
	password := "testpassword123"
	hashedPassword, _ := HashPassword(password)

	tests := []struct {
		name          string
		hashedPwd     []byte
		plainPwd      string
		expectedError bool
	}{
		{
			name:          "correct password",
			hashedPwd:     hashedPassword,
			plainPwd:      password,
			expectedError: false,
		},
		{
			name:          "incorrect password",
			hashedPwd:     hashedPassword,
			plainPwd:      "wrongpassword",
			expectedError: true,
		},
		{
			name:          "empty password",
			hashedPwd:     hashedPassword,
			plainPwd:      "",
			expectedError: true,
		},
		{
			name:          "case sensitive",
			hashedPwd:     hashedPassword,
			plainPwd:      "TestPassword123",
			expectedError: true,
		},
		{
			name:          "extra characters",
			hashedPwd:     hashedPassword,
			plainPwd:      password + "extra",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := VerifyPassword(tt.hashedPwd, tt.plainPwd)

			if tt.expectedError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			}
		})
	}
}

func TestHashPasswordDeterminism(t *testing.T) {
	password := "test123"

	hash1, err1 := HashPassword(password)
	hash2, err2 := HashPassword(password)

	if err1 != nil || err2 != nil {
		t.Fatal("Expected no errors during hashing")
	}

	// Hashes should be different due to random salt
	if string(hash1) == string(hash2) {
		t.Error("Expected different hashes for same password (bcrypt uses random salt)")
	}

	// But both should verify successfully
	if err := VerifyPassword(hash1, password); err != nil {
		t.Error("First hash should verify successfully")
	}

	if err := VerifyPassword(hash2, password); err != nil {
		t.Error("Second hash should verify successfully")
	}
}

func TestTokenManagerClaimsIsolation(t *testing.T) {
	secret := "test-secret"
	tm := NewTokenManager(secret)

	// Generate first token
	token1, _ := tm.GenerateToken("user1", "user1@example.com", 1)

	// Generate second token with different data
	token2, _ := tm.GenerateToken("user2", "user2@example.com", 2)

	// Parse first token again
	tm2 := NewTokenManager(secret)
	err := tm2.ParseToken(token1)
	if err != nil {
		t.Fatalf("Failed to parse first token: %v", err)
	}

	// Verify it contains first user's data, not second
	if tm2.claims.UserID != "user1" {
		t.Errorf("Expected UserID 'user1', got '%s'", tm2.claims.UserID)
	}

	if tm2.claims.Email != "user1@example.com" {
		t.Errorf("Expected Email 'user1@example.com', got '%s'", tm2.claims.Email)
	}

	// Parse second token with new instance
	tm3 := NewTokenManager(secret)
	err = tm3.ParseToken(token2)
	if err != nil {
		t.Fatalf("Failed to parse second token: %v", err)
	}

	// Verify it contains second user's data
	if tm3.claims.UserID != "user2" {
		t.Errorf("Expected UserID 'user2', got '%s'", tm3.claims.UserID)
	}
}

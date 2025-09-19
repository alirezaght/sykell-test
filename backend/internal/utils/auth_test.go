package utils

import (
	"strings"
	"testing"
	"time"

	"sykell-backend/internal/db"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "valid password",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "empty password",
			password: "",
			wantErr:  false, // bcrypt allows empty passwords
		},
		{
			name:     "long password",
			password: strings.Repeat("a", 72), // bcrypt max is 72 bytes
			wantErr:  false,
		},
		{
			name:     "too long password",
			password: strings.Repeat("a", 100), // exceeds bcrypt limit
			wantErr:  true,
		},
		{
			name:     "special characters",
			password: "p@ssw0rd!@#$%^&*()",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)
			
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, hash)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, hash)
				assert.NotEqual(t, tt.password, hash, "hash should not equal original password")
				assert.True(t, len(hash) > 50, "hash should be reasonably long")
				
				// Verify the hash is different each time (salt is random)
				hash2, err2 := HashPassword(tt.password)
				assert.NoError(t, err2)
				assert.NotEqual(t, hash, hash2, "hashes should be different due to random salt")
			}
		})
	}
}

func TestCheckPassword(t *testing.T) {
	password := "testpassword123"
	hash, err := HashPassword(password)
	require.NoError(t, err)

	tests := []struct {
		name           string
		hashedPassword string
		password       string
		wantErr        bool
	}{
		{
			name:           "correct password",
			hashedPassword: hash,
			password:       password,
			wantErr:        false,
		},
		{
			name:           "incorrect password",
			hashedPassword: hash,
			password:       "wrongpassword",
			wantErr:        true,
		},
		{
			name:           "empty password",
			hashedPassword: hash,
			password:       "",
			wantErr:        true,
		},
		{
			name:           "invalid hash",
			hashedPassword: "invalid-hash",
			password:       password,
			wantErr:        true,
		},
		{
			name:           "empty hash",
			hashedPassword: "",
			password:       password,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckPassword(tt.hashedPassword, tt.password)
			
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGenerateJWT(t *testing.T) {
	jwtSecret := []byte("test-secret-key")
	user := db.User{
		ID:    "550e8400-e29b-41d4-a716-446655440000",
		Email: "test@example.com",
	}

	tests := []struct {
		name      string
		user      db.User
		jwtSecret []byte
		wantErr   bool
	}{
		{
			name:      "valid user and secret",
			user:      user,
			jwtSecret: jwtSecret,
			wantErr:   false,
		},
		{
			name: "empty user ID",
			user: db.User{
				ID:    "",
				Email: "test@example.com",
			},
			jwtSecret: jwtSecret,
			wantErr:   false, // Should still work, just with empty ID
		},
		{
			name: "empty email",
			user: db.User{
				ID:    "550e8400-e29b-41d4-a716-446655440000",
				Email: "",
			},
			jwtSecret: jwtSecret,
			wantErr:   false, // Should still work, just with empty email
		},
		{
			name:      "empty secret",
			user:      user,
			jwtSecret: []byte(""),
			wantErr:   false, // JWT library allows empty secret
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, expiresAt, err := GenerateJWT(tt.user.ID, tt.user.Email, tt.jwtSecret)
			
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, token)
				assert.Zero(t, expiresAt)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
				assert.Greater(t, expiresAt, time.Now().Unix())
				
				// Verify token format (should have 3 parts separated by dots)
				parts := strings.Split(token, ".")
				assert.Len(t, parts, 3, "JWT should have 3 parts")
				
				// Verify expiration is approximately 24 hours from now
				expectedExpiry := time.Now().Add(24 * time.Hour).Unix()
				assert.InDelta(t, expectedExpiry, expiresAt, 60, "expiry should be within 1 minute of 24 hours from now")
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	jwtSecret := []byte("test-secret-key")
	user := db.User{
		ID:    "550e8400-e29b-41d4-a716-446655440000",
		Email: "test@example.com",
	}
	
	// Generate a valid token for testing
	validToken, _, err := GenerateJWT(user.ID, user.Email, jwtSecret)
	require.NoError(t, err)
	
	// Generate a token with different secret
	wrongSecretToken, _, err := GenerateJWT(user.ID, user.Email, []byte("wrong-secret"))
	require.NoError(t, err)
	
	// Generate an expired token
	expiredClaims := &JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // Expired 1 hour ago
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			NotBefore: jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			Issuer:    "sykell-backend",
			Subject:   user.Email,
		},
	}
	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	expiredTokenString, err := expiredToken.SignedString(jwtSecret)
	require.NoError(t, err)

	tests := []struct {
		name        string
		tokenString string
		jwtSecret   []byte
		wantErr     bool
		wantClaims  *JWTClaims
	}{
		{
			name:        "valid token",
			tokenString: validToken,
			jwtSecret:   jwtSecret,
			wantErr:     false,
			wantClaims: &JWTClaims{
				UserID: user.ID,
				Email:  user.Email,
			},
		},
		{
			name:        "token with wrong secret",
			tokenString: wrongSecretToken,
			jwtSecret:   jwtSecret,
			wantErr:     true,
			wantClaims:  nil,
		},
		{
			name:        "expired token",
			tokenString: expiredTokenString,
			jwtSecret:   jwtSecret,
			wantErr:     true,
			wantClaims:  nil,
		},
		{
			name:        "invalid token format",
			tokenString: "invalid.token.format",
			jwtSecret:   jwtSecret,
			wantErr:     true,
			wantClaims:  nil,
		},
		{
			name:        "empty token",
			tokenString: "",
			jwtSecret:   jwtSecret,
			wantErr:     true,
			wantClaims:  nil,
		},
		{
			name:        "malformed token",
			tokenString: "not-a-jwt-token",
			jwtSecret:   jwtSecret,
			wantErr:     true,
			wantClaims:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := ValidateJWT(tt.tokenString, tt.jwtSecret)
			
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, claims)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
				
				if tt.wantClaims != nil {
					assert.Equal(t, tt.wantClaims.UserID, claims.UserID)
					assert.Equal(t, tt.wantClaims.Email, claims.Email)
					assert.Equal(t, "sykell-backend", claims.Issuer)
					assert.Equal(t, user.Email, claims.Subject)
					assert.True(t, claims.ExpiresAt.After(time.Now()))
				}
			}
		})
	}
}

func TestJWTRoundTrip(t *testing.T) {
	// Test generating and then validating a JWT token
	jwtSecret := []byte("test-secret-key-for-roundtrip")
	user := db.User{
		ID:    "550e8400-e29b-41d4-a716-446655440000",
		Email: "roundtrip@example.com",
	}

	// Generate token
	token, expiresAt, err := GenerateJWT(user.ID, user.Email, jwtSecret)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.Greater(t, expiresAt, time.Now().Unix())

	// Validate token
	claims, err := ValidateJWT(token, jwtSecret)
	require.NoError(t, err)
	require.NotNil(t, claims)

	// Verify claims
	assert.Equal(t, user.ID, claims.UserID)
	assert.Equal(t, user.Email, claims.Email)
	assert.Equal(t, "sykell-backend", claims.Issuer)
	assert.Equal(t, user.Email, claims.Subject)
	assert.True(t, claims.ExpiresAt.After(time.Now()))
	assert.True(t, claims.IssuedAt.Before(time.Now().Add(1*time.Minute)))
	assert.True(t, claims.NotBefore.Before(time.Now().Add(1*time.Minute)))
}

func TestPasswordHashingRoundTrip(t *testing.T) {
	// Test hashing and then checking a password
	password := "test-password-for-roundtrip"

	// Hash password
	hash, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hash)

	// Check correct password
	err = CheckPassword(hash, password)
	assert.NoError(t, err)

	// Check incorrect password
	err = CheckPassword(hash, "wrong-password")
	assert.Error(t, err)
}


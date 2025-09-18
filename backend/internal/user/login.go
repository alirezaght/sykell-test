package user

import (
	"context"
	"database/sql"
	"errors"
	"sykell-backend/internal/db"
	"sykell-backend/internal/utils"
)

// Login authenticates a user and returns a JWT token
func (s *UserService) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	// Get user by email
	queries := db.New(s.db)
	user, err := queries.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}

	// Check password
	if err := utils.CheckPassword(user.PasswordHash, req.Password); err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Generate JWT token
	token, expiresAt, err := utils.GenerateJWT(user, []byte(s.config.JWTSecret))
	if err != nil {
		return nil, err
	}

	// Convert to response format with proper date handling
	userResponse := ToUserResponse(user)

	return &LoginResponse{
		Token:     token,
		User:      userResponse,
		ExpiresAt: expiresAt,
	}, nil
}
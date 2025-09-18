package user

import (
	"context"
	"database/sql"
	"errors"
	"sykell-backend/internal/db"
	"sykell-backend/internal/utils"
)

// Register registers a new user
func (s *UserService) Register(ctx context.Context, req RegisterRequest) (sql.Result, error) {
	// Check if user already exists
	queries := db.New(s.db)
	_, err := queries.GetUserByEmail(ctx, req.Email)
	if err == nil {
		return nil, errors.New("user already exists")
	}
	if err != sql.ErrNoRows {
		return nil, err
	}

	// Create the user
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	
	result, err := queries.CreateUser(ctx, db.CreateUserParams{
		Email:        req.Email,
		PasswordHash: hashedPassword,
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}
package user

import (
	"context"
	"database/sql"
	"errors"
	"sykell-backend/internal/utils"
)

// Register registers a new user
func (s *UserService) Register(ctx context.Context, req RegisterRequest) error {
	// Check if user already exists
	_, err := s.repo.GetByEmail(ctx, req.Email)
	if err == nil {
		return errors.New("user already exists")
	}
	if err != sql.ErrNoRows {
		return err
	}

	// Create the user
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}
	
	err = s.repo.Create(ctx, req.Email, hashedPassword)
	if err != nil {
		return err
	}

	return nil
}
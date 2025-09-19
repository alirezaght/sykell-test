package user

import (
	"context"
	"database/sql"
	"errors"
)

// GetProfile returns the profile of the authenticated user
func (s *UserService) GetProfile(ctx context.Context, userID string) (*UserResponse, error) {
	
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	
	user.PasswordHash = "" // Clear password hash before sending response

	return &user, nil
}
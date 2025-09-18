package user

import (
	"context"
	"database/sql"
	"errors"
	"sykell-backend/internal/db"
)

// GetProfile returns the profile of the authenticated user
func (s *UserService) GetProfile(ctx context.Context, userID string) (*UserResponse, error) {
	queries := db.New(s.db)
	user, err := queries.GetUser(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	
	// Convert to response format with proper date handling
	userResponse := ToUserResponse(user)

	return &userResponse, nil
}
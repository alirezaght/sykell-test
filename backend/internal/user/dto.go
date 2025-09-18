package user

import (
	"sykell-backend/internal/db"
	"time"
)

// LoginRequest represents the login request payload
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// UserResponse represents the user data in API responses
type UserResponse struct {
	ID        string     `json:"id"`
	Email     string     `json:"email"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Token     string       `json:"token"`
	User      UserResponse `json:"user"`
	ExpiresAt int64        `json:"expires_at"`
}

// RegisterRequest represents the user registration request
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// ToUserResponse converts a db.User to UserResponse, handling sql.NullTime properly
func ToUserResponse(user db.User) UserResponse {
	resp := UserResponse{
		ID:    user.ID,
		Email: user.Email,
	}
	
	// Handle created_at
	if user.CreatedAt.Valid {
		resp.CreatedAt = &user.CreatedAt.Time
	}
	
	// Handle updated_at
	if user.UpdatedAt.Valid {
		resp.UpdatedAt = &user.UpdatedAt.Time
	}
	
	return resp
}
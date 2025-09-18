package user

import (
	"sykell-backend/internal/db"
)

// LoginRequest represents the login request payload
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Token     string  `json:"token"`
	User      db.User `json:"user"`
	ExpiresAt int64   `json:"expires_at"`
}


// RegisterRequest represents the user registration request
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}
package user

import (
	"database/sql"
	"sykell-backend/internal/config"
)

// UserService provides user-related services
type UserService struct {
	db         *sql.DB
	config    *config.Config
}

// NewUserService creates a new UserService
func NewUserService(database *sql.DB, config *config.Config) *UserService {	
	return &UserService{
		db:         database,
		config: config,
	}
}

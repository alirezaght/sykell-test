package user

import (
	"sykell-backend/internal/config"
)

// UserService provides user-related services
type UserService struct {
	repo      Repo
	config    *config.Config
}

// NewUserService creates a new UserService
func NewUserService(repo Repo, config *config.Config) *UserService {	
	return &UserService{
		repo:	repo,
		config: config,
	}
}

package url

import (
	"sykell-backend/internal/config"
)

// Service provides URL-related services
type Service struct {
	repo	Repo
	config	*config.Config
}

// NewService creates a new URL Service
func NewService(repo Repo, config *config.Config) *Service {
	return &Service{
		repo:	repo,
		config:	config,
	}
}




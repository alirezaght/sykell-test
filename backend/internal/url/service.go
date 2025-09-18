package url

import (
	"database/sql"
	"sykell-backend/internal/config"
)

// Service provides URL-related services
type Service struct {
	db      *sql.DB
	config  *config.Config	
}

// NewService creates a new URL Service
func NewService(database *sql.DB, config *config.Config) *Service {
	return &Service{
		db:      database,
		config:  config,		
	}
}




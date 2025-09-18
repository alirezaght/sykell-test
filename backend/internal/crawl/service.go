package crawl

import (
	"database/sql"
	"sykell-backend/internal/config"
)


// CrawlService provides crawl-related services
type CrawlService struct {
	db *sql.DB
	config *config.Config
}


// NewCrawlService creates a new CrawlService
func NewCrawlService(database *sql.DB, config *config.Config) *CrawlService {
	return &CrawlService{
		db: database,
		config: config,
	}
}
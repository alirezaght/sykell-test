package crawl

import (
	"database/sql"
	"sykell-backend/internal/config"
	"go.temporal.io/sdk/client"
)

// CrawlService provides crawl-related services
type CrawlService struct {
	db *sql.DB
	config *config.Config
	temporalClient client.Client
}


// NewCrawlService creates a new CrawlService
func NewCrawlService(database *sql.DB, config *config.Config, temporalClient client.Client) *CrawlService {
	return &CrawlService{
		db: database,
		config: config,
		temporalClient: temporalClient,
	}
}
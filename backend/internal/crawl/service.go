package crawl

import (
	"sykell-backend/internal/config"
	"sykell-backend/internal/temporal"
)

// CrawlService provides crawl-related services
type CrawlService struct {
	repo Repo
	config *config.Config
	temporalService *temporal.Service
}


// NewCrawlService creates a new CrawlService
func NewCrawlService(repo Repo, config *config.Config, temporalService *temporal.Service) *CrawlService {
	return &CrawlService{
		repo: repo,
		config: config,
		temporalService: temporalService,
	}
}
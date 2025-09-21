package crawl

import (
	"sykell-backend/internal/config"
)

// CrawlService provides crawl-related services
type CrawlService struct {
	repo Repo
	config *config.Config
	orchestrator Orchestrator
}


// NewCrawlService creates a new CrawlService
func NewCrawlService(repo Repo, config *config.Config, orchestrator Orchestrator) *CrawlService {
	return &CrawlService{
		repo: repo,
		config: config,
		orchestrator: orchestrator,
	}
}
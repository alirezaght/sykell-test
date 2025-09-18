package crawl



// CrawlHandler handles HTTP requests related to crawling
type CrawlHandler struct {
	crawlService *CrawlService
}


// NewCrawlHandler creates a new CrawlHandler
func NewCrawlHandler(crawlService *CrawlService) *CrawlHandler {
	return &CrawlHandler{
		crawlService: crawlService,
	}
}
package crawl

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

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

// StartCrawl handles starting a new crawl
func (h *CrawlHandler) StartCrawl(c echo.Context) error {
	userID := c.Get("user_id")
	urlID := c.Param("id")
	if urlID == "" {
		return c.JSON(400, map[string]string{
			"error": "Missing URL ID",
		})
	}

	ctx := c.Request().Context()

	err := h.crawlService.StartCrawl(ctx, userID.(string), urlID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve user profile",
		})
	}


	return c.NoContent(http.StatusOK)
}

// StopCrawl handles stopping an active crawl
func (h *CrawlHandler) StopCrawl(c echo.Context) error {
	userID := c.Get("user_id")
	urlID := c.Param("id")
	if urlID == "" {
		return c.JSON(400, map[string]string{
			"error": "Missing URL ID",
		})
	}

	log.Printf("StopCrawl handler called for user: %v, URL ID: %s", userID, urlID)

	ctx := c.Request().Context()

	err := h.crawlService.StopCrawl(ctx, userID.(string), urlID)
	if err != nil {
		log.Printf("Error in StopCrawl handler: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	
	log.Printf("StopCrawl handler completed successfully")
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Crawl stopped successfully",
	})
}
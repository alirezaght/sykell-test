package crawl

import (
	"net/http"
	"sykell-backend/internal/logger"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
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
			"error": "Failed to star crawl: " + err.Error(),
		})
	}


	return c.NoContent(http.StatusOK)
}

// StartCrawl batch handles starting multiple crawls
func (h *CrawlHandler) StartCrawlBatch(c echo.Context) error {
	userID := c.Get("user_id")
	var req StartCrawlBatchRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}
	if len(req.URLIDs) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "No URL IDs provided",
		})
	}

	ctx := c.Request().Context()
	var failedURLs []string

	for _, urlID := range req.URLIDs {
		err := h.crawlService.StartCrawl(ctx, userID.(string), urlID)
		if err != nil {
			logger.Error("Failed to start crawl for URL", 
				zap.String("url_id", urlID), 
				zap.String("user_id", userID.(string)),
				zap.Error(err))
			failedURLs = append(failedURLs, urlID)
		}
	}

	if len(failedURLs) > 0 {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":       "Failed to start crawl for some URLs",
			"failed_urls": failedURLs,
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

	logger.Info("StopCrawl handler called", 
		zap.String("user_id", userID.(string)), 
		zap.String("url_id", urlID))

	ctx := c.Request().Context()

	err := h.crawlService.StopCrawl(ctx, userID.(string), urlID)
	if err != nil {
		logger.Error("Error in StopCrawl handler", 
			zap.Error(err),
			zap.String("user_id", userID.(string)),
			zap.String("url_id", urlID))
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	
	logger.Info("StopCrawl handler completed successfully")
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Crawl stopped successfully",
	})
}

// StopCrawlBatch handles stopping multiple active crawls
func (h *CrawlHandler) StopCrawlBatch(c echo.Context) error {
	userID := c.Get("user_id")
	var req StartCrawlBatchRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}
	if len(req.URLIDs) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "No URL IDs provided",
		})
	}

	ctx := c.Request().Context()
	var failedURLs []string

	for _, urlID := range req.URLIDs {
		err := h.crawlService.StopCrawl(ctx, userID.(string), urlID)
		if err != nil {
			logger.Error("Failed to stop crawl for URL", 
				zap.String("url_id", urlID), 
				zap.String("user_id", userID.(string)),
				zap.Error(err))
			failedURLs = append(failedURLs, urlID)
		}
	}

	if len(failedURLs) > 0 {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":       "Failed to stop crawl for some URLs",
			"failed_urls": failedURLs,
		})
	}

	return c.NoContent(http.StatusOK)
}

// NotifyCrawlUpdate handles internal notifications to trigger SSE updates
func (h *CrawlHandler) NotifyCrawlUpdate(c echo.Context) error {
		var request struct {
			UserID string `json:"user_id"`
			URLID  string `json:"url_id"`
		}
		
		if err := c.Bind(&request); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid request body",
			})
		}
		
		logger.Debug("Received internal notification", 
			zap.String("user_id", request.UserID), 
			zap.String("url_id", request.URLID))
		logger.Debug("About to call NotifyCrawlUpdate", 
			zap.String("user_id", request.UserID), 
			zap.String("url_id", request.URLID))
		
		// Trigger the SSE notification
		NotifyCrawlUpdate(request.UserID, request.URLID)
		
		return c.JSON(http.StatusOK, map[string]string{
			"message": "Notification sent",
		})
	}
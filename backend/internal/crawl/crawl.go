package crawl

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"sykell-backend/internal/config"
	"sykell-backend/internal/db"
	"sykell-backend/internal/utils"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"go.temporal.io/sdk/activity"
	"golang.org/x/net/html"
)

// CrawlURLActivity performs the actual URL crawling and metadata extraction
func CrawlURLActivity(ctx context.Context, input WorlFlowInput) error {
	// Get the activity logger for proper Temporal logging
	logger := activity.GetLogger(ctx)
	logger.Info("Starting crawl activity", "url", input.URL, "crawl_id", input.CrawlID)

	cfg, err := config.Load()
	if err != nil {
		logger.Error("Failed to load configuration", "error", err)
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Connect to database
	dbSQL, err := sql.Open("mysql", cfg.DatabaseURL)
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer dbSQL.Close()

	// Test database connection
	if err := dbSQL.Ping(); err != nil {
		logger.Error("Failed to ping database", "error", err)
		return fmt.Errorf("failed to ping database: %w", err)
	}
	logger.Info("Database connected successfully")

	repo := NewRepo(dbSQL)
	

	err = repo.SetCrawlRunning(ctx, input.CrawlID)
	if err != nil {
		logger.Error("Failed to set crawl running", "error", err, "crawl_id", input.CrawlID)
		repo.SetCrawlError(ctx, input.CrawlID, fmt.Sprintf("Failed to set crawl running: %v", err))
		// Notify SSE that crawl failed
		NotifyCrawlUpdateHTTP(input.UserID, input.URLID)
		return err
	}

	logger.Info("Crawl status set to running", "crawl_id", input.CrawlID)
	// Notify SSE that crawl started
	NotifyCrawlUpdateHTTP(input.UserID, input.URLID)

	// Create HTTP client with longer timeout and proper context
	client := &http.Client{
		Timeout: 2 * time.Minute, // Increased timeout for slow websites
	}

	logger.Info("Fetching URL", "url", input.URL)
	
	// Create request with activity context for cancellation support
	req, err := http.NewRequestWithContext(ctx, "GET", input.URL, nil)
	if err != nil {
		logger.Error("Failed to create HTTP request", "error", err, "url", input.URL)
		repo.SetCrawlError(ctx, input.CrawlID, fmt.Sprintf("Failed to create HTTP request: %v", err))
		NotifyCrawlUpdateHTTP(input.UserID, input.URLID)
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	
	// Set a reasonable User-Agent to avoid blocking
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; SykellBot/1.0)")
	
	// Fetch the URL
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("Failed to fetch URL", "error", err, "url", input.URL)
		repo.SetCrawlError(ctx, input.CrawlID, fmt.Sprintf("Failed to fetch URL: %v", err))
		NotifyCrawlUpdateHTTP(input.UserID, input.URLID)
		return fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	logger.Info("HTTP response received", "status_code", resp.StatusCode, "url", input.URL)
	activity.RecordHeartbeat(ctx, "HTTP response received")
	
	if resp.StatusCode != http.StatusOK {
		logger.Error("HTTP error response", "status_code", resp.StatusCode, "url", input.URL)
		repo.SetCrawlError(ctx, input.CrawlID, fmt.Sprintf("HTTP error: %d", resp.StatusCode))
		// Notify SSE that crawl failed
		NotifyCrawlUpdateHTTP(input.UserID, input.URLID)
		return fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	logger.Info("Parsing HTML content")
	activity.RecordHeartbeat(ctx, "Parsing HTML")
	// Parse HTML
	doc, err := html.Parse(resp.Body)
	if err != nil {
		logger.Error("Failed to parse HTML", "error", err, "url", input.URL)
		repo.SetCrawlError(ctx, input.CrawlID, fmt.Sprintf("Failed to parse HTML: %v", err))
		// Notify SSE that crawl failed
		NotifyCrawlUpdateHTTP(input.UserID, input.URLID)
		return fmt.Errorf("failed to parse HTML: %w", err)
	}

	
	
	logger.Info("Extracting page metadata")
	activity.RecordHeartbeat(ctx, "Extracting metadata")
	// Extract HTML version
	htmlVersion := utils.ExtractHtmlVersion(doc)
	logger.Info("HTML version extracted", "version", htmlVersion)
	

	// Extract page title
	pageTitle := utils.SanitizeText(utils.ExtractTitle(doc), 500)
	logger.Info("Page title extracted", "title", pageTitle)
	

	// Count headings
	headingCounts := utils.CountHeadings(doc)
	h1Count := int32(headingCounts["h1"])
	h2Count := int32(headingCounts["h2"])
	h3Count := int32(headingCounts["h3"])
	h4Count := int32(headingCounts["h4"])
	h5Count := int32(headingCounts["h5"])
	h6Count := int32(headingCounts["h6"])
	logger.Info("Heading counts extracted", "h1", h1Count, "h2", h2Count, "h3", h3Count, "h4", h4Count, "h5", h5Count, "h6", h6Count)
		

	// Count links
	logger.Info("Analyzing links")
	linkAnalysis := utils.CountLinks(doc, input.URL)
	linkCounts := linkAnalysis.Counts
	internalLinksCount := int32(linkCounts["internal"])
	externalLinksCount := int32(linkCounts["external"])
	inaccessibleLinksCount := int32(linkCounts["inaccessible"])
	logger.Info("Link analysis completed", "internal", internalLinksCount, "external", externalLinksCount, "inaccessible", inaccessibleLinksCount, "total_links", len(linkAnalysis.Links))
	
	for i, url := range linkAnalysis.Links {				
		// Send heartbeat every 25 links to show activity is alive
		if i%25 == 0 {
			activity.RecordHeartbeat(ctx, fmt.Sprintf("Processing link %d/%d", i, len(linkAnalysis.Links)))
		}
		
		// Sanitize anchor text to prevent encoding issues
		sanitizedAnchorText := utils.SanitizeText(url.AnchorText, 1024)
						
		err := repo.CreateInaccessibleLink(ctx, input.CrawlID, url.Href, url.AbsoluteURL, url.IsInternal, *url.StatusCode, sanitizedAnchorText)
		if err != nil {
			logger.Error("Error saving link", "href", url.Href, "anchor", sanitizedAnchorText, "error", err)
			log.Printf("Error saving link (href: %s, anchor: %s): %v", url.Href, sanitizedAnchorText, err)
		}
		
		if i%50 == 0 && i > 0 {
			logger.Info("Processed links", "completed", i, "total", len(linkAnalysis.Links))
		}
	}

	// Check for login form
	hasLoginForm := utils.HasLoginForm(doc)
	logger.Info("Login form analysis completed", "has_login_form", hasLoginForm)
		

	logger.Info("Updating crawl results in database")
	err = repo.UpdateCrawlResult(ctx, input.CrawlID, htmlVersion, pageTitle, h1Count, h2Count, h3Count, h4Count, h5Count, h6Count, internalLinksCount, externalLinksCount, inaccessibleLinksCount, hasLoginForm, string(db.CrawlsStatusDone))

	if err != nil {
		logger.Error("Failed to update crawl result", "error", err, "crawl_id", input.CrawlID)
		repo.SetCrawlError(ctx, input.CrawlID, fmt.Sprintf("Failed to update crawl result: %v", err))
		// Notify SSE that crawl failed
		NotifyCrawlUpdateHTTP(input.UserID, input.URLID)
		return fmt.Errorf("failed to update crawl result: %w", err)
	}

	logger.Info("Crawl completed successfully", "crawl_id", input.CrawlID, "url", input.URL)
	// Notify SSE that crawl completed successfully
	NotifyCrawlUpdateHTTP(input.UserID, input.URLID)

	return nil
}
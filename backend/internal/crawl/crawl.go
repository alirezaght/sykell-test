package crawl

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"sykell-backend/internal/config"
	"sykell-backend/internal/db"
	"sykell-backend/internal/utils"
	"time"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
	"go.temporal.io/sdk/activity"
	"golang.org/x/net/html"
)

// CrawlURLActivity performs the actual URL crawling and metadata extraction, it runs in the Temporal worker process
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
	
	// Start keep-alive goroutine to send heartbeats every 30 seconds
	cancelKeepAlive := keepAlive(ctx, 30*time.Second)
	defer cancelKeepAlive()
	
	// Defer function to handle error cases and set crawl status to error
	defer func() {		
		if r := recover(); r != nil {
			logger.Error("Crawl activity panicked", "panic", r, "crawl_id", input.CrawlID)
			bctx, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
			defer cancel()
			repo.SetCrawlError(bctx, input.CrawlID, fmt.Sprintf("Activity panicked: %v", r))
			NotifyCrawlUpdateHTTP(input.UserID, input.URLID)
		}
	}()
	
	// Track if we successfully complete the crawl
	var crawlCompleted bool
	defer func() {
		if !crawlCompleted {
			logger.Error("Crawl did not complete successfully", "crawl_id", input.CrawlID)			
			bctx, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
			defer cancel()
			repo.SetCrawlError(bctx, input.CrawlID, "Crawl failed to complete (timeout, error, or cancellation)")
			NotifyCrawlUpdateHTTP(input.UserID, input.URLID)
		}
	}()
	

	if err = repo.SetCrawlRunning(ctx, input.CrawlID); err != nil {	
		logger.Error("Failed to set crawl running", "error", err, "crawl_id", input.CrawlID)
		return err
	}

	logger.Info("Crawl status set to running", "crawl_id", input.CrawlID)
	// Notify SSE that crawl started
	NotifyCrawlUpdateHTTP(input.UserID, input.URLID)

	// Create HTTP client with longer timeout and proper context
	client := &http.Client{
		Timeout: 20 * time.Second,
	}

	logger.Info("Fetching URL", "url", input.URL)
	
	// Create request with activity context for cancellation support
	req, err := http.NewRequestWithContext(ctx, "GET", input.URL, nil)
	if err != nil {
		logger.Error("Failed to create HTTP request", "error", err, "url", input.URL)
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	
	// Set a reasonable User-Agent to avoid blocking
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; SykellBot/1.0)")
	
	// Fetch the URL
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("Failed to fetch URL", "error", err, "url", input.URL)
		return fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	logger.Info("HTTP response received", "status_code", resp.StatusCode, "url", input.URL)
	activity.RecordHeartbeat(ctx, "HTTP response received")
	
	if resp.StatusCode != http.StatusOK {
		logger.Error("HTTP error response", "status_code", resp.StatusCode, "url", input.URL)
		return fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	logger.Info("Parsing HTML content")
	activity.RecordHeartbeat(ctx, "Parsing HTML")
	// Parse HTML
	doc, err := html.Parse(resp.Body)
	if err != nil {
		logger.Error("Failed to parse HTML", "error", err, "url", input.URL)
		return fmt.Errorf("failed to parse HTML: %w", err)
	}
	activity.RecordHeartbeat(ctx, "HTML parsing completed")

	
	
	logger.Info("Extracting page metadata")
	activity.RecordHeartbeat(ctx, "Starting metadata extraction")
	
	// Extract HTML version
	activity.RecordHeartbeat(ctx, "About to extract HTML version")
	htmlVersion := utils.ExtractHtmlVersion(doc)
	activity.RecordHeartbeat(ctx, "HTML version extraction completed")
	logger.Info("HTML version extracted", "version", htmlVersion)
	

	// Extract page title
	activity.RecordHeartbeat(ctx, "About to extract page title")
	pageTitle := utils.SanitizeText(utils.ExtractTitle(doc), 500)
	activity.RecordHeartbeat(ctx, "Page title extraction completed")
	logger.Info("Page title extracted", "title", pageTitle)
	

	// Count headings
	activity.RecordHeartbeat(ctx, "About to count headings")
	headingCounts := utils.CountHeadings(doc)
	activity.RecordHeartbeat(ctx, "Heading counting completed")
	h1Count := int32(headingCounts["h1"])
	h2Count := int32(headingCounts["h2"])
	h3Count := int32(headingCounts["h3"])
	h4Count := int32(headingCounts["h4"])
	h5Count := int32(headingCounts["h5"])
	h6Count := int32(headingCounts["h6"])
	logger.Info("Heading counts extracted", "h1", h1Count, "h2", h2Count, "h3", h3Count, "h4", h4Count, "h5", h5Count, "h6", h6Count)
		

	// Count links
	logger.Info("Analyzing links")
	activity.RecordHeartbeat(ctx, "About to start link analysis")
	linkAnalysis := utils.CountLinks(doc, input.URL)
	activity.RecordHeartbeat(ctx, "Link analysis function completed")
	linkCounts := linkAnalysis.Counts
	internalLinksCount := int32(linkCounts["internal"])
	externalLinksCount := int32(linkCounts["external"])
	inaccessibleLinksCount := int32(linkCounts["inaccessible"])
	activity.RecordHeartbeat(ctx, "Link counts processed")
	logger.Info("Link analysis completed", "internal", internalLinksCount, "external", externalLinksCount, "inaccessible", inaccessibleLinksCount, "total_links", len(linkAnalysis.Links))
	
	// Skip individual link saving for now to avoid performance issues
	// TODO: Implement efficient link checking in a separate background process
	logger.Info("Skipping individual link saving to improve performance", "total_links", len(linkAnalysis.Links))

	// Check for login form
	activity.RecordHeartbeat(ctx, "Checking for login form")
	hasLoginForm := utils.HasLoginForm(doc)
	logger.Info("Login form analysis completed", "has_login_form", hasLoginForm)
		

	logger.Info("Updating crawl results in database")
	err = repo.UpdateCrawlResult(ctx, input.CrawlID, htmlVersion, pageTitle, h1Count, h2Count, h3Count, h4Count, h5Count, h6Count, internalLinksCount, externalLinksCount, inaccessibleLinksCount, hasLoginForm, string(db.CrawlsStatusDone))

	if err != nil {
		logger.Error("Failed to update crawl result", "error", err, "crawl_id", input.CrawlID)
		return fmt.Errorf("failed to update crawl result: %w", err)
	}

	// Mark crawl as completed successfully
	crawlCompleted = true
	logger.Info("Crawl completed successfully", "crawl_id", input.CrawlID, "url", input.URL)
	// Notify SSE that crawl completed successfully
	NotifyCrawlUpdateHTTP(input.UserID, input.URLID)

	return nil
}

func keepAlive(ctx context.Context, interval time.Duration) (stop func()) {
	done := make(chan struct{})
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				activity.RecordHeartbeat(ctx, "Crawl still in progress")
			case <-ctx.Done():
				return
			case <-done:
				return
			}
		}
	}()
	return func() {
		close(done)
	}
}
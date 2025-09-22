package crawl

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"sykell-backend/internal/config"
	"sykell-backend/internal/db"
	"sykell-backend/internal/logger"
	"sykell-backend/internal/utils"
	"time"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
	"go.uber.org/zap"
	"golang.org/x/net/html"
)
// Heartbeat is a function type for sending heartbeats from the activity
type Heartbeat func(context.Context, string)


// CrawlURLActivity performs the actual URL crawling and metadata extraction, it runs in the Temporal worker process
func CrawlURLActivity(ctx context.Context, input WorlFlowInput, heartBeat Heartbeat) error {
	// Get the activity logger for proper Temporal logging
	
	logger.Info("Starting crawl activity", zap.String("url", input.URL))

	cfg, err := config.Load()
	if err != nil {
		logger.Error("Failed to load configuration", zap.Error(err))
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Connect to database
	dbSQL, err := sql.Open("mysql", cfg.DatabaseURL)
	if err != nil {
		logger.Error("Failed to connect to database", zap.Error(err))
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer dbSQL.Close()

	// Test database connection
	if err := dbSQL.Ping(); err != nil {
		logger.Error("Failed to ping database", zap.Error(err))
		return fmt.Errorf("failed to ping database: %w", err)
	}
	logger.Info("Database connected successfully")

	repo := NewRepo(dbSQL)
	
	// Start keep-alive goroutine to send heartbeats every 30 seconds
	cancelKeepAlive := keepAlive(ctx, 30*time.Second, heartBeat)
	defer cancelKeepAlive()
	
	// Defer function to handle error cases and set crawl status to error
	defer func() {		
		if r := recover(); r != nil {
			logger.Error("Crawl activity panicked", zap.String("crawl_id", input.CrawlID), zap.Any("panic", r))
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
			logger.Error("Crawl did not complete successfully", zap.String("crawl_id", input.CrawlID))			
			bctx, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
			defer cancel()
			repo.SetCrawlError(bctx, input.CrawlID, "Crawl failed to complete (timeout, error, or cancellation)")
			NotifyCrawlUpdateHTTP(input.UserID, input.URLID)
		}
	}()
	

	if err = repo.SetCrawlRunning(ctx, input.CrawlID); err != nil {	
		logger.Error("Failed to set crawl running", zap.Error(err), zap.String("crawl_id", input.CrawlID))
		return err
	}

	logger.Info("Crawl status set to running", zap.String("crawl_id", input.CrawlID))
	// Notify SSE that crawl started
	NotifyCrawlUpdateHTTP(input.UserID, input.URLID)

	// Create HTTP client with longer timeout and proper context
	client := &http.Client{
		Timeout: 20 * time.Second,
	}

	logger.Info("Fetching URL", zap.String("url", input.URL))
	
	// Create request with activity context for cancellation support
	req, err := http.NewRequestWithContext(ctx, "GET", input.URL, nil)
	if err != nil {
		logger.Error("Failed to create HTTP request", zap.Error(err), zap.String("url", input.URL))
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	
	// Set a reasonable User-Agent to avoid blocking
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; SykellBot/1.0)")
	
	// Fetch the URL
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("Failed to fetch URL", zap.Error(err), zap.String("url", input.URL))
		return fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	logger.Info("HTTP response received", zap.Int("status_code", resp.StatusCode), zap.String("url", input.URL))
	heartBeat(ctx, "HTTP response received")
	
	if resp.StatusCode != http.StatusOK {
		logger.Error("HTTP error response", zap.Int("status_code", resp.StatusCode), zap.String("url", input.URL))
		return fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	logger.Info("Parsing HTML content")
	heartBeat(ctx, "Parsing HTML")
	// Parse HTML
	doc, err := html.Parse(resp.Body)
	if err != nil {
		logger.Error("Failed to parse HTML", zap.Error(err), zap.String("url", input.URL))
		return fmt.Errorf("failed to parse HTML: %w", err)
	}
	heartBeat(ctx, "HTML parsing completed")

	
	
	logger.Info("Extracting page metadata")
	heartBeat(ctx, "Starting metadata extraction")
	
	// Extract HTML version
	heartBeat(ctx, "About to extract HTML version")
	htmlVersion := utils.ExtractHtmlVersion(doc)
	heartBeat(ctx, "HTML version extraction completed")
	logger.Info("HTML version extracted", zap.String("html_version", htmlVersion))
	

	// Extract page title
	heartBeat(ctx, "About to extract page title")
	pageTitle := utils.SanitizeText(utils.ExtractTitle(doc), 500)
	heartBeat(ctx, "Page title extraction completed")
	logger.Info("Page title extracted", zap.String("page_title", pageTitle))
	

	// Count headings
	heartBeat(ctx, "About to count headings")	
	headingCounts := utils.CountHeadings(doc)
	heartBeat(ctx, "Heading counting completed")
	h1Count := int32(headingCounts["h1"])
	h2Count := int32(headingCounts["h2"])
	h3Count := int32(headingCounts["h3"])
	h4Count := int32(headingCounts["h4"])
	h5Count := int32(headingCounts["h5"])
	h6Count := int32(headingCounts["h6"])
	logger.Info("Heading counts extracted", zap.Int32("h1", h1Count), zap.Int32("h2", h2Count), zap.Int32("h3", h3Count), zap.Int32("h4", h4Count), zap.Int32("h5", h5Count), zap.Int32("h6", h6Count))
		

	// Count links
	logger.Info("Analyzing links")
	heartBeat(ctx, "About to start link analysis")
	linkAnalysis := utils.CountLinks(ctx, doc, input.URL)
	heartBeat(ctx, "Link analysis function completed")
	linkCounts := linkAnalysis.Counts
	internalLinksCount := int32(linkCounts["internal"])
	externalLinksCount := int32(linkCounts["external"])
	inaccessibleLinksCount := int32(linkCounts["inaccessible"])
	heartBeat(ctx, "Link counts processed")
	logger.Info("Link analysis completed", zap.Int32("internal_links", internalLinksCount), zap.Int32("external_links", externalLinksCount), zap.Int32("inaccessible_links", inaccessibleLinksCount))
	
	// Skip individual link saving for now to avoid performance issues
	// TODO: Implement efficient link checking in a separate background process
	logger.Info("Skipping individual link saving to improve performance", zap.String("crawl_id", input.CrawlID))

	// Check for login form
	heartBeat(ctx, "Checking for login form")
	hasLoginForm := utils.HasLoginForm(doc)
	logger.Info("Login form analysis completed", zap.Bool("has_login_form", hasLoginForm))
		

	logger.Info("Updating crawl results in database")
	err = repo.UpdateCrawlResult(ctx, input.CrawlID, htmlVersion, pageTitle, h1Count, h2Count, h3Count, h4Count, h5Count, h6Count, internalLinksCount, externalLinksCount, inaccessibleLinksCount, hasLoginForm, string(db.CrawlsStatusDone))

	if err != nil {
		logger.Error("Failed to update crawl result", zap.Error(err), zap.String("crawl_id", input.CrawlID))
		return fmt.Errorf("failed to update crawl result: %w", err)
	}

	// Mark crawl as completed successfully
	crawlCompleted = true
	logger.Info("Crawl completed successfully", zap.String("crawl_id", input.CrawlID))
	// Notify SSE that crawl completed successfully
	NotifyCrawlUpdateHTTP(input.UserID, input.URLID)

	return nil
}

func keepAlive(ctx context.Context, interval time.Duration, heartBeat Heartbeat) (stop func()) {
	done := make(chan struct{})
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				heartBeat(ctx, "Crawl still in progress")
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
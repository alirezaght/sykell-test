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
	"golang.org/x/net/html"
)

// CrawlURLActivity performs the actual URL crawling and metadata extraction
func CrawlURLActivity(ctx context.Context, input WorlFlowInput) error {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Connect to database
	dbSQL, err := sql.Open("mysql", cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer dbSQL.Close()

	// Test database connection
	if err := dbSQL.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	log.Println("Database connected successfully")


	queries := db.New(dbSQL)

	err = queries.SetCrawlRunning(ctx, input.CrawlID)
	if err != nil {
		queries.SetCrawlError(ctx, db.SetCrawlErrorParams{
			ErrorMessage: sql.NullString{String: err.Error(), Valid: true},
			ID: input.CrawlID,
		})
		// Notify SSE that crawl failed
		NotifyCrawlUpdateHTTP(input.UserID, input.URLID)
		return err
	}

	// Notify SSE that crawl started
	NotifyCrawlUpdateHTTP(input.UserID, input.URLID)

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Fetch the URL
	resp, err := client.Get(input.URL)
	if err != nil {
		queries.SetCrawlError(ctx, db.SetCrawlErrorParams{
			ErrorMessage: sql.NullString{String: err.Error(), Valid: true},
			ID: input.CrawlID,
		})
		return fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		queries.SetCrawlError(ctx, db.SetCrawlErrorParams{
			ErrorMessage: sql.NullString{String: fmt.Sprintf("HTTP error: %d", resp.StatusCode), Valid: true},
			ID: input.CrawlID,
		})
		// Notify SSE that crawl failed
		NotifyCrawlUpdateHTTP(input.UserID, input.URLID)
		return fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	// Parse HTML
	doc, err := html.Parse(resp.Body)
	if err != nil {
		queries.SetCrawlError(ctx, db.SetCrawlErrorParams{
			ErrorMessage: sql.NullString{String: err.Error(), Valid: true},
			ID: input.CrawlID,
		})
		// Notify SSE that crawl failed
		NotifyCrawlUpdateHTTP(input.UserID, input.URLID)
		return fmt.Errorf("failed to parse HTML: %w", err)
	}

	
	
	// Extract HTML version
	htmlVersion := utils.ExtractHtmlVersion(doc)
	

	// Extract page title
	pageTitle := utils.SanitizeText(utils.ExtractTitle(doc), 500)
	

	// Count headings
	headingCounts := utils.CountHeadings(doc)
	h1Count := int32(headingCounts["h1"])
	h2Count := int32(headingCounts["h2"])
	h3Count := int32(headingCounts["h3"])
	h4Count := int32(headingCounts["h4"])
	h5Count := int32(headingCounts["h5"])
	h6Count := int32(headingCounts["h6"])
		

	// Count links
	linkAnalysis := utils.CountLinks(doc, input.URL)
	linkCounts := linkAnalysis.Counts
	internalLinksCount := int32(linkCounts["internal"])
	externalLinksCount := int32(linkCounts["external"])
	inaccessibleLinksCount := int32(linkCounts["inaccessible"])
	
	for _, url := range linkAnalysis.Links {
		statusCode := sql.NullInt32{}
		if url.StatusCode != nil {
			statusCode = sql.NullInt32{Int32: int32(*url.StatusCode), Valid: true}
		}
		
		// Sanitize anchor text to prevent encoding issues
		sanitizedAnchorText := utils.SanitizeText(url.AnchorText, 1024)
		
		_, err := queries.CreateInaccessibleLink(ctx, db.CreateInaccessibleLinkParams{
			CrawlID:     input.CrawlID,
			Href:        url.Href,
			AbsoluteUrl: url.AbsoluteURL,			
			IsInternal:  url.IsInternal,
			StatusCode:  statusCode,
			AnchorText:  sql.NullString{String: sanitizedAnchorText, Valid: sanitizedAnchorText != ""},
		})
		if err != nil {
			log.Printf("Error saving link (href: %s, anchor: %s): %v", url.Href, sanitizedAnchorText, err)
		}
	}

	// Check for login form
	hasLoginForm := utils.HasLoginForm(doc)
		

	err = queries.UpdateCrawlResult(ctx, db.UpdateCrawlResultParams{
		HtmlVersion:            sql.NullString{String: htmlVersion, Valid: true},
		PageTitle:              sql.NullString{String: pageTitle, Valid: true},
		H1Count:                sql.NullInt32{Int32: h1Count, Valid: true},
		H2Count:                sql.NullInt32{Int32: h2Count, Valid: true},
		H3Count:                sql.NullInt32{Int32: h3Count, Valid: true},
		H4Count:                sql.NullInt32{Int32: h4Count, Valid: true},
		H5Count:                sql.NullInt32{Int32: h5Count, Valid: true},
		H6Count:                sql.NullInt32{Int32: h6Count, Valid: true},
		InternalLinksCount:     sql.NullInt32{Int32: internalLinksCount, Valid: true},
		ExternalLinksCount:     sql.NullInt32{Int32: externalLinksCount, Valid: true},
		InaccessibleLinksCount: sql.NullInt32{Int32: inaccessibleLinksCount, Valid: true},
		HasLoginForm:           hasLoginForm,
		Status:                 db.CrawlsStatusDone,
		ID:                     input.CrawlID,
	})

	if err != nil {
		queries.SetCrawlError(ctx, db.SetCrawlErrorParams{
			ErrorMessage: sql.NullString{String: err.Error(), Valid: true},
			ID: input.CrawlID,
		})
		// Notify SSE that crawl failed
		NotifyCrawlUpdateHTTP(input.UserID, input.URLID)
		return fmt.Errorf("failed to update crawl result: %w", err)
	}

	// Notify SSE that crawl completed successfully
	NotifyCrawlUpdateHTTP(input.UserID, input.URLID)

	return nil
}
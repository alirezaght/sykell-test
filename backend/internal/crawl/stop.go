package crawl

import (
	"context"
	"fmt"
	"log"
	"sykell-backend/internal/db"
)

// StopCrawl stops an active crawl for the specified URL by the user
func (s *CrawlService) StopCrawl(ctx context.Context, userID string, urlID string) error {
	log.Printf("StopCrawl called for userID: %s, urlID: %s", userID, urlID)
	
	// Check if Temporal client is available
	if s.temporalClient == nil {
		log.Printf("Temporal client is nil")
		return fmt.Errorf("crawling functionality is unavailable: Temporal client not connected")
	}
	
	queries := db.New(s.db)
	
	// Verify that the URL belongs to the user
	url, err := queries.GetUrlByIdAndUserId(ctx, db.GetUrlByIdAndUserIdParams{
		ID:   urlID,
		UserID: userID,
	})
	if err != nil {
		log.Printf("Error getting URL by ID and user ID: %v", err)
		return err
	}
	log.Printf("Found URL: %s", url.NormalizedUrl)
	
	activeCrawls, err := queries.GetActiveCrawlsForUrlId(ctx, url.ID)

	if err != nil {
		log.Printf("Error getting active crawls: %v", err)
		return err
	}
	
	log.Printf("Found %d active crawls", len(activeCrawls))

	for _, crawl := range activeCrawls {
		log.Printf("Stopping crawl ID: %s, workflow ID: %s", crawl.ID, crawl.WorkflowID)

		err = queries.SetCrawlStopped(ctx, crawl.ID)
		if err != nil {
			log.Printf("Error updating crawl status: %v", err)
			return fmt.Errorf("failed to update crawl status: %w", err)
		}
		log.Printf("Successfully updated crawl status to stopped for crawl ID: %s", crawl.ID)

		// Signal the workflow to stop	
		err = s.temporalClient.CancelWorkflow(ctx, crawl.WorkflowID, "")
		if err != nil {
			log.Printf("Error canceling workflow: %v", err)
			return fmt.Errorf("failed to cancel workflow: %w", err)
		}
		log.Printf("Successfully canceled workflow: %s", crawl.WorkflowID)
		
	}
	
	log.Printf("StopCrawl completed successfully")
	return nil
}
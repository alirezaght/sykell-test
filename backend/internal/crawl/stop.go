package crawl

import (
	"context"
	"fmt"
	"log"
)

// StopCrawl stops an active crawl for the specified URL by the user
func (s *CrawlService) StopCrawl(ctx context.Context, userID string, urlID string) error {
	log.Printf("StopCrawl called for userID: %s, urlID: %s", userID, urlID)
				
	// Verify that the URL belongs to the user
	url, err := s.repo.GetUrlByIdAndUserId(ctx, urlID, userID)
	if err != nil {
		log.Printf("Error getting URL by ID and user ID: %v", err)
		return err
	}
	log.Printf("Found URL: %s", url.NormalizedUrl)
	
	activeCrawls, err := s.repo.GetActiveCrawlsForUrlId(ctx, url.ID)

	if err != nil {
		log.Printf("Error getting active crawls: %v", err)
		return err
	}
	
	log.Printf("Found %d active crawls", len(activeCrawls))

	for _, crawl := range activeCrawls {
		log.Printf("Stopping crawl ID: %s, workflow ID: %s", crawl.ID, crawl.WorkflowID)

		err = s.repo.SetCrawlStopped(ctx, crawl.ID)
		if err != nil {
			log.Printf("Error updating crawl status: %v", err)
			return fmt.Errorf("failed to update crawl status: %w", err)
		}
		log.Printf("Successfully updated crawl status to stopped for crawl ID: %s", crawl.ID)

		// Signal the workflow to stop		
		err = s.temporalService.GetTemporalClient().CancelWorkflow(ctx, crawl.WorkflowID, "")
		if err != nil {
			log.Printf("Error canceling workflow: %v", err)
			return fmt.Errorf("failed to cancel workflow: %w", err)
		}
		log.Printf("Successfully canceled workflow: %s", crawl.WorkflowID)
		
		// Notify SSE that crawl was stopped
		NotifyCrawlUpdateHTTP(userID, urlID)
		
	}
	
	log.Printf("StopCrawl completed successfully")
	return nil
}
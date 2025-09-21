package crawl

import (
	"context"	
	"github.com/google/uuid"
)

// StartCrawl initiates a crawl for the specified URL by the user
func (s *CrawlService) StartCrawl(ctx context.Context, userID string, urlID string) error {	
		
	// Verify that the URL belongs to the user
	url, err := s.repo.GetUrlByIdAndUserId(ctx, urlID, userID)
	if err != nil {
		return err
	}
	
	// Check if there are active crawls for the user
	activeCrawls, err := s.repo.CountOfActiveCrawlForUrlId(ctx, url.ID)
	if err != nil {
		return err
	}
	if activeCrawls > 0 {
		return nil
	}
	// Enqueue the crawl task
	workflowID := "crawl_" + url.ID + "_" + uuid.New().String()
	if err = s.repo.QueueCrawl(ctx, urlID, workflowID); err != nil {
		return err
	}	

	crawlID, err := s.repo.GetCrawlIDByWorkflowID(ctx, workflowID)
	if err != nil {
		return err
	}
		
	// Start the workflow	
	err = s.orchestrator.Start(ctx, WorlFlowInput{
		URLID: url.ID,
		UserID: userID,
		WorkflowID: workflowID,
		URL: url.NormalizedUrl,
		CrawlID: crawlID,
	})
	
	

	return err
}
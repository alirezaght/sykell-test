package crawl

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
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
	err = s.repo.QueueCrawl(ctx, urlID, workflowID)

	if err != nil {
		return err
	}

	crawlID, err := s.repo.GetCrawlIDByWorkflowID(ctx, workflowID)
	if err != nil {
		return err
	}

	workflowOptions := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: TaskQueueName,
		WorkflowExecutionTimeout: 10 * time.Minute, // Set explicit workflow timeout
		WorkflowTaskTimeout:      time.Minute,      // Set workflow task timeout		
		StartDelay: 3 * time.Second,
	}
	
	// Start the workflow
	
	_, err = s.temporalService.GetTemporalClient().ExecuteWorkflow(ctx, workflowOptions, WorkflowName, WorlFlowInput{
		URLID: url.ID,
		UserID: userID,
		WorkflowID: workflowID,
		URL: url.NormalizedUrl,
		CrawlID: crawlID,
	})
	
	

	return err
}
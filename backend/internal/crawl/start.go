package crawl

import (
	"context"
	"sykell-backend/internal/db"

	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
)

// StartCrawl initiates a crawl for the specified URL by the user
func (s *CrawlService) StartCrawl(ctx context.Context, userID string, urlID string) error {	
	
	queries := db.New(s.db)
	
	// Verify that the URL belongs to the user
	url, err := queries.GetUrlByIdAndUserId(ctx, db.GetUrlByIdAndUserIdParams{
		ID:   urlID,
		UserID: userID,
	})
	if err != nil {
		return err
	}
	
	// Check if there are active crawls for the user
	activeCrawls, err := queries.CountOfActiveCrawlForUrlId(ctx, url.ID)
	if err != nil {
		return err
	}
	if activeCrawls > 0 {
		return nil
	}
	// Enqueue the crawl task
	workflowID := "crawl_" + url.ID + "_" + uuid.New().String()
	_, err = queries.QueueCrawl(ctx, db.QueueCrawlParams{
		UrlID: url.ID,
		WorkflowID: workflowID,
	})

	if err != nil {
		return err
	}

	crawl, err := queries.GetCrawlByWorkflowID(ctx, workflowID)
	if err != nil {
		return err
	}

	workflowOptions := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: TaskQueueName,
	}
	
	// Start the workflow
	
	_, err = s.temporalClient.ExecuteWorkflow(ctx, workflowOptions, WorkflowName, WorlFlowInput{
		URLID: url.ID,
		UserID: userID,
		WorkflowID: workflowID,
		URL: url.NormalizedUrl,
		CrawlID: crawl.ID,
	})
	
	

	return err
}
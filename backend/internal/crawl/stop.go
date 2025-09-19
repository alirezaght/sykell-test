package crawl

import (
	"context"
	"fmt"
	"sykell-backend/internal/logger"

	"go.uber.org/zap"
)

// StopCrawl stops an active crawl for the specified URL by the user
func (s *CrawlService) StopCrawl(ctx context.Context, userID string, urlID string) error {
	logger.Info("StopCrawl called", 
		zap.String("user_id", userID), 
		zap.String("url_id", urlID))
				
	// Verify that the URL belongs to the user
	url, err := s.repo.GetUrlByIdAndUserId(ctx, urlID, userID)
	if err != nil {
		logger.Error("Error getting URL by ID and user ID", zap.Error(err))
		return err
	}
	logger.Info("Found URL", zap.String("normalized_url", url.NormalizedUrl))
	
	activeCrawls, err := s.repo.GetActiveCrawlsForUrlId(ctx, url.ID)

	if err != nil {
		logger.Error("Error getting active crawls", zap.Error(err))
		return err
	}
	
	logger.Info("Found active crawls", zap.Int("count", len(activeCrawls)))

	for _, crawl := range activeCrawls {
		logger.Info("Stopping crawl", 
			zap.String("crawl_id", crawl.ID), 
			zap.String("workflow_id", crawl.WorkflowID))

		if err = s.repo.SetCrawlStopped(ctx, crawl.ID); err != nil {
			logger.Error("Error updating crawl status", zap.Error(err))
			return fmt.Errorf("failed to update crawl status: %w", err)
		}
		
		logger.Info("Successfully updated crawl status to stopped", zap.String("crawl_id", crawl.ID))

		// Signal the workflow to stop		
		if err = s.temporalService.GetTemporalClient().CancelWorkflow(ctx, crawl.WorkflowID, ""); err != nil {		
			logger.Error("Error canceling workflow", zap.Error(err))
			return fmt.Errorf("failed to cancel workflow: %w", err)
		}
		logger.Info("Successfully canceled workflow", zap.String("workflow_id", crawl.WorkflowID))
		
		// Notify SSE that crawl was stopped
		NotifyCrawlUpdateHTTP(userID, urlID)
		
	}
	
	logger.Info("StopCrawl completed successfully")
	return nil
}
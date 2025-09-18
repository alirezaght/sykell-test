package crawl

import "time"

// Constants for crawl workflow configuration
const (
	TaskQueueName = "crawl-task-queue"
	WorkflowName  = "CrawlWorkflow"
)

// WorlFlowInput represents the input parameters for the crawl workflow
type WorlFlowInput struct {
	URLID      string `json:"url_id"`
	UserID     string `json:"user_id"`
	WorkflowID string `json:"workflow_id"`
	URL        string `json:"url,omitempty"`
	CrawlID    string `json:"crawl_id,omitempty"`
}


// SSENotification represents a simple notification to invalidate queries
type SSENotification struct {
	Type      string    `json:"type"`      // "crawl_update"
	URLID     string    `json:"url_id"`    // URL ID that needs to be refetched
	UserID    string    `json:"user_id"`   // User ID (for verification)
	Timestamp time.Time `json:"timestamp"`
}


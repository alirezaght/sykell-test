package crawl

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

package crawl

import "context"

// Orchestrator defines the interface for starting and stopping crawl workflows
type Orchestrator interface {
		Start(ctx context.Context, input WorlFlowInput) error
		Stop(ctx context.Context, workflowID string) error		
}

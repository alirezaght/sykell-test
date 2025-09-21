package temporal

import (
	"context"
	"sykell-backend/internal/config"
	"sykell-backend/internal/crawl"
	"time"

	"go.temporal.io/sdk/client"
)

// Service provides Temporal-related services
type Service struct {
	config *config.Config
	temporalClient client.Client
}

// NewService creates a new Temporal Service
func NewService(config *config.Config) *Service {
	return &Service{
		config: config,
		temporalClient: nil,
	}
}


// Setup initializes the Temporal client
func (s *Service) Setup() {
	s.temporalClient, _ = client.Dial(client.Options{
		HostPort:  s.config.TemporalHostPort,
		Namespace: s.config.Namespace,
		ConnectionOptions: client.ConnectionOptions{
			TLS: nil, // Disable TLS for local development
			KeepAliveTime:  10 * time.Second, // 10 seconds
			KeepAliveTimeout: 20 * time.Second, // 20 seconds			
		},		
	})	
}

// GetTemporalClient returns the Temporal client, initializing it if necessary
func (s *Service) GetTemporalClient() client.Client {
	if s.temporalClient == nil {
		s.Setup()
	}
	return s.temporalClient
}

// Close closes the Temporal client connection
func (s *Service) Close() {
	if s.temporalClient != nil {
		s.temporalClient.Close()
	}
}

type temporalOrchestrator struct {
	service *Service
}

// NewOrchestrator creates a new Orchestrator
func NewOrchestrator(service Service) *temporalOrchestrator {
	return &temporalOrchestrator{
		service: &service,
	}
}

// Constants for Temporal
const (
	TaskQueueName = "crawl-task-queue"
	WorkflowName  = "CrawlWorkflow"
)

func (o *temporalOrchestrator) Start(ctx context.Context, input crawl.WorlFlowInput) error {
	workflowOptions := client.StartWorkflowOptions{
		ID:        input.WorkflowID,
		TaskQueue: TaskQueueName,
		WorkflowExecutionTimeout: 10 * time.Minute, // Set explicit workflow timeout
		WorkflowTaskTimeout:      time.Minute,      // Set workflow task timeout		
		StartDelay: 3 * time.Second, // Small delay to ensure the sse connection is ready
	}
	_, err := o.service.GetTemporalClient().ExecuteWorkflow(ctx, workflowOptions, WorkflowName, input)
	return err
}

func (o *temporalOrchestrator) Stop(ctx context.Context, workflowID string) error {
	return o.service.GetTemporalClient().CancelWorkflow(ctx, workflowID, "")
}
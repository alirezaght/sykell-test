package crawl

import (
	"sykell-backend/internal/config"
	"sykell-backend/internal/logger"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
	"go.uber.org/zap"
)

// CrawlWorkflow is the main workflow for crawling a URL, it orchestrates the crawl activity
func CrawlWorkflow(ctx workflow.Context, input WorlFlowInput) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting crawl workflow", "url", input.URL, "crawl_id", input.CrawlID, "user_id", input.UserID)

	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Minute,
		ScheduleToCloseTimeout: 15 * time.Minute, 
		ScheduleToStartTimeout: 30 * time.Minute,
		HeartbeatTimeout: 1 * time.Minute, 
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute,
			MaximumAttempts:    3,
		},
	}
	
	ctx = workflow.WithActivityOptions(ctx, activityOptions)
	
	logger.Info("Executing crawl activity", "url", input.URL, "crawl_id", input.CrawlID)
	err := workflow.ExecuteActivity(ctx, CrawlURLActivity, input).Get(ctx, nil)
	if err != nil {
		logger.Error("Crawl workflow failed", "error", err, "url", input.URL, "crawl_id", input.CrawlID)
		return err
	}

	logger.Info("Crawl workflow completed successfully", "url", input.URL, "crawl_id", input.CrawlID)
	return nil
}

// StartWorker initializes and starts the Temporal worker to process crawl workflows and activities
func StartWorker(config *config.Config) error {	
	logger.Info("Attempting to connect to Temporal server", zap.String("host_port", config.TemporalHostPort))
	
	// Create Temporal client with better connection settings
	clientOptions := client.Options{
		HostPort:       config.TemporalHostPort,
		Namespace:      config.Namespace,
		ConnectionOptions: client.ConnectionOptions{
			TLS: nil, // Disable TLS for local development
			KeepAliveTime:  10 * time.Second, // seconds
			KeepAliveTimeout: 20 * time.Second, // seconds			
		},
	}
	
	logger.Info("Connecting to Temporal with options", 
		zap.String("host_port", clientOptions.HostPort),
		zap.String("namespace", clientOptions.Namespace))
	
	temporalClient, err := client.Dial(clientOptions)
	if err != nil {
		logger.Error("Failed to create Temporal client", zap.Error(err))
		logger.Info("Make sure Temporal server is running. Start it with: docker compose up -d")
		return err
	}
	
	// Ensure proper cleanup
	defer func() {
		logger.Info("Closing Temporal client connection")
		temporalClient.Close()
	}()

	logger.Info("Successfully connected to Temporal server")

	// Create worker with debug-enabled options
	w := worker.New(temporalClient, TaskQueueName, worker.Options{
		EnableLoggingInReplay: true, // This ensures logs are visible during replay		
		MaxConcurrentActivityExecutionSize: 2,
		
	})

	// Register workflows
	w.RegisterWorkflow(CrawlWorkflow)

	// Register activities
	w.RegisterActivity(CrawlURLActivity)
	
	logger.Info("Starting Temporal worker on task queue", zap.String("task_queue", TaskQueueName))
	
	// Start listening for tasks
	return w.Run(worker.InterruptCh())
}
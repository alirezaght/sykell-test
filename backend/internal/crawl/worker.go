package crawl

import (
	"log"
	"sykell-backend/internal/config"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

func CrawlWorkflow(ctx workflow.Context, input WorlFlowInput) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting crawl workflow", "url", input.URL, "crawl_id", input.CrawlID, "user_id", input.UserID)

	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Minute, // Increased timeout for slow websites
		ScheduleToCloseTimeout: 15 * time.Minute, // Overall timeout including retries
		HeartbeatTimeout: 30 * time.Second, // Add heartbeat for long-running activities
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

func StartWorker(config *config.Config) error {	
	log.Printf("Attempting to connect to Temporal server at %s", config.TemporalHostPort)
	
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
	
	log.Printf("Connecting to Temporal with options: HostPort=%s, Namespace=%s", 
		clientOptions.HostPort, clientOptions.Namespace)
	
	temporalClient, err := client.Dial(clientOptions)
	if err != nil {
		log.Printf("Failed to create Temporal client: %v", err)
		log.Printf("Make sure Temporal server is running. Start it with: docker compose up -d")
		return err
	}
	
	// Ensure proper cleanup
	defer func() {
		log.Printf("Closing Temporal client connection")
		temporalClient.Close()
	}()

	log.Printf("Successfully connected to Temporal server")

	// Create worker
	w := worker.New(temporalClient, TaskQueueName, worker.Options{})

	// Register workflows
	w.RegisterWorkflow(CrawlWorkflow)

	// Register activities
	w.RegisterActivity(CrawlURLActivity)
	
	log.Printf("Starting Temporal worker on task queue: %s", TaskQueueName)
	
	// Start listening for tasks
	return w.Run(worker.InterruptCh())
}
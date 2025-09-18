package main

import (
	"log"
	"os"
	"sykell-backend/internal/config"
	"sykell-backend/internal/crawl"
)

func main() {
	// Get configuration from environment variables
	temporalHostPort := os.Getenv("TEMPORAL_HOST_PORT")
	if temporalHostPort == "" {
		temporalHostPort = "localhost:7233"
	}

	temporalNamespace := os.Getenv("TEMPORAL_NAMESPACE")
	if temporalNamespace == "" {
		temporalNamespace = "default"
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	log.Printf("Starting Temporal worker with config: %+v", cfg)
	
	// Start the worker
	if err := crawl.StartWorker(cfg); err != nil {
		log.Fatalf("Failed to start worker: %v", err)
	}
}
package main

import (
	"sykell-backend/internal/config"
	"sykell-backend/internal/crawl"
	"sykell-backend/internal/logger"

	"go.uber.org/zap"
)

func main() {
	// Load configuration first
	cfg, err := config.Load()
	if err != nil {
		panic("Failed to load configuration: " + err.Error())
	}

	// Initialize logger
	if err := logger.InitLogger(cfg.LogLevel, cfg.LogFormat, cfg.Environment); err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	defer logger.Sync()

	logger.Info("Starting Temporal worker")
	
	// Start the worker
	if err := crawl.StartWorker(cfg); err != nil {
		logger.Fatal("Failed to start worker", zap.Error(err))
	}
}
package main

import (
	"database/sql"
	"log"
	"net/http"

	"sykell-backend/internal/config"
	"sykell-backend/internal/crawl"
	sykellMiddleware "sykell-backend/internal/middleware"
	"sykell-backend/internal/url"
	"sykell-backend/internal/user"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.temporal.io/sdk/client"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Connect to database
	db, err := sql.Open("mysql", cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	log.Println("Database connected successfully")

	// Initialize services
	userService := user.NewUserService(db, cfg)
	userHandler := user.NewUserHandler(userService)

	// Initialize Temporal client with better connection settings
	log.Printf("Connecting to Temporal server at %s", cfg.TemporalHostPort)
	temporalClient, err := client.Dial(client.Options{
		HostPort:  cfg.TemporalHostPort,
		Namespace: cfg.Namespace,
		ConnectionOptions: client.ConnectionOptions{
			TLS: nil, // Disable TLS for local development
		},
	})
	if err != nil {
		log.Printf("Failed to create Temporal client: %v", err)
		log.Printf("Make sure Temporal server is running with: docker compose up -d")
		log.Printf("Temporal connection required for crawling functionality")
		// Don't fatal here - allow server to start without Temporal for other endpoints
		temporalClient = nil
	} else {
		log.Printf("Successfully connected to Temporal server")
	}
	
	// Ensure proper cleanup on shutdown
	if temporalClient != nil {
		defer temporalClient.Close()
	}

	// Initialize Temporal service
	urlService := url.NewService(db, cfg)
	urlHandler := url.NewHandler(urlService)

	var crawlService *crawl.CrawlService
	var crawlHandler *crawl.CrawlHandler
	
	if temporalClient != nil {
		crawlService = crawl.NewCrawlService(db, cfg, temporalClient)
		crawlHandler = crawl.NewCrawlHandler(crawlService)
	} else {
		log.Printf("Warning: Crawling functionality disabled due to Temporal connection failure")
	}

	// Create Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE, echo.OPTIONS},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	// Routes
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "Sykell Backend API",
			"version": "1.0.0",
		})
	})

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "healthy",
		})
	})

	// API routes
	api := e.Group("/api/v1")
	
	// Auth routes (public)
	api.POST("/auth/register", userHandler.Register)
	api.POST("/auth/login", userHandler.Login)	
		
	
	// Protected routes (require JWT)
	protected := api.Group("", sykellMiddleware.JWTMiddleware([]byte(cfg.JWTSecret)))
	// Profile route
	protected.GET("/auth/me", userHandler.GetProfile)
	
	// Url routes
	protected.GET("/urls", urlHandler.ListURLs)
	protected.POST("/urls", urlHandler.AddURL)
	protected.DELETE("/urls/:id", urlHandler.RemoveURL)
	
	// Crawl routes (only if Temporal is available)
	if crawlHandler != nil {
		protected.POST("/crawl/start/:id", crawlHandler.StartCrawl)
		protected.POST("/crawl/stop/:id", crawlHandler.StopCrawl)
	}

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	e.Logger.Fatal(e.Start(":" + cfg.Port))
}
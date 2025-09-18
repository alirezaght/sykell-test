package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

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
			KeepAliveTime:   10 * time.Second, // seconds
			KeepAliveTimeout: 20 * time.Second, // seconds						
		},
	})
	if err != nil {
		log.Printf("Failed to create Temporal client: %v", err)
		log.Printf("Make sure Temporal server is running with: docker compose up -d")
		log.Printf("Temporal connection required for crawling functionality")		
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

	crawlService := crawl.NewCrawlService(db, cfg, temporalClient)
	crawlHandler := crawl.NewCrawlHandler(crawlService)
	
	

	// Create Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	
	// Add custom middleware to log all requests for debugging
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			log.Printf("Request: %s %s from %s", c.Request().Method, c.Request().URL.Path, c.RealIP())
			return next(c)
		}
	})
	
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:5173"}, // Add specific origins for cookie support
		AllowMethods:     []string{echo.GET, echo.PUT, echo.POST, echo.DELETE, echo.OPTIONS},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true, // Enable credentials (cookies) support
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
	
	// Test routes for debugging
	api.GET("/test", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "API is working",
		})
	})

	// Auth routes (public)
	api.POST("/auth/register", userHandler.Register)
	api.POST("/auth/login", userHandler.Login)
	api.POST("/auth/logout", userHandler.Logout)	
		
	
	// Protected routes (require JWT)
	protected := api.Group("", sykellMiddleware.JWTMiddleware([]byte(cfg.JWTSecret)))
	// Profile route
	protected.GET("/auth/me", userHandler.GetProfile)
	
	// Url routes
	protected.GET("/urls", urlHandler.ListURLs)
	protected.POST("/urls", urlHandler.AddURL)
	protected.DELETE("/urls/:id", urlHandler.RemoveURL)
	
	// Crawl routes (only if Temporal is available)
	
	log.Printf("Registering crawl routes...")
	protected.POST("/crawl/start/:id", crawlHandler.StartCrawl)
	protected.POST("/crawl/stop/:id", crawlHandler.StopCrawl)
	protected.GET("/crawl/stream", crawlHandler.StreamCrawlUpdates)
		
	// Internal notification endpoint for Temporal worker to trigger SSE notifications
	api.POST("/internal/notify-crawl-update", crawlHandler.NotifyCrawlUpdate)
	
	log.Printf("Crawl routes registered successfully")
	

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	e.Logger.Fatal(e.Start(":" + cfg.Port))
}
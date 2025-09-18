package main

import (
	"log"
	"net/http"

	"sykell-backend/internal/config"
	"sykell-backend/internal/user"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Create Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

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
	
	// User routes
	userHandler := user.NewUserHandler()
	api.GET("/users", userHandler.ListUsers)
	api.GET("/users/:id", userHandler.GetUser)
	api.POST("/users", userHandler.CreateUser)
	api.PUT("/users/:id", userHandler.UpdateUser)
	api.DELETE("/users/:id", userHandler.DeleteUser)

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	e.Logger.Fatal(e.Start(":" + cfg.Port))
}
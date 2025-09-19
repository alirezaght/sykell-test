package integration

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	
	"sykell-backend/internal/config"
	"sykell-backend/internal/logger"
	sykellMiddleware "sykell-backend/internal/middleware"
	"sykell-backend/internal/temporal"
	"sykell-backend/internal/url"
	"sykell-backend/internal/user"
	
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// TestURLWorkflow_Integration tests the complete workflow: signup -> login -> add URL
func TestURLWorkflow_Integration(t *testing.T) {
	// Skip if running in short mode
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	
	// Initialize test configuration
	cfg := &config.Config{
		JWTSecret:   "test-secret-key-for-integration-testing",
		LogLevel:    "debug",
		LogFormat:   "json",
		Environment: "test",
		Port:        "8080",
		TemporalHostPort: "localhost:7233",
		Namespace: "default",
		DatabaseURL: "sykell_user:sykell_password@tcp(localhost:3306)/sykell_db?charset=utf8mb4&parseTime=True&loc=Local",
	}
	
	// Initialize logger
	if err := logger.InitLogger(cfg.LogLevel, cfg.LogFormat, cfg.Environment); err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()
	
	// Setup database connection
	database, err := sql.Open("mysql", cfg.DatabaseURL)
	require.NoError(t, err)
	defer database.Close()
	
	// Test database connection
	err = database.Ping()
	require.NoError(t, err)
	
	// Generate unique test email to avoid conflicts
	testEmail := "testuser" + time.Now().Format("20060102150405") + "@example.com"
	testPassword := "testpassword123"
	testURL := "https://example.com"
	var testUserID string
	
	// Cleanup function to remove test data
	cleanup := func() {
		// Delete test URLs for the user
		_, err := database.Exec("DELETE FROM urls WHERE user_id = ?", testUserID)
		if err != nil {
			t.Logf("Warning: Failed to clean up test URLs: %v", err)
		}
		
		// Delete test user
		_, err = database.Exec("DELETE FROM users WHERE email = ?", testEmail)
		if err != nil {
			t.Logf("Warning: Failed to clean up test user: %v", err)
		}
		
		logger.Info("Test data cleanup completed",
		zap.String("test_email", testEmail),
		zap.String("test_user_id", testUserID),
	)
}

// Ensure cleanup runs even if test fails
defer cleanup()

// Initialize services
userRepo := user.NewRepo(database)
userService := user.NewUserService(userRepo, cfg)
userHandler := user.NewUserHandler(userService)

urlRepo := url.NewRepo(database)
urlService := url.NewService(urlRepo, cfg)
urlHandler := url.NewHandler(urlService)

// Initialize Temporal (optional for this test)
temporalService := temporal.NewService(cfg)
defer temporalService.Close()

// Create Echo instance with middleware
e := echo.New()
e.Use(middleware.Recover())
e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	AllowOrigins:     []string{"http://localhost:5173"},
	AllowMethods:     []string{echo.GET, echo.PUT, echo.POST, echo.DELETE, echo.OPTIONS},
	AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	AllowCredentials: true,
}))

// Setup routes
api := e.Group("/api/v1")

// Auth routes
api.POST("/auth/register", userHandler.Register)
api.POST("/auth/login", userHandler.Login)

// Protected routes
protected := api.Group("", sykellMiddleware.JWTMiddleware([]byte(cfg.JWTSecret), false))
protected.POST("/urls", urlHandler.AddURL)
protected.GET("/urls", urlHandler.ListURLs)

logger.Info("Starting URL integration test",
zap.String("test_email", testEmail),
zap.String("test_url", testURL),
)

// Step 1: Register user
t.Run("step_1_register_user", func(t *testing.T) {
	registerReq := map[string]interface{}{
		"email":    testEmail,
		"password": testPassword,
	}
	
	reqBody, err := json.Marshal(registerReq)
	require.NoError(t, err)
	
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	
	e.ServeHTTP(rec, req)
	
	assert.Equal(t, http.StatusCreated, rec.Code)
	
	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Equal(t, "User registered successfully", response["message"])
	
	logger.Info("User registration successful",
	zap.String("email", testEmail),
	zap.Int("status_code", rec.Code),
)
})

var authToken string

// Step 2: Login user and get token
t.Run("step_2_login_user", func(t *testing.T) {
	loginReq := map[string]interface{}{
		"email":    testEmail,
		"password": testPassword,
	}
	
	reqBody, err := json.Marshal(loginReq)
	require.NoError(t, err)
	
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	
	e.ServeHTTP(rec, req)
	
	assert.Equal(t, http.StatusOK, rec.Code)
	
	var loginResponse map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &loginResponse)
	require.NoError(t, err)
	
	assert.Contains(t, loginResponse, "token")
	assert.Contains(t, loginResponse, "user")
	
	authToken = loginResponse["token"].(string)
	assert.NotEmpty(t, authToken)
	
	// Verify user information
	userInfo := loginResponse["user"].(map[string]interface{})
	assert.Equal(t, testEmail, userInfo["email"])
	
	// Capture user ID for cleanup
	testUserID = userInfo["id"].(string)
	assert.NotEmpty(t, testUserID)
	
	logger.Info("User login successful",
	zap.String("email", testEmail),
	zap.String("user_id", testUserID),
	zap.Int("token_length", len(authToken)),
	zap.Int("status_code", rec.Code),
)
})

// Step 3: Add URL using the authentication token
t.Run("step_3_add_url", func(t *testing.T) {
	addURLReq := map[string]interface{}{
		"url": testURL,
	}
	
	reqBody, err := json.Marshal(addURLReq)
	require.NoError(t, err)
	
	req := httptest.NewRequest(http.MethodPost, "/api/v1/urls", bytes.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+authToken)
	rec := httptest.NewRecorder()
	
	e.ServeHTTP(rec, req)
	
	assert.Equal(t, http.StatusOK, rec.Code)
	
	logger.Info("URL added successfully",
	zap.String("url", testURL),
	zap.Int("status_code", rec.Code),
)
})

// Step 4: Verify URL was added by listing URLs
t.Run("step_4_verify_url_added", func(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/urls", nil)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+authToken)
	rec := httptest.NewRecorder()
	
	e.ServeHTTP(rec, req)
	
	assert.Equal(t, http.StatusOK, rec.Code)
	
	var urlsResponse map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &urlsResponse)
	require.NoError(t, err)
	
	assert.Contains(t, urlsResponse, "urls")
	assert.Contains(t, urlsResponse, "total_count")
	
	urls := urlsResponse["urls"].([]interface{})
	totalCount := urlsResponse["total_count"].(float64)
	
	assert.Greater(t, int(totalCount), 0, "Should have at least one URL")
	assert.Greater(t, len(urls), 0, "URLs array should not be empty")
	
	// Verify our test URL is in the list (accounting for URL normalization)
	foundTestURL := false
	for _, urlItem := range urls {
		urlData := urlItem.(map[string]interface{})
		// Check both original_url and normalized_url
		if originalURL, ok := urlData["original_url"].(string); ok {
			if originalURL == testURL {
				foundTestURL = true
				break
			}
		}
		if normalizedURL, ok := urlData["normalized_url"].(string); ok {
			// Account for URL normalization (trailing slash)
			if normalizedURL == testURL || normalizedURL == testURL+"/" {
				foundTestURL = true
				break
			}
		}
	}
	
	assert.True(t, foundTestURL, "Test URL should be found in the list")
	
	logger.Info("URL verification successful",
	zap.String("url", testURL),
	zap.Int("total_urls", int(totalCount)),
	zap.Bool("found_test_url", foundTestURL),
)
})

logger.Info("URL integration test completed successfully",
zap.String("test_email", testEmail),
zap.String("test_url", testURL),
)
}

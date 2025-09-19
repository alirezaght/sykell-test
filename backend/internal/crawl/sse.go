package crawl

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sykell-backend/internal/logger"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// SSEManager manages SSE connections and broadcasting
type SSEManager struct {
	clients map[string]chan SSENotification // userID -> channel
	mutex   sync.RWMutex
}

// Global SSE manager instance
var sseManager = &SSEManager{
	clients: make(map[string]chan SSENotification),
}

// StreamCrawlUpdates handles SSE connections for crawl status notifications
func (h *CrawlHandler) StreamCrawlUpdates(c echo.Context) error {
	logger.Debug("SSE connection attempt", zap.String("remote_ip", c.RealIP()))
	
	// Get user ID from middleware context (JWT middleware handles cookie authentication)
	userID := c.Get("user_id").(string)
	logger.Debug("SSE connection for user", zap.String("user_id", userID))
	
	// Set SSE headers
	c.Response().Header().Set("Content-Type", "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")
	c.Response().Header().Set("Access-Control-Allow-Credentials", "true")	
	c.Response().Header().Set("Access-Control-Allow-Headers", "Cache-Control")

	ctx := c.Request().Context()
	
	// Create a channel for this client
	clientChan := make(chan SSENotification, 10)
	
	// Register client
	sseManager.mutex.Lock()
	sseManager.clients[userID] = clientChan
	sseManager.mutex.Unlock()
	
	// Cleanup on disconnect
	defer func() {
		sseManager.mutex.Lock()
		delete(sseManager.clients, userID)
		close(clientChan)
		sseManager.mutex.Unlock()
		logger.Debug("SSE connection closed for user", zap.String("user_id", userID))
	}()

	logger.Info("SSE connection established for user", zap.String("user_id", userID))
	
	// Send initial connection confirmation
	initialNotification := SSENotification{
		Type:      "connection",
		UserID:    userID,
		Timestamp: time.Now(),
	}
	if err := sendSSEEvent(c, initialNotification); err != nil {
		return err
	}

	// Listen for notifications and context cancellation
	for {
		select {
		case <-ctx.Done():
			return nil
		case notification := <-clientChan:
			if err := sendSSEEvent(c, notification); err != nil {
				logger.Error("Error sending SSE event", zap.Error(err))
				return err
			}
		case <-time.After(30 * time.Second):
			// Send keepalive ping every 30 seconds
			pingNotification := SSENotification{
				Type:      "ping",
				UserID:    userID,
				Timestamp: time.Now(),
			}
			if err := sendSSEEvent(c, pingNotification); err != nil {
				return err
			}
		}
	}
}

// NotifyCrawlUpdate sends a notification to invalidate a specific URL's data
func NotifyCrawlUpdate(userID, urlID string) {
	notification := SSENotification{
		Type:      "crawl_update",
		URLID:     urlID,
		UserID:    userID,
		Timestamp: time.Now(),
	}

	sseManager.mutex.RLock()
	defer sseManager.mutex.RUnlock()

	logger.Debug("SSE notification attempt", 
		zap.String("user_id", userID), 
		zap.String("url_id", urlID), 
		zap.Int("total_clients", len(sseManager.clients)))

	if clientChan, exists := sseManager.clients[userID]; exists {
		select {
		case clientChan <- notification:
			logger.Debug("Sent crawl update notification", 
				zap.String("user_id", userID), 
				zap.String("url_id", urlID))
		default:
			logger.Warn("Client channel full, dropping notification", zap.String("user_id", userID))
		}
	} else {
		logger.Debug("No SSE connection found for user", zap.String("user_id", userID))
	}
}



// sendSSEEvent sends an SSE event to the client
func sendSSEEvent(c echo.Context, notification SSENotification) error {
	data, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	logger.Debug("Sending SSE event to client", zap.String("data", string(data)))

	// Write SSE format: data: {json}\n\n
	_, err = fmt.Fprintf(c.Response().Writer, "data: %s\n\n", string(data))
	if err != nil {
		return err
	}

	// Flush the response
	if flusher, ok := c.Response().Writer.(http.Flusher); ok {
		flusher.Flush()
	}

	return nil
}
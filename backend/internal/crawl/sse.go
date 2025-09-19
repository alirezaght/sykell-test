package crawl

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sykell-backend/internal/logger"
	"sync"
	"sync/atomic"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// SSEManager manages SSE connections and broadcasting
type SSEManager struct {
	clients map[string]map[string]chan SSENotification
	mutex   sync.RWMutex
	nextID  uint64
}

// Global SSE manager instance
var sseManager = &SSEManager{
	clients: make(map[string]map[string]chan SSENotification),
}

func (m *SSEManager) addClient(userID string, ch chan SSENotification) string {
	id := strconv.FormatUint(atomic.AddUint64(&m.nextID, 1), 10)
	m.mutex.Lock()
	if _, ok := m.clients[userID]; !ok {
		m.clients[userID] = make(map[string]chan SSENotification)
	}
	m.clients[userID][id] = ch
	m.mutex.Unlock()
	return id
}

func (m *SSEManager) removeClient(userID, connID string) {
	m.mutex.Lock()
	if conns, ok := m.clients[userID]; ok {
		if ch, ok2 := conns[connID]; ok2 {
			delete(conns, connID)
			close(ch)
		}
		if len(conns) == 0 {
			delete(m.clients, userID)
		}
	}
	m.mutex.Unlock()
}

func (m *SSEManager) connectionCount(userID string) int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	if conns, ok := m.clients[userID]; ok {
		return len(conns)
	}
	return 0
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
	
	clientChan := make(chan SSENotification, 10)
	connID := sseManager.addClient(userID, clientChan)
	
	// Cleanup on disconnect
	defer func() {
		sseManager.removeClient(userID, connID)
		logger.Debug("SSE connection closed for user", 
			zap.String("user_id", userID), 
			zap.String("conn_id", connID),
			zap.Int("remaining", sseManager.connectionCount(userID)))
	}()

	logger.Info("SSE connection established for user", 
		zap.String("user_id", userID), 
		zap.String("conn_id", connID),
		zap.Int("total_connections", sseManager.connectionCount(userID)))
	
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
		case notification, ok := <-clientChan:
			if !ok {
				return nil
			}
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
	conns := sseManager.clients[userID]
	var chans []chan SSENotification
	for _, ch := range conns {
		chans = append(chans, ch)
	}
	sseManager.mutex.RUnlock()

	sent := 0
	for _, ch := range chans {
		select {
		case ch <- notification:
			sent++
		default:
		}
	}
	logger.Debug("SSE broadcast attempted", 
		zap.String("user_id", userID), 
		zap.String("url_id", urlID), 
		zap.Int("connections", len(chans)),
		zap.Int("sent", sent))
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
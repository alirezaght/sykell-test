package crawl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"sykell-backend/internal/logger"

	"go.uber.org/zap"
)

// NotifyCrawlUpdateHTTP sends an HTTP request to the main server to trigger SSE notifications
// This is used from the Temporal worker process to communicate with the main server process
func NotifyCrawlUpdateHTTP(userID, urlID string) {			
	// Create the notification request
	request := NotificationRequest{
		UserID: userID,
		URLID:  urlID,
	}
	
	// Marshal to JSON
	jsonData, err := json.Marshal(request)
	if err != nil {
		logger.Error("Error marshaling notification request", zap.Error(err))
		return
	}
	
	// Make HTTP request to the main server
	url := fmt.Sprintf("%s/api/v1/internal/notify-crawl-update", os.Getenv("BACKEND_URL"))
	
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Error("Error sending notification to main server", zap.Error(err))
		return
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		logger.Warn("Notification request failed", zap.Int("status_code", resp.StatusCode))
		return
	}
	
	logger.Info("Successfully sent crawl update notification", 
		zap.String("user_id", userID), 
		zap.String("url_id", urlID))
}
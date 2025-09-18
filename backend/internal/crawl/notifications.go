package crawl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// NotificationRequest represents the payload for internal notification requests
type NotificationRequest struct {
	UserID string `json:"user_id"`
	URLID  string `json:"url_id"`
}

// NotifyCrawlUpdateHTTP sends an HTTP request to the main server to trigger SSE notifications
// This is used from the Temporal worker process to communicate with the main server process
func NotifyCrawlUpdateHTTP(userID, urlID string) {
	// Get the server port from environment, default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	// Create the notification request
	request := NotificationRequest{
		UserID: userID,
		URLID:  urlID,
	}
	
	// Marshal to JSON
	jsonData, err := json.Marshal(request)
	if err != nil {
		log.Printf("Error marshaling notification request: %v", err)
		return
	}
	
	// Make HTTP request to the main server
	url := fmt.Sprintf("http://localhost:%s/api/v1/internal/notify-crawl-update", port)
	
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error sending notification to main server: %v", err)
		return
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		log.Printf("Notification request failed with status: %d", resp.StatusCode)
		return
	}
	
	log.Printf("Successfully sent crawl update notification for user %s, URL %s", userID, urlID)
}
package temporal

import (
	"sykell-backend/internal/config"
	"time"

	"go.temporal.io/sdk/client"
)

// Service provides Temporal-related services
type Service struct {
	config *config.Config
	temporalClient client.Client
}

// NewService creates a new Temporal Service
func NewService(config *config.Config) *Service {
	return &Service{
		config: config,
		temporalClient: nil,
	}
}


// Setup initializes the Temporal client
func (s *Service) Setup() {
	s.temporalClient, _ = client.Dial(client.Options{
		HostPort:  s.config.TemporalHostPort,
		Namespace: s.config.Namespace,
		ConnectionOptions: client.ConnectionOptions{
			TLS: nil, // Disable TLS for local development
			KeepAliveTime:  10 * time.Second, // 10 seconds
			KeepAliveTimeout: 20 * time.Second, // 20 seconds			
		},		
	})	
}

// GetTemporalClient returns the Temporal client, initializing it if necessary
func (s *Service) GetTemporalClient() client.Client {
	if s.temporalClient == nil {
		s.Setup()
	}
	return s.temporalClient
}

// Close closes the Temporal client connection
func (s *Service) Close() {
	if s.temporalClient != nil {
		s.temporalClient.Close()
	}
}
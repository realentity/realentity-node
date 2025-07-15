package impl

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/realentity/realentity-node/internal/services"
)

// EchoRequest represents the payload for echo service
type EchoRequest struct {
	Message string `json:"message"`
}

// EchoResponse represents the response from echo service
type EchoResponse struct {
	Echo      string    `json:"echo"`
	Timestamp time.Time `json:"timestamp"`
	NodeID    string    `json:"nodeId"`
}

// CreateEchoService creates a simple echo service for testing
func CreateEchoService(nodeID string) *services.Service {
	return &services.Service{
		Name:        "echo",
		Description: "Simple echo service for testing connectivity",
		Version:     "1.0.0",
		Metadata: map[string]string{
			"category": "utility",
			"cost":     "free",
		},
		Handler: func(payload []byte) ([]byte, error) {
			var req EchoRequest
			if err := json.Unmarshal(payload, &req); err != nil {
				return nil, fmt.Errorf("invalid echo request: %v", err)
			}

			log.Printf("echo service executing: message='%s', nodeID='%s'", req.Message, nodeID)

			response := EchoResponse{
				Echo:      req.Message,
				Timestamp: time.Now(),
				NodeID:    nodeID,
			}

			result, err := json.Marshal(response)
			if err != nil {
				return nil, err
			}

			log.Printf("echo service completed: echoed='%s'", req.Message)
			return result, nil
		},
	}
}

// Register this service with the global registry
func init() {
	services.GlobalServiceRegistry.RegisterServiceFactory("echo", CreateEchoService)
}

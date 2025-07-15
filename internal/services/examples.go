package services

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
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
func CreateEchoService(nodeID string) *Service {
	return &Service{
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

			response := EchoResponse{
				Echo:      req.Message,
				Timestamp: time.Now(),
				NodeID:    nodeID,
			}

			return json.Marshal(response)
		},
	}
}

// TextProcessRequest represents the payload for text processing
type TextProcessRequest struct {
	Text      string `json:"text"`
	Operation string `json:"operation"` // "uppercase", "lowercase", "reverse"
}

// TextProcessResponse represents the response from text processing
type TextProcessResponse struct {
	Original  string `json:"original"`
	Processed string `json:"processed"`
	Operation string `json:"operation"`
}

// CreateTextProcessService creates a text processing service
func CreateTextProcessService() *Service {
	return &Service{
		Name:        "text.process",
		Description: "Basic text processing operations",
		Version:     "1.0.0",
		Metadata: map[string]string{
			"category":   "text",
			"operations": "uppercase,lowercase,reverse",
			"cost":       "free",
		},
		Handler: func(payload []byte) ([]byte, error) {
			var req TextProcessRequest
			if err := json.Unmarshal(payload, &req); err != nil {
				return nil, fmt.Errorf("invalid text process request: %v", err)
			}

			log.Printf("text.process service executing: operation='%s', text='%s'", req.Operation, req.Text)

			var processed string
			switch req.Operation {
			case "uppercase":
				processed = strings.ToUpper(req.Text)
			case "lowercase":
				processed = strings.ToLower(req.Text)
			case "reverse":
				runes := []rune(req.Text)
				for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
					runes[i], runes[j] = runes[j], runes[i]
				}
				processed = string(runes)
			default:
				return nil, fmt.Errorf("unsupported operation: %s", req.Operation)
			}

			response := TextProcessResponse{
				Original:  req.Text,
				Processed: processed,
				Operation: req.Operation,
			}

			log.Printf("text.process service completed: result='%s'", processed)
			return json.Marshal(response)
		},
	}
}

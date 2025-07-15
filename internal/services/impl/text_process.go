package impl

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/realentity/realentity-node/internal/services"
)

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
func CreateTextProcessService(nodeID string) *services.Service {
	return &services.Service{
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

			result, err := json.Marshal(response)
			if err != nil {
				return nil, err
			}

			log.Printf("text.process service completed: result='%s'", processed)
			return result, nil
		},
	}
}

// Register this service with the global registry
func init() {
	services.GlobalServiceRegistry.RegisterServiceFactory("text.process", CreateTextProcessService)
}

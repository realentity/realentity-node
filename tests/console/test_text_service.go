package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/realentity/realentity-node/internal/services"
)

// Simple program to test the text.process service from console
func main() {
	// Initialize services
	fmt.Println("Initializing services...")

	// Register the text processing service
	textService := services.CreateTextProcessService()
	err := services.GlobalRegistry.RegisterService(textService)
	if err != nil {
		log.Fatalf("Failed to register text service: %v", err)
	}

	// Get command line arguments
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run test_text_service.go <text> <operation>")
		fmt.Println("Operations: uppercase, lowercase, reverse")
		fmt.Println("Example: go run test_text_service.go \"hello world\" uppercase")
		os.Exit(1)
	}

	text := os.Args[1]
	operation := os.Args[2]

	// Validate operation
	validOps := map[string]bool{
		"uppercase": true,
		"lowercase": true,
		"reverse":   true,
	}

	if !validOps[operation] {
		fmt.Printf("Invalid operation: %s\n", operation)
		fmt.Println("Valid operations: uppercase, lowercase, reverse")
		os.Exit(1)
	}

	// Create service request
	request := &services.ServiceRequest{
		Service:   "text.process",
		RequestID: "console-test",
	}

	// Create payload
	payload := map[string]interface{}{
		"text":      text,
		"operation": operation,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("Failed to marshal payload: %v", err)
	}
	request.Payload = json.RawMessage(payloadBytes)

	// Execute service
	fmt.Printf("Processing text: '%s' with operation: '%s'\n", text, operation)
	response := services.GlobalRegistry.ExecuteService(request)

	// Check response
	if !response.Success {
		fmt.Printf("Service execution failed: %s\n", response.Error)
		os.Exit(1)
	}

	// Parse and display result
	var result map[string]interface{}
	if err := json.Unmarshal(response.Result, &result); err != nil {
		fmt.Printf("Failed to parse result: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Original: %s\n", result["original"])
	fmt.Printf("Processed: %s\n", result["processed"])
	fmt.Printf("Operation: %s\n", result["operation"])
}

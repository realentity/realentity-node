package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/realentity/realentity-node/internal/services"
)

func main() {
	// Initialize services
	fmt.Println("Interactive Text Processing Service Tester")
	fmt.Println("============================================")

	// Register the text processing service
	textService := services.CreateTextProcessService()
	err := services.GlobalRegistry.RegisterService(textService)
	if err != nil {
		fmt.Printf("Failed to register text service: %v\n", err)
		return
	}

	fmt.Println("Services initialized successfully!")
	fmt.Println("\nAvailable operations: uppercase, lowercase, reverse")
	fmt.Println("Type 'quit' or 'exit' to stop\n")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Enter text to process: ")
		if !scanner.Scan() {
			break
		}
		text := strings.TrimSpace(scanner.Text())

		if text == "quit" || text == "exit" {
			fmt.Println("Goodbye!")
			break
		}

		if text == "" {
			fmt.Println("️  Please enter some text")
			continue
		}

		fmt.Print("Enter operation (uppercase/lowercase/reverse): ")
		if !scanner.Scan() {
			break
		}
		operation := strings.TrimSpace(strings.ToLower(scanner.Text()))

		validOps := map[string]bool{
			"uppercase": true,
			"lowercase": true,
			"reverse":   true,
		}

		if !validOps[operation] {
			fmt.Printf("Invalid operation: %s\n", operation)
			fmt.Println("Valid operations: uppercase, lowercase, reverse\n")
			continue
		}

		// Execute service
		result := executeTextService(text, operation)
		fmt.Printf("Result: %s\n\n", result)
	}
}

func executeTextService(text, operation string) string {
	request := &services.ServiceRequest{
		Service:   "text.process",
		RequestID: "interactive-test",
	}

	payload := map[string]interface{}{
		"text":      text,
		"operation": operation,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Sprintf("Error marshaling payload: %v", err)
	}
	request.Payload = json.RawMessage(payloadBytes)

	response := services.GlobalRegistry.ExecuteService(request)

	if !response.Success {
		return fmt.Sprintf("Service execution failed: %s", response.Error)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(response.Result, &result); err != nil {
		return fmt.Sprintf("Failed to parse result: %v", err)
	}

	return fmt.Sprintf("%s", result["processed"])
}

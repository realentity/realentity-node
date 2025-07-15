package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"strings"
)

// ServiceConfig represents configuration for a service
type ServiceConfig struct {
	Name    string                 `json:"name"`
	Enabled bool                   `json:"enabled"`
	Version string                 `json:"version,omitempty"`
	Config  map[string]interface{} `json:"config,omitempty"`
}

// ServicesConfig represents the configuration file structure
type ServicesConfig struct {
	Services []ServiceConfig `json:"services"`
}

// ServiceLoader handles loading and initializing services
type ServiceLoader struct {
	manager *ServiceManager
}

// NewServiceLoader creates a new service loader
func NewServiceLoader(manager *ServiceManager) *ServiceLoader {
	return &ServiceLoader{
		manager: manager,
	}
}

// LoadAllServices loads all available services with default configurations
func (sl *ServiceLoader) LoadAllServices(nodeID string) error {
	log.Printf("Loading all services for node %s", nodeID)

	// Load echo service
	echoService := NewEchoServiceImpl(nodeID)
	if err := sl.manager.RegisterService(echoService); err != nil {
		log.Printf("Failed to register echo service: %v", err)
	} else {
		echoService.Start()
	}

	// Load text processing service
	textService := NewTextProcessServiceImpl()
	if err := sl.manager.RegisterService(textService); err != nil {
		log.Printf("Failed to register text service: %v", err)
	} else {
		textService.Start()
	}

	log.Printf("All services loaded successfully")
	return nil
}

// LoadServicesFromConfig loads services from a configuration file
func (sl *ServiceLoader) LoadServicesFromConfig(configPath string, nodeID string) error {
	// Check if config file exists
	if !fileExists(configPath) {
		return fmt.Errorf("config file %s not found", configPath)
	}

	// Read configuration file
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}

	// Parse configuration
	var config ServicesConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse config file: %v", err)
	}

	// Load services based on configuration
	for _, serviceConfig := range config.Services {
		service, err := sl.createServiceByName(serviceConfig.Name, nodeID)
		if err != nil {
			log.Printf("Failed to create service %s: %v", serviceConfig.Name, err)
			continue
		}

		// Apply configuration
		if len(serviceConfig.Config) > 0 {
			if err := service.SetConfig(serviceConfig.Config); err != nil {
				log.Printf("Failed to set config for service %s: %v", serviceConfig.Name, err)
			}
		}

		// Register service
		if err := sl.manager.RegisterService(service); err != nil {
			log.Printf("Failed to register service %s: %v", serviceConfig.Name, err)
			continue
		}

		// Start service if enabled
		if serviceConfig.Enabled {
			if err := service.Start(); err != nil {
				log.Printf("Failed to start service %s: %v", serviceConfig.Name, err)
			}
		}
	}

	return nil
}

// createServiceByName creates a service instance by name
func (sl *ServiceLoader) createServiceByName(name string, nodeID string) (ServiceProvider, error) {
	switch name {
	case "echo":
		return NewEchoServiceImpl(nodeID), nil
	case "text.process":
		return NewTextProcessServiceImpl(), nil
	default:
		return nil, fmt.Errorf("unknown service: %s", name)
	}
}

// LoadServiceFromDirectory loads services from a directory structure
func (sl *ServiceLoader) LoadServiceFromDirectory(dir string, nodeID string) error {
	// This could be extended to automatically discover services in directories
	// For now, we'll use the predefined services
	return sl.LoadAllServices(nodeID)
}

// CreateDefaultConfig creates a default configuration file
func (sl *ServiceLoader) CreateDefaultConfig(configPath string) error {
	defaultConfig := ServicesConfig{
		Services: []ServiceConfig{
			{
				Name:    "echo",
				Enabled: true,
				Version: "1.0.0",
				Config: map[string]interface{}{
					"prefix": "Echo: ",
				},
			},
			{
				Name:    "text.process",
				Enabled: true,
				Version: "1.0.0",
				Config: map[string]interface{}{
					"max_length": 1000,
					"operations": []string{"uppercase", "lowercase", "reverse", "word_count"},
				},
			},
		},
	}

	data, err := json.MarshalIndent(defaultConfig, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal default config: %v", err)
	}

	if err := ioutil.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	log.Printf("Default configuration created at %s", configPath)
	return nil
}

// ReloadServices reloads all services from configuration
func (sl *ServiceLoader) ReloadServices(configPath string, nodeID string) error {
	// Stop all current services
	if err := sl.manager.StopAllServices(); err != nil {
		log.Printf("Warning: failed to stop some services during reload: %v", err)
	}

	// Clear all services (this would need to be implemented in ServiceManager)
	// For now, we'll just load new services

	// Reload from config
	return sl.LoadServicesFromConfig(configPath, nodeID)
}

// fileExists checks if a file exists
func fileExists(filename string) bool {
	_, err := filepath.Abs(filename)
	if err != nil {
		return false
	}
	_, err = ioutil.ReadFile(filename)
	return err == nil
}

// Global service loader instance
var GlobalServiceLoader = NewServiceLoader(GlobalServiceManager)

// EchoServiceImpl implements a simple echo service
type EchoServiceImpl struct {
	*BaseService
	nodeID string
}

// NewEchoServiceImpl creates a new echo service instance
func NewEchoServiceImpl(nodeID string) *EchoServiceImpl {
	return &EchoServiceImpl{
		BaseService: NewBaseService(
			"echo",
			"1.0.0",
			"Simple echo service that returns the input message with optional prefix",
		),
		nodeID: nodeID,
	}
}

// Execute processes the echo request
func (e *EchoServiceImpl) Execute(request ServiceRequest) (*ServiceResponse, error) {
	if !e.IsEnabled() {
		return &ServiceResponse{
			RequestID: request.RequestID,
			Success:   false,
			Error:     "Echo service is not enabled",
		}, nil
	}

	// Parse the payload
	var data map[string]interface{}
	if len(request.Payload) > 0 {
		if err := json.Unmarshal(request.Payload, &data); err != nil {
			return &ServiceResponse{
				RequestID: request.RequestID,
				Success:   false,
				Error:     "Invalid payload format",
			}, nil
		}
	} else {
		data = make(map[string]interface{})
	}

	// Get the input message
	message, ok := data["message"]
	if !ok {
		return &ServiceResponse{
			RequestID: request.RequestID,
			Success:   false,
			Error:     "Missing 'message' parameter",
		}, nil
	}

	messageStr, ok := message.(string)
	if !ok {
		return &ServiceResponse{
			RequestID: request.RequestID,
			Success:   false,
			Error:     "Message parameter must be a string",
		}, nil
	}

	// Get prefix from config if available
	config := e.GetConfig()
	prefix := ""
	if prefixVal, exists := config["prefix"]; exists {
		if prefixStr, ok := prefixVal.(string); ok {
			prefix = prefixStr
		}
	}

	// Process the message
	result := prefix + messageStr

	// Add node ID if requested
	if addNodeID, exists := data["add_node_id"]; exists {
		if addNodeIDBool, ok := addNodeID.(bool); ok && addNodeIDBool {
			result = fmt.Sprintf("[Node: %s] %s", e.nodeID, result)
		}
	}

	// Handle special commands
	if strings.HasPrefix(messageStr, "/") {
		switch strings.ToLower(messageStr) {
		case "/status":
			result = fmt.Sprintf("Echo service is running on node %s", e.nodeID)
		case "/help":
			result = "Available commands: /status, /help, /version. Or send any message to echo it back."
		case "/version":
			result = fmt.Sprintf("Echo service version %s", e.GetVersion())
		default:
			result = fmt.Sprintf("Unknown command: %s. Type /help for available commands.", messageStr)
		}
	}

	// Prepare response data
	responseData := map[string]interface{}{
		"original": messageStr,
		"echo":     result,
		"node_id":  e.nodeID,
		"service":  "echo",
	}

	resultJSON, err := json.Marshal(responseData)
	if err != nil {
		return &ServiceResponse{
			RequestID: request.RequestID,
			Success:   false,
			Error:     "Failed to marshal response data",
		}, nil
	}

	return &ServiceResponse{
		RequestID: request.RequestID,
		Success:   true,
		Result:    resultJSON,
	}, nil
}

// Start starts the echo service
func (e *EchoServiceImpl) Start() error {
	if err := e.BaseService.Start(); err != nil {
		return err
	}

	// Set default configuration if not already set
	config := e.GetConfig()
	if len(config) == 0 {
		defaultConfig := map[string]interface{}{
			"prefix": "Echo: ",
		}
		e.SetConfig(defaultConfig)
	}

	return nil
}

// TextProcessServiceImpl implements text processing operations
type TextProcessServiceImpl struct {
	*BaseService
}

// NewTextProcessServiceImpl creates a new text processing service instance
func NewTextProcessServiceImpl() *TextProcessServiceImpl {
	return &TextProcessServiceImpl{
		BaseService: NewBaseService(
			"text.process",
			"1.0.0",
			"Advanced text processing service with multiple operations",
		),
	}
}

// Execute processes the text processing request
func (t *TextProcessServiceImpl) Execute(request ServiceRequest) (*ServiceResponse, error) {
	if !t.IsEnabled() {
		return &ServiceResponse{
			RequestID: request.RequestID,
			Success:   false,
			Error:     "Text processing service is not enabled",
		}, nil
	}

	// Parse the payload
	var data map[string]interface{}
	if len(request.Payload) > 0 {
		if err := json.Unmarshal(request.Payload, &data); err != nil {
			return &ServiceResponse{
				RequestID: request.RequestID,
				Success:   false,
				Error:     "Invalid payload format",
			}, nil
		}
	} else {
		data = make(map[string]interface{})
	}

	// Get the input text
	text, ok := data["text"]
	if !ok {
		return &ServiceResponse{
			RequestID: request.RequestID,
			Success:   false,
			Error:     "Missing 'text' parameter",
		}, nil
	}

	textStr, ok := text.(string)
	if !ok {
		return &ServiceResponse{
			RequestID: request.RequestID,
			Success:   false,
			Error:     "Text parameter must be a string",
		}, nil
	}

	// Check text length limit
	config := t.GetConfig()
	maxLength := 1000 // default
	if maxLenVal, exists := config["max_length"]; exists {
		if maxLenFloat, ok := maxLenVal.(float64); ok {
			maxLength = int(maxLenFloat)
		} else if maxLenInt, ok := maxLenVal.(int); ok {
			maxLength = maxLenInt
		}
	}

	if len(textStr) > maxLength {
		return &ServiceResponse{
			RequestID: request.RequestID,
			Success:   false,
			Error:     fmt.Sprintf("Text length exceeds maximum of %d characters", maxLength),
		}, nil
	}

	// Get the operation to perform
	operation, ok := data["operation"]
	if !ok {
		operation = "info" // default operation
	}

	operationStr, ok := operation.(string)
	if !ok {
		return &ServiceResponse{
			RequestID: request.RequestID,
			Success:   false,
			Error:     "Operation parameter must be a string",
		}, nil
	}

	// Perform the requested operation
	result, err := t.processText(textStr, operationStr, data)
	if err != nil {
		return &ServiceResponse{
			RequestID: request.RequestID,
			Success:   false,
			Error:     err.Error(),
		}, nil
	}

	// Prepare response data
	responseData := map[string]interface{}{
		"original":  textStr,
		"operation": operationStr,
		"result":    result,
		"service":   "text.process",
	}

	resultJSON, err := json.Marshal(responseData)
	if err != nil {
		return &ServiceResponse{
			RequestID: request.RequestID,
			Success:   false,
			Error:     "Failed to marshal response data",
		}, nil
	}

	return &ServiceResponse{
		RequestID: request.RequestID,
		Success:   true,
		Result:    resultJSON,
	}, nil
}

// processText performs the actual text processing
func (t *TextProcessServiceImpl) processText(text, operation string, data map[string]interface{}) (interface{}, error) {
	switch strings.ToLower(operation) {
	case "uppercase":
		return strings.ToUpper(text), nil

	case "lowercase":
		return strings.ToLower(text), nil

	case "reverse":
		runes := []rune(text)
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		return string(runes), nil

	case "word_count":
		words := strings.Fields(text)
		return map[string]interface{}{
			"words":      len(words),
			"characters": len(text),
			"lines":      len(strings.Split(text, "\n")),
		}, nil

	case "capitalize":
		return strings.Title(text), nil

	case "trim":
		return strings.TrimSpace(text), nil

	case "replace":
		// Get find and replace parameters
		find, findOk := data["find"]
		replace, replaceOk := data["replace"]

		if !findOk || !replaceOk {
			return nil, fmt.Errorf("replace operation requires 'find' and 'replace' parameters")
		}

		findStr, findStrOk := find.(string)
		replaceStr, replaceStrOk := replace.(string)

		if !findStrOk || !replaceStrOk {
			return nil, fmt.Errorf("find and replace parameters must be strings")
		}

		return strings.ReplaceAll(text, findStr, replaceStr), nil

	case "extract_emails":
		emailRegex := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
		emails := emailRegex.FindAllString(text, -1)
		return emails, nil

	case "extract_urls":
		urlRegex := regexp.MustCompile(`https?://[^\s]+`)
		urls := urlRegex.FindAllString(text, -1)
		return urls, nil

	case "split":
		// Get delimiter parameter
		delimiter, delimOk := data["delimiter"]
		if !delimOk {
			delimiter = " " // default to space
		}

		delimiterStr, delimStrOk := delimiter.(string)
		if !delimStrOk {
			return nil, fmt.Errorf("delimiter parameter must be a string")
		}

		return strings.Split(text, delimiterStr), nil

	case "info":
		words := strings.Fields(text)
		lines := strings.Split(text, "\n")
		return map[string]interface{}{
			"length":        len(text),
			"word_count":    len(words),
			"line_count":    len(lines),
			"char_count":    len(text),
			"has_numbers":   regexp.MustCompile(`\d`).MatchString(text),
			"has_uppercase": regexp.MustCompile(`[A-Z]`).MatchString(text),
			"has_lowercase": regexp.MustCompile(`[a-z]`).MatchString(text),
			"has_special":   regexp.MustCompile(`[^a-zA-Z0-9\s]`).MatchString(text),
		}, nil

	case "clean":
		// Remove extra whitespace and normalize
		cleaned := regexp.MustCompile(`\s+`).ReplaceAllString(strings.TrimSpace(text), " ")
		return cleaned, nil

	default:
		// List available operations
		operations := []string{
			"uppercase", "lowercase", "reverse", "word_count", "capitalize",
			"trim", "replace", "extract_emails", "extract_urls", "split",
			"info", "clean",
		}
		return nil, fmt.Errorf("unknown operation '%s'. Available operations: %s",
			operation, strings.Join(operations, ", "))
	}
}

// Start starts the text processing service
func (t *TextProcessServiceImpl) Start() error {
	if err := t.BaseService.Start(); err != nil {
		return err
	}

	// Set default configuration if not already set
	config := t.GetConfig()
	if len(config) == 0 {
		defaultConfig := map[string]interface{}{
			"max_length": 1000,
			"operations": []string{"uppercase", "lowercase", "reverse", "word_count"},
		}
		t.SetConfig(defaultConfig)
	}

	return nil
}

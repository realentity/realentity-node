package main

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/realentity/realentity-node/internal/config"
	"github.com/realentity/realentity-node/internal/discovery"
	"github.com/realentity/realentity-node/internal/node"
	"github.com/realentity/realentity-node/internal/protocol"
	"github.com/realentity/realentity-node/internal/services"
)

// TestInitializeServices tests the service initialization function
func TestInitializeServices(t *testing.T) {
	// Create a temporary registry for testing
	originalRegistry := services.GlobalRegistry
	services.GlobalRegistry = services.NewRegistry()
	defer func() {
		services.GlobalRegistry = originalRegistry
	}()

	// Test service initialization
	nodeID := "test-node-12345"
	initializeServices(nodeID)

	// Check if services were registered
	registeredServices := services.GlobalRegistry.ListServices()

	// Should have at least echo and text services
	if len(registeredServices) < 2 {
		t.Errorf("Expected at least 2 services, got %d", len(registeredServices))
	}

	// Check if echo service exists
	echoService, exists := services.GlobalRegistry.GetService("echo")
	if !exists {
		t.Error("Echo service should be registered")
	}
	if echoService != nil && echoService.Name != "echo" {
		t.Errorf("Expected echo service name to be 'echo', got '%s'", echoService.Name)
	}

	// Check if text service exists
	textService, exists := services.GlobalRegistry.GetService("text.process")
	if !exists {
		t.Error("Text processing service should be registered")
	}
	if textService != nil && textService.Name != "text.process" {
		t.Errorf("Expected text service name to be 'text.process', got '%s'", textService.Name)
	}
}

// TestServiceExecution tests that services can be executed properly
func TestServiceExecution(t *testing.T) {
	// Create a temporary registry for testing
	originalRegistry := services.GlobalRegistry
	services.GlobalRegistry = services.NewRegistry()
	defer func() {
		services.GlobalRegistry = originalRegistry
	}()

	// Initialize services
	nodeID := "test-node-12345"
	initializeServices(nodeID)

	// Test echo service
	echoRequest := services.ServiceRequest{
		Service:   "echo",
		RequestID: "test-request-1",
	}

	// Create echo payload
	echoPayload := map[string]string{"message": "Hello, World!"}
	payloadBytes, err := json.Marshal(echoPayload)
	if err != nil {
		t.Fatalf("Failed to marshal echo payload: %v", err)
	}
	echoRequest.Payload = json.RawMessage(payloadBytes)

	// Execute echo service
	response := services.GlobalRegistry.ExecuteService(&echoRequest)
	if !response.Success {
		t.Errorf("Echo service execution failed: %s", response.Error)
	}
	if response.RequestID != "test-request-1" {
		t.Errorf("Expected request ID 'test-request-1', got '%s'", response.RequestID)
	}

	// Test text processing service
	textRequest := services.ServiceRequest{
		Service:   "text.process",
		RequestID: "test-request-2",
	}

	// Create text processing payload
	textPayload := map[string]interface{}{
		"text":      "hello world",
		"operation": "uppercase",
	}
	textPayloadBytes, err := json.Marshal(textPayload)
	if err != nil {
		t.Fatalf("Failed to marshal text payload: %v", err)
	}
	textRequest.Payload = json.RawMessage(textPayloadBytes)

	// Execute text service
	textResponse := services.GlobalRegistry.ExecuteService(&textRequest)
	if !textResponse.Success {
		t.Errorf("Text service execution failed: %s", textResponse.Error)
	}
	if textResponse.RequestID != "test-request-2" {
		t.Errorf("Expected request ID 'test-request-2', got '%s'", textResponse.RequestID)
	}
}

// TestHostCreation tests the host creation functionality
func TestHostCreation(t *testing.T) {
	ctx := context.Background()

	// Test default host creation
	host, err := node.CreateHost(ctx)
	if err != nil {
		t.Fatalf("Failed to create host: %v", err)
	}
	defer host.Close()

	if host.ID() == "" {
		t.Error("Host should have a valid peer ID")
	}

	// Check if host has addresses
	addrs := host.Addrs()
	if len(addrs) == 0 {
		t.Error("Host should have at least one address")
	}
}

// TestDiscoveryManagerCreation tests discovery manager initialization
func TestDiscoveryManagerCreation(t *testing.T) {
	ctx := context.Background()

	// Create a test host
	host, err := node.CreateHost(ctx)
	if err != nil {
		t.Fatalf("Failed to create host: %v", err)
	}
	defer host.Close()

	// Create discovery manager
	dm := discovery.NewDiscoveryManager(host)
	if dm == nil {
		t.Error("Discovery manager should not be nil")
	}

	// Test getting peers (should be empty initially)
	peers := dm.GetPeers()
	if peers == nil {
		t.Error("GetPeers should return a valid slice, not nil")
	}
	if len(peers) != 0 {
		t.Errorf("Expected 0 initial peers, got %d", len(peers))
	}

	// Test getting connectable peers (can be nil if no peers)
	connectable := dm.GetConnectablePeers()
	if len(connectable) != 0 {
		t.Errorf("Expected 0 initial connectable peers, got %d", len(connectable))
	}
}

// TestProtocolRegistration tests protocol handler registration
func TestProtocolRegistration(t *testing.T) {
	ctx := context.Background()

	// Create a test host
	host, err := node.CreateHost(ctx)
	if err != nil {
		t.Fatalf("Failed to create host: %v", err)
	}
	defer host.Close()

	// Register protocol handler
	protocol.RegisterHandler(host, "/realentity/1.0.0")

	// Check if protocol is registered by getting supported protocols
	protocols := host.Mux().Protocols()
	found := false
	for _, p := range protocols {
		if string(p) == "/realentity/1.0.0" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Protocol /realentity/1.0.0 should be registered")
	}
}

// TestConfigIntegration tests configuration loading with default values
func TestConfigIntegration(t *testing.T) {
	// Create a temporary config file for testing
	tempConfig := "test_config.json"
	defer os.Remove(tempConfig)

	// Create test configuration
	testConfig := config.NodeConfig{
		Discovery: config.DiscoveryConfig{
			EnableMDNS:      true,
			EnableBootstrap: false,
			EnableDHT:       false,
			MDNSServiceTag:  "test-realentity",
			MDNSQuietMode:   true,
			BootstrapPeers:  []string{},
			DHTRendezvous:   "test-realentity",
		},
		Server: config.ServerConfig{
			BindAddress: "127.0.0.1",
			Port:        0, // Use random port for testing
			PublicIP:    "",
		},
	}

	// Write test config to file
	configData, err := json.MarshalIndent(testConfig, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal test config: %v", err)
	}

	err = os.WriteFile(tempConfig, configData, 0644)
	if err != nil {
		t.Fatalf("Failed to write test config file: %v", err)
	}

	// Load configuration
	cfg, err := config.LoadConfig(tempConfig)
	if err != nil {
		t.Fatalf("Failed to load test config: %v", err)
	}

	// Verify configuration values
	if !cfg.Discovery.EnableMDNS {
		t.Error("Expected mDNS to be enabled in test config")
	}
	if cfg.Discovery.EnableBootstrap {
		t.Error("Expected bootstrap to be disabled in test config")
	}
	if cfg.Discovery.MDNSServiceTag != "test-realentity" {
		t.Errorf("Expected mDNS service tag 'test-realentity', got '%s'", cfg.Discovery.MDNSServiceTag)
	}
}

// TestEndToEndFlow tests a simplified end-to-end flow without external dependencies
func TestEndToEndFlow(t *testing.T) {
	// Create a temporary registry for testing
	originalRegistry := services.GlobalRegistry
	services.GlobalRegistry = services.NewRegistry()
	defer func() {
		services.GlobalRegistry = originalRegistry
	}()

	ctx := context.Background()

	// Create host
	host, err := node.CreateHost(ctx)
	if err != nil {
		t.Fatalf("Failed to create host: %v", err)
	}
	defer host.Close()

	// Initialize services
	nodeID := host.ID().String()
	initializeServices(nodeID)

	// Create discovery manager
	dm := discovery.NewDiscoveryManager(host)

	// Register protocol handler
	protocol.RegisterHandler(host, "/realentity/1.0.0")

	// Verify everything is working
	services := services.GlobalRegistry.ListServices()
	if len(services) < 2 {
		t.Errorf("Expected at least 2 services after initialization, got %d", len(services))
	}

	// Test that discovery manager is ready
	peers := dm.GetPeers()
	if peers == nil {
		t.Error("Discovery manager should return valid peer list")
	}
}

// Benchmark for service execution
func BenchmarkServiceExecution(b *testing.B) {
	// Create a temporary registry for testing
	originalRegistry := services.GlobalRegistry
	services.GlobalRegistry = services.NewRegistry()
	defer func() {
		services.GlobalRegistry = originalRegistry
	}()

	// Initialize services
	nodeID := "benchmark-node"
	initializeServices(nodeID)

	// Prepare echo request
	echoRequest := services.ServiceRequest{
		Service:   "echo",
		RequestID: "benchmark-request",
	}

	echoPayload := map[string]string{"message": "Benchmark test"}
	payloadBytes, _ := json.Marshal(echoPayload)
	echoRequest.Payload = json.RawMessage(payloadBytes)

	// Run benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		response := services.GlobalRegistry.ExecuteService(&echoRequest)
		if !response.Success {
			b.Fatalf("Service execution failed: %s", response.Error)
		}
	}
}

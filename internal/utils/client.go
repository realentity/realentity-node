package utils

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	host "github.com/libp2p/go-libp2p/core/host"
	peer "github.com/libp2p/go-libp2p/core/peer"
	protocol "github.com/libp2p/go-libp2p/core/protocol"
	"github.com/realentity/realentity-node/internal/services"
)

// ServiceClient handles communication with remote nodes
type ServiceClient struct {
	host host.Host
}

// NewServiceClient creates a new service client
func NewServiceClient(h host.Host) *ServiceClient {
	return &ServiceClient{host: h}
}

// CallService calls a service on a remote peer
func (c *ServiceClient) CallService(ctx context.Context, peerID peer.ID, serviceName string, payload interface{}) (*services.ServiceResponse, error) {
	// Marshal the payload
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %v", err)
	}

	// Create service request
	request := services.ServiceRequest{
		Service:   serviceName,
		Payload:   json.RawMessage(payloadBytes),
		RequestID: uuid.New().String(),
	}

	// Open stream to peer
	stream, err := c.host.NewStream(ctx, peerID, protocol.ID("/realentity/1.0.0"))
	if err != nil {
		return nil, fmt.Errorf("failed to open stream: %v", err)
	}
	defer stream.Close()

	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	// Send request
	if err := json.NewEncoder(rw).Encode(request); err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	rw.Flush()

	// Read response
	var response services.ServiceResponse
	if err := json.NewDecoder(rw).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	return &response, nil
}

// TestEcho tests the echo service on a remote peer
func (c *ServiceClient) TestEcho(ctx context.Context, peerID peer.ID, message string) error {
	log.Printf("Testing echo service on peer %s with message: %s\n", peerID.String(), message)

	echoReq := services.EchoRequest{Message: message}

	response, err := c.CallService(ctx, peerID, "echo", echoReq)
	if err != nil {
		return fmt.Errorf("echo test failed: %v", err)
	}

	if !response.Success {
		return fmt.Errorf("echo service returned error: %s", response.Error)
	}

	var echoResp services.EchoResponse
	if err := json.Unmarshal(response.Result, &echoResp); err != nil {
		return fmt.Errorf("failed to unmarshal echo response: %v", err)
	}

	log.Printf("Echo response: %s (from node: %s, time: %v)\n",
		echoResp.Echo, echoResp.NodeID, echoResp.Timestamp)
	return nil
}

// TestTextProcess tests the text processing service
func (c *ServiceClient) TestTextProcess(ctx context.Context, peerID peer.ID, text, operation string) error {
	log.Printf("Testing text.process service on peer %s: %s -> %s\n", peerID.String(), operation, text)

	textReq := services.TextProcessRequest{
		Text:      text,
		Operation: operation,
	}

	response, err := c.CallService(ctx, peerID, "text.process", textReq)
	if err != nil {
		return fmt.Errorf("text process test failed: %v", err)
	}

	if !response.Success {
		return fmt.Errorf("text service returned error: %s", response.Error)
	}

	var textResp services.TextProcessResponse
	if err := json.Unmarshal(response.Result, &textResp); err != nil {
		return fmt.Errorf("failed to unmarshal text response: %v", err)
	}

	log.Printf("Text processing result: '%s' -> '%s' (operation: %s)\n",
		textResp.Original, textResp.Processed, textResp.Operation)
	return nil
}

// AutoTestServices automatically tests services when a new peer is discovered
func (c *ServiceClient) AutoTestServices(peerID peer.ID) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Printf("Auto-testing services on newly discovered peer: %s\n", peerID.String())

	// Test echo service
	if err := c.TestEcho(ctx, peerID, "Hello from "+c.host.ID().String()); err != nil {
		log.Printf("Echo test failed: %v\n", err)
	}

	// Test text processing
	if err := c.TestTextProcess(ctx, peerID, "Hello World", "uppercase"); err != nil {
		log.Printf("Text process test failed: %v\n", err)
	}
}

// FormatPeerID returns a shortened version of peer ID for logging
func FormatPeerID(peerID peer.ID) string {
	str := peerID.String()
	// if len(str) > 12 {
	// 	return str[:8] + "..." + str[len(str)-4:]
	// }
	return str
}

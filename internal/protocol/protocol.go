package protocol

import (
	"bufio"
	"encoding/json"
	"log"

	host "github.com/libp2p/go-libp2p/core/host"
	network "github.com/libp2p/go-libp2p/core/network"
	protocol "github.com/libp2p/go-libp2p/core/protocol"
	"github.com/realentity/realentity-node/internal/services"
)

// Legacy Request struct for backward compatibility
type Request struct {
	Service string `json:"service"`
	Payload string `json:"payload"`
}

func HandleStream(stream network.Stream) {
	log.Println("New stream opened")
	defer stream.Close()

	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	// Try to parse as ServiceRequest first, then fall back to legacy Request
	var serviceReq services.ServiceRequest
	decoder := json.NewDecoder(rw)

	if err := decoder.Decode(&serviceReq); err != nil {
		log.Println("Invalid request:", err)
		errorResponse := services.ServiceResponse{
			Success: false,
			Error:   "Invalid request format",
		}
		json.NewEncoder(rw).Encode(errorResponse)
		rw.Flush()
		return
	}

	log.Printf("Received service request: %s (ID: %s)\n", serviceReq.Service, serviceReq.RequestID)

	// Execute the service
	response := services.GlobalRegistry.ExecuteService(&serviceReq)

	// Send response
	if err := json.NewEncoder(rw).Encode(response); err != nil {
		log.Printf("Failed to send response: %v\n", err)
		return
	}

	rw.Flush()
	log.Printf("Response sent for request %s\n", serviceReq.RequestID)
}

func RegisterHandler(h host.Host, protocolID string) {
	h.SetStreamHandler(protocol.ID(protocolID), HandleStream)
	log.Printf("Protocol handler registered for: %s\n", protocolID)
}

package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/realentity/realentity-node/internal/discovery"
	"github.com/realentity/realentity-node/internal/services"
)

// Server represents the HTTP API server
type Server struct {
	host        host.Host
	discovery   *discovery.DiscoveryManager
	port        int
	httpsPort   int
	certFile    string
	keyFile     string
	server      *http.Server
	httpsServer *http.Server
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string          `json:"status"`
	Timestamp time.Time       `json:"timestamp"`
	PeerID    string          `json:"peer_id"`
	Peers     int             `json:"connected_peers"`
	Services  []string        `json:"services"`
	Uptime    string          `json:"uptime"`
	Discovery map[string]bool `json:"discovery"`
	Version   string          `json:"version,omitempty"`
}

// NodeInfoResponse represents detailed node information
type NodeInfoResponse struct {
	PeerID      string          `json:"peer_id"`
	Addresses   []string        `json:"addresses"`
	Peers       []string        `json:"peers"`
	Services    []string        `json:"services"`
	Discovery   map[string]bool `json:"discovery"`
	Protocols   []string        `json:"protocols"`
	Connections int             `json:"connections"`
}

var startTime = time.Now()

// NewServer creates a new HTTP API server
func NewServer(h host.Host, dm *discovery.DiscoveryManager, port int, httpsPort int, certFile, keyFile string) *Server {
	return &Server{
		host:      h,
		discovery: dm,
		port:      port,
		httpsPort: httpsPort,
		certFile:  certFile,
		keyFile:   keyFile,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	mux := http.NewServeMux()

	// Health endpoint
	mux.HandleFunc("/health", s.handleHealth)

	// Node info endpoint
	mux.HandleFunc("/api/node", s.handleNodeInfo)

	// Peers endpoint
	mux.HandleFunc("/api/peers", s.handlePeers)

	// Services endpoint
	mux.HandleFunc("/api/services", s.handleServices)

	// Service execution endpoint
	mux.HandleFunc("/api/services/execute", s.handleServiceExecution)

	// Start HTTP server
	if s.port > 0 {
		s.server = &http.Server{
			Addr:    fmt.Sprintf(":%d", s.port),
			Handler: mux,
		}

		log.Printf("Starting HTTP API server on port %d", s.port)
		log.Printf("Health endpoint: http://localhost:%d/health", s.port)

		go func() {
			if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Printf("HTTP server error: %v", err)
			}
		}()
	}

	// Start HTTPS server if configured
	if s.httpsPort > 0 && s.certFile != "" && s.keyFile != "" {
		s.httpsServer = &http.Server{
			Addr:    fmt.Sprintf(":%d", s.httpsPort),
			Handler: mux,
		}

		log.Printf("Starting HTTPS API server on port %d", s.httpsPort)
		log.Printf("Health endpoint: https://localhost:%d/health", s.httpsPort)

		return s.httpsServer.ListenAndServeTLS(s.certFile, s.keyFile)
	}

	// If only HTTP is configured, block on it
	if s.server != nil {
		return s.server.ListenAndServe()
	}

	return fmt.Errorf("no HTTP or HTTPS server configured")
}

// Stop stops the HTTP server
func (s *Server) Stop() error {
	var err error
	if s.server != nil {
		if e := s.server.Close(); e != nil {
			err = e
		}
	}
	if s.httpsServer != nil {
		if e := s.httpsServer.Close(); e != nil {
			err = e
		}
	}
	return err
}

// handleHealth handles the /health endpoint
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	peers := s.discovery.GetPeers()
	services := services.GlobalRegistry.ListServices()

	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		PeerID:    s.host.ID().String(),
		Peers:     len(peers),
		Services:  services,
		Uptime:    time.Since(startTime).String(),
		Discovery: map[string]bool{
			"mdns":      true, // TODO: Get actual status from discovery manager
			"bootstrap": len(peers) > 0,
			"dht":       false, // TODO: Get actual DHT status
		},
		Version: "1.0.0", // TODO: Get from build info
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// handleNodeInfo handles the /api/node endpoint
func (s *Server) handleNodeInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	peers := s.discovery.GetPeers()
	peerStrings := make([]string, 0, len(peers))
	for peerID := range peers {
		peerStrings = append(peerStrings, peerID.String())
	}

	addresses := make([]string, len(s.host.Addrs()))
	for i, addr := range s.host.Addrs() {
		addresses[i] = fmt.Sprintf("%s/p2p/%s", addr, s.host.ID())
	}

	protocols := s.host.Mux().Protocols()
	protocolStrings := make([]string, len(protocols))
	for i, p := range protocols {
		protocolStrings[i] = string(p)
	}

	response := NodeInfoResponse{
		PeerID:    s.host.ID().String(),
		Addresses: addresses,
		Peers:     peerStrings,
		Services:  services.GlobalRegistry.ListServices(),
		Discovery: map[string]bool{
			"mdns":      true,
			"bootstrap": len(peers) > 0,
			"dht":       false,
		},
		Protocols:   protocolStrings,
		Connections: len(s.host.Network().Conns()),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// handlePeers handles the /api/peers endpoint
func (s *Server) handlePeers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	peers := s.discovery.GetPeers()
	peerInfo := make([]map[string]interface{}, 0, len(peers))

	for peerID, info := range peers {
		peerInfo = append(peerInfo, map[string]interface{}{
			"peer_id":   peerID.String(),
			"last_seen": info.LastSeen,
			"source":    info.Source,
			"status":    info.Status,
			"services":  info.Services,
			"connected": s.host.Network().Connectedness(peerID),
		})
	}

	response := map[string]interface{}{
		"total_peers": len(peers),
		"peers":       peerInfo,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// handleServices handles the /api/services endpoint
func (s *Server) handleServices(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	servicesList := services.GlobalRegistry.ListServices()

	response := map[string]interface{}{
		"total_services": len(servicesList),
		"services":       servicesList,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// ServiceExecutionRequest represents a request to execute a service
type ServiceExecutionRequest struct {
	Service   string      `json:"service"`
	Payload   interface{} `json:"payload"`
	RequestID string      `json:"requestId,omitempty"`
}

// handleServiceExecution handles the /api/services/execute endpoint
func (s *Server) handleServiceExecution(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Only allow POST requests
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Only POST method is allowed",
		})
		return
	}

	// Parse request body
	var execReq ServiceExecutionRequest
	if err := json.NewDecoder(r.Body).Decode(&execReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Invalid request body: %v", err),
		})
		return
	}

	// Validate required fields
	if execReq.Service == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Service name is required",
		})
		return
	}

	// Generate request ID if not provided
	if execReq.RequestID == "" {
		execReq.RequestID = fmt.Sprintf("api-req-%d", time.Now().UnixNano())
	}

	// Convert payload to JSON
	payloadBytes, err := json.Marshal(execReq.Payload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Invalid payload: %v", err),
		})
		return
	}

	// Create service request
	serviceReq := &services.ServiceRequest{
		Service:   execReq.Service,
		RequestID: execReq.RequestID,
		Payload:   json.RawMessage(payloadBytes),
	}

	// Execute service using the original GlobalRegistry
	response := services.GlobalRegistry.ExecuteService(serviceReq)

	// Return response
	if response.Success {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(response)
}

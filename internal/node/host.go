package node

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"

	libp2p "github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p/core/crypto"
	host "github.com/libp2p/go-libp2p/core/host"
	peer "github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

var HostInstance host.Host

// HostConfig contains configuration for libp2p host creation
type HostConfig struct {
	ListenPort   int
	ExternalIP   string
	ForcePort    bool
	EnableRelay  bool
	EnableNATSvc bool
}

// DefaultHostConfig returns sensible defaults for VPS deployment
func DefaultHostConfig() *HostConfig {
	return &HostConfig{
		ListenPort:   4001, // Standard libp2p port
		ExternalIP:   "",   // Auto-detect if empty
		ForcePort:    true, // Force specific port for VPS
		EnableRelay:  true, // Enable relay for NAT traversal
		EnableNATSvc: true, // Enable NAT service
	}
}

// CreateHost creates a libp2p host optimized for VPS deployment
func CreateHost(ctx context.Context) (host.Host, error) {
	return CreateHostWithConfig(ctx, DefaultHostConfig())
}

// CreateHostWithConfig creates a libp2p host with custom configuration
func CreateHostWithConfig(ctx context.Context, config *HostConfig) (host.Host, error) {
	// Generate cryptographic identity
	priv, _, err := crypto.GenerateKeyPair(crypto.Ed25519, -1)
	if err != nil {
		return nil, fmt.Errorf("failed to generate key pair: %v", err)
	}

	// Build listen addresses
	var listenAddrs []multiaddr.Multiaddr

	// Add IPv4 TCP listener
	tcpAddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", config.ListenPort))
	if err != nil {
		return nil, fmt.Errorf("failed to create TCP multiaddr: %v", err)
	}
	listenAddrs = append(listenAddrs, tcpAddr)

	// Add QUIC listener for better performance
	quicAddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/udp/%d/quic-v1", config.ListenPort))
	if err != nil {
		log.Printf("Warning: Failed to create QUIC listener: %v", err)
	} else {
		listenAddrs = append(listenAddrs, quicAddr)
	}

	// Configure libp2p options
	opts := []libp2p.Option{
		libp2p.Identity(priv),
		libp2p.ListenAddrs(listenAddrs...),
		libp2p.EnableRelay(), // Important for VPS connectivity
	}

	// Add NAT port mapping for VPS environments
	if config.EnableNATSvc {
		opts = append(opts, libp2p.EnableNATService())
	}

	// Force external address if provided (for VPS with known public IP)
	if config.ExternalIP != "" {
		extAddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", config.ExternalIP, config.ListenPort))
		if err != nil {
			log.Printf("Warning: Invalid external IP %s: %v", config.ExternalIP, err)
		} else {
			opts = append(opts, libp2p.AddrsFactory(func([]multiaddr.Multiaddr) []multiaddr.Multiaddr {
				return []multiaddr.Multiaddr{extAddr}
			}))
			log.Printf("Announcing external address: %s", extAddr)
		}
	}

	// Create the host
	h, err := libp2p.New(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create libp2p host: %v", err)
	}

	HostInstance = h

	log.Printf("Host created: %s", h.ID().String())
	for _, addr := range h.Addrs() {
		log.Printf("Listening on: %s", addr)
	}

	// Log the full multiaddr for bootstrap purposes
	if len(h.Addrs()) > 0 {
		fullAddr := h.Addrs()[0].Encapsulate(multiaddr.StringCast("/p2p/" + h.ID().String()))
		log.Printf("Full multiaddr for bootstrap: %s", fullAddr)

		// Also log external address if configured
		if config.ExternalIP != "" {
			extFullAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d/p2p/%s",
				config.ExternalIP, config.ListenPort, h.ID().String()))
			log.Printf("External multiaddr: %s", extFullAddr)
		}
	}

	return h, nil
}

// CreateVPSHost creates a host specifically optimized for VPS deployment
func CreateVPSHost(ctx context.Context, externalIP string, port int) (host.Host, error) {
	config := &HostConfig{
		ListenPort:   port,
		ExternalIP:   externalIP,
		ForcePort:    true,
		EnableRelay:  true,
		EnableNATSvc: false, // Usually not needed on VPS
	}
	return CreateHostWithConfig(ctx, config)
}

// CreateHostWithPeerID creates a host with a specific peer ID (from private key)
func CreateHostWithPeerID(ctx context.Context, peerIDStr string) (host.Host, error) {
	config := DefaultHostConfig()

	// Decode the peer ID to get the private key
	priv, err := DecodePrivateKeyFromPeerID(peerIDStr)
	if err != nil {
		return nil, fmt.Errorf("failed to decode peer ID: %v", err)
	}

	return CreateHostWithPrivateKey(ctx, config, priv)
}

// CreateHostWithPrivateKey creates a host with a specific private key
func CreateHostWithPrivateKey(ctx context.Context, config *HostConfig, priv crypto.PrivKey) (host.Host, error) {
	// Build listen addresses
	var listenAddrs []multiaddr.Multiaddr

	// Add IPv4 TCP listener
	tcpAddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", config.ListenPort))
	if err != nil {
		return nil, fmt.Errorf("failed to create TCP multiaddr: %v", err)
	}
	listenAddrs = append(listenAddrs, tcpAddr)

	// Add QUIC listener for better performance
	quicAddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/udp/%d/quic-v1", config.ListenPort))
	if err != nil {
		log.Printf("Warning: Failed to create QUIC listener: %v", err)
	} else {
		listenAddrs = append(listenAddrs, quicAddr)
	}

	// Configure libp2p options
	opts := []libp2p.Option{
		libp2p.Identity(priv),
		libp2p.ListenAddrs(listenAddrs...),
		libp2p.EnableRelay(), // Important for VPS connectivity
		libp2p.EnableNATService(),
	}

	// Force external address if provided (for VPS with known public IP)
	if config.ExternalIP != "" {
		extAddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", config.ExternalIP, config.ListenPort))
		if err != nil {
			log.Printf("Warning: Invalid external IP %s: %v", config.ExternalIP, err)
		} else {
			opts = append(opts, libp2p.AddrsFactory(func([]multiaddr.Multiaddr) []multiaddr.Multiaddr {
				return []multiaddr.Multiaddr{extAddr}
			}))
			log.Printf("Announcing external address: %s", extAddr)
		}
	}

	// Create the host
	h, err := libp2p.New(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create libp2p host: %v", err)
	}

	log.Printf("Host created: %s\n", h.ID().String())
	for _, addr := range h.Addrs() {
		log.Printf("Listening on: %s\n", addr)
	}

	// Log the full multiaddr for bootstrap purposes
	if config.ExternalIP != "" {
		log.Printf("Full multiaddr for bootstrap: /ip4/%s/tcp/%d/p2p/%s\n",
			config.ExternalIP, config.ListenPort, h.ID().String())
	}

	return h, nil
}

// GeneratePrivateKeyBase64 generates a new private key and returns it as base64
func GeneratePrivateKeyBase64() (string, peer.ID, error) {
	priv, _, err := crypto.GenerateKeyPair(crypto.Ed25519, -1)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate key pair: %v", err)
	}

	// Get peer ID
	peerID, err := peer.IDFromPrivateKey(priv)
	if err != nil {
		return "", "", fmt.Errorf("failed to get peer ID: %v", err)
	}

	// Encode private key to base64
	privBytes, err := crypto.MarshalPrivateKey(priv)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal private key: %v", err)
	}

	privBase64 := base64.StdEncoding.EncodeToString(privBytes)
	return privBase64, peerID, nil
}

// DecodePrivateKeyFromBase64 decodes a base64 private key
func DecodePrivateKeyFromBase64(privBase64 string) (crypto.PrivKey, error) {
	privBytes, err := base64.StdEncoding.DecodeString(privBase64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %v", err)
	}

	priv, err := crypto.UnmarshalPrivateKey(privBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal private key: %v", err)
	}

	return priv, nil
}

// DecodePrivateKeyFromPeerID is a legacy function - we need the private key, not peer ID
func DecodePrivateKeyFromPeerID(peerIDStr string) (crypto.PrivKey, error) {
	return nil, fmt.Errorf("cannot derive private key from peer ID - please use private_key field instead")
}

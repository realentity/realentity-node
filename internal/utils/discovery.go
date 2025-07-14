package utils

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

// DiscoveryUtils provides utilities for discovery management
type DiscoveryUtils struct {
	host host.Host
}

// NewDiscoveryUtils creates a new discovery utilities instance
func NewDiscoveryUtils(h host.Host) *DiscoveryUtils {
	return &DiscoveryUtils{
		host: h,
	}
}

// GetMultiaddr returns the first multiaddr for this node
func (du *DiscoveryUtils) GetMultiaddr() string {
	addrs := du.host.Addrs()
	if len(addrs) == 0 {
		return ""
	}

	// Find a good public address (prefer non-localhost)
	for _, addr := range addrs {
		addrStr := addr.String()
		if !strings.Contains(addrStr, "127.0.0.1") && !strings.Contains(addrStr, "::1") {
			return fmt.Sprintf("%s/p2p/%s", addrStr, du.host.ID().String())
		}
	}

	// Fallback to first address
	return fmt.Sprintf("%s/p2p/%s", addrs[0].String(), du.host.ID().String())
}

// ConnectToPeer manually connects to a peer given their multiaddr
func (du *DiscoveryUtils) ConnectToPeer(ctx context.Context, addrStr string) error {
	addr, err := multiaddr.NewMultiaddr(addrStr)
	if err != nil {
		return fmt.Errorf("invalid multiaddr: %v", err)
	}

	addrInfo, err := peer.AddrInfoFromP2pAddr(addr)
	if err != nil {
		return fmt.Errorf("invalid peer address: %v", err)
	}

	log.Printf("Attempting to connect to peer: %s\n", FormatPeerID(addrInfo.ID))

	if err := du.host.Connect(ctx, *addrInfo); err != nil {
		return fmt.Errorf("connection failed: %v", err)
	}

	log.Printf("Successfully connected to peer: %s\n", FormatPeerID(addrInfo.ID))

	return nil
}

// ListConnectedPeers lists all currently connected peers
func (du *DiscoveryUtils) ListConnectedPeers() {
	peers := du.host.Network().Peers()

	if len(peers) == 0 {
		log.Println("No connected peers")
		return
	}

	log.Printf("Connected peers (%d total):\n", len(peers))
	for _, peerID := range peers {
		addrs := du.host.Peerstore().Addrs(peerID)
		log.Printf("%s (%d addresses)\n", FormatPeerID(peerID), len(addrs))
	}
}

// TestAllConnectedPeers tests services on all connected peers
func (du *DiscoveryUtils) TestAllConnectedPeers() {
	peers := du.host.Network().Peers()

	if len(peers) == 0 {
		log.Println("No connected peers to test")
		return
	}

	log.Printf("Testing services on %d connected peers...\n", len(peers))

	client := NewServiceClient(du.host)

	for _, peerID := range peers {
		go func(id peer.ID) {
			log.Printf("Testing services on peer: %s\n", FormatPeerID(id))
			client.AutoTestServices(id)
		}(peerID)
	}
}

// GetConnectionInfo returns connection information for sharing
func (du *DiscoveryUtils) GetConnectionInfo() map[string]interface{} {
	return map[string]interface{}{
		"peer_id":   du.host.ID().String(),
		"multiaddr": du.GetMultiaddr(),
		"services":  []string{"echo", "text.process"}, // Could be dynamic
		"node_info": map[string]interface{}{
			"version":    "1.0.0",
			"node_type":  "realentity-node",
			"started_at": time.Now().Format(time.RFC3339),
		},
	}
}

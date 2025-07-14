package discovery

import (
	"context"
	"log"
	"time"

	host "github.com/libp2p/go-libp2p/core/host"
	peer "github.com/libp2p/go-libp2p/core/peer"
	"github.com/realentity/realentity-node/internal/utils"
)

// EnhancedNotifee handles peer discovery with the discovery manager
type EnhancedNotifee struct {
	host             host.Host
	client           *utils.ServiceClient
	discoveryManager *DiscoveryManager
}

func (n *EnhancedNotifee) HandlePeerFound(pi peer.AddrInfo) {
	// Don't connect to ourselves
	if pi.ID == n.host.ID() {
		return
	}

	log.Printf("Discovered peer: %s\n", utils.FormatPeerID(pi.ID))

	// Connect to the peer
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := n.host.Connect(ctx, pi); err != nil {
		log.Printf("Failed to connect to peer %s: %v\n", utils.FormatPeerID(pi.ID), err)

		// Update peer status in discovery manager
		if n.discoveryManager != nil {
			n.discoveryManager.peerStore.UpdatePeerStatus(pi.ID, PeerStatusUnreachable, err)
		}
		return
	}

	log.Printf("Connected to peer: %s\n", utils.FormatPeerID(pi.ID))

	// Update peer status in discovery manager
	if n.discoveryManager != nil {
		n.discoveryManager.peerStore.UpdatePeerStatus(pi.ID, PeerStatusConnected, nil)
	}

	// Auto-test services on the new peer
	go n.client.AutoTestServices(pi.ID)
}

// SetupEnhancedMDNS sets up mDNS with the discovery manager integration
func SetupEnhancedMDNS(ctx context.Context, h host.Host, serviceTag string, dm *DiscoveryManager) error {
	client := utils.NewServiceClient(h)
	notifee := &EnhancedNotifee{
		host:             h,
		client:           client,
		discoveryManager: dm,
	}

	// Create and add mDNS discovery mechanism
	mdnsDisc := NewMDNSDiscovery(h, serviceTag)
	dm.AddMechanism(mdnsDisc)

	// Set up the callback to handle found peers
	dm.SetPeerFoundCallback(func(addrInfo peer.AddrInfo) {
		notifee.HandlePeerFound(addrInfo)
	})

	log.Printf("Enhanced mDNS discovery configured with service tag: %s\n", serviceTag)
	return nil
}

package discovery

import (
	"context"
	"log"
	"time"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

// BootstrapDiscovery finds peers through bootstrap nodes
type BootstrapDiscovery struct {
	host           host.Host
	bootstrapPeers []peer.AddrInfo
	connected      map[peer.ID]bool
}

// NewBootstrapDiscovery creates a new bootstrap discovery mechanism
func NewBootstrapDiscovery(h host.Host, bootstrapAddrs []string) (*BootstrapDiscovery, error) {
	var bootstrapPeers []peer.AddrInfo

	for _, addrStr := range bootstrapAddrs {
		addr, err := multiaddr.NewMultiaddr(addrStr)
		if err != nil {
			log.Printf("Invalid bootstrap address %s: %v\n", addrStr, err)
			continue
		}

		addrInfo, err := peer.AddrInfoFromP2pAddr(addr)
		if err != nil {
			log.Printf("Invalid bootstrap peer address %s: %v\n", addrStr, err)
			continue
		}

		bootstrapPeers = append(bootstrapPeers, *addrInfo)
	}

	return &BootstrapDiscovery{
		host:           h,
		bootstrapPeers: bootstrapPeers,
		connected:      make(map[peer.ID]bool),
	}, nil
}

func (bd *BootstrapDiscovery) Name() string {
	return "bootstrap"
}

func (bd *BootstrapDiscovery) Start(ctx context.Context) error {
	log.Printf("Bootstrap discovery starting with %d bootstrap peers\n", len(bd.bootstrapPeers))

	// Connect to bootstrap peers immediately
	go bd.connectToBootstrapPeers(ctx)

	return nil
}

func (bd *BootstrapDiscovery) Stop() error {
	log.Println("Bootstrap discovery stopped")
	return nil
}

func (bd *BootstrapDiscovery) FindPeers(ctx context.Context, limit int) ([]peer.AddrInfo, error) {
	var found []peer.AddrInfo

	// First, ensure we're connected to bootstrap peers
	bd.connectToBootstrapPeers(ctx)

	// Get peers from the host's peerstore (peers we've learned about)
	peers := bd.host.Peerstore().Peers()

	for _, peerID := range peers {
		if peerID == bd.host.ID() {
			continue // Skip ourselves
		}

		if len(found) >= limit {
			break
		}

		addrs := bd.host.Peerstore().Addrs(peerID)
		if len(addrs) > 0 {
			found = append(found, peer.AddrInfo{
				ID:    peerID,
				Addrs: addrs,
			})
		}
	}

	return found, nil
}

func (bd *BootstrapDiscovery) connectToBootstrapPeers(ctx context.Context) {
	for _, addrInfo := range bd.bootstrapPeers {
		if bd.connected[addrInfo.ID] {
			continue // Already connected
		}

		go func(ai peer.AddrInfo) {
			// Retry logic for VPS environments where connections might be unstable
			maxRetries := 3
			backoff := time.Second

			for attempt := 0; attempt < maxRetries; attempt++ {
				connectCtx, cancel := context.WithTimeout(ctx, 15*time.Second)

				if err := bd.host.Connect(connectCtx, ai); err != nil {
					cancel()
					log.Printf("Failed to connect to bootstrap peer %s (attempt %d/%d): %v",
						ai.ID.String(), attempt+1, maxRetries, err)

					if attempt < maxRetries-1 {
						time.Sleep(backoff)
						backoff *= 2 // Exponential backoff
						continue
					}
					return
				}

				cancel()
				bd.connected[ai.ID] = true
				log.Printf("Connected to bootstrap peer: %s (attempt %d)", ai.ID.String(), attempt+1)

				// Verify connection is stable by checking network reachability
				go bd.monitorBootstrapConnection(ai.ID)
				return
			}
		}(addrInfo)
	}
}

// monitorBootstrapConnection monitors bootstrap peer connections and reconnects if needed
func (bd *BootstrapDiscovery) monitorBootstrapConnection(peerID peer.ID) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Check if peer is still connected
			connectedness := bd.host.Network().Connectedness(peerID)
			if connectedness != network.Connected {
				log.Printf("Lost connection to bootstrap peer %s, attempting reconnect...", peerID)
				bd.connected[peerID] = false

				// Find the peer info and reconnect
				for _, ai := range bd.bootstrapPeers {
					if ai.ID == peerID {
						go func() {
							ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
							defer cancel()

							if err := bd.host.Connect(ctx, ai); err == nil {
								bd.connected[ai.ID] = true
								log.Printf("Reconnected to bootstrap peer: %s", ai.ID)
							}
						}()
						break
					}
				}
			}
		}
	}
}

// GetBootstrapPeers returns the list of bootstrap peer addresses
func (bd *BootstrapDiscovery) GetBootstrapPeers() []peer.AddrInfo {
	return bd.bootstrapPeers
}

// AddBootstrapPeer adds a new bootstrap peer
func (bd *BootstrapDiscovery) AddBootstrapPeer(addrInfo peer.AddrInfo) {
	bd.bootstrapPeers = append(bd.bootstrapPeers, addrInfo)
	log.Printf("Added bootstrap peer: %s\n", addrInfo.ID.String())
}

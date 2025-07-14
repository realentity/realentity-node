package discovery

import (
	"context"
	"log"
	"time"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	discovery "github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

// MDNSDiscovery wraps mDNS discovery as a DiscoveryMechanism
type MDNSDiscovery struct {
	host       host.Host
	service    discovery.Service
	serviceTag string
	found      chan peer.AddrInfo
}

// MDNSNotifee handles mDNS peer discovery notifications
type MDNSNotifee struct {
	found chan peer.AddrInfo
	host  host.Host
}

func (n *MDNSNotifee) HandlePeerFound(pi peer.AddrInfo) {
	// Don't notify about ourselves
	if pi.ID == n.host.ID() {
		return
	}

	select {
	case n.found <- pi:
	default:
		// Channel full, skip this peer
	}
}

// NewMDNSDiscovery creates a new mDNS discovery mechanism
func NewMDNSDiscovery(h host.Host, serviceTag string) *MDNSDiscovery {
	found := make(chan peer.AddrInfo, 100)

	notifee := &MDNSNotifee{
		found: found,
		host:  h,
	}

	service := discovery.NewMdnsService(h, serviceTag, notifee)

	return &MDNSDiscovery{
		host:       h,
		service:    service,
		serviceTag: serviceTag,
		found:      found,
	}
}

func (md *MDNSDiscovery) Name() string {
	return "mdns"
}

func (md *MDNSDiscovery) Start(ctx context.Context) error {
	log.Printf("Starting mDNS discovery with service tag: %s\n", md.serviceTag)
	return md.service.Start()
}

func (md *MDNSDiscovery) Stop() error {
	if md.service != nil {
		return md.service.Close()
	}
	return nil
}

func (md *MDNSDiscovery) FindPeers(ctx context.Context, limit int) ([]peer.AddrInfo, error) {
	var found []peer.AddrInfo
	timeout := time.After(5 * time.Second)

	for len(found) < limit {
		select {
		case peer := <-md.found:
			found = append(found, peer)
		case <-timeout:
			return found, nil
		case <-ctx.Done():
			return found, ctx.Err()
		}
	}

	return found, nil
}

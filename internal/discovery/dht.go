package discovery

import (
	"context"
	"log"
	"time"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/routing"
	"github.com/libp2p/go-libp2p/p2p/discovery/util"
)

// DHTDiscovery uses DHT for peer discovery
type DHTDiscovery struct {
	host        host.Host
	dht         *dht.IpfsDHT
	routingDisc *routing.RoutingDiscovery
	rendezvous  string
}

// NewDHTDiscovery creates a new DHT-based discovery mechanism
func NewDHTDiscovery(h host.Host, rendezvous string, bootstrapPeers []peer.AddrInfo) (*DHTDiscovery, error) {
	// Create DHT
	kademliaDHT, err := dht.New(context.Background(), h)
	if err != nil {
		return nil, err
	}

	// Create routing discovery
	routingDiscovery := routing.NewRoutingDiscovery(kademliaDHT)

	dd := &DHTDiscovery{
		host:        h,
		dht:         kademliaDHT,
		routingDisc: routingDiscovery,
		rendezvous:  rendezvous,
	}

	return dd, nil
}

func (dd *DHTDiscovery) Name() string {
	return "dht"
}

func (dd *DHTDiscovery) Start(ctx context.Context) error {
	log.Println("Starting DHT discovery...")

	// Bootstrap the DHT
	if err := dd.dht.Bootstrap(ctx); err != nil {
		return err
	}

	// Start advertising our presence
	go dd.advertise(ctx)

	log.Printf("DHT discovery started with rendezvous: %s\n", dd.rendezvous)
	return nil
}

func (dd *DHTDiscovery) Stop() error {
	if dd.dht != nil {
		return dd.dht.Close()
	}
	return nil
}

func (dd *DHTDiscovery) FindPeers(ctx context.Context, limit int) ([]peer.AddrInfo, error) {
	log.Printf("Searching for peers via DHT (rendezvous: %s)\n", dd.rendezvous)

	peerChan, err := dd.routingDisc.FindPeers(ctx, dd.rendezvous)
	if err != nil {
		return nil, err
	}

	var found []peer.AddrInfo
	timeout := time.After(10 * time.Second)

	for {
		select {
		case peer, ok := <-peerChan:
			if !ok {
				return found, nil
			}
			if peer.ID == dd.host.ID() {
				continue // Skip ourselves
			}
			found = append(found, peer)
			if len(found) >= limit {
				return found, nil
			}

		case <-timeout:
			return found, nil

		case <-ctx.Done():
			return found, ctx.Err()
		}
	}
}

func (dd *DHTDiscovery) advertise(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	// Advertise immediately
	dd.doAdvertise(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			dd.doAdvertise(ctx)
		}
	}
}

func (dd *DHTDiscovery) doAdvertise(ctx context.Context) {
	log.Printf("Advertising presence on DHT (rendezvous: %s)\n", dd.rendezvous)

	// Use util.Advertise for easier advertising
	util.Advertise(ctx, dd.routingDisc, dd.rendezvous)
}

// GetDHT returns the underlying DHT instance
func (dd *DHTDiscovery) GetDHT() *dht.IpfsDHT {
	return dd.dht
}

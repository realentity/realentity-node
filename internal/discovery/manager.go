package discovery

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

// DiscoveryManager coordinates different discovery mechanisms
type DiscoveryManager struct {
	host        host.Host
	mechanisms  []DiscoveryMechanism
	peerStore   *PeerStore
	ctx         context.Context
	cancel      context.CancelFunc
	mu          sync.RWMutex
	onPeerFound func(peer.AddrInfo)
}

// DiscoveryMechanism interface for different discovery methods
type DiscoveryMechanism interface {
	Name() string
	Start(ctx context.Context) error
	Stop() error
	FindPeers(ctx context.Context, limit int) ([]peer.AddrInfo, error)
}

// PeerStore manages discovered peers with metadata
type PeerStore struct {
	peers       map[peer.ID]*PeerInfo
	mu          sync.RWMutex
	maxPeers    int
	cleanupTime time.Duration
}

// PeerInfo contains metadata about discovered peers
type PeerInfo struct {
	AddrInfo     peer.AddrInfo
	LastSeen     time.Time
	Source       string // Which discovery mechanism found this peer
	Status       PeerStatus
	Services     []string
	Reliability  float64 // 0.0 to 1.0
	LastError    error
	ConnectCount int
}

type PeerStatus int

const (
	PeerStatusUnknown PeerStatus = iota
	PeerStatusConnectable
	PeerStatusUnreachable
	PeerStatusConnected
)

// NewDiscoveryManager creates a new discovery manager
func NewDiscoveryManager(h host.Host) *DiscoveryManager {
	ctx, cancel := context.WithCancel(context.Background())

	dm := &DiscoveryManager{
		host:       h,
		mechanisms: make([]DiscoveryMechanism, 0),
		peerStore:  NewPeerStore(1000, 10*time.Minute),
		ctx:        ctx,
		cancel:     cancel,
	}

	return dm
}

// NewPeerStore creates a new peer store
func NewPeerStore(maxPeers int, cleanupTime time.Duration) *PeerStore {
	return &PeerStore{
		peers:       make(map[peer.ID]*PeerInfo),
		maxPeers:    maxPeers,
		cleanupTime: cleanupTime,
	}
}

// AddMechanism adds a discovery mechanism
func (dm *DiscoveryManager) AddMechanism(mechanism DiscoveryMechanism) {
	dm.mu.Lock()
	defer dm.mu.Unlock()
	dm.mechanisms = append(dm.mechanisms, mechanism)
	log.Printf("Added discovery mechanism: %s\n", mechanism.Name())
}

// SetPeerFoundCallback sets the callback for when peers are found
func (dm *DiscoveryManager) SetPeerFoundCallback(callback func(peer.AddrInfo)) {
	dm.mu.Lock()
	defer dm.mu.Unlock()
	dm.onPeerFound = callback
}

// Start begins all discovery mechanisms
func (dm *DiscoveryManager) Start() error {
	dm.mu.RLock()
	mechanisms := make([]DiscoveryMechanism, len(dm.mechanisms))
	copy(mechanisms, dm.mechanisms)
	dm.mu.RUnlock()

	for _, mechanism := range mechanisms {
		if err := mechanism.Start(dm.ctx); err != nil {
			log.Printf("Failed to start discovery mechanism %s: %v\n", mechanism.Name(), err)
			continue
		}
		log.Printf("Started discovery mechanism: %s\n", mechanism.Name())
	}

	// Start periodic discovery
	go dm.periodicDiscovery()

	// Start peer store cleanup
	go dm.peerStore.startCleanup(dm.ctx)

	return nil
}

// Stop stops all discovery mechanisms
func (dm *DiscoveryManager) Stop() error {
	dm.cancel()

	dm.mu.RLock()
	defer dm.mu.RUnlock()

	for _, mechanism := range dm.mechanisms {
		if err := mechanism.Stop(); err != nil {
			log.Printf("Error stopping discovery mechanism %s: %v\n", mechanism.Name(), err)
		}
	}

	return nil
}

// periodicDiscovery runs continuous peer discovery
func (dm *DiscoveryManager) periodicDiscovery() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-dm.ctx.Done():
			return
		case <-ticker.C:
			dm.discoverPeers()
		}
	}
}

// discoverPeers actively searches for peers using all mechanisms
func (dm *DiscoveryManager) discoverPeers() {
	dm.mu.RLock()
	mechanisms := make([]DiscoveryMechanism, len(dm.mechanisms))
	copy(mechanisms, dm.mechanisms)
	dm.mu.RUnlock()

	for _, mechanism := range mechanisms {
		go func(m DiscoveryMechanism) {
			ctx, cancel := context.WithTimeout(dm.ctx, 10*time.Second)
			defer cancel()

			peers, err := m.FindPeers(ctx, 10)
			if err != nil {
				log.Printf("Discovery mechanism %s failed: %v\n", m.Name(), err)
				return
			}

			for _, peerInfo := range peers {
				dm.handleFoundPeer(peerInfo, m.Name())
			}
		}(mechanism)
	}
}

// handleFoundPeer processes a newly found peer
func (dm *DiscoveryManager) handleFoundPeer(addrInfo peer.AddrInfo, source string) {
	// Don't add ourselves
	if addrInfo.ID == dm.host.ID() {
		return
	}

	// Add to peer store
	dm.peerStore.AddPeer(addrInfo, source)

	// Notify callback
	dm.mu.RLock()
	callback := dm.onPeerFound
	dm.mu.RUnlock()

	if callback != nil {
		callback(addrInfo)
	}
}

// HandleFoundPeer makes the handleFoundPeer method accessible
func (dm *DiscoveryManager) HandleFoundPeer(addrInfo peer.AddrInfo, source string) {
	dm.handleFoundPeer(addrInfo, source)
}

// GetPeers returns all known peers
func (dm *DiscoveryManager) GetPeers() map[peer.ID]*PeerInfo {
	return dm.peerStore.GetAllPeers()
}

// GetConnectablePeers returns peers that are likely connectable
func (dm *DiscoveryManager) GetConnectablePeers() []*PeerInfo {
	return dm.peerStore.GetConnectablePeers()
}

// AddPeer adds a peer to the store
func (ps *PeerStore) AddPeer(addrInfo peer.AddrInfo, source string) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	existing, exists := ps.peers[addrInfo.ID]
	if exists {
		// Update existing peer
		existing.LastSeen = time.Now()
		existing.AddrInfo.Addrs = append(existing.AddrInfo.Addrs, addrInfo.Addrs...)
		// Remove duplicates
		existing.AddrInfo.Addrs = removeDuplicateAddrs(existing.AddrInfo.Addrs)
	} else {
		// Add new peer
		ps.peers[addrInfo.ID] = &PeerInfo{
			AddrInfo:     addrInfo,
			LastSeen:     time.Now(),
			Source:       source,
			Status:       PeerStatusUnknown,
			Services:     make([]string, 0),
			Reliability:  0.5, // Start with neutral reliability
			ConnectCount: 0,
		}
	}
}

// GetAllPeers returns all peers
func (ps *PeerStore) GetAllPeers() map[peer.ID]*PeerInfo {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	result := make(map[peer.ID]*PeerInfo)
	for id, info := range ps.peers {
		result[id] = info
	}
	return result
}

// GetConnectablePeers returns peers that are likely to be connectable
func (ps *PeerStore) GetConnectablePeers() []*PeerInfo {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	var connectable []*PeerInfo
	for _, info := range ps.peers {
		if info.Status == PeerStatusConnectable || info.Status == PeerStatusConnected {
			connectable = append(connectable, info)
		}
	}
	return connectable
}

// UpdatePeerStatus updates the status of a peer
func (ps *PeerStore) UpdatePeerStatus(peerID peer.ID, status PeerStatus, err error) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if info, exists := ps.peers[peerID]; exists {
		info.Status = status
		info.LastError = err
		info.LastSeen = time.Now()

		// Update reliability based on connection success/failure
		if status == PeerStatusConnected {
			info.ConnectCount++
			info.Reliability = min(info.Reliability+0.1, 1.0)
		} else if status == PeerStatusUnreachable {
			info.Reliability = max(info.Reliability-0.2, 0.0)
		}
	}
}

// startCleanup periodically removes old/unreliable peers
func (ps *PeerStore) startCleanup(ctx context.Context) {
	ticker := time.NewTicker(ps.cleanupTime)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			ps.cleanup()
		}
	}
}

// cleanup removes old or unreliable peers
func (ps *PeerStore) cleanup() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	cutoff := time.Now().Add(-ps.cleanupTime)
	removed := 0

	for id, info := range ps.peers {
		// Remove peers that haven't been seen recently and have low reliability
		if info.LastSeen.Before(cutoff) && info.Reliability < 0.3 {
			delete(ps.peers, id)
			removed++
		}
	}

	if removed > 0 {
		log.Printf("Cleaned up %d old/unreliable peers\n", removed)
	}
}

// Helper functions
func removeDuplicateAddrs(addrs []multiaddr.Multiaddr) []multiaddr.Multiaddr {
	seen := make(map[string]bool)
	result := make([]multiaddr.Multiaddr, 0)

	for _, addr := range addrs {
		addrStr := addr.String()
		if !seen[addrStr] {
			seen[addrStr] = true
			result = append(result, addr)
		}
	}

	return result
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

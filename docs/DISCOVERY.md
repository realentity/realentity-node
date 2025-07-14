# RealEntity Node Discovery Architecture

##  Overview

The RealEntity Node now supports **multiple discovery mechanisms** for finding other nodes across different network environments. This enhanced discovery system allows nodes to find each other locally (mDNS), through bootstrap peers, and potentially through DHT (future).

## ️ Discovery Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                 DiscoveryManager                            │
├─────────────────────────────────────────────────────────────┤
│  • Coordinates multiple discovery mechanisms               │
│  • Manages peer store with metadata                        │
│  • Handles peer lifecycle (discovery → connection → test)  │
│  • Provides reliability scoring                            │
└─────────────────────┬───────────────────────────────────────┘
                      │
        ┌─────────────┼─────────────┐
        │             │             │
   ┌────▼────┐  ┌─────▼─────┐  ┌───▼───┐
   │  mDNS   │  │ Bootstrap │  │  DHT  │
   │Discovery│  │ Discovery │  │(Future)│
   └─────────┘  └───────────┘  └───────┘
```

##  Discovery Mechanisms

### 1. mDNS Discovery (Local Network)
- **Purpose**: Find peers on the same local network
- **Range**: Local subnet only  
- **Configuration**: `enable_mdns: true`
- **Service Tag**: Configurable (default: "realentity-mdns")

### 2. Bootstrap Discovery (Internet)
- **Purpose**: Connect to known bootstrap peers to discover the wider network
- **Range**: Global (internet)
- **Configuration**: `enable_bootstrap: true` + bootstrap peer list
- **Peer List**: Configurable in `config.json`

### 3. DHT Discovery (Future)
- **Purpose**: Distributed peer discovery using Kademlia DHT
- **Range**: Global (internet)
- **Status**: Framework ready, implementation pending
- **Configuration**: `enable_dht: false` (disabled by default)

## ️ Peer Store Features

### Peer Metadata
Each discovered peer includes:
- **Connection Status**: Unknown, Connectable, Unreachable, Connected
- **Discovery Source**: Which mechanism found this peer
- **Reliability Score**: 0.0 to 1.0 based on connection success
- **Last Seen**: Timestamp of last contact
- **Services**: Available services (future enhancement)
- **Address List**: All known multiaddresses

### Automatic Cleanup
- Removes peers not seen recently with low reliability
- Configurable cleanup intervals
- Prevents memory bloat from stale peers

## ️ Configuration System

### Config File (`config.json`)
```json
{
  "discovery": {
    "enable_mdns": true,
    "enable_bootstrap": true,
    "enable_dht": false,
    "mdns_service_tag": "realentity-mdns",
    "bootstrap_peers": [
      "/ip4/203.0.113.1/tcp/4001/p2p/12D3KooW...",
      "/ip4/198.51.100.1/tcp/4001/p2p/12D3KooW..."
    ],
    "dht_rendezvous": "realentity-dht"
  },
  "log_level": "info"
}
```

### Automatic Config Creation
- Creates default config if none exists
- Sensible defaults for immediate functionality
- Easily customizable for different environments

##  Discovery Workflow

### 1. Startup Phase
```
Node Start → Load Config → Create Discovery Manager → Add Mechanisms → Start Discovery
```

### 2. Peer Discovery Phase
```
Mechanism Finds Peer → Add to Peer Store → Callback → Attempt Connection → Update Status
```

### 3. Service Testing Phase
```
Connection Success → Auto-test Services → Update Reliability → Periodic Re-testing
```

### 4. Maintenance Phase
```
Periodic Discovery → Cleanup Old Peers → Update Reliability Scores → Log Statistics
```

##  Discovery Statistics

The system provides real-time stats:
- **Total Known Peers**: All peers ever discovered
- **Connectable Peers**: Peers likely to accept connections
- **Connected Peers**: Currently active connections
- **Reliability Metrics**: Success rates per peer

## ️ Discovery Utilities

### Manual Connection
```go
utils := NewDiscoveryUtils(host)
err := utils.ConnectToPeer(ctx, "/ip4/192.168.1.100/tcp/4001/p2p/12D3KooW...")
```

### Peer Management
```go
utils.ListConnectedPeers()           // Show connected peers
utils.TestAllConnectedPeers()        // Test services on all peers
info := utils.GetConnectionInfo()    // Get shareable connection info
```

##  Network Topologies Supported

### 1. Local Development
- **mDNS only**: Automatic discovery on same network
- **Perfect for**: Testing, development, local demos

### 2. Internet Deployment
- **Bootstrap + mDNS**: Global discovery with local optimization
- **Perfect for**: Production VPS deployments

### 3. Hybrid Networks
- **All mechanisms**: Maximum discovery coverage
- **Perfect for**: Complex network environments

##  Future Enhancements

### Planned Features
1. **DHT Integration**: Full Kademlia DHT support
2. **Service Registry Protocol**: Advertise available services
3. **Reputation System**: Advanced peer scoring
4. **Network Partitioning**: Handle network splits gracefully
5. **Mobile Discovery**: Optimizations for mobile networks

### Discovery Protocol Extensions
1. **Service Discovery**: "Who has service X?"
2. **Load Balancing**: "Which peer can handle my request?"
3. **Geographic Awareness**: "Find peers near me"
4. **Capability Matching**: "Find peers with specific capabilities"

##  Performance Characteristics

### mDNS Discovery
- **Latency**: < 1 second (local network)
- **Overhead**: Minimal multicast traffic
- **Reliability**: High (same network)

### Bootstrap Discovery
- **Latency**: 2-10 seconds (depends on bootstrap peer)
- **Overhead**: Initial connection cost
- **Reliability**: Medium (depends on bootstrap peer availability)

### Expected DHT Performance
- **Latency**: 5-30 seconds (network-dependent)
- **Overhead**: Moderate (DHT maintenance)
- **Reliability**: High (distributed, no single point of failure)

##  Usage Examples

### Adding Bootstrap Peers
1. Edit `config.json` to add bootstrap peer addresses
2. Restart node or use dynamic configuration updates
3. Monitor logs for bootstrap connection attempts

### Manual Peer Connection
1. Get peer's multiaddress: `utils.GetMultiaddr()`
2. Share address with other nodes
3. Connect manually: `utils.ConnectToPeer(ctx, addr)`

### Monitoring Discovery
1. Watch startup logs for active mechanisms
2. Check periodic discovery statistics
3. Use utility functions to inspect peer state

##  Troubleshooting

### Common Issues
1. **mDNS warnings**: Normal on some network configurations
2. **No bootstrap peers**: Add valid bootstrap addresses to config
3. **Connection failures**: Check firewall settings and peer addresses

### Debug Commands
```bash
# Check current peers
go run . # Monitor logs for peer discovery

# Test connectivity
# Use utility functions in code to manually test connections
```

---

This discovery system provides a solid foundation for building a truly decentralized peer-to-peer network that can operate effectively across different network environments and scale from local development to global deployment.

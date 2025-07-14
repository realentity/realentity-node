#  VPS Deployment Guide for RealEntity Nodes

##  Overview

When deploying RealEntity nodes across different VPS servers, discovery works through **Bootstrap Discovery** - nodes connect to known peer addresses across the internet.

## ️ Architecture for VPS Deployment

```
VPS 1 (Bootstrap Node)          VPS 2 (Peer Node)              VPS 3 (Peer Node)
┌─────────────────────┐        ┌─────────────────────┐        ┌─────────────────────┐
│ realentity-node     │◄──────►│ realentity-node     │◄──────►│ realentity-node     │
│ 198.51.100.1:4001  │        │ 203.0.113.5:4001   │        │ 192.0.2.10:4001    │
│ (Bootstrap)         │        │ (Connects to VPS1)  │        │ (Connects to VPS1)  │
└─────────────────────┘        └─────────────────────┘        └─────────────────────┘
```

##  Step-by-Step Deployment

### Step 1: Deploy Bootstrap Node (First VPS)

1. **Upload your binary** to VPS 1
2. **Run the bootstrap setup**:
   ```bash
   # Linux/Unix
   chmod +x deploy_bootstrap.sh
   ./deploy_bootstrap.sh
   
   # Windows
   deploy_bootstrap.bat
   ```

3. **Start the node**:
   ```bash
   ./realentity-node
   ```

4. **Note the multiaddr** from the logs (example):
   ```
   2025/07/12 Host created: 12D3KooWExample123...
   Listening on: /ip4/198.51.100.1/tcp/4001
   ```
   
   Your bootstrap address will be:
   ```
   /ip4/198.51.100.1/tcp/4001/p2p/12D3KooWExample123...
   ```

### Step 2: Deploy Peer Nodes (Additional VPS)

1. **Upload your binary** to VPS 2, 3, etc.
2. **Run the peer setup** with bootstrap address:
   ```bash
   # Linux/Unix
   chmod +x deploy_peer.sh
   ./deploy_peer.sh "/ip4/198.51.100.1/tcp/4001/p2p/12D3KooWExample123..."
   
   # Windows
   deploy_peer.bat "/ip4/198.51.100.1/tcp/4001/p2p/12D3KooWExample123..."
   ```

3. **Start the node**:
   ```bash
   ./realentity-node
   ```

##  How Discovery Works Across VPS

### 1. **Bootstrap Connection Process**
```
Peer Node                    Bootstrap Node
    │                             │
    │──── TCP Connection ────────▶│ (Initial contact)
    │◄─── Peer Exchange ─────────│ (Share known peers)
    │──── Service Discovery ────▶│ (Discover services)
    │◄─── Service Testing ──────│ (Auto-test services)
```

### 2. **Peer-to-Peer Mesh Formation**
- Once connected to bootstrap, nodes learn about other peers
- Direct peer-to-peer connections form between all nodes
- Services can be called directly between any two nodes

### 3. **Service Discovery Flow**
```
VPS 1: [echo, text.process]
  │
  ├─ VPS 2: [echo, text.process, ai.infer]
  │
  └─ VPS 3: [echo, image.resize, data.process]
```

##  Configuration Examples

### Bootstrap Node Config (VPS 1)
```json
{
  "discovery": {
    "enable_mdns": false,        ← Disabled for VPS
    "enable_bootstrap": true,     ← Enabled to accept connections
    "enable_dht": false,
    "bootstrap_peers": [],       ← Empty for bootstrap node
    "mdns_quiet_mode": true
  },
  "log_level": "info"
}
```

### Peer Node Config (VPS 2, 3, ...)
```json
{
  "discovery": {
    "enable_mdns": false,
    "enable_bootstrap": true,
    "enable_dht": false,
    "bootstrap_peers": [
      "/ip4/198.51.100.1/tcp/4001/p2p/12D3KooWExample123..."
    ],
    "mdns_quiet_mode": true
  },
  "log_level": "info"
}
```

##  Network Requirements

### **Firewall Rules**
Each VPS needs these ports open:
- **TCP 4001** (or whatever port libp2p chooses)
- **UDP ports** for QUIC (libp2p will log the actual ports)

### **Public IP Addresses**
- Nodes need to be reachable on their public IPs
- libp2p automatically detects and announces public addresses

##  Expected Behavior

### **Bootstrap Node Logs**
```
2025/07/12 Node started with ID: 12D3KooW...Bootstrap
2025/07/12 Listening on: /ip4/198.51.100.1/tcp/4001
2025/07/12 Node is ready! Registered services: [echo text.process]
2025/07/12 Discovery mechanisms active:
2025/07/12   - mDNS: false
2025/07/12   - Bootstrap: true (0 peers)
2025/07/12   - DHT: false
2025/07/12 Waiting for connections...

// When peer connects:
2025/07/12 Discovered peer: 12D3KooW...Peer2
2025/07/12 Connected to peer: 12D3KooW...Peer2
2025/07/12 Auto-testing services on newly discovered peer...
```

### **Peer Node Logs**
```
2025/07/12 Node started with ID: 12D3KooW...Peer2
2025/07/12 Listening on: /ip4/203.0.113.5/tcp/4001
2025/07/12 Node is ready! Registered services: [echo text.process]
2025/07/12 Discovery mechanisms active:
2025/07/12   - mDNS: false
2025/07/12   - Bootstrap: true (1 peers)
2025/07/12   - DHT: false
2025/07/12 Connecting to bootstrap peers...
2025/07/12 Connected to bootstrap peer: 12D3KooW...Bootstrap
2025/07/12 Auto-testing services on bootstrap peer...
```

##  Scaling to Many Nodes

### **Multiple Bootstrap Nodes**
For redundancy, configure multiple bootstrap peers:
```json
{
  "bootstrap_peers": [
    "/ip4/198.51.100.1/tcp/4001/p2p/12D3KooW...Bootstrap1",
    "/ip4/203.0.113.10/tcp/4001/p2p/12D3KooW...Bootstrap2"
  ]
}
```

### **DHT for Large Networks**
When you have 10+ nodes, enable DHT discovery:
```json
{
  "enable_dht": true,
  "dht_rendezvous": "realentity-global"
}
```

##  Testing Your Deployment

1. **Deploy bootstrap node** on VPS 1
2. **Deploy 2-3 peer nodes** on different VPS
3. **Watch the logs** for automatic discovery and service testing
4. **Verify peer mesh** by checking discovery stats every 60 seconds

Your nodes should automatically find each other and start testing services within seconds of connection!

##  Troubleshooting

- **Connection failures**: Check firewall rules and public IP accessibility
- **No peer discovery**: Verify bootstrap peer address format
- **Service test failures**: Ensure protocol handlers are properly registered

The beauty of this system is that once the bootstrap connection is made, the network becomes self-organizing and resilient!
